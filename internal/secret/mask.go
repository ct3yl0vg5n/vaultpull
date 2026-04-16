package secret

import "strings"

// MaskOptions controls how secret values are masked in output.
type MaskOptions struct {
	Enabled    bool
	MaskChar   string
	RevealLast int
}

// DefaultMaskOptions returns sensible defaults for masking.
func DefaultMaskOptions() MaskOptions {
	return MaskOptions{
		Enabled:    true,
		MaskChar:   "*",
		RevealLast: 4,
	}
}

// Mask replaces a secret value with a masked representation.
func Mask(value string, opts MaskOptions) string {
	if !opts.Enabled || value == "" {
		return value
	}
	reveal := opts.RevealLast
	if reveal < 0 {
		reveal = 0
	}
	if reveal >= len(value) {
		return strings.Repeat(opts.MaskChar, len(value))
	}
	masked := strings.Repeat(opts.MaskChar, len(value)-reveal)
	return masked + value[len(value)-reveal:]
}

// MaskMap applies Mask to all values in a map.
func MaskMap(secrets map[string]string, opts MaskOptions) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = Mask(v, opts)
	}
	return out
}
