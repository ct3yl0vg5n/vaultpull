package secret

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func tempAnnotatePath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "annotations.json")
}

func TestLoadAnnotations_NotExist(t *testing.T) {
	store, err := LoadAnnotations("/nonexistent/path/annotations.json")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(store) != 0 {
		t.Fatalf("expected empty store, got %d entries", len(store))
	}
}

func TestAnnotate_AddsEntry(t *testing.T) {
	store := make(AnnotationStore)
	store = Annotate(store, "DB_PASSWORD", "rotated quarterly", "alice")
	a, ok := store["DB_PASSWORD"]
	if !ok {
		t.Fatal("expected annotation to be present")
	}
	if a.Note != "rotated quarterly" {
		t.Errorf("unexpected note: %s", a.Note)
	}
	if a.Owner != "alice" {
		t.Errorf("unexpected owner: %s", a.Owner)
	}
	if a.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestAnnotate_EmptyKey_NoOp(t *testing.T) {
	store := make(AnnotationStore)
	store = Annotate(store, "", "some note", "bob")
	if len(store) != 0 {
		t.Fatalf("expected no entries for empty key, got %d", len(store))
	}
}

func TestAnnotate_UpdatesExisting(t *testing.T) {
	store := make(AnnotationStore)
	store = Annotate(store, "API_KEY", "first note", "alice")
	before := store["API_KEY"].UpdatedAt
	time.Sleep(2 * time.Millisecond)
	store = Annotate(store, "API_KEY", "updated note", "bob")
	if store["API_KEY"].Note != "updated note" {
		t.Errorf("note not updated")
	}
	if !store["API_KEY"].UpdatedAt.After(before) {
		t.Errorf("UpdatedAt should be newer after update")
	}
}

func TestRemoveAnnotation(t *testing.T) {
	store := make(AnnotationStore)
	store = Annotate(store, "SECRET", "note", "")
	store = RemoveAnnotation(store, "SECRET")
	if _, ok := store["SECRET"]; ok {
		t.Error("expected annotation to be removed")
	}
}

func TestSaveAndLoadAnnotations(t *testing.T) {
	path := tempAnnotatePath(t)
	store := make(AnnotationStore)
	store = Annotate(store, "TOKEN", "service account token", "ops")
	if err := SaveAnnotations(path, store); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadAnnotations(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded["TOKEN"].Note != "service account token" {
		t.Errorf("unexpected note after reload: %s", loaded["TOKEN"].Note)
	}
}

func TestLoadAnnotations_InvalidJSON(t *testing.T) {
	path := tempAnnotatePath(t)
	_ = os.WriteFile(path, []byte("not json{"), 0o600)
	_, err := LoadAnnotations(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestFormatAnnotationReport_Empty(t *testing.T) {
	out := FormatAnnotationReport(make(AnnotationStore))
	if !strings.Contains(out, "no annotations") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatAnnotationReport_WithEntries(t *testing.T) {
	store := make(AnnotationStore)
	store = Annotate(store, "DB_URL", "primary database", "dba")
	out := FormatAnnotationReport(store)
	if !strings.Contains(out, "DB_URL") {
		t.Errorf("expected key in output: %s", out)
	}
	if !strings.Contains(out, "primary database") {
		t.Errorf("expected note in output: %s", out)
	}
}
