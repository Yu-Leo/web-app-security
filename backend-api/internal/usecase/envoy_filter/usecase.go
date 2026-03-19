package envoy_filter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"

	authpb "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/Yu-Leo/web-app-security/backend-api/internal/models"
	"github.com/Yu-Leo/web-app-security/backend-api/internal/monitoring"
)

const (
	actionAllow = "allow"
	actionBlock = "block"
)

var errNoResources = errors.New("no resources configured")
var errNoResourceMatch = errors.New("no matching resource")

type UsecaseImpl struct {
	resources        ResourceRepository
	securityProfiles SecurityProfileRepository
	trafficProfiles  TrafficProfileRepository
	securityRules    SecurityRuleRepository
	trafficRules     TrafficRuleRepository
	mlModels         MLModelRepository
	mlScorer         MLScorer
	requestLogs      RequestLogRepository
	eventLogs        EventLogRepository
	rateLimiter      *InMemoryRateLimiter
}

func NewUsecase(
	resources ResourceRepository,
	securityProfiles SecurityProfileRepository,
	trafficProfiles TrafficProfileRepository,
	securityRules SecurityRuleRepository,
	trafficRules TrafficRuleRepository,
	mlModels MLModelRepository,
	mlScorer MLScorer,
	requestLogs RequestLogRepository,
	eventLogs EventLogRepository,
) *UsecaseImpl {
	if mlScorer == nil {
		mlScorer = NewUnavailableMLScorer()
	}

	return &UsecaseImpl{
		resources:        resources,
		securityProfiles: securityProfiles,
		trafficProfiles:  trafficProfiles,
		securityRules:    securityRules,
		trafficRules:     trafficRules,
		mlModels:         mlModels,
		mlScorer:         mlScorer,
		requestLogs:      requestLogs,
		eventLogs:        eventLogs,
		rateLimiter:      NewInMemoryRateLimiter(),
	}
}

func (u *UsecaseImpl) Check(ctx context.Context, req *authpb.CheckRequest) (*Decision, error) {
	startedAt := time.Now()
	var finalDecision *Decision
	defer func() {
		if finalDecision == nil {
			monitoring.ObserveEnvoyRequest(nil, time.Since(startedAt))
			return
		}
		allowed := finalDecision.Allowed
		monitoring.ObserveEnvoyRequest(&allowed, time.Since(startedAt))
	}()

	requestCtx := buildRequestContext(req)
	requestCtx.MLFeatureVector = extractMLFeatureVector(requestCtx)

	resource, resourceMatched, err := u.matchResource(ctx, requestCtx.Path)
	if err != nil {
		if errors.Is(err, errNoResources) || errors.Is(err, errNoResourceMatch) {
			reason := "no_resources_configured"
			if errors.Is(err, errNoResourceMatch) {
				reason = "no_matching_resource"
			}
			decision := &Decision{
				Allowed:    true,
				Action:     actionAllow,
				StatusCode: 200,
				DecidedAt:  time.Now(),
				Reason:     strPtr(reason),
			}
			u.persistLogs(ctx, nil, false, requestCtx, decision)
			finalDecision = decision
			return decision, nil
		}
		return nil, err
	}

	decision := &Decision{
		Allowed:    true,
		Action:     actionAllow,
		StatusCode: 200,
		DecidedAt:  time.Now(),
	}

	if resource.SecurityProfileID != nil {
		securityProfile, err := u.securityProfiles.Get(ctx, *resource.SecurityProfileID)
		if err != nil {
			return nil, err
		}
		securityEval := u.applySecurityRules(ctx, securityProfile, requestCtx)
		if securityEval.Rule != nil {
			decision.RuleID = &securityEval.Rule.ID
			decision.ProfileID = &securityEval.Rule.ProfileID
			decision.MatchedRule = &securityEval.Rule.Name
			kind := "security"
			decision.MatchedRuleType = &kind
			decision.DryRun = securityEval.DryRun
		}

		decision.DecisionSource = strPtr(securityEval.DecisionSource)
		decision.MLModelID = securityEval.MLModelID
		decision.MLModelName = securityEval.MLModelName
		decision.MLThreshold = securityEval.MLThreshold
		decision.MLScore = securityEval.MLScore
		decision.MLScorePercent = securityEval.MLScorePercent
		decision.MLError = securityEval.MLError

		if securityEval.Action == actionBlock {
			if securityEval.DryRun {
				decision.Action = actionAllow
				decision.Reason = strPtr("security_rule_block_dry_run")
			} else {
				decision.Allowed = false
				decision.Action = actionBlock
				decision.StatusCode = 403
				if decision.Reason == nil {
					decision.Reason = strPtr("security_rule_block")
				}
			}
		}
	}

	if decision.Allowed && resource.TrafficProfileID != nil {
		trafficProfile, err := u.trafficProfiles.Get(ctx, *resource.TrafficProfileID)
		if err != nil {
			return nil, err
		}
		trafficRule, trafficAction, trafficDryRun := u.applyTrafficRules(ctx, trafficProfile, resource.ID, requestCtx)
		if trafficRule != nil {
			decision.RuleID = &trafficRule.ID
			decision.ProfileID = &trafficRule.ProfileID
			decision.MatchedRule = &trafficRule.Name
			kind := "traffic"
			decision.MatchedRuleType = &kind
			decision.DryRun = trafficDryRun
		}
		if trafficAction == actionBlock {
			decision.Allowed = false
			decision.Action = actionBlock
			decision.StatusCode = 429
			decision.Reason = strPtr("traffic_rule_limit")
		} else if trafficRule != nil && trafficDryRun {
			decision.Action = actionAllow
			decision.Reason = strPtr("traffic_rule_limit_dry_run")
		}
	}

	u.persistLogs(ctx, &resource, resourceMatched, requestCtx, decision)
	finalDecision = decision

	return decision, nil
}

