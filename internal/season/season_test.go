package season

import (
	"testing"
	"time"
)

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

func TestFormatRowLengths(t *testing.T) {
	dates := []time.Time{
		time.Date(2026, 3, 25, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 6, 25, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 9, 25, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 12, 25, 0, 0, 0, 0, time.UTC),
	}
	for _, d := range dates {
		lines := Format(d)
		if lines[0] == "" || lines[1] == "" || lines[2] == "" {
			t.Errorf("%v: got empty row: %v", d.Format("Jan 2"), lines)
		}
	}
}

func TestPhaseLabelCoverage(t *testing.T) {
	seasons := []Season{Spring, Summer, Fall, Winter}
	progs := []float64{0.05, 0.25, 0.50, 0.75, 0.92}
	for _, s := range seasons {
		for _, p := range progs {
			label := phaseLabel(s, p)
			if label == "" {
				t.Errorf("season=%d progress=%.2f returned empty label", s, p)
			}
		}
	}
}
