package envoy_filter

import (
	"math"
	"net"
	"net/url"
	"regexp"
	"strings"
	"unicode"
)

var (
	pathEncodedCharsRegex = regexp.MustCompile(`%[0-9a-fA-F]{2}`)

	pathTraversalTokens = []string{
		"../",
		"..\\",
		"%2e%2e",
		"%2f",
		"%5c",
	}
	pathSuspiciousKeywords = []string{
		"admin",
		"phpmyadmin",
		"wp-admin",
		"debug",
		"actuator",
	}
	sqliTokens = []string{
		"union",
		"select",
		"or 1=1",
		"sleep(",
		"benchmark(",
		"--",
		"/*",
	}
	xssTokens = []string{
		"<script",
		"onerror=",
		"javascript:",
	}
	uaAutomationTokens = []string{
		"curl",
		"wget",
		"python-requests",
		"sqlmap",
		"nikto",
		"nmap",
	}

	featureOrder = []string{
		"method_get",
		"method_post",
		"method_put",
		"method_delete",
		"method_patch",
		"method_other",
		"path_len",
		"path_depth",
		"path_segment_max_len",
		"path_digit_ratio",
		"path_special_char_ratio",
		"path_has_encoded_chars",
		"path_has_traversal_pattern",
		"path_has_suspicious_keywords",
		"query_len",
		"query_param_count",
		"query_key_count",
		"query_value_avg_len",
		"query_special_char_ratio",
		"query_has_sqli_tokens",
		"query_has_xss_tokens",
		"host_len",
		"host_is_ip_literal",
		"scheme_http",
		"scheme_https",
		"scheme_other",
		"protocol_http_1_1",
		"protocol_http_2",
		"protocol_other",
		"headers_count",
		"header_name_count_unique",
		"has_authorization_header",
		"has_cookie_header",
		"has_content_type_header",
		"content_type_json",
		"content_type_form",
		"content_type_multipart",
		"content_type_xml",
		"content_type_text",
		"content_type_other",
		"content_type_missing",
		"accept_json",
		"accept_html",
		"accept_any",
		"accept_other",
		"accept_missing",
		"user_agent_len",
		"user_agent_missing",
		"user_agent_token_count",
		"user_agent_has_automation_tokens",
		"x_forwarded_for_hops_count",
		"client_ip_is_private",
		"client_ip_is_loopback",
		"client_port_present",
		"body_size",
		"body_present",
		"body_entropy",
		"body_non_printable_ratio",
		"body_has_sqli_or_xss_tokens",
	}
)

