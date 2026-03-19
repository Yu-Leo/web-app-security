package http

import (
	"encoding/json"

	"github.com/Yu-Leo/web-app-security/backend-api/internal/generated/service"
)

var securityConditionArrayKeys = map[string]struct{}{
	"source_ip_cidr": {},
	"uri_regex":      {},
	"host_regex":     {},
	"method_regex":   {},
}

func normalizeSecurityConditions(raw json.RawMessage) (*service.SecurityRuleConditions, error) {
	if len(raw) == 0 {
		return nil, nil
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}

	normalized := map[string]interface{}{}
	for key, value := range payload {
		if _, ok := securityConditionArrayKeys[key]; ok {
			normalized[key] = normalizeStringArray(value)
			continue
		}

		switch key {
		case "headers":
			normalized[key] = normalizeSecurityHeaders(value)
		default:
			normalized[key] = value
		}
	}

	encoded, err := json.Marshal(normalized)
	if err != nil {
		return nil, err
	}

	var result service.SecurityRuleConditions
	if err := json.Unmarshal(encoded, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func normalizeStringArray(value interface{}) []string {
	switch typed := value.(type) {
	case string:
		return []string{typed}
	case []interface{}:
		result := make([]string, 0, len(typed))
		for _, item := range typed {
			if str, ok := item.(string); ok {
				result = append(result, str)
			}
		}
		return result
	default:
		return []string{}
	}
}

func normalizeSecurityHeaders(value interface{}) []map[string]interface{} {
	items, ok := value.([]interface{})
	if !ok || items == nil {
		return []map[string]interface{}{}
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		typed, ok := item.(map[string]interface{})
		if !ok || typed == nil {
			continue
		}

		header := map[string]interface{}{}
		if name, ok := typed["name"].(string); ok {
			header["name"] = name
		}
		header["value_regex"] = normalizeStringArray(typed["value_regex"])
		result = append(result, header)
	}

	return result
}
