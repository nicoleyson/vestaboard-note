package moonphase

import (
	"strings"
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

// countTiles counts logical display tiles in an encoded line.
// {N} sequences count as 1; each other rune counts as 1.
func countTiles(s string) int {
	count := 0
	for len(s) > 0 {
		if s[0] == '{' {
			end := strings.IndexByte(s, '}')
			if end > 0 {
				s = s[end+1:]
				count++
				continue
			}
		}
		// consume one rune
		for i, _ := range s {
			if i > 0 {
				s = s[i:]
				break
			}
			if i == 0 && len(s) == 1 {
				s = ""
			}
		}
		count++
	}
	return count
}
