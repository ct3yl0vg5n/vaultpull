package secret

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// HistoryEntry records a snapshot of a secret value at a point in time.
type HistoryEntry struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

// HistoryLog is a collection of history entries.
type HistoryLog struct {
	Entries []HistoryEntry `json:"entries"`
}

// LoadHistory reads a history log from a JSON file.
func LoadHistory(path string) (*HistoryLog, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &HistoryLog{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read history: %w", err)
	}
	var log HistoryLog
	if err := json.Unmarshal(data, &log); err != nil {
		return nil, fmt.Errorf("parse history: %w", err)
	}
	return &log, nil
}

// Append adds a new entry for each key/value pair and saves the log.
func (h *HistoryLog) Append(path string, secrets map[string]string) error {
	now := time.Now().UTC()
	for k, v := range secrets {
		h.Entries = append(h.Entries, HistoryEntry{Key: k, Value: v, Timestamp: now})
	}
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal history: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// FilterByKey returns all entries matching the given key.
func (h *HistoryLog) FilterByKey(key string) []HistoryEntry {
	var out []HistoryEntry
	for _, e := range h.Entries {
		if e.Key == key {
			out = append(out, e)
		}
	}
	return out
}
