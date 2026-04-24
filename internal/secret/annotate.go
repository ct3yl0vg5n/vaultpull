package secret

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// AnnotateOptions controls annotation behavior.
type AnnotateOptions struct {
	Path string
}

// DefaultAnnotateOptions returns sensible defaults.
func DefaultAnnotateOptions() AnnotateOptions {
	return AnnotateOptions{
		Path: ".vaultpull_annotations.json",
	}
}

// Annotation holds metadata attached to a secret key.
type Annotation struct {
	Key       string    `json:"key"`
	Note      string    `json:"note"`
	Owner     string    `json:"owner,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AnnotationStore maps secret keys to their annotations.
type AnnotationStore map[string]Annotation

// LoadAnnotations reads annotations from disk, returning an empty store if the file does not exist.
func LoadAnnotations(path string) (AnnotationStore, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return make(AnnotationStore), nil
	}
	if err != nil {
		return nil, fmt.Errorf("annotate: read %s: %w", path, err)
	}
	var store AnnotationStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("annotate: parse %s: %w", path, err)
	}
	return store, nil
}

// SaveAnnotations persists the annotation store to disk.
func SaveAnnotations(path string, store AnnotationStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("annotate: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("annotate: write %s: %w", path, err)
	}
	return nil
}

// Annotate adds or updates an annotation for the given key.
func Annotate(store AnnotationStore, key, note, owner string) AnnotationStore {
	if key == "" {
		return store
	}
	store[key] = Annotation{
		Key:       key,
		Note:      note,
		Owner:     owner,
		UpdatedAt: time.Now().UTC(),
	}
	return store
}

// RemoveAnnotation deletes the annotation for the given key.
func RemoveAnnotation(store AnnotationStore, key string) AnnotationStore {
	delete(store, key)
	return store
}

// FormatAnnotationReport returns a human-readable summary of all annotations.
func FormatAnnotationReport(store AnnotationStore) string {
	if len(store) == 0 {
		return "no annotations found\n"
	}
	out := ""
	for _, a := range store {
		owner := a.Owner
		if owner == "" {
			owner = "(none)"
		}
		out += fmt.Sprintf("%-24s  owner=%-16s  note=%s  updated=%s\n",
			a.Key, owner, a.Note, a.UpdatedAt.Format(time.RFC3339))
	}
	return out
}
