package secret

import (
	"strings"
	"testing"
)

func TestDetectChanges_NoChanges(t *testing.T) {
	prev := map[string]string{"A": "1", "B": "2"}
	curr := map[string]string{"A": "1", "B": "2"}
	events := DetectChanges(prev, curr, DefaultWatchOptions())
	if len(events) != 0 {
		t.Fatalf("expected 0 events, got %d", len(events))
	}
}

func TestDetectChanges_Added(t *testing.T) {
	prev := map[string]string{}
	curr := map[string]string{"NEW_KEY": "hello"}
	events := DetectChanges(prev, curr, DefaultWatchOptions())
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Kind != "added" || events[0].Key != "NEW_KEY" {
		t.Errorf("unexpected event: %+v", events[0])
	}
}

func TestDetectChanges_Removed(t *testing.T) {
	prev := map[string]string{"OLD_KEY": "val"}
	curr := map[string]string{}
	events := DetectChanges(prev, curr, DefaultWatchOptions())
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Kind != "removed" || events[0].Key != "OLD_KEY" {
		t.Errorf("unexpected event: %+v", events[0])
	}
}

func TestDetectChanges_Modified(t *testing.T) {
	prev := map[string]string{"X": "old"}
	curr := map[string]string{"X": "new"}
	events := DetectChanges(prev, curr, DefaultWatchOptions())
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e := events[0]
	if e.Kind != "modified" || e.OldValue != "old" || e.NewValue != "new" {
		t.Errorf("unexpected event: %+v", e)
	}
}

func TestDetectChanges_AlertOnNewDisabled(t *testing.T) {
	opts := DefaultWatchOptions()
	opts.AlertOnNew = false
	prev := map[string]string{}
	curr := map[string]string{"K": "v"}
	events := DetectChanges(prev, curr, opts)
	if len(events) != 0 {
		t.Fatalf("expected 0 events with AlertOnNew=false, got %d", len(events))
	}
}

func TestFormatWatchReport_NoChanges(t *testing.T) {
	out := FormatWatchReport(nil)
	if out != "no changes detected" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatWatchReport_WithEvents(t *testing.T) {
	events := []WatchEvent{
		{Key: "A", Kind: "added"},
		{Key: "B", Kind: "removed"},
		{Key: "C", Kind: "modified"},
	}
	out := FormatWatchReport(events)
	if !strings.Contains(out, "3 change(s)") {
		t.Errorf("expected change count in output, got: %s", out)
	}
	if !strings.Contains(out, "+ A") || !strings.Contains(out, "- B") || !strings.Contains(out, "~ C") {
		t.Errorf("missing expected symbols in output: %s", out)
	}
}
