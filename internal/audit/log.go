package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Operation string    `json:"operation"`
	Mount     string    `json:"mount"`
	Path      string    `json:"path"`
	EnvFile   string    `json:"env_file"`
	DryRun    bool      `json:"dry_run"`
	Added     int       `json:"added"`
	Removed   int       `json:"removed"`
	Modified  int       `json:"modified"`
	Error     string    `json:"error,omitempty"`
}

// Logger writes audit entries to a file.
type Logger struct {
	path string
}

// NewLogger creates a Logger that appends to the given file path.
// Pass an empty string to disable logging (no-op logger).
func NewLogger(path string) *Logger {
	return &Logger{path: path}
}

// Write appends the entry as a JSON line to the audit log file.
// If the logger has no path configured, Write is a no-op.
func (l *Logger) Write(e Entry) error {
	if l.path == "" {
		return nil
	}

	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}

	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("audit: open log file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	if err := enc.Encode(e); err != nil {
		return fmt.Errorf("audit: encode entry: %w", err)
	}
	return nil
}
