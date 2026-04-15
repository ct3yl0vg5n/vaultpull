package diff

import "fmt"

// ChangeType represents the type of change between two env maps.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
	Unchanged ChangeType = "unchanged"
)

// Change represents a single key-level difference.
type Change struct {
	Key      string
	Type     ChangeType
	OldValue string
	NewValue string
}

// Compare computes the diff between a local env map and a remote env map.
// local is the existing .env contents; remote is the incoming Vault secrets.
func Compare(local, remote map[string]string) []Change {
	var changes []Change

	// Keys present in remote
	for key, newVal := range remote {
		oldVal, exists := local[key]
		if !exists {
			changes = append(changes, Change{
				Key:      key,
				Type:     Added,
				NewValue: newVal,
			})
		} else if oldVal != newVal {
			changes = append(changes, Change{
				Key:      key,
				Type:     Modified,
				OldValue: oldVal,
				NewValue: newVal,
			})
		}
	}

	// Keys only in local (removed in remote)
	for key, oldVal := range local {
		if _, exists := remote[key]; !exists {
			changes = append(changes, Change{
				Key:      key,
				Type:     Removed,
				OldValue: oldVal,
			})
		}
	}

	return changes
}

// Format returns a human-readable summary of a slice of Changes.
func Format(changes []Change) string {
	if len(changes) == 0 {
		return "No changes detected."
	}
	var out string
	for _, c := range changes {
		switch c.Type {
		case Added:
			out += fmt.Sprintf("+ %s=%s\n", c.Key, c.NewValue)
		case Removed:
			out += fmt.Sprintf("- %s=%s\n", c.Key, c.OldValue)
		case Modified:
			out += fmt.Sprintf("~ %s: %s -> %s\n", c.Key, c.OldValue, c.NewValue)
		}
	}
	return out
}

// HasChanges returns true if any non-unchanged entries exist.
func HasChanges(changes []Change) bool {
	return len(changes) > 0
}