func (u *UsecaseImpl) persistLogs(
	ctx context.Context,
	resource *models.Resource,
	resourceMatched bool,
	requestCtx RequestContext,
	decision *Decision,
) {
	metadata := map[string]interface{}{
		"resource_matched": resourceMatched,
	}
	if resource != nil {
		if resource.SecurityProfileID != nil {
			metadata["security_profile_id"] = *resource.SecurityProfileID
		}
		if resource.TrafficProfileID != nil {
			metadata["traffic_profile_id"] = *resource.TrafficProfileID
		}
	}
	if decision.MatchedRuleType != nil {
		metadata["matched_rule_type"] = *decision.MatchedRuleType
	}
	if decision.DryRun {
		metadata["dry_run"] = true
	}
	if decision.Reason != nil {
		metadata["reason"] = *decision.Reason
	}
	if decision.DecisionSource != nil {
		metadata["decision_source"] = *decision.DecisionSource
	}
	if decision.MLModelID != nil {
		metadata["ml_model_id"] = *decision.MLModelID
	}
	if decision.MLModelName != nil {
		metadata["ml_model_name"] = *decision.MLModelName
	}
	if decision.MLThreshold != nil {
		metadata["ml_threshold"] = *decision.MLThreshold
	}
	if decision.MLScore != nil {
		metadata["ml_score"] = *decision.MLScore
	}
	if decision.MLScorePercent != nil {
		metadata["ml_score_percent"] = *decision.MLScorePercent
	}
	if decision.MLError != nil {
		metadata["ml_error"] = *decision.MLError
	}
	metadataJSON, _ := json.Marshal(metadata)

	var resourceID *int64
	if resource != nil {
		resourceID = &resource.ID
	}

	_, _ = u.requestLogs.Create(ctx, models.RequestLog{
		ResourceID:         resourceID,
		OccurredAt:         requestCtx.OccurredAt,
		ClientIP:           requestCtx.ClientIP,
		Method:             requestCtx.Method,
		Path:               requestCtx.Path,
		StatusCode:         decision.StatusCode,
		Action:             decision.Action,
		RuleID:             decision.RuleID,
		ProfileID:          decision.ProfileID,
		UserAgent:          nullableString(requestCtx.UserAgent),
		RequestID:          nullableString(requestCtx.RequestID),
		Metadata:           metadataJSON,
		Host:               nullableString(requestCtx.Host),
		Scheme:             nullableString(requestCtx.Scheme),
		Protocol:           nullableString(requestCtx.Protocol),
		Authority:          nullableString(requestCtx.Authority),
		Query:              nullableString(requestCtx.Query),
		SourcePort:         nullableInt32(requestCtx.SourcePort),
		DestinationIP:      nullableString(requestCtx.DestinationIP),
		DestinationPort:    nullableInt32(requestCtx.DestinationPort),
		SourcePrincipal:    nullableString(requestCtx.SourcePrincipal),
		SourceService:      nullableString(requestCtx.SourceService),
		SourceLabels:       mapToJSON(requestCtx.SourceLabels),
		DestinationService: nullableString(requestCtx.DestinationService),
		DestinationLabels:  mapToJSON(requestCtx.DestinationLabels),
		RequestHTTPID:      nullableString(requestCtx.RequestHTTPID),
		Fragment:           nullableString(requestCtx.Fragment),
		RequestHeaders:     headersToJSON(requestCtx.Headers),
		RequestBodySize:    nullableInt32(requestCtx.RequestBodySize),
		RequestBody:        nullableString(requestCtx.RequestBody),
		ContextExtensions:  mapToJSON(requestCtx.ContextExtensions),
		MetadataContext:    requestCtx.MetadataContext,
		RouteMetadataCtx:   requestCtx.RouteMetadataCtx,
	})

	if resource == nil || (decision.Action == actionAllow && !decision.DryRun) {
		return
	}

	eventType := decision.Action
	message := fmt.Sprintf("request %s", decision.Action)
	severity := "info"
	if decision.DryRun {
		eventType = "dry_run"
		message = "request dry_run"
	} else if decision.Action == actionBlock {
		severity = "high"
	}
	_, _ = u.eventLogs.Create(ctx, models.EventLog{
		ResourceID: resource.ID,
		OccurredAt: decision.DecidedAt,
		EventType:  eventType,
		Severity:   severity,
		Message:    message,
		RuleID:     decision.RuleID,
		ProfileID:  decision.ProfileID,
		Metadata:   metadataJSON,
		RequestID:  nullableString(requestCtx.RequestID),
		ClientIP:   nullableString(requestCtx.ClientIP),
		Method:     nullableString(requestCtx.Method),
		Path:       nullableString(requestCtx.Path),
	})
}

