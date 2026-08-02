package pattern

import (
	"regexp"
	"testing"

	"github.com/nicoleyson/vestaboard-note/internal/season"
)

var tileReC = regexp.MustCompile(`\{\d+\}|.`)

func countTilesC(s string) int {
	return len(tileReC.FindAllString(s, -1))
}

func assertCurrentLines(t *testing.T, name string, lines [3]string) {
	t.Helper()
	for i, line := range lines {
		n := countTilesC(line)
		if n != cols {
			t.Errorf("%s: row %d has %d tiles, want %d: %q", name, i, n, cols, line)
		}
	}
}

func TestStripesP_tileCounts(t *testing.T) {
	assertCurrentLines(t, "stripesP", stripesP([]int{cRed, cWhite}).toLines())
}

func TestCheckerP_tileCounts(t *testing.T) {
	assertCurrentLines(t, "checkerP", checkerP([]int{cBlue, cWhite}).toLines())
}

func TestBarsP_tileCounts(t *testing.T) {
	assertCurrentLines(t, "barsP", barsP([]int{cGreen, cWhite, cRed}).toLines())
}

func TestFadeP_tileCounts(t *testing.T) {
	assertCurrentLines(t, "fadeP", fadeP([]int{cYellow, cOrange, cRed}).toLines())
}

func TestDiagonalP_tileCounts(t *testing.T) {
	assertCurrentLines(t, "diagonalP", diagonalP([]int{cBlue, cViolet}).toLines())
}

func TestHeartsP_tileCounts(t *testing.T) {
	assertCurrentLines(t, "heartsP", heartsP([]int{cRed, cHeart, cWhite}).toLines())
}

func TestConfettiP_tileCounts(t *testing.T) {
	assertCurrentLines(t, "confettiP", confettiP([]int{cRed, cYellow, cGreen, cBlue, cViolet}).toLines())
}

func TestSparkleP_tileCounts(t *testing.T) {
	assertCurrentLines(t, "sparkleP", sparkleP([]int{cWhite, cViolet}).toLines())
}

func TestSparkleP_singleColor(t *testing.T) {
	assertCurrentLines(t, "sparkleP-single", sparkleP([]int{cWhite}).toLines())
}

func TestRainbowP_tileCounts(t *testing.T) {
	assertCurrentLines(t, "rainbowP", rainbowP([]int{cViolet, cYellow, cGreen}).toLines())
}

func TestPulseP_tileCounts(t *testing.T) {
	assertCurrentLines(t, "pulseP", pulseP([]int{cBlue, cWhite}).toLines())
}

func TestRenderWithTheme_allPatternNames(t *testing.T) {
	pal := []int{cRed, cWhite, cBlue}
	for _, name := range []string{
		"stripes", "checker", "bars", "fade", "diagonal",
		"hearts", "confetti", "sparkle", "rainbow", "pulse",
		"unknown_default",
	} {
		assertCurrentLines(t, "renderWithTheme/"+name, renderWithTheme(theme{palette: pal, patterns: []string{name}}))
	}
}

func TestThemeForProgress_bandCoverage(t *testing.T) {
	for _, s := range []season.Season{season.Spring, season.Summer, season.Fall, season.Winter} {
		for _, p := range []float64{0.1, 0.5, 0.9} {
			th := themeForProgress(s, p)
			if len(th.palette) == 0 || len(th.patterns) == 0 {
				t.Errorf("themeForProgress(season=%d, progress=%.1f): empty theme", s, p)
			}
		}
	}
}

func TestHolidayThemes_allProduceValidOutput(t *testing.T) {
	for _, ht := range holidayThemes {
		assertCurrentLines(t, "holidayTheme/"+ht.keyword, renderWithTheme(ht.t))
	}
}

func TestCurrent_rowLengths(t *testing.T) {
	lines, err := Current(37.7, -122.4)
	if err != nil {
		t.Fatalf("Current returned error: %v", err)
	}
	assertCurrentLines(t, "Current", lines)
}
