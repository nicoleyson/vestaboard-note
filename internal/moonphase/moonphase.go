package moonphase

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/nicoleyson/vestaboard-note/internal/layout"
)

const (
	colorLit  = 65
	colorDark = 70
)

type Phase struct {
	Name         string
	Illumination float64
	Waxing       bool
}

// Calculate uses the synodic period (29.53 days) anchored to the known
// new moon on 2000-01-06 18:14 UTC to determine the current phase.
// Illumination = (1 - cos(2π·cycle/period)) / 2, ranging 0 (new) → 1 (full).
func Calculate(t time.Time) Phase {
	ref := time.Date(2000, 1, 6, 18, 14, 0, 0, time.UTC)
	synodicDays := 29.53058867

	elapsed := t.UTC().Sub(ref).Hours() / 24
	cycle := math.Mod(elapsed, synodicDays)
	if cycle < 0 {
		cycle += synodicDays
	}

	angleRad := 2 * math.Pi * cycle / synodicDays
	illum := (1 - math.Cos(angleRad)) / 2

	return Phase{
		Name:         phaseName(cycle, synodicDays),
		Illumination: illum,
		Waxing:       cycle < synodicDays/2,
	}
}

func phaseName(cycle, period float64) string {
	frac := cycle / period
	switch {
	case frac < 0.02 || frac >= 0.98:
		return "NEW MOON"
	case frac < 0.23:
		return "WAX CRESCENT"
	case frac < 0.27:
		return "FIRST QUARTER"
	case frac < 0.48:
		return "WAX GIBBOUS"
	case frac < 0.52:
		return "FULL MOON"
	case frac < 0.73:
		return "WAN GIBBOUS"
	case frac < 0.77:
		return "LAST QUARTER"
	default:
		return "WAN CRESCENT"
	}
}

func Format(t time.Time) [3]string {
	return render(Calculate(t))
}

func tileColor(i, litCols int, waxing bool) int {
	if waxing {
		if i < litCols {
			return colorLit
		}
		return colorDark
	}
	darkCols := layout.Cols - litCols
	if i < darkCols {
		return colorDark
	}
	return colorLit
}

func render(p Phase) [3]string {
	litCols := int(math.Round(p.Illumination * float64(layout.Cols)))

	colorLine := func() string {
		var sb strings.Builder
		for i := 0; i < layout.Cols; i++ {
			sb.WriteString(fmt.Sprintf("{%d}", tileColor(i, litCols, p.Waxing)))
		}
		return sb.String()
	}

	row := colorLine()
	return [3]string{row, nameOverlay(p.Name, litCols, p.Waxing), row}
}

// nameOverlay builds the middle row. Blank positions emit a {color} tile
// (one tile each in encodeLines); named characters emit their rune directly
// (also one tile each). This keeps the total at exactly 15 encoded tiles.
func nameOverlay(name string, litCols int, waxing bool) string {
	centered := layout.Center(name, layout.Cols)
	var sb strings.Builder
	for i, r := range []rune(centered) {
		col := tileColor(i, litCols, waxing)
		if r == ' ' {
			sb.WriteString(fmt.Sprintf("{%d}", col))
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}
