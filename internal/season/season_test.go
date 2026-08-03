package season

import (
	"regexp"
	"testing"
	"time"
)

var placeholderRe = regexp.MustCompile(`\{[0-9]+\}`)

// tileCount counts rendered tiles: each {N} placeholder is 1 tile, each other rune is 1 tile.
func tileCount(s string) int {
	rest := placeholderRe.ReplaceAllString(s, "X")
	return len([]rune(rest))
}

func TestCurrentSpring(t *testing.T) {
	d := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	s, prog := Current(d)
	if s != Spring {
		t.Errorf("Apr 1 should be Spring, got %d", s)
	}
	if prog < 0 || prog > 1 {
		t.Errorf("progress out of range: %f", prog)
	}
}

func TestCurrentSummer(t *testing.T) {
	d := time.Date(2026, 7, 15, 0, 0, 0, 0, time.UTC)
	s, _ := Current(d)
	if s != Summer {
		t.Errorf("Jul 15 should be Summer, got %d", s)
	}
}

func TestCurrentFall(t *testing.T) {
	d := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	s, _ := Current(d)
	if s != Fall {
		t.Errorf("Oct 1 should be Fall, got %d", s)
	}
}

func TestCurrentWinter(t *testing.T) {
	d := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	s, _ := Current(d)
	if s != Winter {
		t.Errorf("Jan 15 should be Winter, got %d", s)
	}
}

func TestCurrentWinterDecember(t *testing.T) {
	d := time.Date(2026, 12, 25, 0, 0, 0, 0, time.UTC)
	s, _ := Current(d)
	if s != Winter {
		t.Errorf("Dec 25 should be Winter, got %d", s)
	}
}

func TestFormatTileCount(t *testing.T) {
	dates := []time.Time{
		time.Date(2026, 3, 25, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 6, 25, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 9, 25, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 12, 25, 0, 0, 0, 0, time.UTC),
	}
	for _, d := range dates {
		lines := Format(d)
		for i, l := range lines {
			n := tileCount(l)
			if n != 15 {
				t.Errorf("%v row %d: got %d tiles, want 15: %q", d.Format("Jan 2"), i, n, l)
			}
		}
	}
}

func TestProgressLabel_midSeason(t *testing.T) {
	label := progressLabel(Fall, 0.44)
	if label != "44% INTO FALL" {
		t.Errorf("got %q, want %q", label, "44% INTO FALL")
	}
}

func TestProgressLabel_lastDay(t *testing.T) {
	label := progressLabel(Summer, 0.999)
	if label != "LAST DAY" {
		t.Errorf("got %q, want %q", label, "LAST DAY")
	}
}

func TestProgressLabel_allSeasonsLessThan16(t *testing.T) {
	seasons := []Season{Spring, Summer, Fall, Winter}
	progs := []float64{0.01, 0.50, 0.98}
	for _, s := range seasons {
		for _, p := range progs {
			label := progressLabel(s, p)
			if len([]rune(label)) > 15 {
				t.Errorf("season=%d progress=%.2f label %q is %d runes, want ≤15", s, p, label, len([]rune(label)))
			}
			if label == "" {
				t.Errorf("season=%d progress=%.2f returned empty label", s, p)
			}
		}
	}
}