func (u *UsecaseImpl) matchResource(ctx context.Context, path string) (models.Resource, bool, error) {
	resources, err := u.resources.List(ctx)
	if err != nil {
		return models.Resource{}, false, err
	}
	if len(resources) == 0 {
		return models.Resource{}, false, errNoResources
	}

	for _, resource := range resources {
		re, err := regexp.Compile(resource.URLPattern)
		if err != nil {
			continue
		}
		if re.MatchString(path) {
			return resource, true, nil
		}
	}

	return models.Resource{}, false, errNoResourceMatch
}

type securityEvaluation struct {
	Rule           *models.SecurityRule
	Action         string
	DryRun         bool
	DecisionSource string
	MLModelID      *int64
	MLModelName    *string
	MLThreshold    *int16
	MLScore        *float32
	MLScorePercent *float32
	MLError        *string
}

func (u *UsecaseImpl) applySecurityRules(ctx context.Context, profile models.SecurityProfile, requestCtx RequestContext) securityEvaluation {
	baseAction := mapBaseAction(profile.BaseAction)
	if !profile.IsEnabled {
		monitoring.IncSecurityBaseAction(profile.ID, actionAllow)
		return securityEvaluation{Action: actionAllow, DecisionSource: "base_action"}
	}

	rules, err := u.securityRules.List(ctx)
	if err != nil {
		monitoring.IncSecurityBaseAction(profile.ID, actionAllow)
		return securityEvaluation{Action: actionAllow, DecisionSource: "base_action"}
	}
	filtered := make([]models.SecurityRule, 0)
	for _, rule := range rules {
		if rule.ProfileID == profile.ID && rule.IsEnabled {
			filtered = append(filtered, rule)
		}
	}
	if len(filtered) == 0 {
		monitoring.IncSecurityBaseAction(profile.ID, baseAction)
		return securityEvaluation{Action: baseAction, DecisionSource: "base_action"}
	}

	sort.Slice(filtered, func(i, j int) bool { return filtered[i].Priority < filtered[j].Priority })
	var lastMLError *string
	for _, rule := range filtered {
		if !securityConditionsMatch(rule.Conditions, requestCtx) {
			continue
		}

		if rule.RuleType == models.SecurityRuleTypeDeterministic {
			monitoring.IncSecurityRuleMatch(rule.ID, rule.ProfileID, string(rule.Action), string(rule.RuleType), rule.DryRun)
			if rule.DryRun {
				monitoring.IncDryRun("security", string(rule.Action))
			}
			return securityEvaluation{
				Rule:           &rule,
				Action:         string(rule.Action),
				DryRun:         rule.DryRun,
				DecisionSource: "rule",
				MLError:        lastMLError,
			}
		}

		if rule.RuleType != models.SecurityRuleTypeML {
			continue
		}

		if rule.MLModelID == nil || rule.MLThreshold == nil {
			lastMLError = strPtr("invalid_ml_rule_config")
			continue
		}

		model, err := u.mlModels.Get(ctx, *rule.MLModelID)
		if err != nil {
			lastMLError = strPtr(fmt.Sprintf("ml_model_load_error: %v", err))
			monitoring.IncMLError(rule.MLModelID, nil)
			continue
		}
		monitoring.IncMLInference(model.ID, model.Name)

		score, err := u.mlScorer.Score(ctx, model.Name, model.ModelData, requestCtx.MLFeatureVector)
		if err != nil {
			lastMLError = strPtr(fmt.Sprintf("ml_inference_error: %v", err))
			modelID := model.ID
			modelName := model.Name
			monitoring.IncMLError(&modelID, &modelName)
			continue
		}
		monitoring.ObserveMLScore(model.ID, model.Name, float64(score))

		scorePercent := score * 100
		modelName := model.Name
		if scorePercent >= float32(*rule.MLThreshold) {
			monitoring.IncMLThresholdPass(model.ID, model.Name)
			monitoring.IncSecurityRuleMatch(rule.ID, rule.ProfileID, string(rule.Action), string(rule.RuleType), rule.DryRun)
			if rule.DryRun {
				monitoring.IncDryRun("security", string(rule.Action))
			}
			return securityEvaluation{
				Rule:           &rule,
				Action:         string(rule.Action),
				DryRun:         rule.DryRun,
				DecisionSource: "ml",
				MLModelID:      rule.MLModelID,
				MLModelName:    &modelName,
				MLThreshold:    rule.MLThreshold,
				MLScore:        &score,
				MLScorePercent: &scorePercent,
				MLError:        lastMLError,
			}
		}
		monitoring.IncMLThresholdMiss(model.ID, model.Name)
	}

	monitoring.IncSecurityBaseAction(profile.ID, baseAction)
	return securityEvaluation{
		Action:         baseAction,
		DecisionSource: "base_action",
		MLError:        lastMLError,
	}
}

