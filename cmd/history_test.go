package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/vaultpull/internal/secret"
)

func writeHistoryFile(t *testing.T, entries []secret.HistoryEntry) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "history.json")
	log := secret.HistoryLog{Entries: entries}
	data, _ := json.MarshalIndent(log, "", "  ")
	os.WriteFile(path, data, 0600)
	return path
}

func runHistoryCmd(args []string) (string, error) {
	buf := &bytes.Buffer{}
	historyCmd.SetOut(buf)
	historyCmd.SetErr(buf)
	historyCmd.SetArgs(args)
	err := historyCmd.Execute()
	return buf.String(), err
}

func TestHistoryCmd_NoEntries(t *testing.T) {
	path := writeHistoryFile(t, nil)
	out, err := runHistoryCmd([]string{"--history-file", path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No history") {
		t.Errorf("expected no-history message, got: %s", out)
	}
}

func TestHistoryCmd_ShowsEntries(t *testing.T) {
	entries := []secret.HistoryEntry{
		{Key: "FOO", Value: "bar", Timestamp: time.Now().UTC()},
	}
	path := writeHistoryFile(t, entries)
	out, err := runHistoryCmd([]string{"--history-file", path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "FOO") {
		t.Errorf("expected FOO in output, got: %s", out)
	}
}

func TestHistoryCmd_FilterByKey(t *testing.T) {
	entries := []secret.HistoryEntry{
		{Key: "FOO", Value: "v1", Timestamp: time.Now().UTC()},
		{Key: "BAR", Value: "v2", Timestamp: time.Now().UTC()},
	}
	path := writeHistoryFile(t, entries)
	out, err := runHistoryCmd([]string{"--history-file", path, "--key", "BAR"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "FOO") {
		t.Errorf("FOO should be filtered out, got: %s", out)
	}
	if !strings.Contains(out, "BAR") {
		t.Errorf("expected BAR in output, got: %s", out)
	}
}
