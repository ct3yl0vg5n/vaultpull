package cache

import "time"

// Options holds cache configuration passed in from CLI flags or config.
type Options struct {
	Enabled  bool
	Path     string
	TTL      time.Duration
}

// DefaultOptions returns sensible defaults for cache behaviour.
func DefaultOptions() Options {
	return Options{
		Enabled: false,
		Path:    ".vaultpull.cache.json",
		TTL:     5 * time.Minute,
	}
}
