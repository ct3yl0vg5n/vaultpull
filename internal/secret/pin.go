package secret

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// PinEntry records a pinned version of a secret key.
type PinEntry struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	PinnedAt  time.Time `json:"pinned_at"`
	PinnedBy  string    `json:"pinned_by,omitempty"`
}

// PinStore holds all pinned secrets.
type PinStore struct {
	Entries []PinEntry `json:"entries"`
}

// LoadPins reads a pin store from disk. Returns empty store if file not found.
func LoadPins(path string) (*PinStore, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &PinStore{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read pin file: %w", err)
	}
	var store PinStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("parse pin file: %w", err)
	}
	return &store, nil
}

// SavePins writes the pin store to disk.
func SavePins(path string, store *PinStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal pins: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// Pin adds or updates a pinned entry for the given key.
func Pin(store *PinStore, key, value, pinnedBy string) {
	for i, e := range store.Entries {
		if e.Key == key {
			store.Entries[i] = PinEntry{Key: key, Value: value, PinnedAt: time.Now(), PinnedBy: pinnedBy}
			return
		}
	}
	store.Entries = append(store.Entries, PinEntry{Key: key, Value: value, PinnedAt: time.Now(), PinnedBy: pinnedBy})
}

// Unpin removes a pinned entry by key. Returns true if removed.
func Unpin(store *PinStore, key string) bool {
	for i, e := range store.Entries {
		if e.Key == key {
			store.Entries = append(store.Entries[:i], store.Entries[i+1:]...)
			return true
		}
	}
	return false
}

// GetPinned returns the pinned value for a key, and whether it exists.
func GetPinned(store *PinStore, key string) (string, bool) {
	for _, e := range store.Entries {
		if e.Key == key {
			return e.Value, true
		}
	}
	return "", false
}