func extractMLFeatureVector(requestCtx RequestContext) []float32 {
	vector := make([]float32, 0, len(featureOrder))

	path := requestCtx.Path
	pathLen := len(path)
	pathLower := strings.ToLower(path)
	pathDepth := calcPathDepth(path)
	pathSegmentMaxLen := calcPathSegmentMaxLen(path)
	pathDigitRatio := ratio(countDigits(path), pathLen)
	pathSpecialRatio := ratio(countPathSpecialChars(path), pathLen)
	pathHasEncodedChars := pathEncodedCharsRegex.MatchString(path)
	pathHasTraversal := containsAny(pathLower, pathTraversalTokens)
	pathHasSuspiciousKeywords := containsAny(pathLower, pathSuspiciousKeywords)

	query := requestCtx.Query
	queryLen := len(query)
	queryLower := strings.ToLower(query)
	queryParamCount, queryKeyCount, queryValueAvgLen := parseQueryStats(query)
	querySpecialRatio := ratio(countQuerySpecialChars(query), queryLen)
	queryHasSQLITokens := containsAny(queryLower, sqliTokens)
	queryHasXSSTokens := containsAny(queryLower, xssTokens)

	host := requestCtx.Host
	hostLen := len(host)
	hostIsIPLiteral := isIPLiteralHost(host)

	headersCount := len(requestCtx.Headers)
	headerNameCountUnique := len(requestCtx.Headers)

	hasAuthorization := hasHeader(requestCtx.Headers, "authorization")
	hasCookie := hasHeader(requestCtx.Headers, "cookie")
	hasContentType := hasHeader(requestCtx.Headers, "content-type")

	contentTypeBucket := classifyContentType(requestCtx.Headers["content-type"])
	acceptBucket := classifyAccept(requestCtx.Headers["accept"])

	userAgent := requestCtx.UserAgent
	userAgentLower := strings.ToLower(userAgent)
	userAgentLen := len(userAgent)
	userAgentMissing := userAgent == ""
	userAgentTokenCount := countUserAgentTokens(userAgent)
	userAgentHasAutomationTokens := containsAny(userAgentLower, uaAutomationTokens)

	xForwardedForHopsCount := countXForwardedForHops(requestCtx.Headers["x-forwarded-for"])

	clientIPIsPrivate, clientIPIsLoopback := classifyClientIP(requestCtx.ClientIP)
	clientPortPresent := requestCtx.SourcePort > 0

	body := requestCtx.RequestBody
	bodyLower := strings.ToLower(body)
	bodySize := requestCtx.RequestBodySize
	bodyPresent := bodySize > 0 || body != ""
	bodyEntropy := calcEntropy(body)
	bodyNonPrintableRatio := ratio(countNonPrintable(body), len(body))
	bodyHasSQLIOrXSS := containsAny(bodyLower, sqliTokens) || containsAny(bodyLower, xssTokens)

	vector = appendOneHot(vector, normalizeMethod(requestCtx.Method), []string{"get", "post", "put", "delete", "patch", "other"})
	vector = append(vector,
		float32(pathLen),
		float32(pathDepth),
		float32(pathSegmentMaxLen),
		pathDigitRatio,
		pathSpecialRatio,
		boolToFloat(pathHasEncodedChars),
		boolToFloat(pathHasTraversal),
		boolToFloat(pathHasSuspiciousKeywords),
		float32(queryLen),
		float32(queryParamCount),
		float32(queryKeyCount),
		queryValueAvgLen,
		querySpecialRatio,
		boolToFloat(queryHasSQLITokens),
		boolToFloat(queryHasXSSTokens),
		float32(hostLen),
		boolToFloat(hostIsIPLiteral),
	)

	vector = appendOneHot(vector, normalizeScheme(requestCtx.Scheme), []string{"http", "https", "other"})
	vector = appendOneHot(vector, normalizeProtocol(requestCtx.Protocol), []string{"http/1.1", "http/2", "other"})

	vector = append(vector,
		float32(headersCount),
		float32(headerNameCountUnique),
		boolToFloat(hasAuthorization),
		boolToFloat(hasCookie),
		boolToFloat(hasContentType),
	)

	vector = appendOneHot(vector, contentTypeBucket, []string{"json", "form", "multipart", "xml", "text", "other", "missing"})
	vector = appendOneHot(vector, acceptBucket, []string{"json", "html", "*/*", "other", "missing"})

	vector = append(vector,
		float32(userAgentLen),
		boolToFloat(userAgentMissing),
		float32(userAgentTokenCount),
		boolToFloat(userAgentHasAutomationTokens),
		float32(xForwardedForHopsCount),
		boolToFloat(clientIPIsPrivate),
		boolToFloat(clientIPIsLoopback),
		boolToFloat(clientPortPresent),
		float32(bodySize),
		boolToFloat(bodyPresent),
		bodyEntropy,
		bodyNonPrintableRatio,
		boolToFloat(bodyHasSQLIOrXSS),
	)

	return vector
}

func mlFeatureOrder() []string {
	order := make([]string, len(featureOrder))
	copy(order, featureOrder)
	return order
}

func appendOneHot(target []float32, value string, dictionary []string) []float32 {
	for _, item := range dictionary {
		target = append(target, boolToFloat(value == item))
	}
	return target
}

func boolToFloat(value bool) float32 {
	if value {
		return 1
	}
	return 0
}

func ratio(numerator int, denominator int) float32 {
	if denominator <= 0 {
		return 0
	}
	return float32(numerator) / float32(denominator)
}

func calcPathDepth(path string) int {
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return 0
	}
	return strings.Count(trimmed, "/") + 1
}

func calcPathSegmentMaxLen(path string) int {
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return 0
	}
	maxLen := 0
	for _, segment := range strings.Split(trimmed, "/") {
		if l := len(segment); l > maxLen {
			maxLen = l
		}
	}
	return maxLen
}

func countDigits(value string) int {
	total := 0
	for _, r := range value {
		if unicode.IsDigit(r) {
			total++
		}
	}
	return total
}

func countPathSpecialChars(path string) int {
	total := 0
	for _, r := range path {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r):
			continue
		case r == '/', r == '-', r == '_', r == '.', r == '~':
			continue
		default:
			total++
		}
	}
	return total
}

func countQuerySpecialChars(query string) int {
	total := 0
	for _, r := range query {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r):
			continue
		case r == '&', r == '=', r == '-', r == '_', r == '.', r == '~', r == '%':
			continue
		default:
			total++
		}
	}
	return total
}

func parseQueryStats(query string) (int, int, float32) {
	if query == "" {
		return 0, 0, 0
	}

	parsed, err := url.ParseQuery(query)
	if err != nil {
		return fallbackQueryStats(query)
	}

	paramCount := 0
	valueLenSum := 0
	valueCount := 0

	for _, values := range parsed {
		for _, value := range values {
			paramCount++
			valueCount++
			valueLenSum += len(value)
		}
	}

	avg := float32(0)
	if valueCount > 0 {
		avg = float32(valueLenSum) / float32(valueCount)
	}

	return paramCount, len(parsed), avg
}

