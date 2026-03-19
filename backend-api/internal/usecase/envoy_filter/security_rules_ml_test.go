package envoy_filter

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/Yu-Leo/web-app-security/backend-api/internal/models"
)

func TestApplySecurityRules_MLThresholdMatched(t *testing.T) {
	modelID := int64(10)
	threshold := int16(70)

	u := &UsecaseImpl{
		securityRules: stubSecurityRulesRepo{
			rules: []models.SecurityRule{{
				ID:          1,
				ProfileID:   100,
				Name:        "ml-block",
				Priority:    1,
				RuleType:    models.SecurityRuleTypeML,
				Action:      models.SecurityRuleActionBlock,
				Conditions:  json.RawMessage(`{"method_regex":["POST"]}`),
				MLModelID:   &modelID,
				MLThreshold: &threshold,
				IsEnabled:   true,
			}},
		},
		mlModels: stubMLModelRepo{
			model: models.MLModel{ID: modelID, Name: "fraud-v1", ModelData: []byte("model")},
		},
		mlScorer: stubMLScorer{
			score: 0.82,
		},
	}

	result := u.applySecurityRules(context.Background(), models.SecurityProfile{
		ID:         100,
		IsEnabled:  true,
		BaseAction: models.SecurityProfileBaseActionAllow,
	}, RequestContext{
		Method:          "POST",
		MLFeatureVector: []float32{0.1, 0.2, 0.3},
	})

	if result.Rule == nil {
		t.Fatalf("expected matched rule")
	}
	if result.Action != actionBlock {
		t.Fatalf("unexpected action: %s", result.Action)
	}
	if result.DecisionSource != "ml" {
		t.Fatalf("unexpected decision source: %s", result.DecisionSource)
	}
	if result.MLScorePercent == nil || *result.MLScorePercent < 70 {
		t.Fatalf("unexpected score percent: %v", result.MLScorePercent)
	}
}

func TestApplySecurityRules_MLInferenceErrorSkipsRule(t *testing.T) {
	modelID := int64(10)
	threshold := int16(50)

	u := &UsecaseImpl{
		securityRules: stubSecurityRulesRepo{
			rules: []models.SecurityRule{
				{
					ID:          1,
					ProfileID:   100,
					Name:        "ml-block",
					Priority:    1,
					RuleType:    models.SecurityRuleTypeML,
					Action:      models.SecurityRuleActionBlock,
					MLModelID:   &modelID,
					MLThreshold: &threshold,
					IsEnabled:   true,
				},
				{
					ID:        2,
					ProfileID: 100,
					Name:      "allow-next",
					Priority:  2,
					RuleType:  models.SecurityRuleTypeDeterministic,
					Action:    models.SecurityRuleActionAllow,
					IsEnabled: true,
				},
			},
		},
		mlModels: stubMLModelRepo{
			model: models.MLModel{ID: modelID, Name: "fraud-v1", ModelData: []byte("model")},
		},
		mlScorer: stubMLScorer{
			err: errors.New("runtime failed"),
		},
	}

	result := u.applySecurityRules(context.Background(), models.SecurityProfile{
		ID:         100,
		IsEnabled:  true,
		BaseAction: models.SecurityProfileBaseActionBlock,
	}, RequestContext{
		MLFeatureVector: []float32{0.1, 0.2, 0.3},
	})

	if result.Rule == nil || result.Rule.ID != 2 {
		t.Fatalf("expected next rule to match after ml error, got %+v", result.Rule)
	}
	if result.Action != actionAllow {
		t.Fatalf("unexpected action: %s", result.Action)
	}
	if result.DecisionSource != "rule" {
		t.Fatalf("unexpected decision source: %s", result.DecisionSource)
	}
	if result.MLError == nil {
		t.Fatalf("expected ml error to be preserved in result")
	}
}

func TestApplySecurityRules_MLBelowThresholdSkipsRule(t *testing.T) {
	modelID := int64(10)
	threshold := int16(90)

	u := &UsecaseImpl{
		securityRules: stubSecurityRulesRepo{
			rules: []models.SecurityRule{
				{
					ID:          1,
					ProfileID:   100,
					Name:        "ml-block",
					Priority:    1,
					RuleType:    models.SecurityRuleTypeML,
					Action:      models.SecurityRuleActionBlock,
					MLModelID:   &modelID,
					MLThreshold: &threshold,
					IsEnabled:   true,
				},
				{
					ID:        2,
					ProfileID: 100,
					Name:      "block-next",
					Priority:  2,
					RuleType:  models.SecurityRuleTypeDeterministic,
					Action:    models.SecurityRuleActionBlock,
					IsEnabled: true,
				},
			},
		},
		mlModels: stubMLModelRepo{
			model: models.MLModel{ID: modelID, Name: "fraud-v1", ModelData: []byte("model")},
		},
		mlScorer: stubMLScorer{
			score: 0.82,
		},
	}

	result := u.applySecurityRules(context.Background(), models.SecurityProfile{
		ID:         100,
		IsEnabled:  true,
		BaseAction: models.SecurityProfileBaseActionAllow,
	}, RequestContext{
		MLFeatureVector: []float32{0.1, 0.2, 0.3},
	})

	if result.Rule == nil || result.Rule.ID != 2 {
		t.Fatalf("expected next rule to match after low ml score, got %+v", result.Rule)
	}
	if result.Action != actionBlock {
		t.Fatalf("unexpected action: %s", result.Action)
	}
}

