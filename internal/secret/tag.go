package secret

import (
	"fmt"
	"strings"
)

// TagOptions controls tagging behavior.
type TagOptions struct {
	Tags      map[string][]string // key -> list of tags
	Separator string
}

// DefaultTagOptions returns sensible defaults.
func DefaultTagOptions() TagOptions {
	return TagOptions{
		Tags:      make(map[string][]string),
		Separator: ",",
	}
}

// TagEntry holds a key and its associated tags.
type TagEntry struct {
	Key  string
	Tags []string
}

// AddTag associates a tag with a key.
func AddTag(opts TagOptions, key, tag string) TagOptions {
	tag = strings.TrimSpace(tag)
	if key == "" || tag == "" {
		return opts
	}
	opts.Tags[key] = append(opts.Tags[key], tag)
	return opts
}

// GetTags returns tags for a given key.
func GetTags(opts TagOptions, key string) []string {
	return opts.Tags[key]
}

// FilterByTag returns keys that have the given tag.
func FilterByTag(opts TagOptions, tag string) []string {
	var result []string
	for key, tags := range opts.Tags {
		for _, t := range tags {
			if t == tag {
				result = append(result, key)
				break
			}
		}
	}
	return result
}

// FormatTagReport returns a human-readable tag summary.
func FormatTagReport(opts TagOptions) string {
	if len(opts.Tags) == 0 {
		return "no tags defined"
	}
	var sb strings.Builder
	for key, tags := range opts.Tags {
		sb.WriteString(fmt.Sprintf("  %s: %s\n", key, strings.Join(tags, opts.Separator+" ")))
	}
	return sb.String()
}
