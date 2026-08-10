package weekprogress

import (
	"fmt"
	"strings"
	"time"

	"github.com/nicoleyson/vestaboard-note/internal/layout"
)

const (
	cGreen = 66
	cWhite = 69
	cols   = layout.Cols
)

var weekdays = map[time.Weekday]int{
	time.Monday:    0,
	time.Tuesday:   1,
	time.Wednesday: 2,
	time.Thursday:  3,
	time.Friday:    4,
	time.Saturday:  5,
	time.Sunday:    6,
}

func progress(t time.Time) float64 {
	dayIndex := float64(weekdays[t.Weekday()])
	hourFraction := float64(t.Hour())/24.0 + float64(t.Minute())/(24.0*60.0)
	return (dayIndex + hourFraction) / 7.0
}

func progressBar(elapsedTiles int) string {
	var b strings.Builder
	for i := 0; i < cols; i++ {
		if i < elapsedTiles {
			fmt.Fprintf(&b, "{%d}", cGreen)
		} else {
			fmt.Fprintf(&b, "{%d}", cWhite)
		}
	}
	return b.String()
}

func Format(t time.Time) [3]string {
	p := progress(t)
	elapsedTiles := int(p*float64(cols) + 0.5)
	if elapsedTiles < 0 {
		elapsedTiles = 0
	}
	if elapsedTiles > cols {
		elapsedTiles = cols
	}

	dayName := strings.ToUpper(t.Weekday().String())

	return [3]string{
		layout.ColorRow(cGreen),
		layout.Center(dayName, cols),
		progressBar(elapsedTiles),
	}
}
