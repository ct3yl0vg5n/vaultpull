package secret

import (
	"os"
	"path/filepath"
	"testing"
)

func tempPinPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "pins.json")
}

func TestLoadPins_NotExist(t *testing.T) {
	store, err := LoadPins("/nonexistent/pins.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(store.Entries) != 0 {
		t.Errorf("expected empty store")
	}
}

func TestPin_AddsEntry(t *testing.T) {
	store := &PinStore{}
	Pin(store, "DB_PASS", "secret123", "alice")
	if len(store.Entries) != 1 {
		t.Fatalf("expected 1 entry")
	}
	if store.Entries[0].Key != "DB_PASS" || store.Entries[0].Value != "secret123" {
		t.Errorf("unexpected entry: %+v", store.Entries[0])
	}
}

func TestPin_UpdatesExisting(t *testing.T) {
	store := &PinStore{}
	Pin(store, "DB_PASS", "old", "alice")
	Pin(store, "DB_PASS", "new", "bob")
	if len(store.Entries) != 1 {
		t.Fatalf("expected 1 entry after update")
	}
	if store.Entries[0].Value != "new" {
		t.Errorf("expected updated value")
	}
}

func TestUnpin_RemovesEntry(t *testing.T) {
	store := &PinStore{}
	Pin(store, "API_KEY", "val", "")
	removed := Unpin(store, "API_KEY")
	if !removed {
		t.Errorf("expected true")
	}
	if len(store.Entries) != 0 {
		t.Errorf("expected empty store")
	}
}

func TestUnpin_NotFound(t *testing.T) {
	store := &PinStore{}
	if Unpin(store, "MISSING") {
		t.Errorf("expected false for missing key")
	}
}

func TestGetPinned_Found(t *testing.T) {
	store := &PinStore{}
	Pin(store, "TOKEN", "abc", "")
	val, ok := GetPinned(store, "TOKEN")
	if !ok || val != "abc" {
		t.Errorf("expected abc, got %q ok=%v", val, ok)
	}
}

func TestSaveAndLoadPins(t *testing.T) {
	path := tempPinPath(t)
	store := &PinStore{}
	Pin(store, "SECRET", "xyz", "ci")
	if err := SavePins(path, store); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadPins(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Entries) != 1 || loaded.Entries[0].Value != "xyz" {
		t.Errorf("unexpected loaded store: %+v", loaded)
	}
}

func TestLoadPins_InvalidJSON(t *testing.T) {
	path := tempPinPath(t)
	os.WriteFile(path, []byte("not json"), 0600)
	_, err := LoadPins(path)
	if err == nil {
		t.Errorf("expected error for invalid JSON")
	}
}
