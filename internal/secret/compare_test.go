package secret

import (
	"testing"
)

func TestCompareMap_NoChanges(t *testing.T) {
	old := map[string]string{"A": "1", "B": "2"}
	new := map[string]string{"A": "1", "B": "2"}
	results := CompareMap(old, new, DefaultCompareOptions())
	for _, r := range results {
		if r.Changed {
			t.Errorf("expected no change for key %s", r.Key)
		}
	}
}

func TestCompareMap_Modified(t *testing.T) {
	old := map[string]string{"A": "old"}
	new := map[string]string{"A": "new"}
	results := CompareMap(old, new, DefaultCompareOptions())
	if len(results) != 1 || !results[0].Changed {
		t.Fatal("expected A to be changed")
	}
	if results[0].OldValue != "old" || results[0].NewValue != "new" {
		t.Errorf("unexpected values: %+v", results[0])
	}
}

func TestCompareMap_Added(t *testing.T) {
	old := map[string]string{}
	new := map[string]string{"X": "val"}
	results := CompareMap(old, new, DefaultCompareOptions())
	if len(results) != 1 || !results[0].Changed {
		t.Fatal("expected X to be marked changed (added)")
	}
}

func TestCompareMap_Removed(t *testing.T) {
	old := map[string]string{"Z": "gone"}
	new := map[string]string{}
	results := CompareMap(old, new, DefaultCompareOptions())
	if len(results) != 1 || !results[0].Changed {
		t.Fatal("expected Z to be marked changed (removed)")
	}
	if results[0].NewValue != "" {
		t.Errorf("expected empty new value for removed key")
	}
}

func TestCompareMap_IgnoreKeys(t *testing.T) {
	old := map[string]string{"A": "1", "SKIP": "old"}
	new := map[string]string{"A": "1", "SKIP": "new"}
	opts := DefaultCompareOptions()
	opts.IgnoreKeys = []string{"SKIP"}
	results := CompareMap(old, new, opts)
	for _, r := range results {
		if r.Key == "SKIP" {
			t.Error("SKIP key should have been ignored")
		}
	}
}
