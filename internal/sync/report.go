package sync

import (
	"fmt"
	"io"

	"github.com/user/vaultpull/internal/diff"
)

// PrintReport writes a human-readable sync report to w.
func PrintReport(w io.Writer, result *Result) {
	if !diff.HasChanges(result.Changes) {
		fmt.Fprintln(w, "✔ No changes detected. Local env file is up to date.")
		return
	}

	fmt.Fprintf(w, "Diff for %s:\n", result.EnvFile)
	fmt.Fprintln(w, diff.Format(result.Changes))

	if result.Applied {
		fmt.Fprintf(w, "✔ Applied %d change(s) to %s\n", countChanges(result.Changes), result.EnvFile)
	} else {
		fmt.Fprintf(w, "ℹ Dry-run mode: %d change(s) not written.\n", countChanges(result.Changes))
	}
}

// countChanges returns the number of non-unchanged entries.
func countChanges(changes []diff.Change) int {
	count := 0
	for _, c := range changes {
		if c.Type != diff.Unchanged {
			count++
		}
	}
	return count
}
