package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeResolveEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeResolveEnvFile: %v", err)
	}
	return p
}

func runResolveCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}
	resolveCmd.ResetFlags()
	resolveCmd.Flags().StringVarP(&resolveFile, "file", "f", ".env", "env file to resolve")
	resolveCmd.Flags().StringVar(&resolveRefPrefix, "prefix", "ref:", "prefix that marks a value as a reference")
	resolveCmd.Flags().BoolVar(&resolveStrict, "strict", true, "fail if a referenced key is missing")
	resolveCmd.SetOut(buf)
	resolveCmd.SetErr(buf)
	resolveCmd.SetArgs(args)
	_, err := resolveCmd.ExecuteC()
	return buf.String(), err
}

func TestResolveCmd_NoRefs(t *testing.T) {
	p := writeResolveEnvFile(t, "A=hello\nB=world\n")
	out, err := runResolveCmd(t, "--file", p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "0 reference(s)") {
		t.Errorf("expected 0 references in output, got: %s", out)
	}
}

func TestResolveCmd_ResolvesRef(t *testing.T) {
	p := writeResolveEnvFile(t, "BASE=https://example.com\nAPI=ref:BASE\n")
	_, err := runResolveCmd(t, "--file", p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(p)
	if !strings.Contains(string(data), "API=https://example.com") {
		t.Errorf("expected resolved value in file, got:\n%s", string(data))
	}
}

func TestResolveCmd_MissingRef_StrictFails(t *testing.T) {
	p := writeResolveEnvFile(t, "A=ref:MISSING\n")
	_, err := runResolveCmd(t, "--file", p, "--strict=true")
	if err == nil {
		t.Fatal("expected error for missing ref in strict mode")
	}
}

func TestResolveCmd_MissingFile(t *testing.T) {
	_, err := runResolveCmd(t, "--file", "/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestResolveCmd_FlagDefaults(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().StringVarP(&resolveFile, "file", "f", ".env", "")
	cmd.Flags().StringVar(&resolveRefPrefix, "prefix", "ref:", "")
	cmd.Flags().BoolVar(&resolveStrict, "strict", true, "")
	if resolveFile != ".env" {
		t.Errorf("expected default file .env, got %q", resolveFile)
	}
	if resolveRefPrefix != "ref:" {
		t.Errorf("expected default prefix ref:, got %q", resolveRefPrefix)
	}
}
