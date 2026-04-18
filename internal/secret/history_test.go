package secret

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempHistoryPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "history.json")
}

func TestLoadHistory_NotExist(t *testing.T) {
	h, err := LoadHistory("/nonexistent/path/history.json")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(h.Entries) != 0 {
		t.Errorf("expected empty log")
	}
}

func TestAppend_WritesEntries(t *testing.T) {
	path := tempHistoryPath(t)
	h := &HistoryLog{}
	err := h.Append(path, map[string]string{"FOO": "bar", "BAZ": "qux"})
	if err != nil {
		t.Fatalf("append failed: %v", err)
	}
	if len(h.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(h.Entries))
	}
}

func TestLoadHistory_ReadsFile(t *testing.T) {
	path := tempHistoryPath(t)
	log := HistoryLog{Entries: []HistoryEntry{
		{Key: "A", Value: "1", Timestamp: time.Now().UTC()},
	}}
	data, _ := json.MarshalIndent(log, "", "  ")
	os.WriteFile(path, data, 0600)

	h, err := LoadHistory(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if len(h.Entries) != 1 || h.Entries[0].Key != "A" {
		t.Errorf("unexpected entries: %+v", h.Entries)
	}
}

func TestFilterByKey(t *testing.T) {
	h := &HistoryLog{Entries: []HistoryEntry{
		{Key: "FOO", Value: "v1"},
		{Key: "BAR", Value: "v2"},
		{Key: "FOO", Value: "v3"},
	}}
	result := h.FilterByKey("FOO")
	if len(result) != 2 {
		t.Errorf("expected 2, got %d", len(result))
	}
}

func TestLoadHistory_InvalidJSON(t *testing.T) {
	path := tempHistoryPath(t)
	os.WriteFile(path, []byte("not json"), 0600)
	_, err := LoadHistory(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
