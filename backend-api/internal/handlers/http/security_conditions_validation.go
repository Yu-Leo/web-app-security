package http

import (
	"errors"
	"fmt"
	"net/netip"
	"regexp"
	"strings"

	"github.com/Yu-Leo/web-app-security/backend-api/internal/generated/service"
)

func validateSecurityConditions(conditions *service.SecurityRuleConditions) error {
	if conditions == nil {
		return nil
	}

	if conditions.SourceIpCidr != nil {
		for _, cidr := range *conditions.SourceIpCidr {
			if _, err := netip.ParsePrefix(cidr); err != nil {
				return fmt.Errorf("invalid source_ip_cidr value: %s", cidr)
			}
		}
	}

	if conditions.UriRegex != nil {
		for _, pattern := range *conditions.UriRegex {
			if err := validateRegexp(pattern); err != nil {
				return fmt.Errorf("invalid uri_regex: %w", err)
			}
		}
	}
	if conditions.HostRegex != nil {
		for _, pattern := range *conditions.HostRegex {
			if err := validateRegexp(pattern); err != nil {
				return fmt.Errorf("invalid host_regex: %w", err)
			}
		}
	}
	if conditions.MethodRegex != nil {
		for _, pattern := range *conditions.MethodRegex {
			if err := validateRegexp(pattern); err != nil {
				return fmt.Errorf("invalid method_regex: %w", err)
			}
		}
	}

	if conditions.Headers != nil {
		for _, header := range *conditions.Headers {
			if strings.TrimSpace(header.Name) == "" {
				return errors.New("header condition name must not be empty")
			}
			if len(header.ValueRegex) == 0 {
				return fmt.Errorf("header %q must contain at least one value_regex", header.Name)
			}
			for _, pattern := range header.ValueRegex {
				if err := validateRegexp(pattern); err != nil {
					return fmt.Errorf("invalid header %q value_regex: %w", header.Name, err)
				}
			}
		}
	}

	return nil
}

func validateRegexp(pattern string) error {
	if _, err := regexp.Compile(pattern); err != nil {
		return err
	}
	return nil
}
