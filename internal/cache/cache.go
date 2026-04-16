package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Entry holds cached secrets with metadata.
type Entry struct {
	Secrets   map[string]string `json:"secrets"`
	FetchedAt time.Time         `json:"fetched_at"`
	TTL       time.Duration     `json:"ttl"`
}

// IsExpired returns true if the cache entry is past its TTL.
func (e *Entry) IsExpired() bool {
	if e.TTL <= 0 {
		return true
	}
	return time.Since(e.FetchedAt) > e.TTL
}

// ExpiresAt returns the absolute time when the cache entry expires.
func (e *Entry) ExpiresAt() time.Time {
	return e.FetchedAt.Add(e.TTL)
}

// Store writes an entry to the given file path as JSON.
func Store(path string, secrets map[string]string, ttl time.Duration) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	entry := Entry{
		Secrets:   secrets,
		FetchedAt: time.Now(),
		TTL:       ttl,
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(entry)
}

// Load reads a cache entry from disk. Returns nil, nil if file does not exist.
func Load(path string) (*Entry, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var entry Entry
	if err := json.NewDecoder(f).Decode(&entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

// Invalidate removes the cache file if it exists.
func Invalidate(path string) error {
	err := os.Remove(path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
