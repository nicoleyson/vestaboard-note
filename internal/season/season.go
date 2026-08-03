package season

import (
	"fmt"
	"math"
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

type seasonInfo struct {
	color int
	name  string
}

var seasonData = map[Season]seasonInfo{
	Spring: {color: 66, name: "SPRING"},
	Summer: {color: 65, name: "SUMMER"},
	Fall:   {color: 64, name: "FALL"},
	Winter: {color: 69, name: "WINTER"},
}

func progressLabel(s Season, progress float64) string {
	pct := int(math.Round(progress * 100))
	if pct > 99 {
		pct = 99
	}
	return fmt.Sprintf("%d%% INTO %s", pct, seasonData[s].name)
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
	info := seasonData[s]
	color := info.color

	if progress >= 0.988 {
		label := fmt.Sprintf("LAST DAY %s", info.name)
		return [3]string{
			layout.ColorRow(color),
			layout.Center(label, layout.Cols),
			layout.PadRight("", layout.Cols),
		}
	}

	return [3]string{
		layout.ColorRow(color),
		layout.Center(progressLabel(s, progress), layout.Cols),
		layout.ColorRow(color),
	}
}
