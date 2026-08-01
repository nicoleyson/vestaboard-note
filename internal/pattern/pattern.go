package pattern

import (
	"fmt"
	"math/rand"
	"strings"
)

const (
	rows = 3
	cols = 15

	cRed    = 63
	cOrange = 64
	cYellow = 65
	cGreen  = 66
	cBlue   = 67
	cViolet = 68
	cWhite  = 69
	cBlack  = 70
	cHeart  = 62
)

var allColors = []int{cRed, cOrange, cYellow, cGreen, cBlue, cViolet, cWhite, cBlack}
var brightColors = []int{cRed, cOrange, cYellow, cGreen, cBlue, cViolet}
var warmColors = []int{cRed, cOrange, cYellow}
var coolColors = []int{cGreen, cBlue, cViolet}

var gradientPalettes = [][]int{
	{cYellow, cOrange, cRed},
	{cBlue, cViolet, cRed},
	{cGreen, cBlue, cViolet},
	{cYellow, cGreen, cBlue},
	{cRed, cViolet, cBlue},
	{cOrange, cYellow, cGreen},
}

var Names = []string{
	"stripes", "checker", "bars", "fade", "diagonal",
	"hearts", "confetti", "sparkle", "pulse", "rainbow",
}

type grid [rows][cols]int

func (g grid) toLines() [3]string {
	var lines [3]string
	for r := 0; r < rows; r++ {
		var b strings.Builder
		for c := 0; c < cols; c++ {
			code := g[r][c]
			if code >= 1 && code <= 26 {
				b.WriteRune(rune('A' + code - 1))
			} else if code >= 27 && code <= 35 {
				b.WriteRune(rune('1' + code - 27))
			} else if code == 36 {
				b.WriteRune('0')
			} else {
				b.WriteString(fmt.Sprintf("{%d}", code))
			}
		}
		lines[r] = b.String()
	}
	return lines
}

func pick(palette []int) int {
	return palette[rand.Intn(len(palette))]
}

func twoDistinct(palette []int) (int, int) {
	a := pick(palette)
	b := pick(palette)
	for b == a && len(palette) > 1 {
		b = pick(palette)
	}
	return a, b
}

func stripes() grid {
	var g grid
	numColors := 2 + rand.Intn(3)
	palette := make([]int, numColors)
	used := pick(brightColors)
	palette[0] = used
	for i := 1; i < numColors; i++ {
		c := pick(brightColors)
		for c == palette[i-1] {
			c = pick(brightColors)
		}
		palette[i] = c
	}
	width := 1 + rand.Intn(3)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			g[r][c] = palette[(c/width)%numColors]
		}
	}
	return g
}

func checker() grid {
	var g grid
	a, b := twoDistinct(brightColors)
	size := 1 + rand.Intn(2)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if ((r/size)+(c/size))%2 == 0 {
				g[r][c] = a
			} else {
				g[r][c] = b
			}
		}
	}
	return g
}

func bars() grid {
	var g grid
	colors := make([]int, rows)
	colors[0] = pick(brightColors)
	for i := 1; i < rows; i++ {
		c := pick(brightColors)
		for c == colors[i-1] {
			c = pick(brightColors)
		}
		colors[i] = c
	}
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			g[r][c] = colors[r]
		}
	}
	return g
}

func fade() grid {
	var g grid
	palette := gradientPalettes[rand.Intn(len(gradientPalettes))]
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			idx := c * len(palette) / cols
			g[r][c] = palette[idx]
		}
	}
	return g
}

func diagonal() grid {
	var g grid
	a, b := twoDistinct(brightColors)
	width := 2 + rand.Intn(3)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if ((r+c)/width)%2 == 0 {
				g[r][c] = a
			} else {
				g[r][c] = b
			}
		}
	}
	return g
}

func hearts() grid {
	var g grid
	bg := pick(brightColors)
	density := 3 + rand.Intn(5)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			g[r][c] = bg
		}
	}
	placed := 0
	attempts := 0
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

func confetti() grid {
	var g grid
	letters := []int{}
	for i := 1; i <= 26; i++ {
		letters = append(letters, i)
	}
	all := append(append([]int{}, brightColors...), letters...)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			g[r][c] = all[rand.Intn(len(all))]
		}
	}
	return g
}

func sparkle() grid {
	var g grid
	bg := cBlack
	if rand.Intn(2) == 0 {
		bg = cViolet
	}
	density := 4 + rand.Intn(6)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			g[r][c] = bg
		}
	}
	placed := 0
	attempts := 0
	sparklePalette := []int{cYellow, cWhite, cOrange, cHeart}
	for placed < density && attempts < 100 {
		r := rand.Intn(rows)
		c := rand.Intn(cols)
		if g[r][c] == bg {
			g[r][c] = sparklePalette[rand.Intn(len(sparklePalette))]
			placed++
		}
		attempts++
	}
	return g
}

func pulse() grid {
	var g grid
	center := pick(brightColors)
	mid := pick(brightColors)
	for mid == center {
		mid = pick(brightColors)
	}
	outer := pick(brightColors)
	for outer == center || outer == mid {
		outer = pick(brightColors)
	}
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			distFromCenter := abs(c - cols/2)
			switch {
			case distFromCenter <= 2:
				g[r][c] = center
			case distFromCenter <= 5:
				g[r][c] = mid
			default:
				g[r][c] = outer
			}
		}
	}
	return g
}

func rainbow() grid {
	var g grid
	spectrum := []int{cRed, cOrange, cYellow, cGreen, cBlue, cViolet}
	offset := rand.Intn(len(spectrum))
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			g[r][c] = spectrum[(c/3+offset)%len(spectrum)]
		}
	}
	return g
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Generate(name string) ([3]string, error) {
	var g grid
	switch name {
	case "stripes":
		g = stripes()
	case "checker":
		g = checker()
	case "bars":
		g = bars()
	case "fade":
		g = fade()
	case "diagonal":
		g = diagonal()
	case "hearts":
		g = hearts()
	case "confetti":
		g = confetti()
	case "sparkle":
		g = sparkle()
	case "pulse":
		g = pulse()
	case "rainbow":
		g = rainbow()
	default:
		return [3]string{}, fmt.Errorf("unknown pattern %q — available: %s", name, strings.Join(Names, ", "))
	}
	return g.toLines(), nil
}

func Random() [3]string {
	lines, _ := Generate(Names[rand.Intn(len(Names))])
	return lines
}
