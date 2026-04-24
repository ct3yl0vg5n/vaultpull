package secret

import (
	"fmt"
	"strings"
)

// RollbackReport wraps a RollbackResult with metadata for display.
type RollbackReport struct {
	Result  RollbackResult
	Total   int
	Changed int
}

// NewRollbackReport builds a RollbackReport from a RollbackResult.
func NewRollbackReport(r RollbackResult) RollbackReport {
	return RollbackReport{
		Result:  r,
		Total:   len(r.Applied) + len(r.Skipped),
		Changed: len(r.Applied),
	}
}

// Summary returns a one-line summary of the rollback operation.
func (rr RollbackReport) Summary() string {
	if rr.Changed == 0 {
		return "rollback: nothing to restore"
	}
	action := "restored"
	if rr.Result.DryRun {
		action = "would restore"
	}
	parts := []string{fmt.Sprintf("%s %d key(s)", action, rr.Changed)}
	if len(rr.Result.Skipped) > 0 {
		parts = append(parts, fmt.Sprintf("skipped %d", len(rr.Result.Skipped)))
	}
	return "rollback: " + strings.Join(parts, ", ")
}
