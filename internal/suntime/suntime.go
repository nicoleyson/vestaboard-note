package suntime

import (
	"fmt"
	"time"

	"github.com/nicoleyson/vestaboard-note/internal/layout"
	"github.com/nicoleyson/vestaboard-note/internal/sunapi"
)

func Fetch(lat, lon float64) ([3]string, error) {
	now := time.Now()

	rise, set, err := sunapi.FetchTimes(lat, lon, now)
	if err != nil {
		return [3]string{}, err
	}

	if now.After(set) {
		tomorrow := now.AddDate(0, 0, 1)
		rise, set, err = sunapi.FetchTimes(lat, lon, tomorrow)
		if err != nil {
			return [3]string{}, err
		}
	}

	var label, timeStr string
	var color int
	var eventTime time.Time

	if now.Before(rise) {
		label = "SUNRISE"
		color = 65
		eventTime = rise
	} else {
		label = "SUNSET"
		color = 64
		eventTime = set
	}

	timeStr = eventTime.Local().Format("3:04PM")
	dateStr := eventTime.Local().Format("Mon Jan 2")

	row2 := layout.Center(fmt.Sprintf("%s %s", label, timeStr), layout.Cols)
	row3 := layout.Center(dateStr, layout.Cols)

	return [3]string{
		layout.ColorRow(color),
		row2,
		row3,
	}, nil
}
