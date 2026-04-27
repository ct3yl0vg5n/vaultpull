package secret

import (
	"errors"
	"strings"
	"testing"
)

func uppercaseStep(src map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = strings.ToUpper(v)
	}
	return out, nil
}

func failStep(_ map[string]string) (map[string]string, error) {
	return nil, errors.New("step error")
}

func prefixStep(src map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = "pre_" + v
	}
	return out, nil
}

func TestRunChain_SingleStep(t *testing.T) {
	src := map[string]string{"KEY": "value"}
	steps := []ChainStep{{Name: "uppercase", Action: uppercaseStep}}
	out, results, err := RunChain(src, steps, DefaultChainOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "VALUE" {
		t.Errorf("expected VALUE, got %s", out["KEY"])
	}
	if len(results) != 1 || !results[0].Applied {
		t.Errorf("expected step to be applied")
	}
}

func TestRunChain_MultipleSteps(t *testing.T) {
	src := map[string]string{"K": "val"}
	steps := []ChainStep{
		{Name: "uppercase", Action: uppercaseStep},
		{Name: "prefix", Action: prefixStep},
	}
	out, _, err := RunChain(src, steps, DefaultChainOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["K"] != "pre_VAL" {
		t.Errorf("expected pre_VAL, got %s", out["K"])
	}
}

func TestRunChain_StopOnError(t *testing.T) {
	src := map[string]string{"K": "v"}
	steps := []ChainStep{
		{Name: "fail", Action: failStep},
		{Name: "uppercase", Action: uppercaseStep},
	}
	_, results, err := RunChain(src, steps, DefaultChainOptions())
	if err == nil {
		t.Fatal("expected error")
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestRunChain_ContinueOnError(t *testing.T) {
	src := map[string]string{"K": "v"}
	steps := []ChainStep{
		{Name: "fail", Action: failStep},
		{Name: "uppercase", Action: uppercaseStep},
	}
	opts := ChainOptions{StopOnError: false}
	out, results, err := RunChain(src, steps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["K"] != "V" {
		t.Errorf("expected V, got %s", out["K"])
	}
	if results[0].Err == nil {
		t.Error("expected first step to have error")
	}
}

func TestRunChain_DryRun_NoMutation(t *testing.T) {
	src := map[string]string{"K": "original"}
	steps := []ChainStep{{Name: "uppercase", Action: uppercaseStep}}
	opts := ChainOptions{DryRun: true}
	out, results, err := RunChain(src, steps, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["K"] != "original" {
		t.Errorf("expected original value in dry run, got %s", out["K"])
	}
	if results[0].Applied {
		t.Error("step should not be marked applied in dry run")
	}
}

func TestFormatChainReport_Mixed(t *testing.T) {
	results := []ChainResult{
		{Step: "step1", Applied: true},
		{Step: "step2", Applied: false, Err: errors.New("oops")},
		{Step: "step3", Applied: false},
	}
	report := FormatChainReport(results)
	if !strings.Contains(report, "[OK]") {
		t.Error("expected OK marker")
	}
	if !strings.Contains(report, "[FAIL]") {
		t.Error("expected FAIL marker")
	}
	if !strings.Contains(report, "[SKIP]") {
		t.Error("expected SKIP marker")
	}
}
