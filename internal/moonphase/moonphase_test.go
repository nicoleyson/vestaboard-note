package moonphase

import (
	"regexp"
	"testing"
	"time"
)

func TestFormat_tileCount(t *testing.T) {
	// Use a known date: 2024-01-17 (known waxing gibbous, ~77% lit)
	ts := time.Date(2024, 1, 17, 12, 0, 0, 0, time.UTC)
	lines := Format(ts)
	for i, line := range lines {
		// Count tiles: each {N} = 1 tile, each plain rune = 1 tile
		count := countTiles(line)
		if count != 15 {
			t.Errorf("row %d: got %d tiles, want 15 (line: %q)", i+1, count, line)
		}
	}
}

func TestFormat_newMoon(t *testing.T) {
	// Reference new moon: 2000-01-06 18:14 UTC
	ts := time.Date(2000, 1, 6, 18, 14, 0, 0, time.UTC)
	p := Calculate(ts)
	if p.Name != "NEW MOON" {
		t.Errorf("got phase %q, want NEW MOON", p.Name)
	}
	if p.Illumination > 0.01 {
		t.Errorf("got illumination %.4f, want ~0", p.Illumination)
	}
}

func TestFormat_fullMoon(t *testing.T) {
	// ~14.76 days after reference new moon = full moon
	ts := time.Date(2000, 1, 21, 0, 0, 0, 0, time.UTC)
	p := Calculate(ts)
	if p.Name != "FULL MOON" {
		t.Errorf("got phase %q, want FULL MOON", p.Name)
	}
	if p.Illumination < 0.95 {
		t.Errorf("got illumination %.4f, want ~1.0", p.Illumination)
	}
}

func TestPhaseName(t *testing.T) {
	period := 29.53058867
	cases := []struct {
		frac float64
		want string
	}{
		{0.00, "NEW MOON"},
		{0.01, "NEW MOON"},
		{0.015, "NEW MOON"},
		{0.02, "WAX CRESCENT"},
		{0.10, "WAX CRESCENT"},
		{0.22, "WAX CRESCENT"},
		{0.23, "FIRST QUARTER"},
		{0.25, "FIRST QUARTER"},
		{0.269, "FIRST QUARTER"},
		{0.27, "WAX GIBBOUS"},
		{0.40, "WAX GIBBOUS"},
		{0.479, "WAX GIBBOUS"},
		{0.48, "FULL MOON"},
		{0.50, "FULL MOON"},
		{0.519, "FULL MOON"},
		{0.52, "WAN GIBBOUS"},
		{0.60, "WAN GIBBOUS"},
		{0.729, "WAN GIBBOUS"},
		{0.73, "LAST QUARTER"},
		{0.75, "LAST QUARTER"},
		{0.769, "LAST QUARTER"},
		{0.77, "WAN CRESCENT"},
		{0.90, "WAN CRESCENT"},
		{0.979, "WAN CRESCENT"},
		{0.98, "NEW MOON"},
		{0.999, "NEW MOON"},
	}
	for _, tc := range cases {
		cycle := tc.frac * period
		got := phaseName(cycle, period)
		if got != tc.want {
			t.Errorf("frac=%.3f: got %q, want %q", tc.frac, got, tc.want)
		}
	}
}

// countTiles counts logical display tiles in an encoded line.
// {N} sequences count as 1; each other rune counts as 1.
var tileRe = regexp.MustCompile(`\{\d+\}|.`)

func countTiles(s string) int {
	return len(tileRe.FindAllString(s, -1))
}
