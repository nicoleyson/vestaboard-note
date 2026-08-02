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

	row2 := layout.Center(next.Label, layout.Cols)
	var row3 string
	switch {
	case days < 0:
		row3 = layout.Center(fmt.Sprintf("%d DAYS AGO", -days), layout.Cols)
	case days == 0:
		row3 = layout.Center("TODAY", layout.Cols)
	case days == 1:
		row3 = layout.Center("TOMORROW", layout.Cols)
	default:
		row3 = layout.Center(fmt.Sprintf("IN %d DAYS", days), layout.Cols)
	}

	return [3]string{
		layout.Center("COUNTDOWN", layout.Cols),
		row2,
		row3,
	}
}

func nearest(events []Event, now time.Time) Event {
	var bestFuture *Event
	bestFutureDiff := math.MaxFloat64
	var bestPast *Event
	bestPastDiff := math.MaxFloat64

	for i := range events {
		e := &events[i]
		diff := e.Date.Sub(now).Hours()
		if diff >= 0 {
			if diff < bestFutureDiff {
				bestFutureDiff = diff
				bestFuture = e
			}
		} else {
			absDiff := -diff
			if absDiff < bestPastDiff {
				bestPastDiff = absDiff
				bestPast = e
			}
		}
	}

	if bestFuture != nil {
		return *bestFuture
	}
	return *bestPast
}
