package secret

import (
	"strings"
	"testing"
)

func TestGroup_NoSecrets(t *testing.T) {
	result := Group(map[string]string{}, DefaultGroupOptions())
	if len(result.Groups) != 0 {
		t.Fatalf("expected 0 groups, got %d", len(result.Groups))
	}
}

func TestGroup_AllDefault(t *testing.T) {
	secrets := map[string]string{
		"NOPREFIXKEY": "val",
	}
	result := Group(secrets, DefaultGroupOptions())
	if _, ok := result.Groups["default"]; !ok {
		t.Fatal("expected 'default' group")
	}
	if len(result.Groups["default"]) != 1 {
		t.Fatalf("expected 1 key in default group, got %d", len(result.Groups["default"]))
	}
}

func TestGroup_PrefixRouting(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST":   "localhost",
		"DB_PORT":   "5432",
		"APP_DEBUG": "true",
	}
	result := Group(secrets, DefaultGroupOptions())

	if len(result.Groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(result.Groups))
	}
	if len(result.Groups["DB"]) != 2 {
		t.Fatalf("expected 2 keys in DB group, got %d", len(result.Groups["DB"]))
	}
	if len(result.Groups["APP"]) != 1 {
		t.Fatalf("expected 1 key in APP group, got %d", len(result.Groups["APP"]))
	}
}

func TestGroup_MaxDepth2(t *testing.T) {
	secrets := map[string]string{
		"AWS_S3_BUCKET": "my-bucket",
		"AWS_S3_REGION": "us-east-1",
		"AWS_EC2_KEY":   "key-123",
	}
	opts := GroupOptions{Delimiter: "_", MaxDepth: 2}
	result := Group(secrets, opts)

	if _, ok := result.Groups["AWS_S3"]; !ok {
		t.Fatal("expected 'AWS_S3' group")
	}
	if _, ok := result.Groups["AWS_EC2"]; !ok {
		t.Fatal("expected 'AWS_EC2' group")
	}
}

func TestGroup_OrderIsSorted(t *testing.T) {
	secrets := map[string]string{
		"Z_KEY": "1",
		"A_KEY": "2",
		"M_KEY": "3",
	}
	result := Group(secrets, DefaultGroupOptions())
	if result.Order[0] != "A" || result.Order[1] != "M" || result.Order[2] != "Z" {
		t.Fatalf("expected sorted order A,M,Z got %v", result.Order)
	}
}

func TestFormatGroupReport_NoSecrets(t *testing.T) {
	result := Group(map[string]string{}, DefaultGroupOptions())
	out := FormatGroupReport(result)
	if !strings.Contains(out, "no secrets") {
		t.Fatalf("expected 'no secrets' message, got: %s", out)
	}
}

func TestFormatGroupReport_ShowsGroups(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	result := Group(secrets, DefaultGroupOptions())
	out := FormatGroupReport(result)
	if !strings.Contains(out, "[DB]") {
		t.Fatalf("expected '[DB]' in output, got: %s", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Fatalf("expected 'DB_HOST' in output, got: %s", out)
	}
}
