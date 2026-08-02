package tearoff

import (
	"regexp"
	"strings"
	"testing"
	"time"
)

var tileRe = regexp.MustCompile(`\{(\d+)\}|[A-Z0-9]`)

func tileCount(s string) int {
	return len(tileRe.FindAllString(s, -1))
}

func TestFormatRowLengths(t *testing.T) {
	cases := []struct {
		day int
	}{
		{1}, {9}, {10}, {15}, {28}, {31},
	}
	for _, tc := range cases {
		date := time.Date(2025, 8, tc.day, 0, 0, 0, 0, time.UTC)
		lines := Format(date)
		for r, line := range lines {
			got := tileCount(line)
			if got != 15 {
				t.Errorf("day %d row %d: want 15 tiles, got %d (line=%q)", tc.day, r, got, line)
			}
		}
	}
}

func TestRow1Pattern(t *testing.T) {
	date := time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)
	lines := Format(date)
	want := "{63}{63}{63}{65}{65}{63}{63}{63}{63}{63}{65}{65}{63}{63}{63}"
	if lines[0] != want {
		t.Errorf("row1: want %q, got %q", want, lines[0])
	}
}

func TestRow3AllWhite(t *testing.T) {
	date := time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)
	lines := Format(date)
	want := strings.Repeat("{69}", 15)
	if lines[2] != want {
		t.Errorf("row3: want all white, got %q", lines[2])
	}
}

func TestRow2ContainsDayNumber(t *testing.T) {
	cases := []struct {
		day    int
		dayStr string
	}{
		{1, "1"},
		{15, "15"},
		{31, "31"},
	}
	for _, tc := range cases {
		date := time.Date(2025, 8, tc.day, 0, 0, 0, 0, time.UTC)
		lines := Format(date)
		if !strings.Contains(lines[1], tc.dayStr) {
			t.Errorf("day %d: row2 %q does not contain %q", tc.day, lines[1], tc.dayStr)
		}
	}
}

func TestRow2Centered(t *testing.T) {
	date := time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)
	lines := Format(date)
	left := strings.Index(lines[1], "1")
	right := strings.LastIndex(lines[1], "1")
	if left != right {
		t.Skip("multiple 1s in row2, centering test skipped")
	}
	leftTiles := tileCount(lines[1][:left])
	rightTilesStart := left + 1
	rightTiles := tileCount(lines[1][rightTilesStart:])
	diff := leftTiles - rightTiles
	if diff < -1 || diff > 1 {
		t.Errorf("day 1 row2 not centered: %d left tiles, %d right tiles in %q", leftTiles, rightTiles, lines[1])
	}
}
