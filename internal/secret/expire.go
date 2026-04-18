package secret

import (
	"fmt"
	"time"
)

// ExpireOptions controls expiration checking behavior.
type ExpireOptions struct {
	WarnWithin time.Duration // warn if secret expires within this window
	NowFunc    func() time.Time
}

// ExpireEntry holds metadata about a single secret's expiration.
type ExpireEntry struct {
	Key       string
	ExpiresAt time.Time
	Expired   bool
	Warn      bool
}

// ExpireReport summarises expiration checks.
type ExpireReport struct {
	Entries []ExpireEntry
	Expired int
	Warning int
}

// DefaultExpireOptions returns sensible defaults.
func DefaultExpireOptions() ExpireOptions {
	return ExpireOptions{
		WarnWithin: 7 * 24 * time.Hour,
		NowFunc:    time.Now,
	}
}

// CheckExpiration evaluates a map of key→expiresAt timestamps.
func CheckExpiration(secrets map[string]time.Time, opts ExpireOptions) ExpireReport {
	if opts.NowFunc == nil {
		opts.NowFunc = time.Now
	}
	now := opts.NowFunc()
	report := ExpireReport{}
	for k, exp := range secrets {
		e := ExpireEntry{Key: k, ExpiresAt: exp}
		if now.After(exp) {
			e.Expired = true
			report.Expired++
		} else if exp.Sub(now) <= opts.WarnWithin {
			e.Warn = true
			report.Warning++
		}
		report.Entries = append(report.Entries, e)
	}
	return report
}

// FormatExpireReport returns a human-readable summary.
func FormatExpireReport(r ExpireReport) string {
	if r.Expired == 0 && r.Warning == 0 {
		return "All secrets are valid and not expiring soon."
	}
	return fmt.Sprintf("%d expired, %d expiring soon", r.Expired, r.Warning)
}
