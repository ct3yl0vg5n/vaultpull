package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeScopeEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func runScopeCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}
	// Re-create command to avoid flag pollution between tests.
	cmd := &cobra.Command{Use: "scope", RunE: runScope}
	cmd.Flags().StringVarP(&scopeFile, "file", "f", ".env", "")
	cmd.Flags().StringArrayVarP(&scopeDefinitions, "scope", "s", nil, "")
	cmd.Flags().StringVar(&scopeDefault, "default", "default", "")
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	os.Stdout = os.Stdout // keep stdout for FormatScopeReport
	cmd.SetArgs(args)
	err := cmd.Execute()
	return buf.String(), err
}

func TestScopeCmd_AllDefault(t *testing.T) {
	p := writeScopeEnvFile(t, "FOO=1\nBAR=2\n")
	_, err := runScopeCmd(t, "--file", p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestScopeCmd_WithScopeFlag(t *testing.T) {
	p := writeScopeEnvFile(t, "DB_HOST=localhost\nDB_PORT=5432\nAPP_KEY=secret\n")
	// Reset slice between tests
	scopeDefinitions = nil
	_, err := runScopeCmd(t, "--file", p, "--scope", "database=DB_", "--scope", "app=APP_")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestScopeCmd_MissingFile(t *testing.T) {
	scopeDefinitions = nil
	_, err := runScopeCmd(t, "--file", "/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestScopeCmd_InvalidScopeDefinition(t *testing.T) {
	p := writeScopeEnvFile(t, "FOO=1\n")
	scopeDefinitions = nil
	_, err := runScopeCmd(t, "--file", p, "--scope", "badscopenoequals")
	if err == nil {
		t.Fatal("expected error for invalid scope definition")
	}
	if !strings.Contains(err.Error(), "invalid scope definition") {
		t.Errorf("unexpected error message: %v", err)
	}
}
