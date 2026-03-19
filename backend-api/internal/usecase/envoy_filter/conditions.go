package envoy_filter

import (
	"net"
	"regexp"
	"strings"
)

type Conditions struct {
	Methods        []string            `json:"methods,omitempty"`
	PathRegex      []string            `json:"path_regex,omitempty"`
	PathPrefix     []string            `json:"path_prefix,omitempty"`
	QueryRegex     []string            `json:"query_regex,omitempty"`
	Hosts          []string            `json:"hosts,omitempty"`
	HostRegex      []string            `json:"host_regex,omitempty"`
	Headers        map[string][]string `json:"headers,omitempty"`
	HeadersRegex   map[string][]string `json:"headers_regex,omitempty"`
	IPCIDR         []string            `json:"ip_cidr,omitempty"`
	UserAgentRegex []string            `json:"user_agent_regex,omitempty"`
}

type SecurityConditions struct {
	SourceIPCIDR []string                  `json:"source_ip_cidr,omitempty"`
	URIRegex     []string                  `json:"uri_regex,omitempty"`
	HostRegex    []string                  `json:"host_regex,omitempty"`
	MethodRegex  []string                  `json:"method_regex,omitempty"`
	Headers      []SecurityHeaderCondition `json:"headers,omitempty"`
}

type SecurityHeaderCondition struct {
	Name       string   `json:"name"`
	ValueRegex []string `json:"value_regex,omitempty"`
}

func (c Conditions) Match(ctx RequestContext) bool {
	if len(c.Methods) > 0 && !matchStringSlice(ctx.Method, c.Methods, true) {
		return false
	}
	if len(c.PathPrefix) > 0 && !matchPrefixSlice(ctx.Path, c.PathPrefix) {
		return false
	}
	if len(c.PathRegex) > 0 && !matchRegexSlice(ctx.Path, c.PathRegex) {
		return false
	}
	if len(c.QueryRegex) > 0 && !matchRegexSlice(ctx.Query, c.QueryRegex) {
		return false
	}
	if len(c.Hosts) > 0 && !matchStringSlice(ctx.Host, c.Hosts, true) {
		return false
	}
	if len(c.HostRegex) > 0 && !matchRegexSlice(ctx.Host, c.HostRegex) {
		return false
	}
	if len(c.Headers) > 0 && !matchHeadersExact(ctx.Headers, c.Headers) {
		return false
	}
	if len(c.HeadersRegex) > 0 && !matchHeadersRegex(ctx.Headers, c.HeadersRegex) {
		return false
	}
	if len(c.IPCIDR) > 0 && !matchIPCIDR(ctx.ClientIP, c.IPCIDR) {
		return false
	}
	if len(c.UserAgentRegex) > 0 && !matchRegexSlice(ctx.UserAgent, c.UserAgentRegex) {
		return false
	}

	return true
}

func (c SecurityConditions) Match(ctx RequestContext) bool {
	if len(c.SourceIPCIDR) > 0 && !matchIPCIDR(ctx.ClientIP, c.SourceIPCIDR) {
		return false
	}
	if len(c.URIRegex) > 0 && !matchRegexSlice(ctx.Path, c.URIRegex) {
		return false
	}
	if len(c.HostRegex) > 0 && !matchRegexSlice(ctx.Host, c.HostRegex) {
		return false
	}
	if len(c.MethodRegex) > 0 && !matchRegexSlice(ctx.Method, c.MethodRegex) {
		return false
	}
	if len(c.Headers) > 0 && !matchSecurityHeaders(ctx.Headers, c.Headers) {
		return false
	}

	return true
}

func matchStringSlice(value string, candidates []string, caseInsensitive bool) bool {
	if caseInsensitive {
		value = strings.ToLower(value)
	}
	for _, candidate := range candidates {
		c := candidate
		if caseInsensitive {
			c = strings.ToLower(candidate)
		}
		if value == c {
			return true
		}
	}
	return false
}

func matchPrefixSlice(value string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(value, prefix) {
			return true
		}
	}
	return false
}

func matchRegexSlice(value string, regexes []string) bool {
	for _, pattern := range regexes {
		re, err := regexp.Compile(pattern)
		if err != nil {
			continue
		}
		if re.MatchString(value) {
			return true
		}
	}
	return false
}

func matchHeadersExact(headers map[string]string, rules map[string][]string) bool {
	for key, allowed := range rules {
		value, ok := headers[strings.ToLower(key)]
		if !ok {
			return false
		}
		if !matchStringSlice(value, allowed, false) {
			return false
		}
	}
	return true
}

func matchHeadersRegex(headers map[string]string, rules map[string][]string) bool {
	for key, patterns := range rules {
		value, ok := headers[strings.ToLower(key)]
		if !ok {
			return false
		}
		if !matchRegexSlice(value, patterns) {
			return false
		}
	}
	return true
}

func matchIPCIDR(ip string, cidrs []string) bool {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return false
	}
	for _, cidr := range cidrs {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if network.Contains(parsed) {
			return true
		}
	}
	return false
}

func matchSecurityHeaders(headers map[string]string, rules []SecurityHeaderCondition) bool {
	for _, rule := range rules {
		value, ok := headers[strings.ToLower(rule.Name)]
		if !ok {
			return false
		}
		if !matchRegexSlice(value, rule.ValueRegex) {
			return false
		}
	}
	return true
}
