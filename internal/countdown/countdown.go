package countdown

import (
	"fmt"
	"math"
	"time"

	"github.com/nicoleyson/vestaboard-note/internal/layout"
)

type Event struct {
	Label string    `yaml:"label"`
	Date  time.Time `yaml:"date"`
}

func Format(events []Event, now time.Time) [3]string {
	if len(events) == 0 {
		return [3]string{
			layout.Center("COUNTDOWN", layout.Cols),
			layout.Center("NO EVENTS", layout.Cols),
			layout.Center("SET IN CONFIG", layout.Cols),
		}
	}

	next := nearest(events, now)
	days := int(math.Ceil(next.Date.Sub(now).Hours() / 24))

	var row2, row3 string
	switch {
	case days < 0:
		row2 = layout.Center(next.Label, layout.Cols)
		row3 = layout.Center(fmt.Sprintf("%d DAYS AGO", -days), layout.Cols)
	case days == 0:
		row2 = layout.Center(next.Label, layout.Cols)
		row3 = layout.Center("TODAY", layout.Cols)
	case days == 1:
		row2 = layout.Center(next.Label, layout.Cols)
		row3 = layout.Center("TOMORROW", layout.Cols)
	default:
		row2 = layout.Center(next.Label, layout.Cols)
		row3 = layout.Center(fmt.Sprintf("IN %d DAYS", days), layout.Cols)
	}

	return [3]string{
		layout.Center("COUNTDOWN", layout.Cols),
		row2,
		row3,
	}
}

func nearest(events []Event, now time.Time) Event {
	best := events[0]
	bestDiff := math.Abs(events[0].Date.Sub(now).Hours())
	for _, e := range events[1:] {
		diff := math.Abs(e.Date.Sub(now).Hours())
		if diff < bestDiff {
			bestDiff = diff
			best = e
		}
	}
	return best
}
