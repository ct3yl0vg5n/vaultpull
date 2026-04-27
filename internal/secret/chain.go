package secret

import (
	"fmt"
	"strings"
)

// ChainStep represents a single transformation step in a pipeline.
type ChainStep struct {
	Name   string
	Action func(map[string]string) (map[string]string, error)
}

// ChainOptions configures the transformation chain.
type ChainOptions struct {
	StopOnError bool
	DryRun      bool
}

// DefaultChainOptions returns sensible defaults.
func DefaultChainOptions() ChainOptions {
	return ChainOptions{
		StopOnError: true,
		DryRun:      false,
	}
}

// ChainResult captures the outcome of each step.
type ChainResult struct {
	Step    string
	Applied bool
	Err     error
}

// RunChain executes a series of transformation steps sequentially.
func RunChain(src map[string]string, steps []ChainStep, opts ChainOptions) (map[string]string, []ChainResult, error) {
	current := copyMap(src)
	results := make([]ChainResult, 0, len(steps))

	for _, step := range steps {
		if opts.DryRun {
			results = append(results, ChainResult{Step: step.Name, Applied: false})
			continue
		}
		out, err := step.Action(current)
		if err != nil {
			results = append(results, ChainResult{Step: step.Name, Applied: false, Err: err})
			if opts.StopOnError {
				return current, results, fmt.Errorf("chain step %q failed: %w", step.Name, err)
			}
			continue
		}
		current = out
		results = append(results, ChainResult{Step: step.Name, Applied: true})
	}
	return current, results, nil
}

// FormatChainReport returns a human-readable summary of chain execution.
func FormatChainReport(results []ChainResult) string {
	if len(results) == 0 {
		return "no steps executed"
	}
	var sb strings.Builder
	sb.WriteString("chain results:\n")
	for _, r := range results {
		if r.Err != nil {
			fmt.Fprintf(&sb, "  [FAIL] %s: %v\n", r.Step, r.Err)
		} else if r.Applied {
			fmt.Fprintf(&sb, "  [OK]   %s\n", r.Step)
		} else {
			fmt.Fprintf(&sb, "  [SKIP] %s\n", r.Step)
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
