package secret

import (
	"strings"
	"testing"
)

func TestAddTag_Basic(t *testing.T) {
	opts := DefaultTagOptions()
	opts = AddTag(opts, "DB_PASS", "sensitive")
	tags := GetTags(opts, "DB_PASS")
	if len(tags) != 1 || tags[0] != "sensitive" {
		t.Fatalf("expected [sensitive], got %v", tags)
	}
}

func TestAddTag_EmptyKeyOrTag_NoOp(t *testing.T) {
	opts := DefaultTagOptions()
	opts = AddTag(opts, "", "sensitive")
	opts = AddTag(opts, "KEY", "")
	if len(opts.Tags) != 0 {
		t.Fatalf("expected no tags, got %v", opts.Tags)
	}
}

func TestAddTag_MultipleTags(t *testing.T) {
	opts := DefaultTagOptions()
	opts = AddTag(opts, "API_KEY", "sensitive")
	opts = AddTag(opts, "API_KEY", "rotate")
	tags := GetTags(opts, "API_KEY")
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
}

func TestFilterByTag(t *testing.T) {
	opts := DefaultTagOptions()
	opts = AddTag(opts, "DB_PASS", "sensitive")
	opts = AddTag(opts, "API_KEY", "sensitive")
	opts = AddTag(opts, "APP_ENV", "config")

	keys := FilterByTag(opts, "sensitive")
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
}

func TestFilterByTag_NoMatch(t *testing.T) {
	opts := DefaultTagOptions()
	opts = AddTag(opts, "KEY", "config")
	keys := FilterByTag(opts, "sensitive")
	if len(keys) != 0 {
		t.Fatalf("expected 0 keys, got %d", len(keys))
	}
}

func TestFormatTagReport_Empty(t *testing.T) {
	opts := DefaultTagOptions()
	out := FormatTagReport(opts)
	if out != "no tags defined" {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestFormatTagReport_WithTags(t *testing.T) {
	opts := DefaultTagOptions()
	opts = AddTag(opts, "DB_PASS", "sensitive")
	out := FormatTagReport(opts)
	if !strings.Contains(out, "DB_PASS") || !strings.Contains(out, "sensitive") {
		t.Fatalf("unexpected output: %s", out)
	}
}
