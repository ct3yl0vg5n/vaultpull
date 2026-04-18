package secret

import (
	"testing"
	"time"
)

func fixedNow() time.Time {
	return time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
}

func TestCheckExpiration_NoEntries(t *testing.T) {
	opts := DefaultExpireOptions()
	opts.NowFunc = fixedNow
	r := CheckExpiration(map[string]time.Time{}, opts)
	if r.Expired != 0 || r.Warning != 0 {
		t.Fatalf("expected no issues, got %+v", r)
	}
}

func TestCheckExpiration_Expired(t *testing.T) {
	opts := DefaultExpireOptions()
	opts.NowFunc = fixedNow
	secrets := map[string]time.Time{
		"OLD_KEY": fixedNow().Add(-24 * time.Hour),
	}
	r := CheckExpiration(secrets, opts)
	if r.Expired != 1 {
		t.Fatalf("expected 1 expired, got %d", r.Expired)
	}
	if !r.Entries[0].Expired {
		t.Error("entry should be marked expired")
	}
}

func TestCheckExpiration_WarningSoon(t *testing.T) {
	opts := DefaultExpireOptions()
	opts.NowFunc = fixedNow
	secrets := map[string]time.Time{
		"SOON_KEY": fixedNow().Add(3 * 24 * time.Hour),
	}
	r := CheckExpiration(secrets, opts)
	if r.Warning != 1 {
		t.Fatalf("expected 1 warning, got %d", r.Warning)
	}
	if !r.Entries[0].Warn {
		t.Error("entry should be marked warn")
	}
}

func TestCheckExpiration_Healthy(t *testing.T) {
	opts := DefaultExpireOptions()
	opts.NowFunc = fixedNow
	secrets := map[string]time.Time{
		"HEALTHY": fixedNow().Add(30 * 24 * time.Hour),
	}
	r := CheckExpiration(secrets, opts)
	if r.Expired != 0 || r.Warning != 0 {
		t.Fatalf("expected healthy, got %+v", r)
	}
}

func TestFormatExpireReport_AllGood(t *testing.T) {
	r := ExpireReport{}
	if msg := FormatExpireReport(r); msg == "" {
		t.Error("expected non-empty message")
	}
}

func TestFormatExpireReport_Issues(t *testing.T) {
	r := ExpireReport{Expired: 2, Warning: 1}
	msg := FormatExpireReport(r)
	if msg == "" {
		t.Error("expected non-empty message")
	}
}
