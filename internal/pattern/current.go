package pattern

import (
	"math/rand"
	"strings"
	"time"

	"github.com/nicoleyson/vestaboard-note/internal/holiday"
	"github.com/nicoleyson/vestaboard-note/internal/season"
)

type theme struct {
	palette  []int
	patterns []string
}

var holidayThemes = []struct {
	keyword string
	t       theme
}{
	{"LUNAR NEW YEAR", theme{[]int{cRed, cYellow}, []string{"stripes", "checker", "diagonal", "bars"}}},
	{"CHINESE NEW YEAR", theme{[]int{cRed, cYellow}, []string{"stripes", "checker", "diagonal", "bars"}}},
	{"NEW YEAR", theme{[]int{cViolet, cYellow, cWhite}, []string{"confetti", "rainbow", "sparkle"}}},
	{"CHRISTMAS", theme{[]int{cRed, cGreen}, []string{"checker", "stripes", "diagonal", "bars"}}},
	{"HALLOWEEN", theme{[]int{cOrange, cBlack}, []string{"checker", "diagonal", "stripes"}}},
	{"VALENTINE", theme{[]int{cRed, cHeart}, []string{"hearts", "stripes", "checker"}}},
	{"PATRICK", theme{[]int{cGreen, cWhite}, []string{"stripes", "checker", "bars"}}},
	{"EASTER", theme{[]int{cRed, cGreen, cYellow}, []string{"checker", "confetti", "fade"}}},
	{"DIWALI", theme{[]int{cYellow, cOrange}, []string{"confetti", "sparkle", "rainbow"}}},
	{"DEEPAVALI", theme{[]int{cYellow, cOrange}, []string{"confetti", "sparkle", "rainbow"}}},
	{"HANUKKAH", theme{[]int{cBlue, cWhite}, []string{"stripes", "bars", "fade"}}},
	{"CHANUKAH", theme{[]int{cBlue, cWhite}, []string{"stripes", "bars", "fade"}}},
	{"EID", theme{[]int{cGreen, cWhite}, []string{"stripes", "checker", "bars"}}},
	{"INDEPENDENCE", theme{[]int{cRed, cWhite, cBlue}, []string{"stripes", "bars", "diagonal"}}},
	{"NATIONAL DAY", theme{[]int{cRed, cWhite, cBlue}, []string{"stripes", "bars", "diagonal"}}},
	{"THANKSGIVING", theme{[]int{cOrange, cYellow}, []string{"fade", "diagonal", "stripes"}}},
	{"MEMORIAL", theme{[]int{cRed, cWhite, cBlue}, []string{"stripes", "bars"}}},
	{"VETERANS", theme{[]int{cRed, cWhite, cBlue}, []string{"stripes", "bars"}}},
	{"BASTILLE", theme{[]int{cBlue, cWhite, cRed}, []string{"stripes", "bars", "diagonal"}}},
	{"ANZAC", theme{[]int{cRed, cWhite}, []string{"stripes", "bars", "checker"}}},
	{"REMEMBRANCE", theme{[]int{cRed, cWhite}, []string{"stripes", "bars", "checker"}}},
	{"GUY FAWKES", theme{[]int{cOrange, cBlack}, []string{"checker", "diagonal", "stripes"}}},
	{"MIDSUMMER", theme{[]int{cYellow, cGreen, cWhite}, []string{"stripes", "fade", "bars"}}},
	{"MIDSOMMAR", theme{[]int{cYellow, cGreen, cWhite}, []string{"stripes", "fade", "bars"}}},
	{"HOLI", theme{[]int{cRed, cYellow, cGreen, cBlue, cViolet}, []string{"confetti", "rainbow", "stripes"}}},
	{"NOWRUZ", theme{[]int{cGreen, cWhite}, []string{"stripes", "fade", "bars"}}},
	{"VESAK", theme{[]int{cOrange, cWhite}, []string{"stripes", "bars", "fade"}}},
	{"WESAK", theme{[]int{cOrange, cWhite}, []string{"stripes", "bars", "fade"}}},
	{"BUDDHA", theme{[]int{cOrange, cWhite}, []string{"stripes", "bars", "fade"}}},
	{"CARNIVAL", theme{[]int{cViolet, cYellow, cGreen}, []string{"confetti", "rainbow", "diagonal"}}},
	{"MARDI GRAS", theme{[]int{cViolet, cYellow, cGreen}, []string{"confetti", "rainbow", "diagonal"}}},
	{"JUNETEENTH", theme{[]int{cRed, cBlack, cGreen}, []string{"stripes", "bars", "checker"}}},
	{"KWANZAA", theme{[]int{cRed, cBlack, cGreen}, []string{"stripes", "bars", "checker"}}},
	{"SONGKRAN", theme{[]int{cBlue, cWhite}, []string{"stripes", "checker", "fade"}}},
	{"CHUSEOK", theme{[]int{cOrange, cYellow}, []string{"fade", "stripes", "bars"}}},
	{"OKTOBERFEST", theme{[]int{cBlue, cWhite}, []string{"checker", "stripes", "bars"}}},
	{"RAMADAN", theme{[]int{cGreen, cWhite, cYellow}, []string{"stripes", "fade", "bars"}}},
	{"DAY OF THE DEAD", theme{[]int{cOrange, cViolet, cYellow}, []string{"confetti", "sparkle", "checker"}}},
	{"DIA DE", theme{[]int{cOrange, cViolet, cYellow}, []string{"confetti", "sparkle", "checker"}}},
	{"CHILDREN", theme{[]int{cRed, cYellow, cGreen, cBlue}, []string{"confetti", "rainbow", "stripes"}}},
	{"ONAM", theme{[]int{cYellow, cGreen, cWhite}, []string{"stripes", "fade", "bars"}}},
	{"AUSTRALIA", theme{[]int{cGreen, cYellow}, []string{"stripes", "checker", "bars"}}},
	{"CANADA", theme{[]int{cRed, cWhite}, []string{"stripes", "bars", "checker"}}},
	{"LIBERATION", theme{[]int{cRed, cWhite}, []string{"stripes", "bars"}}},
	{"FREEDOM", theme{[]int{cRed, cWhite, cBlue}, []string{"stripes", "bars", "diagonal"}}},
	{"CONSTITUTION", theme{[]int{cBlue, cWhite, cRed}, []string{"stripes", "bars", "diagonal"}}},
	{"REPUBLIC", theme{[]int{cOrange, cWhite, cGreen}, []string{"stripes", "bars", "diagonal"}}},
	{"BONFIRE", theme{[]int{cOrange, cRed, cBlack}, []string{"checker", "diagonal", "stripes"}}},
	{"ARMISTICE", theme{[]int{cRed, cWhite}, []string{"stripes", "bars", "checker"}}},
	{"MARTIN LUTHER", theme{[]int{cBlue, cWhite}, []string{"stripes", "bars", "checker"}}},
	{"PRESIDENTS", theme{[]int{cRed, cWhite, cBlue}, []string{"stripes", "bars", "diagonal"}}},
	{"LABOR", theme{[]int{cBlue, cWhite}, []string{"stripes", "bars", "checker"}}},
	{"WORKERS", theme{[]int{cBlue, cWhite}, []string{"stripes", "bars", "checker"}}},
	{"NOCHE BUENA", theme{[]int{cRed, cGreen}, []string{"checker", "stripes", "diagonal"}}},
	{"NOCHEBUENA", theme{[]int{cRed, cGreen}, []string{"checker", "stripes", "diagonal"}}},
	{"MOTHER", theme{[]int{cRed, cHeart, cWhite}, []string{"hearts", "stripes", "fade"}}},
	{"FATHER", theme{[]int{cBlue, cWhite}, []string{"stripes", "bars", "checker"}}},
}

