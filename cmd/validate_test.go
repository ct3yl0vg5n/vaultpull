package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func writeValidateEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func runValidateCmd(args []string) (string, error) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(append([]string{"validate"}, args...))
	err := rootCmd.Execute()
	return buf.String(), err
}

func resetValidateCmd() {
	for _, c := range rootCmd.Commands() {
		if c.Use == "validate" {
			c.ResetFlags()
			init()
			return
		}
	}
}

func TestValidateCmd_AllValid(t *testing.T) {
	p := writeValidateEnvFile(t, "MY_KEY=value\nOTHER=data\n")
	out, err := runValidateCmd([]string{"--file", p})
	if err != nil {
		t.Fatalf("unexpected error: %v, out: %s", err, out)
	}
	if out == "" {
		t.Log("no output but no error — ok")
	}
}

func TestValidateCmd_InvalidKey(t *testing.T) {
	p := writeValidateEnvFile(t, "lower_key=value\n")
	_, err := runValidateCmd([]string{"--file", p, "--strict-keys=true"})
	if err == nil {
		t.Fatal("expected error for invalid key")
	}
}

func TestValidateCmd_MissingFile(t *testing.T) {
	_, err := runValidateCmd([]string{"--file", "/nonexistent/.env"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestValidateCmd_NoStrictKeys(t *testing.T) {
	p := writeValidateEnvFile(t, "lower_key=value\n")
	_, err := runValidateCmd([]string{"--file", p, "--strict-keys=false"})
	if err != nil {
		t.Fatalf("unexpected error with strict-keys=false: %v", err)
	}
}

var _ *cobra.Command = validateCmd