func (u *UsecaseImpl) applyTrafficRules(ctx context.Context, profile models.TrafficProfile, resourceID int64, requestCtx RequestContext) (*models.TrafficRule, string, bool) {
	if !profile.IsEnabled {
		return nil, actionAllow, false
	}

	rules, err := u.trafficRules.List(ctx)
	if err != nil {
		return nil, actionAllow, false
	}
	filtered := make([]models.TrafficRule, 0)
	for _, rule := range rules {
		if rule.ProfileID == profile.ID && rule.IsEnabled {
			filtered = append(filtered, rule)
		}
	}
	if len(filtered) == 0 {
		return nil, actionAllow, false
	}

	sort.Slice(filtered, func(i, j int) bool { return filtered[i].Priority < filtered[j].Priority })
	for _, rule := range filtered {
		if !rule.MatchAll && !securityConditionsMatch(rule.Conditions, requestCtx) {
			continue
		}

		key := fmt.Sprintf("res:%d:rule:%d:ip:%s", resourceID, rule.ID, requestCtx.ClientIP)
		period := time.Duration(rule.PeriodSeconds) * time.Second
		allowed := u.rateLimiter.Allow(key, rule.RequestsLimit, period, time.Now())
		if !allowed {
			monitoring.IncTrafficRuleMatch(rule.ID, rule.ProfileID, rule.DryRun)
			if rule.DryRun {
				monitoring.IncDryRun("traffic", actionBlock)
				return &rule, actionAllow, true
			}
			return &rule, actionBlock, false
		}
	}

	return nil, actionAllow, false
}

func conditionsMatch(raw json.RawMessage, requestCtx RequestContext) bool {
	if len(raw) == 0 {
		return true
	}
	var cond Conditions
	if err := json.Unmarshal(raw, &cond); err != nil {
		return false
	}
	return cond.Match(requestCtx)
}

func securityConditionsMatch(raw json.RawMessage, requestCtx RequestContext) bool {
	if len(raw) == 0 {
		return true
	}
	var cond SecurityConditions
	if err := json.Unmarshal(raw, &cond); err != nil {
		return false
	}
	return cond.Match(requestCtx)
}

func mapBaseAction(action models.SecurityProfileBaseAction) string {
	switch action {
	case models.SecurityProfileBaseActionBlock:
		return actionBlock
	default:
		return actionAllow
	}
}

