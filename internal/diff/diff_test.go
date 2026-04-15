package diff

import (
	"strings"
	"testing"
)

func TestCompare_Added(t *testing.T) {
	local := map[string]string{}
	remote := map[string]string{"NEW_KEY": "value1"}

	changes := Compare(local, remote)

	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != Added {
		t.Errorf("expected Added, got %s", changes[0].Type)
	}
	if changes[0].Key != "NEW_KEY" {
		t.Errorf("unexpected key: %s", changes[0].Key)
	}
}

func TestCompare_Removed(t *testing.T) {
	local := map[string]string{"OLD_KEY": "oldval"}
	remote := map[string]string{}

	changes := Compare(local, remote)

	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != Removed {
		t.Errorf("expected Removed, got %s", changes[0].Type)
	}
}

func TestCompare_Modified(t *testing.T) {
	local := map[string]string{"KEY": "old"}
	remote := map[string]string{"KEY": "new"}

	changes := Compare(local, remote)

	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != Modified {
		t.Errorf("expected Modified, got %s", changes[0].Type)
	}
	if changes[0].OldValue != "old" || changes[0].NewValue != "new" {
		t.Errorf("unexpected values: old=%s new=%s", changes[0].OldValue, changes[0].NewValue)
	}
}

func TestCompare_Unchanged(t *testing.T) {
	local := map[string]string{"KEY": "same"}
	remote := map[string]string{"KEY": "same"}

	changes := Compare(local, remote)

	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestFormat_NoChanges(t *testing.T) {
	out := Format([]Change{})
	if out != "No changes detected." {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormat_WithChanges(t *testing.T) {
	changes := []Change{
		{Key: "A", Type: Added, NewValue: "1"},
		{Key: "B", Type: Removed, OldValue: "2"},
		{Key: "C", Type: Modified, OldValue: "old", NewValue: "new"},
	}
	out := Format(changes)
	if !strings.Contains(out, "+ A=1") {
		t.Errorf("missing added line in output: %s", out)
	}
	if !strings.Contains(out, "- B=2") {
		t.Errorf("missing removed line in output: %s", out)
	}
	if !strings.Contains(out, "~ C: old -> new") {
		t.Errorf("missing modified line in output: %s", out)
	}
}

func TestHasChanges(t *testing.T) {
	if HasChanges([]Change{}) {
		t.Error("expected false for empty changes")
	}
	if !HasChanges([]Change{{Key: "X", Type: Added}}) {
		t.Error("expected true for non-empty changes")
	}
}
