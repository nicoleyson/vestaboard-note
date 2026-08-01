package clock

import (
	"testing"
	"time"

	"github.com/nicoleyson/vestaboard-note/internal/layout"
)

func TestFormatLength(t *testing.T) {
	fixed := time.Date(2024, time.August, 1, 15, 30, 0, 0, time.UTC)
	lines := Format(fixed)
	for i, line := range lines {
		if len([]rune(line)) != layout.Cols {
			t.Errorf("Format row %d length = %d, want %d: %q", i, len([]rune(line)), layout.Cols, line)
		}
	}
}

func TestFormatContent(t *testing.T) {
	fixed := time.Date(2024, time.August, 1, 15, 4, 0, 0, time.UTC)
	lines := Format(fixed)

	if lines[0] != "   THU AUG 1   " {
		t.Errorf("Format row 0 = %q, want %q", lines[0], "   THU AUG 1   ")
	}
	if lines[1] != "    3:04 PM    " {
		t.Errorf("Format row 1 = %q, want %q", lines[1], "    3:04 PM    ")
	}
	if lines[2] != "     2024      " {
		t.Errorf("Format row 2 = %q, want %q", lines[2], "     2024      ")
	}
}

func TestFormatMidnight(t *testing.T) {
	midnight := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
	lines := Format(midnight)
	for i, line := range lines {
		if len([]rune(line)) != layout.Cols {
			t.Errorf("Format midnight row %d length = %d, want %d", i, len([]rune(line)), layout.Cols)
		}
	}
}
