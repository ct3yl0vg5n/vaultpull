package secret

import (
	"fmt"
	"strings"
)

// PolicyReport holds the structured outcome of a policy check.
type PolicyReport struct {
	Passed     bool
	Violations []PolicyViolation
	Total      int
	Checked    int
}

// NewPolicyReport builds a PolicyReport from violations and the number of secrets checked.
func NewPolicyReport(violations []PolicyViolation, checked int) PolicyReport {
	return PolicyReport{
		Passed:     len(violations) == 0,
		Violations: violations,
		Total:      len(violations),
		Checked:    checked,
	}
}

// Summary returns a short one-line summary of the report.
func (r PolicyReport) Summary() string {
	if r.Passed {
		return fmt.Sprintf("OK: %d secrets checked, no violations", r.Checked)
	}
	return fmt.Sprintf("FAIL: %d violation(s) across %d secrets", r.Total, r.Checked)
}

// FormatPolicyReportDetailed renders a full multi-line report.
func FormatPolicyReportDetailed(r PolicyReport) string {
	var sb strings.Builder
	sb.WriteString(r.Summary() + "\n")
	for _, v := range r.Violations {
		sb.WriteString(fmt.Sprintf("  • [%s] %s: %s\n", v.Rule, v.Key, v.Msg))
	}
	return sb.String()
}
