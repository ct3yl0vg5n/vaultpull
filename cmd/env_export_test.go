package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeExportEnvFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "export-*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func runEnvExportCmd(t *testing.T, args []string) (string, error) {
	t.Helper()
	var out strings.Builder
	envExportCmd.SetOut(&out)
	envExportCmd.SetErr(&out)
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs(append([]string{"env-export"}, args...))
	err := rootCmd.Execute()
	// reset flags for next run
	envExportCmd.Flags().VisitAll(func(f *cobra.Flag) { f.Changed = false })
	return out.String(), err
}

func TestEnvExportCmd_BasicOutput(t *testing.T) {
	path := writeExportEnvFile(t, "FOO=bar\nBAZ=qux\n")
	out, err := runEnvExportCmd(t, []string{"--file", path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "FOO=bar") || !strings.Contains(out, "BAZ=qux") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestEnvExportCmd_WithExportFlag(t *testing.T) {
	path := writeExportEnvFile(t, "TOKEN=secret\n")
	out, err := runEnvExportCmd(t, []string{"--file", path, "--export"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "export TOKEN=secret") {
		t.Errorf("expected export declaration, got: %s", out)
	}
}

func TestEnvExportCmd_WithPrefix(t *testing.T) {
	path := writeExportEnvFile(t, "KEY=val\n")
	out, err := runEnvExportCmd(t, []string{"--file", path, "--prefix", "MY_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "MY_KEY=val") {
		t.Errorf("expected prefixed key, got: %s", out)
	}
}

func TestEnvExportCmd_MissingFile(t *testing.T) {
	_, err := runEnvExportCmd(t, []string{"--file", "/nonexistent/path.env"})
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestEnvExportCmd_SkipEmpty(t *testing.T) {
	path := writeExportEnvFile(t, "PRESENT=yes\nEMPTY=\n")
	out, err := runEnvExportCmd(t, []string{"--file", path, "--skip-empty"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "EMPTY=") {
		t.Errorf("expected EMPTY to be skipped, got: %s", out)
	}
	if !strings.Contains(out, "PRESENT=yes") {
		t.Errorf("expected PRESENT in output, got: %s", out)
	}
}
