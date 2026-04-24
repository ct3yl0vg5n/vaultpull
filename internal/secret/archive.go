package secret

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// ArchiveEntry represents a single archived snapshot of secrets.
type ArchiveEntry struct {
	Label     string            `json:"label"`
	Timestamp time.Time         `json:"timestamp"`
	Secrets   map[string]string `json:"secrets"`
}

// DefaultArchiveOptions returns sensible defaults.
func DefaultArchiveOptions() ArchiveOptions {
	return ArchiveOptions{
		MaxEntries: 10,
	}
}

// ArchiveOptions controls archive behaviour.
type ArchiveOptions struct {
	MaxEntries int
}

// LoadArchive reads all archive entries from the given path.
// Returns an empty slice if the file does not exist.
func LoadArchive(path string) ([]ArchiveEntry, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []ArchiveEntry{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("archive: read %q: %w", path, err)
	}
	var entries []ArchiveEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("archive: parse %q: %w", path, err)
	}
	return entries, nil
}

// AppendArchive adds a new entry and trims to MaxEntries, then persists.
func AppendArchive(path, label string, secrets map[string]string, opts ArchiveOptions) error {
	entries, err := LoadArchive(path)
	if err != nil {
		return err
	}
	copy := make(map[string]string, len(secrets))
	for k, v := range secrets {
		copy[k] = v
	}
	entries = append(entries, ArchiveEntry{
		Label:     label,
		Timestamp: time.Now().UTC(),
		Secrets:   copy,
	})
	if opts.MaxEntries > 0 && len(entries) > opts.MaxEntries {
		entries = entries[len(entries)-opts.MaxEntries:]
	}
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("archive: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("archive: write %q: %w", path, err)
	}
	return nil
}

// FormatArchiveSummary returns a human-readable summary of all archive entries.
func FormatArchiveSummary(entries []ArchiveEntry) string {
	if len(entries) == 0 {
		return "no archive entries found\n"
	}
	out := fmt.Sprintf("%-4s  %-20s  %-16s  %s\n", "#", "timestamp", "label", "keys")
	for i, e := range entries {
		keys := make([]string, 0, len(e.Secrets))
		for k := range e.Secrets {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		out += fmt.Sprintf("%-4d  %-20s  %-16s  %d\n",
			i+1,
			e.Timestamp.Format(time.RFC3339),
			e.Label,
			len(keys),
		)
	}
	return out
}
