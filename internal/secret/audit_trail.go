package secret

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// AuditAction represents the type of operation performed on a secret.
type AuditAction string

const (
	AuditActionRead     AuditAction = "read"
	AuditActionWrite    AuditAction = "write"
	AuditActionDelete   AuditAction = "delete"
	AuditActionPromote  AuditAction = "promote"
	AuditActionRollback AuditAction = "rollback"
	AuditActionMerge    AuditAction = "merge"
)

// AuditEntry records a single auditable event for a secret key.
type AuditEntry struct {
	Timestamp time.Time   `json:"timestamp"`
	Action    AuditAction `json:"action"`
	Key       string      `json:"key"`
	Actor     string      `json:"actor,omitempty"`
	Note      string      `json:"note,omitempty"`
}

// AuditTrailOptions configures audit trail behaviour.
type AuditTrailOptions struct {
	// Path is the file where audit entries are persisted.
	Path string
	// Actor is the identity to record on each entry (e.g. $USER).
	Actor string
}

// DefaultAuditTrailOptions returns sensible defaults.
func DefaultAuditTrailOptions() AuditTrailOptions {
	return AuditTrailOptions{
		Path:  ".vaultpull_audit.jsonl",
		Actor: os.Getenv("USER"),
	}
}

// AppendAuditEntry appends a single audit entry to the JSONL file at opts.Path.
// If opts.Path is empty the call is a no-op.
func AppendAuditEntry(opts AuditTrailOptions, action AuditAction, key, note string) error {
	if opts.Path == "" {
		return nil
	}

	entry := AuditEntry{
		Timestamp: time.Now().UTC(),
		Action:    action,
		Key:       key,
		Actor:     opts.Actor,
		Note:      note,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}

	f, err := os.OpenFile(opts.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("audit: open file: %w", err)
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "%s\n", data)
	if err != nil {
		return fmt.Errorf("audit: write entry: %w", err)
	}
	return nil
}

// LoadAuditTrail reads all entries from the JSONL file at path.
// Returns an empty slice if the file does not exist.
func LoadAuditTrail(path string) ([]AuditEntry, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []AuditEntry{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("audit: read file: %w", err)
	}

	var entries []AuditEntry
	for _, line := range splitLines(data) {
		if len(line) == 0 {
			continue
		}
		var e AuditEntry
		if err := json.Unmarshal(line, &e); err != nil {
			return nil, fmt.Errorf("audit: parse entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// FormatAuditTrail returns a human-readable summary of the provided entries.
func FormatAuditTrail(entries []AuditEntry) string {
	if len(entries) == 0 {
		return "no audit entries found\n"
	}
	out := fmt.Sprintf("%-30s %-10s %-30s %s\n", "TIMESTAMP", "ACTION", "KEY", "ACTOR")
	for _, e := range entries {
		ts := e.Timestamp.Format(time.RFC3339)
		out += fmt.Sprintf("%-30s %-10s %-30s %s\n", ts, e.Action, e.Key, e.Actor)
		if e.Note != "" {
			out += fmt.Sprintf("  note: %s\n", e.Note)
		}
	}
	return out
}

// splitLines splits raw JSONL bytes into individual lines.
func splitLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, b := range data {
		if b == '\n' {
			lines = append(lines, data[start:i])
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}
