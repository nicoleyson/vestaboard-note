package weekprogress

import (
	"regexp"
	"strings"
	"testing"
	"time"
)

func makeTime(weekday time.Weekday, hour, minute int) time.Time {
	// Start from a known Monday (2024-01-01 = Monday) and advance to the desired weekday.
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	for base.Weekday() != time.Monday {
		base = base.Add(24 * time.Hour)
	}
	offset := int(weekday) - int(time.Monday)
	if offset < 0 {
		offset += 7
	}
	return base.AddDate(0, 0, offset).Add(time.Duration(hour)*time.Hour + time.Duration(minute)*time.Minute)
}

func countTiles(s string) int {
	re := regexp.MustCompile(`\{\d+\}|.`)
	return len(re.FindAllString(s, -1))
}

func TestFormat_rowLengths(t *testing.T) {
	now := makeTime(time.Wednesday, 14, 30)
	lines := Format(now)
	for i, line := range lines {
		if n := countTiles(line); n != 15 {
			t.Errorf("row %d: got %d tiles, want 15 (content: %q)", i, n, line)
		}
	}
}

func TestFormat_row2ContainsDayName(t *testing.T) {
	cases := []time.Weekday{
		time.Monday, time.Tuesday, time.Wednesday, time.Thursday,
		time.Friday, time.Saturday, time.Sunday,
	}
	for _, wd := range cases {
		t.Run(wd.String(), func(t *testing.T) {
			lines := Format(makeTime(wd, 12, 0))
			want := strings.ToUpper(wd.String())
			if !strings.Contains(lines[1], want) {
				t.Errorf("row 2 %q does not contain %q", lines[1], want)
			}
		})
	}
}

func TestFormat_row1IsFullGreen(t *testing.T) {
	lines := Format(makeTime(time.Wednesday, 9, 0))
	re := regexp.MustCompile(`\{66\}`)
	matches := re.FindAllString(lines[0], -1)
	if len(matches) != 15 {
		t.Errorf("row 1: got %d green tiles, want 15", len(matches))
	}
}

func TestProgress_mondayMidnight(t *testing.T) {
	p := progress(makeTime(time.Monday, 0, 0))
	if p != 0.0 {
		t.Errorf("Monday midnight: got %v, want 0.0", p)
	}
}

func TestProgress_sundayMidnight(t *testing.T) {
	p := progress(makeTime(time.Sunday, 0, 0))
	want := 6.0 / 7.0
	if p != want {
		t.Errorf("Sunday midnight: got %v, want %v", p, want)
	}
}

func TestFormat_progressBarAdvances(t *testing.T) {
	mondayMorning := makeTime(time.Monday, 6, 0)
	fridayEvening := makeTime(time.Friday, 20, 0)

	linesEarly := Format(mondayMorning)
	linesLate := Format(fridayEvening)

	reGreen := regexp.MustCompile(`\{66\}`)
	earlyGreen := len(reGreen.FindAllString(linesEarly[2], -1))
	lateGreen := len(reGreen.FindAllString(linesLate[2], -1))

	if lateGreen <= earlyGreen {
		t.Errorf("Friday evening (%d green) should have more green tiles than Monday morning (%d green)", lateGreen, earlyGreen)
	}
}

func TestFormat_progressBarTilesSumTo15(t *testing.T) {
	now := makeTime(time.Thursday, 15, 0)
	lines := Format(now)
	reGreen := regexp.MustCompile(`\{66\}`)
	reWhite := regexp.MustCompile(`\{69\}`)
	green := len(reGreen.FindAllString(lines[2], -1))
	white := len(reWhite.FindAllString(lines[2], -1))
	if green+white != 15 {
		t.Errorf("progress bar: green(%d) + white(%d) = %d, want 15", green, white, green+white)
	}
}

func TestProgressBar_allGreenAtEnd(t *testing.T) {
	bar := progressBar(15)
	reGreen := regexp.MustCompile(`\{66\}`)
	if n := len(reGreen.FindAllString(bar, -1)); n != 15 {
		t.Errorf("progressBar(15): got %d green tiles, want 15", n)
	}
}

func TestProgressBar_allWhiteAtStart(t *testing.T) {
	bar := progressBar(0)
	reWhite := regexp.MustCompile(`\{69\}`)
	if n := len(reWhite.FindAllString(bar, -1)); n != 15 {
		t.Errorf("progressBar(0): got %d white tiles, want 15", n)
	}
}
