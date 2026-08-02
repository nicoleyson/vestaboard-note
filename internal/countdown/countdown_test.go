package countdown

import (
	"testing"
	"time"
)

func checkRows(t *testing.T, lines [3]string) {
	t.Helper()
	for i, l := range lines {
		if len([]rune(l)) != 15 {
			t.Errorf("row %d wrong length %d: %q", i, len([]rune(l)), l)
		}
	}
}

func TestFormat_inDays(t *testing.T) {
	now := time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)
	events := []Event{
		{Label: "VACATION", Date: time.Date(2025, 7, 4, 0, 0, 0, 0, time.UTC)},
	}
	lines := Format(events, now)
	checkRows(t, lines)

	wantR1 := "   COUNTDOWN   "
	if lines[0] != wantR1 {
		t.Errorf("row1: got %q, want %q", lines[0], wantR1)
	}
	wantR2 := "   VACATION    "
	if lines[1] != wantR2 {
		t.Errorf("row2: got %q, want %q", lines[1], wantR2)
	}
	// Jun 1 12:00 → Jul 4 00:00 = 32.5 hours → ceil = 33 days
	wantR3 := "  IN 33 DAYS   "
	if lines[2] != wantR3 {
		t.Errorf("row3: got %q, want %q", lines[2], wantR3)
	}
}

func TestFormat_today(t *testing.T) {
	// Event is at midnight; now is 6 hours later same day → days = ceil(-6/24) = 0
	now := time.Date(2025, 7, 4, 6, 0, 0, 0, time.UTC)
	events := []Event{
		{Label: "VACATION", Date: time.Date(2025, 7, 4, 0, 0, 0, 0, time.UTC)},
	}
	lines := Format(events, now)
	checkRows(t, lines)

	wantR3 := "     TODAY     "
	if lines[2] != wantR3 {
		t.Errorf("row3: got %q, want %q", lines[2], wantR3)
	}
}

func TestFormat_tomorrow(t *testing.T) {
	now := time.Date(2025, 7, 3, 12, 0, 0, 0, time.UTC)
	events := []Event{
		{Label: "VACATION", Date: time.Date(2025, 7, 4, 12, 0, 0, 0, time.UTC)},
	}
	lines := Format(events, now)
	checkRows(t, lines)

	wantR3 := "   TOMORROW    "
	if lines[2] != wantR3 {
		t.Errorf("row3: got %q, want %q", lines[2], wantR3)
	}
}

func TestFormat_daysAgo(t *testing.T) {
	now := time.Date(2025, 7, 10, 12, 0, 0, 0, time.UTC)
	events := []Event{
		{Label: "BIRTHDAY", Date: time.Date(2025, 7, 4, 0, 0, 0, 0, time.UTC)},
	}
	lines := Format(events, now)
	checkRows(t, lines)

	// Jul 4 00:00 → Jul 10 12:00 = -156 hours → ceil(-156/24) = -6 → 6 DAYS AGO
	wantR3 := "  6 DAYS AGO   "
	if lines[2] != wantR3 {
		t.Errorf("row3: got %q, want %q", lines[2], wantR3)
	}
}

func TestFormat_empty(t *testing.T) {
	lines := Format(nil, time.Now())
	checkRows(t, lines)

	found := false
	for _, l := range lines {
		if l == "   NO EVENTS   " {
			found = true
		}
	}
	if !found {
		t.Errorf("empty events: expected NO EVENTS row, got %q %q %q", lines[0], lines[1], lines[2])
	}
}

func TestFormat_nearestPicked(t *testing.T) {
	now := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
	events := []Event{
		{Label: "FAR", Date: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)},
		{Label: "NEAR", Date: time.Date(2025, 7, 10, 0, 0, 0, 0, time.UTC)},
	}
	lines := Format(events, now)
	checkRows(t, lines)

	found := false
	for _, l := range lines {
		if len([]rune(l)) == 15 {
			runes := []rune(l)
			s := string(runes)
			if containsStr(s, "NEAR") {
				found = true
			}
		}
	}
	if !found {
		t.Errorf("nearest event not picked: got %q %q %q", lines[0], lines[1], lines[2])
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsRune(s, sub))
}

func containsRune(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