var seasonThemes = map[season.Season][]theme{
	season.Spring: {
		{[]int{cGreen, cWhite, cYellow}, []string{"hearts", "fade", "stripes"}},
		{[]int{cGreen, cYellow, cHeart}, []string{"hearts", "confetti", "checker"}},
		{[]int{cGreen, cWhite}, []string{"stripes", "bars", "fade"}},
	},
	season.Summer: {
		{[]int{cYellow, cOrange}, []string{"fade", "rainbow", "bars"}},
		{[]int{cYellow, cOrange, cRed}, []string{"stripes", "diagonal", "fade"}},
		{[]int{cBlue, cYellow}, []string{"checker", "stripes", "diagonal"}},
	},
	season.Fall: {
		{[]int{cOrange, cYellow}, []string{"fade", "diagonal", "stripes"}},
		{[]int{cOrange, cRed, cYellow}, []string{"diagonal", "checker", "bars"}},
		{[]int{cOrange, cBlack}, []string{"diagonal", "checker", "stripes"}},
	},
	season.Winter: {
		{[]int{cWhite, cViolet}, []string{"sparkle", "bars", "stripes"}},
		{[]int{cBlue, cWhite}, []string{"stripes", "checker", "fade"}},
		{[]int{cWhite, cBlue, cViolet}, []string{"sparkle", "fade", "diagonal"}},
	},
}

