package season

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/nicoleyson/vestaboard-note/internal/layout"
)

type Season int

const (
	Spring Season = iota
	Summer
	Fall
	Winter
)

func colorFor(s Season) int {
	switch s {
	case Spring:
		return 66
	case Summer:
		return 65
	case Fall:
		return 64
	default:
		return 69
	}
}

func nameFor(s Season) string {
	switch s {
	case Spring:
		return "SPRING"
	case Summer:
		return "SUMMER"
	case Fall:
		return "FALL"
	default:
		return "WINTER"
	}
}

func phaseLabel(s Season, progress float64) string {
	switch {
	case progress < 0.15:
		return "IS HERE"
	case progress < 0.40:
		switch s {
		case Spring:
			return "IN BLOOM"
		case Summer:
			return "HEATING UP"
		case Fall:
			return "TURNING"
		default:
			return "SETTING IN"
		}
	case progress < 0.65:
		return "IN FULL SWING"
	case progress < 0.85:
		return "WINDING DOWN"
	default:
		switch s {
		case Spring:
			return "ALMOST SUMMER"
		case Summer:
			return "ALMOST FALL"
		case Fall:
			return "ALMOST WINTER"
		default:
			return "ALMOST SPRING"
		}
	}
}

// solsticeEquinox returns the approximate Unix timestamp (days since J2000.0)
// for a given year's astronomical event using the Jean Meeus algorithm (simplified).
// event: 0=March equinox, 1=June solstice, 2=September equinox, 3=December solstice
func eventJDE(year int, event int) float64 {
	// Table 27.a from Meeus "Astronomical Algorithms" (approximate)
	y := float64(year-2000) / 1000.0
	var jde0 float64
	switch event {
	case 0: // March equinox
		jde0 = 2451623.80984 + 365242.37404*y + 0.05169*y*y - 0.00411*y*y*y - 0.00057*y*y*y*y
	case 1: // June solstice
		jde0 = 2451716.56767 + 365241.62603*y + 0.00325*y*y + 0.00888*y*y*y - 0.00030*y*y*y*y
	case 2: // September equinox
		jde0 = 2451810.21715 + 365242.01767*y - 0.11575*y*y + 0.00337*y*y*y + 0.00078*y*y*y*y
	case 3: // December solstice
		jde0 = 2451900.05952 + 365242.74049*y - 0.06223*y*y - 0.00823*y*y*y + 0.00032*y*y*y*y
	}
	return jde0
}

// jdeToTime converts a Julian Day Number to a time.Time (UTC).
func jdeToTime(jde float64) time.Time {
	// JD 2440587.5 = Unix epoch (1970-01-01 00:00:00 UTC)
	unixSec := (jde - 2440587.5) * 86400
	return time.Unix(int64(unixSec), 0).UTC()
}

func Current(t time.Time) (Season, float64) {
	year := t.Year()

	march := jdeToTime(eventJDE(year, 0))
	june := jdeToTime(eventJDE(year, 1))
	sept := jdeToTime(eventJDE(year, 2))
	dec := jdeToTime(eventJDE(year, 3))
	nextMarch := jdeToTime(eventJDE(year+1, 0))

	var s Season
	var start, end time.Time

	switch {
	case t.Before(march):
		prevDec := jdeToTime(eventJDE(year-1, 3))
		s = Winter
		start = prevDec
		end = march
	case t.Before(june):
		s = Spring
		start = march
		end = june
	case t.Before(sept):
		s = Summer
		start = june
		end = sept
	case t.Before(dec):
		s = Fall
		start = sept
		end = dec
	default:
		s = Winter
		start = dec
		end = nextMarch
	}

	elapsed := t.Sub(start).Hours()
	total := end.Sub(start).Hours()
	progress := elapsed / total
	progress = math.Max(0, math.Min(1, progress))

	return s, progress
}

func Format(t time.Time) [3]string {
	s, progress := Current(t)
	color := colorFor(s)
	name := nameFor(s)
	phase := phaseLabel(s, progress)

	tile := fmt.Sprintf("{%d}", color)
	inner := tile + name + tile
	innerTiles := 2 + len([]rune(name))
	total := layout.Cols - innerTiles
	left := total / 2
	right := total - left
	row2 := strings.Repeat(" ", left) + inner + strings.Repeat(" ", right)

	return [3]string{
		layout.ColorRow(color),
		row2,
		layout.Center(phase, layout.Cols),
	}
}
