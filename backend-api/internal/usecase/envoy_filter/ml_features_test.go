package envoy_filter

import "testing"

func TestExtractMLFeatureVector_LengthMatchesOrder(t *testing.T) {
	vector := extractMLFeatureVector(RequestContext{})
	if len(vector) != len(mlFeatureOrder()) {
		t.Fatalf("vector length mismatch: got %d, want %d", len(vector), len(mlFeatureOrder()))
	}
}

func TestExtractMLFeatureVector_EncodesExpectedValues(t *testing.T) {
	requestCtx := RequestContext{
		Method:          "POST",
		Path:            "/admin/v1/users/123",
		Query:           "id=1&search=<script>alert(1)</script>&q=union+select",
		Scheme:          "https",
		Protocol:        "HTTP/2",
		Host:            "127.0.0.1:8080",
		Headers:         map[string]string{"content-type": "application/json", "accept": "application/json", "x-forwarded-for": "10.0.0.1, 172.16.0.5", "authorization": "Bearer token", "user-agent": "curl/8.7.1"},
		UserAgent:       "curl/8.7.1",
		ClientIP:        "10.0.0.1",
		SourcePort:      8081,
		RequestBodySize: 19,
		RequestBody:     "{\"q\":\"or 1=1 --\"}",
	}

	vector := extractMLFeatureVector(requestCtx)
	byName := vectorByFeatureName(t, vector)

	assertFeatureEq(t, byName, "method_post", 1)
	assertFeatureEq(t, byName, "method_get", 0)
	assertFeatureEq(t, byName, "path_has_suspicious_keywords", 1)
	assertFeatureEq(t, byName, "query_has_sqli_tokens", 1)
	assertFeatureEq(t, byName, "query_has_xss_tokens", 1)
	assertFeatureEq(t, byName, "host_is_ip_literal", 1)
	assertFeatureEq(t, byName, "scheme_https", 1)
	assertFeatureEq(t, byName, "protocol_http_2", 1)
	assertFeatureEq(t, byName, "content_type_json", 1)
	assertFeatureEq(t, byName, "accept_json", 1)
	assertFeatureEq(t, byName, "has_authorization_header", 1)
	assertFeatureEq(t, byName, "x_forwarded_for_hops_count", 2)
	assertFeatureEq(t, byName, "client_ip_is_private", 1)
	assertFeatureEq(t, byName, "client_port_present", 1)
	assertFeatureEq(t, byName, "body_present", 1)
	assertFeatureEq(t, byName, "body_has_sqli_or_xss_tokens", 1)
	assertFeatureEq(t, byName, "user_agent_has_automation_tokens", 1)
}

func vectorByFeatureName(t *testing.T, vector []float32) map[string]float32 {
	t.Helper()
	order := mlFeatureOrder()
	if len(order) != len(vector) {
		t.Fatalf("feature order length mismatch: got %d, vector %d", len(order), len(vector))
	}

	result := make(map[string]float32, len(order))
	for i, name := range order {
		result[name] = vector[i]
	}
	return result
}

func assertFeatureEq(t *testing.T, byName map[string]float32, name string, want float32) {
	t.Helper()
	got, ok := byName[name]
	if !ok {
		t.Fatalf("missing feature: %s", name)
	}
	if got != want {
		t.Fatalf("feature %s mismatch: got %v, want %v", name, got, want)
	}
}
