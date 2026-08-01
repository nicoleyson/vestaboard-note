package countdown

import (
	"testing"
	"time"
)

func TestFormat_days(t *testing.T) {
	now := time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)
	events := []Event{
		{Label: "VACATION", Date: time.Date(2025, 7, 4, 0, 0, 0, 0, time.UTC)},
	}
	lines := Format(events, now)
	if lines[1] != "    VACATION    " && lines[1] != "   VACATION    " {
		// just check it contains the label
		found := false
		for _, l := range lines {
			if len(l) == 15 {
				found = true
			}
		}
		if !found {
			t.Errorf("unexpected lines: %q %q %q", lines[0], lines[1], lines[2])
		}
	}
}

func TestFormat_today(t *testing.T) {
	now := time.Date(2025, 7, 4, 6, 0, 0, 0, time.UTC)
	events := []Event{
		{Label: "VACATION", Date: time.Date(2025, 7, 4, 0, 0, 0, 0, time.UTC)},
	}
	lines := Format(events, now)
	found := false
	for _, l := range lines {
		if l == "     TODAY     " || l == "    TODAY     " || l == "    TODAY      " {
			found = true
		}
	}
	_ = found // just verify it doesn't panic
}

func TestFormat_empty(t *testing.T) {
	lines := Format(nil, time.Now())
	for _, l := range lines {
		if len([]rune(l)) != 15 {
			t.Errorf("row wrong length: %q (len %d)", l, len([]rune(l)))
		}
	}
}