func themeForProgress(s season.Season, progress float64) theme {
	themes := seasonThemes[s]
	var idx int
	switch {
	case progress < 0.33:
		idx = 0
	case progress < 0.67:
		idx = 1
	default:
		idx = 2
	}
	return themes[idx]
}

func renderWithTheme(t theme) [3]string {
	patName := t.patterns[rand.Intn(len(t.patterns))]
	pal := t.palette
	var g grid
	switch patName {
	case "stripes":
		g = stripesP(pal)
	case "checker":
		g = checkerP(pal)
	case "bars":
		g = barsP(pal)
	case "fade":
		g = fadeP(pal)
	case "diagonal":
		g = diagonalP(pal)
	case "hearts":
		g = heartsP(pal)
	case "confetti":
		g = confettiP(pal)
	case "sparkle":
		g = sparkleP(pal)
	case "rainbow":
		g = rainbowP(pal)
	case "pulse":
		g = pulseP(pal)
	default:
		g = stripesP(pal)
	}
	return g.toLines()
}

func Current(lat, lon float64) ([3]string, error) {
	name := holiday.TodaysName(lat, lon)
	if name != "" {
		for _, ht := range holidayThemes {
			if strings.Contains(name, ht.keyword) {
				return renderWithTheme(ht.t), nil
			}
		}
	}

	s, progress := season.Current(time.Now())
	return renderWithTheme(themeForProgress(s, progress)), nil
}

func stripesP(pal []int) grid {
	var g grid
	width := 1 + rand.Intn(3)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			g[r][c] = pal[(c/width)%len(pal)]
		}
	}
	return g
}

func checkerP(pal []int) grid {
	var g grid
	size := 1 + rand.Intn(2)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			g[r][c] = pal[((r/size)+(c/size))%len(pal)]
		}
	}
	return g
}

func barsP(pal []int) grid {
	var g grid
	for r := 0; r < rows; r++ {
		c := pal[r%len(pal)]
		for col := 0; col < cols; col++ {
			g[r][col] = c
		}
	}
	return g
}

func fadeP(pal []int) grid {
	var g grid
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			idx := c * len(pal) / cols
			g[r][c] = pal[idx]
		}
	}
	return g
}

func diagonalP(pal []int) grid {
	var g grid
	width := 2 + rand.Intn(3)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			g[r][c] = pal[((r+c)/width)%len(pal)]
		}
	}
	return g
}

func heartsP(pal []int) grid {
	var g grid
	bg := pal[0]
	density := 3 + rand.Intn(5)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			g[r][c] = bg
		}
	}
	placed, attempts := 0, 0
	for placed < density && attempts < 100 {
		r := rand.Intn(rows)
		c := rand.Intn(cols)
		if g[r][c] == bg {
			g[r][c] = cHeart
			placed++
		}
		attempts++
	}
	return g
}

func confettiP(pal []int) grid {
	var g grid
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			g[r][c] = pal[rand.Intn(len(pal))]
		}
	}
	return g
}

func sparkleP(pal []int) grid {
	if len(pal) < 2 {
		return stripesP(pal)
	}
	var g grid
	bg := pal[len(pal)-1]
	density := 4 + rand.Intn(6)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			g[r][c] = bg
		}
	}
	placed, attempts := 0, 0
	for placed < density && attempts < 100 {
		r := rand.Intn(rows)
		c := rand.Intn(cols)
		if g[r][c] == bg {
			g[r][c] = pal[rand.Intn(len(pal)-1)]
			placed++
		}
		attempts++
	}
	return g
}

// rainbowP cycles through the palette left-to-right in stripes of width 3.
func rainbowP(pal []int) grid {
	var g grid
	offset := rand.Intn(len(pal))
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			g[r][c] = pal[(c/3+offset)%len(pal)]
		}
	}
	return g
}

func pulseP(pal []int) grid {
	var g grid
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			dist := abs(c - cols/2)
			ring := dist / 3
			g[r][c] = pal[ring%len(pal)]
		}
	}
	return g
}
