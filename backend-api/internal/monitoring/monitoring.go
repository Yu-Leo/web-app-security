package monitoring

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	envoyRequestsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "was_envoy_requests_total",
		Help: "Total number of authz requests received from Envoy.",
	})
	envoyRequestsAllowedTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "was_envoy_requests_allowed_total",
		Help: "Total number of requests allowed by backend authz flow.",
	})
	envoyRequestsBlockedTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "was_envoy_requests_blocked_total",
		Help: "Total number of requests blocked by backend authz flow.",
	})
	envoyRequestDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "was_envoy_request_duration_seconds",
		Help:    "Duration of backend authz request processing.",
		Buckets: prometheus.DefBuckets,
	})
	securityRuleMatchesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "was_security_rule_matches_total",
		Help: "Total number of matched security rules.",
	}, []string{"rule_id", "profile_id", "action", "rule_type", "dry_run"})
	trafficRuleMatchesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "was_traffic_rule_matches_total",
		Help: "Total number of matched traffic rules.",
	}, []string{"rule_id", "profile_id", "dry_run"})
	securityBaseActionTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "was_security_base_action_total",
		Help: "Total number of times security profile base action was used.",
	}, []string{"profile_id", "action"})
	mlInferenceTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "was_ml_inference_total",
		Help: "Total number of ML inference executions.",
	}, []string{"model_id", "model_name"})
	mlInferenceErrorsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "was_ml_inference_errors_total",
		Help: "Total number of ML errors.",
	}, []string{"model_id", "model_name"})
	mlThresholdPassTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "was_ml_threshold_pass_total",
		Help: "Total number of ML evaluations with score above or equal to threshold.",
	}, []string{"model_id", "model_name"})
	mlThresholdMissTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "was_ml_threshold_miss_total",
		Help: "Total number of ML evaluations with score below threshold.",
	}, []string{"model_id", "model_name"})
	mlScoreHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "was_ml_score",
		Help:    "Distribution of ML scores.",
		Buckets: []float64{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1},
	}, []string{"model_id", "model_name"})
	dryRunTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "was_dry_run_total",
		Help: "Total number of dry-run activations.",
	}, []string{"rule_type", "action"})
	httpAPIRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "was_http_api_requests_total",
		Help: "Total number of management API HTTP requests.",
	}, []string{"method", "route", "status_code"})
	httpAPIRequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "was_http_api_request_duration_seconds",
		Help:    "Duration of management API HTTP requests.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "route"})
)

func init() {
	prometheus.MustRegister(
		envoyRequestsTotal,
		envoyRequestsAllowedTotal,
		envoyRequestsBlockedTotal,
		envoyRequestDuration,
		securityRuleMatchesTotal,
		trafficRuleMatchesTotal,
		securityBaseActionTotal,
		mlInferenceTotal,
		mlInferenceErrorsTotal,
		mlThresholdPassTotal,
		mlThresholdMissTotal,
		mlScoreHistogram,
		dryRunTotal,
		httpAPIRequestsTotal,
		httpAPIRequestDuration,
	)
}

func Handler() http.Handler {
	return promhttp.Handler()
}

func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		startedAt := time.Now()
		c.Next()

		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}
		if route == "" {
			route = "unknown"
		}

		statusCode := strconv.Itoa(c.Writer.Status())
		httpAPIRequestsTotal.WithLabelValues(c.Request.Method, route, statusCode).Inc()
		httpAPIRequestDuration.WithLabelValues(c.Request.Method, route).Observe(time.Since(startedAt).Seconds())
	}
}

func ObserveEnvoyRequest(allowed *bool, duration time.Duration) {
	envoyRequestsTotal.Inc()
	envoyRequestDuration.Observe(duration.Seconds())
	if allowed == nil {
		return
	}
	if *allowed {
		envoyRequestsAllowedTotal.Inc()
		return
	}
	envoyRequestsBlockedTotal.Inc()
}

func IncSecurityRuleMatch(ruleID int64, profileID int64, action string, ruleType string, dryRun bool) {
	securityRuleMatchesTotal.WithLabelValues(
		strconv.FormatInt(ruleID, 10),
		strconv.FormatInt(profileID, 10),
		action,
		ruleType,
		strconv.FormatBool(dryRun),
	).Inc()
}

func IncTrafficRuleMatch(ruleID int64, profileID int64, dryRun bool) {
	trafficRuleMatchesTotal.WithLabelValues(
		strconv.FormatInt(ruleID, 10),
		strconv.FormatInt(profileID, 10),
		strconv.FormatBool(dryRun),
	).Inc()
}

func IncSecurityBaseAction(profileID int64, action string) {
	securityBaseActionTotal.WithLabelValues(strconv.FormatInt(profileID, 10), action).Inc()
}

func IncMLInference(modelID int64, modelName string) {
	mlInferenceTotal.WithLabelValues(strconv.FormatInt(modelID, 10), modelName).Inc()
}

func IncMLError(modelID *int64, modelName *string) {
	mlInferenceErrorsTotal.WithLabelValues(labelInt64(modelID), labelString(modelName)).Inc()
}

func IncMLThresholdPass(modelID int64, modelName string) {
	mlThresholdPassTotal.WithLabelValues(strconv.FormatInt(modelID, 10), modelName).Inc()
}

func IncMLThresholdMiss(modelID int64, modelName string) {
	mlThresholdMissTotal.WithLabelValues(strconv.FormatInt(modelID, 10), modelName).Inc()
}

func ObserveMLScore(modelID int64, modelName string, score float64) {
	mlScoreHistogram.WithLabelValues(strconv.FormatInt(modelID, 10), modelName).Observe(score)
}

func IncDryRun(ruleType string, action string) {
	dryRunTotal.WithLabelValues(ruleType, action).Inc()
}

func labelInt64(value *int64) string {
	if value == nil {
		return "unknown"
	}
	return strconv.FormatInt(*value, 10)
}

func labelString(value *string) string {
	if value == nil || *value == "" {
		return "unknown"
	}
	return *value
}
