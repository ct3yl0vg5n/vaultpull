package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writePolicyEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func runPolicyCmd(args []string) (string, error) {
	buf := &bytes.Buffer{}
	policyCmd.SetOut(buf)
	policyCmd.SetErr(buf)
	policyCmd.ResetFlags()
	policyCmd.Flags().StringVar(&policyFile, "file", ".env", "path to .env file")
	policyCmd.Flags().StringArrayVar(&policyRequireKeys, "require", nil, "required keys")
	policyCmd.Flags().StringArrayVar(&policyDenyValues, "deny-value", nil, "deny patterns")
	policyCmd.SetArgs(args)
	err := policyCmd.Execute()
	return buf.String(), err
}

func TestPolicyCmd_NoViolations(t *testing.T) {
	p := writePolicyEnvFile(t, "DB_PASSWORD=xK9#mP2$qR\n")
	out, err := runPolicyCmd([]string{"--file", p})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "passed") {
		t.Errorf("expected passed message, got: %s", out)
	}
}

func TestPolicyCmd_DenyViolation(t *testing.T) {
	p := writePolicyEnvFile(t, "DB_PASSWORD=password\n")
	out, _ := runPolicyCmd([]string{"--file", p})
	if !strings.Contains(out, "violation") {
		t.Errorf("expected violation in output, got: %s", out)
	}
}

func TestPolicyCmd_RequireKeyMissing(t *testing.T) {
	p := writePolicyEnvFile(t, "FOO=bar\n")
	out, _ := runPolicyCmd([]string{"--file", p, "--require", "API_KEY"})
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in output, got: %s", out)
	}
}

func TestPolicyCmd_RequireKeyPresent(t *testing.T) {
	p := writePolicyEnvFile(t, "API_KEY=abc123\n")
	out, err := runPolicyCmd([]string{"--file", p, "--require", "API_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "passed") {
		t.Errorf("expected passed, got: %s", out)
	}
}

func TestPolicyCmd_MissingFile(t *testing.T) {
	_, err := runPolicyCmd([]string{"--file", "/nonexistent/.env"})
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func resetPolicyCmd() {
	policyCmd = &cobra.Command{
		Use:   "policy",
		Short: "Enforce secret policy rules against a .env file",
		RunE:  runPolicy,
	}
}
