package audit

import (
	"fmt"
	"os"
	"time"
)

// RotateOptions controls log rotation behaviour.
type RotateOptions struct {
	// MaxBytes is the maximum file size before rotation. 0 means no limit.
	MaxBytes int64
	// Keep is the number of rotated files to retain (oldest deleted first).
	// 0 means keep all.
	Keep int
}

// Rotate renames the current log file to a timestamped backup and creates a
// fresh empty file in its place. It returns the backup path.
// If the file does not exist or is smaller than opts.MaxBytes the call is a
// no-op and returns an empty string.
func Rotate(path string, opts RotateOptions) (string, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("audit rotate: stat: %w", err)
	}

	if opts.MaxBytes > 0 && info.Size() < opts.MaxBytes {
		return "", nil
	}

	backup := fmt.Sprintf("%s.%s", path, time.Now().UTC().Format("20060102T150405Z"))
	if err := os.Rename(path, backup); err != nil {
		return "", fmt.Errorf("audit rotate: rename: %w", err)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return backup, fmt.Errorf("audit rotate: create new file: %w", err)
	}
	f.Close()

	return backup, nil
}
