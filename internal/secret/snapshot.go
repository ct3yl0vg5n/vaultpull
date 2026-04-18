package secret

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a point-in-time capture of secrets.
type Snapshot struct {
	Timestamp time.Time         `json:"timestamp"`
	Label     string            `json:"label,omitempty"`
	Secrets   map[string]string `json:"secrets"`
}

// SaveSnapshot writes a snapshot of the given secrets map to a JSON file.
func SaveSnapshot(path, label string, secrets map[string]string) error {
	s := Snapshot{
		Timestamp: time.Now().UTC(),
		Label:     label,
		Secrets:   secrets,
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("snapshot write: %w", err)
	}
	return nil
}

// LoadSnapshot reads a snapshot from a JSON file.
func LoadSnapshot(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("snapshot read: %w", err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("snapshot parse: %w", err)
	}
	return &s, nil
}

// FormatSnapshotSummary returns a human-readable summary of a snapshot.
func FormatSnapshotSummary(s *Snapshot) string {
	if s == nil {
		return "no snapshot found"
	}
	label := s.Label
	if label == "" {
		label = "(unlabeled)"
	}
	return fmt.Sprintf("snapshot: %s | label: %s | keys: %d", s.Timestamp.Format(time.RFC3339), label, len(s.Secrets))
}
