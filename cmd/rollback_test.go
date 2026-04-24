package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/example/vaultpull/internal/secret"
)

func writeRollbackEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func writeRollbackSnapshot(t *testing.T, data map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "snap.json")
	b, _ := json.Marshal(data)
	if err := os.WriteFile(p, b, 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func runRollbackCmd(args ...string) (string, error) {
	buf := &bytes.Buffer{}
	rollbackCmd.SetOut(buf)
	rollbackCmd.SetErr(buf)
	rollbackCmd.ResetFlags()
	init() // re-register flags
	rollbackCmd.SetArgs(args)
	err := rollbackCmd.Execute()
	return buf.String(), err
}

func TestRollbackCmd_DryRun_ShowsChanges(t *testing.T) {
	envFile := writeRollbackEnvFile(t, "KEY=new\n")
	snapFile := writeRollbackSnapshot(t, map[string]string{"KEY": "old"})

	var buf bytes.Buffer
	rollbackCmd.SetOut(&buf)
	rollbackCmd.SetArgs([]string{"--file", envFile, "--snapshot", snapFile, "--dry-run"})
	if err := rollbackCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "dry-run") {
		t.Errorf("expected dry-run in output, got: %s", out)
	}
	if !strings.Contains(out, "KEY") {
		t.Errorf("expected KEY in output, got: %s", out)
	}
}

func TestRollbackCmd_MissingSnapshotFlag(t *testing.T) {
	envFile := writeRollbackEnvFile(t, "KEY=val\n")
	rollbackCmd.SetArgs([]string{"--file", envFile})
	if err := rollbackCmd.Execute(); err == nil {
		t.Error("expected error when --snapshot is missing")
	}
}

func TestRollbackCmd_AppliesChanges(t *testing.T) {
	envFile := writeRollbackEnvFile(t, "KEY=new\n")
	snapFile := writeRollbackSnapshot(t, map[string]string{"KEY": "old"})

	_ = secret.DefaultRollbackOptions() // ensure package used

	var buf bytes.Buffer
	rollbackCmd.SetOut(&buf)
	rollbackCmd.SetArgs([]string{"--file", envFile, "--snapshot", snapFile})
	if err := rollbackCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(envFile)
	if !strings.Contains(string(data), "old") {
		t.Errorf("expected env file to contain restored value 'old', got: %s", data)
	}
}
