package secret

import "strings"

// RedactOptions controls how values are redacted in output.
type RedactOptions struct {
	Enabled    bool
	MaskChar   string
	RevealChars int // number of trailing chars to reveal
}

// DefaultRedactOptions returns sensible defaults.
func DefaultRedactOptions() RedactOptions {
	return RedactOptions{
		Enabled:     true,
		MaskChar:    "*",
		RevealChars: 4,
	}
}

// Redact masks a secret value according to the given options.
func Redact(value string, opts RedactOptions) string {
	if !opts.Enabled || value == "" {
		return value
	}
	mask := opts.MaskChar
	if mask == "" {
		mask = "*"
	}
	revealed := opts.RevealChars
	if revealed < 0 {
		revealed = 0
	}
	if revealed >= len(value) {
		return strings.Repeat(mask, len(value))
	}
	masked := strings.Repeat(mask, len(value)-revealed)
	return masked + value[len(value)-revealed:]
}

// RedactMap applies Redact to every value in a map.
func RedactMap(secrets map[string]string, opts RedactOptions) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = Redact(v, opts)
	}
	return out
}
