package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func runGenerateCmd(args ...string) (string, error) {
	buf := &bytes.Buffer{}
	RootCmd.SetOut(buf)
	RootCmd.SetErr(buf)
	RootCmd.SetArgs(append([]string{"generate"}, args...))
	err := RootCmd.Execute()
	return buf.String(), err
}

func TestGenerateCmd_DefaultOutput(t *testing.T) {
	out, err := runGenerateCmd()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := strings.TrimSpace(out)
	if len(result) == 0 {
		t.Error("expected non-empty output")
	}
}

func TestGenerateCmd_WithKey(t *testing.T) {
	out, err := runGenerateCmd("--key", "MY_SECRET")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := strings.TrimSpace(out)
	if !strings.HasPrefix(result, "MY_SECRET=") {
		t.Errorf("expected MY_SECRET= prefix, got %q", result)
	}
}

func TestGenerateCmd_CustomLength(t *testing.T) {
	out, err := runGenerateCmd("--length", "16")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := strings.TrimSpace(out)
	if len(result) != 16 {
		t.Errorf("expected length 16, got %d", len(result))
	}
}

func TestGenerateCmd_ZeroLength_Fails(t *testing.T) {
	_, err := runGenerateCmd("--length", "0")
	if err == nil {
		t.Error("expected error for zero length")
	}
}