func fallbackQueryStats(query string) (int, int, float32) {
	parts := strings.Split(query, "&")
	paramCount := 0
	keySet := make(map[string]struct{})
	valueCount := 0
	valueLenSum := 0

	for _, part := range parts {
		if strings.TrimSpace(part) == "" {
			continue
		}
		paramCount++
		kv := strings.SplitN(part, "=", 2)
		key := strings.TrimSpace(kv[0])
		if key != "" {
			keySet[key] = struct{}{}
		}
		if len(kv) == 2 {
			valueCount++
			valueLenSum += len(kv[1])
		}
	}

	avg := float32(0)
	if valueCount > 0 {
		avg = float32(valueLenSum) / float32(valueCount)
	}
	return paramCount, len(keySet), avg
}

func isIPLiteralHost(host string) bool {
	if host == "" {
		return false
	}
	candidate := strings.TrimSpace(host)

	if strings.HasPrefix(candidate, "[") && strings.Contains(candidate, "]") {
		end := strings.Index(candidate, "]")
		if end > 1 {
			candidate = candidate[1:end]
		}
	} else if strings.Count(candidate, ":") == 1 {
		name, _, err := net.SplitHostPort(candidate)
		if err == nil {
			candidate = name
		}
	}

	candidate = strings.Trim(candidate, "[]")
	return net.ParseIP(candidate) != nil
}

func hasHeader(headers map[string]string, key string) bool {
	if len(headers) == 0 {
		return false
	}
	value, ok := headers[strings.ToLower(key)]
	return ok && strings.TrimSpace(value) != ""
}

func normalizeMethod(method string) string {
	switch strings.ToUpper(strings.TrimSpace(method)) {
	case "GET":
		return "get"
	case "POST":
		return "post"
	case "PUT":
		return "put"
	case "DELETE":
		return "delete"
	case "PATCH":
		return "patch"
	default:
		return "other"
	}
}

func normalizeScheme(scheme string) string {
	switch strings.ToLower(strings.TrimSpace(scheme)) {
	case "http":
		return "http"
	case "https":
		return "https"
	default:
		return "other"
	}
}

func normalizeProtocol(protocol string) string {
	normalized := strings.ToLower(strings.TrimSpace(protocol))
	switch normalized {
	case "http/1.1":
		return "http/1.1"
	case "http/2", "h2":
		return "http/2"
	default:
		return "other"
	}
}

func classifyContentType(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return "missing"
	}
	switch {
	case strings.Contains(normalized, "application/json"), strings.Contains(normalized, "+json"):
		return "json"
	case strings.Contains(normalized, "application/x-www-form-urlencoded"):
		return "form"
	case strings.Contains(normalized, "multipart/form-data"):
		return "multipart"
	case strings.Contains(normalized, "application/xml"), strings.Contains(normalized, "text/xml"), strings.Contains(normalized, "+xml"):
		return "xml"
	case strings.Contains(normalized, "text/"):
		return "text"
	default:
		return "other"
	}
}

func classifyAccept(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return "missing"
	}
	switch {
	case strings.Contains(normalized, "application/json"), strings.Contains(normalized, "+json"):
		return "json"
	case strings.Contains(normalized, "text/html"):
		return "html"
	case strings.Contains(normalized, "*/*"):
		return "*/*"
	default:
		return "other"
	}
}

func countUserAgentTokens(userAgent string) int {
	return len(strings.FieldsFunc(userAgent, func(r rune) bool {
		switch r {
		case ' ', '/', ';':
			return true
		default:
			return false
		}
	}))
}

func countXForwardedForHops(value string) int {
	if strings.TrimSpace(value) == "" {
		return 0
	}
	parts := strings.Split(value, ",")
	total := 0
	for _, part := range parts {
		if strings.TrimSpace(part) != "" {
			total++
		}
	}
	return total
}

func classifyClientIP(clientIP string) (bool, bool) {
	parsed := net.ParseIP(strings.TrimSpace(clientIP))
	if parsed == nil {
		return false, false
	}
	return parsed.IsPrivate(), parsed.IsLoopback()
}

func calcEntropy(value string) float32 {
	if len(value) == 0 {
		return 0
	}
	counts := make(map[byte]int)
	data := []byte(value)
	for _, b := range data {
		counts[b]++
	}

	entropy := 0.0
	total := float64(len(data))
	for _, count := range counts {
		probability := float64(count) / total
		entropy -= probability * math.Log2(probability)
	}
	return float32(entropy)
}

func countNonPrintable(value string) int {
	total := 0
	for _, r := range value {
		if unicode.IsPrint(r) {
			continue
		}
		if r == '\n' || r == '\r' || r == '\t' {
			continue
		}
		total++
	}
	return total
}

func containsAny(value string, tokens []string) bool {
	for _, token := range tokens {
		if strings.Contains(value, token) {
			return true
		}
	}
	return false
}