func TestApplyTrafficRules_UsesSecurityConditionsSchema(t *testing.T) {
	u := &UsecaseImpl{
		trafficRules: stubTrafficRulesRepo{
			rules: []models.TrafficRule{
				{
					ID:            1,
					ProfileID:     77,
					Name:          "skip-by-host",
					Priority:      1,
					RequestsLimit: 1,
					PeriodSeconds: 60,
					Conditions:    json.RawMessage(`{"host_regex":["^admin\\.example\\.com$"]}`),
					IsEnabled:     true,
				},
				{
					ID:            2,
					ProfileID:     77,
					Name:          "match-by-method",
					Priority:      2,
					RequestsLimit: 1,
					PeriodSeconds: 60,
					Conditions:    json.RawMessage(`{"method_regex":["^POST$"]}`),
					IsEnabled:     true,
				},
			},
		},
		rateLimiter: NewInMemoryRateLimiter(),
	}

	profile := models.TrafficProfile{
		ID:        77,
		IsEnabled: true,
	}

	requestCtx := RequestContext{
		Method:   "POST",
		Host:     "app.example.com",
		ClientIP: "10.0.0.10",
	}

	if rule, action, _ := u.applyTrafficRules(context.Background(), profile, 42, requestCtx); rule != nil || action != actionAllow {
		t.Fatalf("first request should pass limiter, got rule=%+v action=%s", rule, action)
	}

	rule, action, dryRun := u.applyTrafficRules(context.Background(), profile, 42, requestCtx)
	if rule == nil || rule.ID != 2 {
		t.Fatalf("expected second traffic rule to match, got %+v", rule)
	}
	if action != actionBlock {
		t.Fatalf("unexpected action: %s", action)
	}
	if dryRun {
		t.Fatalf("expected non-dry-run rule")
	}
}

type stubSecurityRulesRepo struct {
	rules []models.SecurityRule
	err   error
}

func (s stubSecurityRulesRepo) List(_ context.Context) ([]models.SecurityRule, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.rules, nil
}

type stubMLModelRepo struct {
	model models.MLModel
	err   error
}

func (s stubMLModelRepo) Get(_ context.Context, _ int64) (models.MLModel, error) {
	if s.err != nil {
		return models.MLModel{}, s.err
	}
	return s.model, nil
}

type stubMLScorer struct {
	score float32
	err   error
}

func (s stubMLScorer) Score(_ context.Context, _ string, _ []byte, featureVector []float32) (float32, error) {
	if len(featureVector) == 0 {
		return 0, errors.New("empty feature vector")
	}
	if s.err != nil {
		return 0, s.err
	}
	return s.score, nil
}

type stubResourceRepo struct{}
type stubSecurityProfileRepo struct{}
type stubTrafficProfileRepo struct{}
type stubTrafficRulesRepo struct {
	rules []models.TrafficRule
	err   error
}
type stubRequestLogRepo struct{}
type stubEventLogRepo struct{}

func (stubResourceRepo) List(context.Context) ([]models.Resource, error) { return nil, nil }
func (stubSecurityProfileRepo) Get(context.Context, int64) (models.SecurityProfile, error) {
	return models.SecurityProfile{}, nil
}
func (stubTrafficProfileRepo) Get(context.Context, int64) (models.TrafficProfile, error) {
	return models.TrafficProfile{}, nil
}
func (s stubTrafficRulesRepo) List(context.Context) ([]models.TrafficRule, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.rules, nil
}
func (stubRequestLogRepo) Create(context.Context, models.RequestLog) (models.RequestLog, error) {
	return models.RequestLog{}, nil
}
func (stubEventLogRepo) Create(context.Context, models.EventLog) (models.EventLog, error) {
	return models.EventLog{}, nil
}

func TestNewUsecase_DefaultMLScorerIsInitialized(t *testing.T) {
	usecase := NewUsecase(
		stubResourceRepo{},
		stubSecurityProfileRepo{},
		stubTrafficProfileRepo{},
		stubSecurityRulesRepo{},
		stubTrafficRulesRepo{},
		stubMLModelRepo{},
		nil,
		stubRequestLogRepo{},
		stubEventLogRepo{},
	)

	if usecase.mlScorer == nil {
		t.Fatalf("expected default ml scorer to be configured")
	}
}

func TestUnavailableMLScorer_ReturnsError(t *testing.T) {
	scorer := NewUnavailableMLScorer()
	_, err := scorer.Score(context.Background(), "model", []byte("data"), []float32{1})
	if err == nil {
		t.Fatalf("expected scorer error")
	}
}

func TestDecisionFieldsForMLMetadata(t *testing.T) {
	now := time.Now()
	score := float32(0.91)
	scorePercent := score * 100
	modelID := int64(7)
	threshold := int16(80)
	modelName := "fraud-v2"
	source := "ml"

	decision := Decision{
		DecidedAt:      now,
		DecisionSource: &source,
		MLModelID:      &modelID,
		MLModelName:    &modelName,
		MLThreshold:    &threshold,
		MLScore:        &score,
		MLScorePercent: &scorePercent,
	}

	if decision.MLModelID == nil || *decision.MLModelID != 7 {
		t.Fatalf("unexpected model id")
	}
	if decision.DecisionSource == nil || *decision.DecisionSource != "ml" {
		t.Fatalf("unexpected decision source")
	}
}

var _ Usecase = (*UsecaseImpl)(nil)