func buildRequestContext(req *authpb.CheckRequest) RequestContext {
	ctx := RequestContext{OccurredAt: time.Now(), Headers: map[string]string{}}
	if req == nil || req.Attributes == nil || req.Attributes.Request == nil || req.Attributes.Request.Http == nil {
		return ctx
	}

	if req.Attributes.Request.Time != nil {
		ctx.OccurredAt = req.Attributes.Request.Time.AsTime()
	}

	httpReq := req.Attributes.Request.Http
	ctx.RequestHTTPID = httpReq.Id
	ctx.Method = httpReq.Method
	ctx.Host = httpReq.Host
	ctx.Scheme = httpReq.Scheme
	ctx.Path = httpReq.Path
	ctx.Fragment = httpReq.Fragment
	ctx.Protocol = httpReq.Protocol
	ctx.Headers = normalizeHeaders(httpReq.Headers)
	ctx.Authority = ctx.Headers[":authority"]
	ctx.UserAgent = ctx.Headers["user-agent"]
	ctx.RequestID = httpReq.Id
	if ctx.RequestID == "" {
		ctx.RequestID = ctx.Headers["x-request-id"]
	}
	ctx.ClientIP = extractClientIP(ctx.Headers)
	ctx.RequestBodySize = int32(httpReq.Size)
	ctx.RequestBody = httpReq.Body

	parsedPath, query := splitQuery(httpReq.Path)
	ctx.Path = parsedPath
	ctx.Query = query

	if req.Attributes.Source != nil {
		ctx.SourceService = req.Attributes.Source.Service
		ctx.SourceLabels = req.Attributes.Source.Labels
		ctx.SourcePrincipal = req.Attributes.Source.Principal
		if req.Attributes.Source.Address != nil {
			if sock := req.Attributes.Source.Address.GetSocketAddress(); sock != nil {
				if ctx.ClientIP == "" {
					ctx.ClientIP = sock.Address
				}
				ctx.SourcePort = int32(sock.GetPortValue())
			}
		}
	}

	if req.Attributes.Destination != nil {
		ctx.DestinationService = req.Attributes.Destination.Service
		ctx.DestinationLabels = req.Attributes.Destination.Labels
		if req.Attributes.Destination.Address != nil {
			if sock := req.Attributes.Destination.Address.GetSocketAddress(); sock != nil {
				ctx.DestinationIP = sock.Address
				ctx.DestinationPort = int32(sock.GetPortValue())
			}
		}
	}

	ctx.ContextExtensions = req.Attributes.ContextExtensions
	ctx.MetadataContext = protoToJSON(req.Attributes.MetadataContext)
	ctx.RouteMetadataCtx = protoToJSON(req.Attributes.RouteMetadataContext)

	return ctx
}

func normalizeHeaders(headers map[string]string) map[string]string {
	result := make(map[string]string, len(headers))
	for k, v := range headers {
		result[strings.ToLower(k)] = v
	}
	return result
}

func extractClientIP(headers map[string]string) string {
	forwarded := headers["x-forwarded-for"]
	if forwarded == "" {
		return ""
	}
	parts := strings.Split(forwarded, ",")
	if len(parts) == 0 {
		return ""
	}
	return strings.TrimSpace(parts[0])
}

func splitQuery(path string) (string, string) {
	if path == "" {
		return "", ""
	}
	parsed, err := url.Parse(path)
	if err != nil {
		if idx := strings.Index(path, "?"); idx != -1 {
			return path[:idx], path[idx+1:]
		}
		return path, ""
	}
	return parsed.Path, parsed.RawQuery
}

func headersToJSON(headers map[string]string) json.RawMessage {
	if len(headers) == 0 {
		return nil
	}
	data, err := json.Marshal(headers)
	if err != nil {
		return nil
	}
	return data
}

func mapToJSON(value map[string]string) json.RawMessage {
	if len(value) == 0 {
		return nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	return data
}

func protoToJSON(msg proto.Message) json.RawMessage {
	if msg == nil {
		return nil
	}
	data, err := protojson.MarshalOptions{
		UseProtoNames:   true,
		EmitUnpopulated: false,
	}.Marshal(msg)
	if err != nil {
		return nil
	}
	return data
}

func nullableString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func nullableInt32(value int32) *int32 {
	if value == 0 {
		return nil
	}
	return &value
}

func strPtr(value string) *string {
	return &value
}
