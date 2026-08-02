package satellites

import (
	"testing"
)

func TestPickBest_empty(t *testing.T) {
	if pickBest(nil) != nil {
		t.Error("expected nil for empty slice")
	}
}

func TestPickBest_skipsUninteresting(t *testing.T) {
	sats := []satellite{
		{Name: "STARLINK-1", Category: "STARLINK", ElevationDeg: 80},
		{Name: "DEBRIS-X", Category: "DEBRIS", ElevationDeg: 60},
	}
	if pickBest(sats) != nil {
		t.Error("expected nil when only STARLINK/DEBRIS present")
	}
}

func TestPickBest_picksHighestElevation(t *testing.T) {
	sats := []satellite{
		{Name: "GPS-1", Category: "GPS", ElevationDeg: 40},
		{Name: "GPS-2", Category: "GPS", ElevationDeg: 75},
		{Name: "ISS", Category: "OTHER", ElevationDeg: 50},
	}
	got := pickBest(sats)
	if got == nil {
		t.Fatal("expected non-nil result")
	}
	if got.Name != "GPS-2" {
		t.Errorf("want GPS-2, got %q", got.Name)
	}
}

func TestPickBest_prefersInterestingOverHigherDebris(t *testing.T) {
	sats := []satellite{
		{Name: "STARLINK-999", Category: "STARLINK", ElevationDeg: 89},
		{Name: "ISS", Category: "OTHER", ElevationDeg: 45},
	}
	got := pickBest(sats)
	if got == nil || got.Name != "ISS" {
		t.Errorf("want ISS, got %v", got)
	}
}

func TestCleanName_stripsParenSuffix(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"ISS (ZARYA)", "ISS"},
		{"NOAA 19 (GOES)", "NOAA 19"},
		{"GPS IIF-3", "GPS IIF-3"},
		{"", ""},
	}
	for _, tc := range cases {
		got := cleanName(tc.input)
		if got != tc.want {
			t.Errorf("cleanName(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestCleanName_truncatesLong(t *testing.T) {
	long := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	got := cleanName(long)
	if len(got) > 15 {
		t.Errorf("cleanName did not truncate: len=%d, got=%q", len(got), got)
	}
}

func TestColorForCategory(t *testing.T) {
	cases := []struct {
		cat  string
		want int
	}{
		{"GPS", 66},
		{"IRIDIUM", 67},
		{"OTHER", 68},
		{"WEATHER", 68},
		{"", 68},
	}
	for _, tc := range cases {
		got := colorForCategory(tc.cat)
		if got != tc.want {
			t.Errorf("colorForCategory(%q) = %d, want %d", tc.cat, got, tc.want)
		}
	}
}
