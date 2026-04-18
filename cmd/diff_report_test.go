package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeDiffEnvFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func runDiffReportCmd(args []string) (string, error) {
	buf := &bytes.Buffer{}
	diffReportCmd.SetOut(buf)
	diffReportCmd.SetErr(buf)
	diffReportCmd.SetArgs(args)
	err := diffReportCmd.Execute()
	return buf.String(), err
}

func TestDiffReportCmd_AddedKey(t *testing.T) {
	dir := t.TempDir()
	base := writeDiffEnvFile(t, dir, "base.env", "EXISTING=val\n")
	next := writeDiffEnvFile(t, dir, "next.env", "EXISTING=val\nNEW_KEY=secret\n")

	out, err := runDiffReportCmd([]string{"--file", base, "--compare", next, "--redact=false"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "+ NEW_KEY") {
		t.Errorf("expected added key in output, got: %s", out)
	}
}

func TestDiffReportCmd_RemovedKey(t *testing.T) {
	dir := t.TempDir()
	base := writeDiffEnvFile(t, dir, "base.env", "GONE=old\n")
	next := writeDiffEnvFile(t, dir, "next.env", "")

	out, err := runDiffReportCmd([]string{"--file", base, "--compare", next, "--redact=false"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "- GONE") {
		t.Errorf("expected removed key in output, got: %s", out)
	}
}

func TestDiffReportCmd_MissingCompareFlag(t *testing.T) {
	_, err := runDiffReportCmd([]string{})
	if err == nil {
		t.Error("expected error when --compare flag is missing")
	}
}

func TestDiffReportCmd_NoChanges(t *testing.T) {
	dir := t.TempDir()
	base := writeDiffEnvFile(t, dir, "base.env", "KEY=value\n")
	next := writeDiffEnvFile(t, dir, "next.env", "KEY=value\n")

	out, err := runDiffReportCmd([]string{"--file", base, "--compare", next, "--redact=false"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "KEY") {
		t.Errorf("expected unchanged key in output, got: %s", out)
	}
}
