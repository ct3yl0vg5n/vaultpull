package secret

import (
	"fmt"
	"regexp"
	"strings"
)

// PolicyRule defines a single enforcement rule for secret keys or values.
type PolicyRule struct {
	Name        string
	Description string
	KeyPattern  string // optional regex that key must match
	DenyPattern string // optional regex that value must NOT match
	Required    bool   // key must be present
}

// PolicyOptions configures policy enforcement.
type PolicyOptions struct {
	Rules []PolicyRule
}

// PolicyViolation describes a single rule breach.
type PolicyViolation struct {
	Key  string
	Rule string
	Msg  string
}

// DefaultPolicyOptions returns a sensible default policy.
func DefaultPolicyOptions() PolicyOptions {
	return PolicyOptions{
		Rules: []PolicyRule{
			{
				Name:        "no-plaintext-password",
				Description: "Values for PASSWORD keys must not be plaintext words",
				KeyPattern:  `(?i)password`,
				DenyPattern: `^(password|pass|secret|admin|root)$`,
			},
		},
	}
}

// EnforcePolicy checks secrets against all rules and returns violations.
func EnforcePolicy(secrets map[string]string, opts PolicyOptions) []PolicyViolation {
	var violations []PolicyViolation

	for _, rule := range opts.Rules {
		if rule.Required {
			if _, ok := secrets[rule.Name]; !ok {
				violations = append(violations, PolicyViolation{
					Key:  rule.Name,
					Rule: rule.Name,
					Msg:  "required key is missing",
				})
			}
		}

		for key, val := range secrets {
			if rule.KeyPattern != "" {
				matched, err := regexp.MatchString(rule.KeyPattern, key)
				if err != nil || !matched {
					continue
				}
			}
			if rule.DenyPattern != "" {
				matched, err := regexp.MatchString(rule.DenyPattern, strings.TrimSpace(val))
				if err == nil && matched {
					violations = append(violations, PolicyViolation{
						Key:  key,
						Rule: rule.Name,
						Msg:  fmt.Sprintf("value matches denied pattern: %s", rule.DenyPattern),
					})
				}
			}
		}
	}
	return violations
}

// FormatPolicyReport renders violations as a human-readable string.
func FormatPolicyReport(violations []PolicyViolation) string {
	if len(violations) == 0 {
		return "policy check passed: no violations found\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("policy violations (%d):\n", len(violations)))
	for _, v := range violations {
		sb.WriteString(fmt.Sprintf("  [%s] %s — %s\n", v.Rule, v.Key, v.Msg))
	}
	return sb.String()
}
