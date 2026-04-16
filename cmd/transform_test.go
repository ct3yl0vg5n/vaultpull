package cmd_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultpull/cmd"
)

func writeEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func runTransformCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	buf := new(strings.Builder)
	root := &cobra.Command{Use: "vaultpull"}
	cmd.RegisterTransformCmd(root)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(append([]string{"transform"}, args...))
	err := root.Execute()
	return buf.String(), err
}

func TestTransformCmd_Uppercase(t *testing.T) {
	p := writeEnvFile(t, "db_host=localhost\ndb_port=5432\n")
	out, err := runTransformCmd(t, "--file", p, "--uppercase")
	if err != nil {
		t.Fatalf("unexpected error: %v — output: %s", err, out)
	}
	data, _ := os.ReadFile(p)
	if !strings.Contains(string(data), "DB_HOST") {
		t.Errorf("expected uppercase keys, got:\n%s", string(data))
	}
}

func TestTransformCmd_MissingFile(t *testing.T) {
	_, err := runTransformCmd(t, "--file", "/nonexistent/.env", "--uppercase")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestTransformCmd_DefaultFile(t *testing.T) {
	// Ensure default flag value is .env
	root := &cobra.Command{Use: "vaultpull"}
	cmd.RegisterTransformCmd(root)
	root.SetArgs([]string{"transform", "--help"})
	buf := new(strings.Builder)
	root.SetOut(buf)
	_ = root.Execute()
	if !strings.Contains(buf.String(), ".env") {
		t.Error("expected default file .env in help output")
	}
}
