package calendar

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/bounoable/ical"
	"github.com/nicoleyson/vestaboard-note/internal/layout"
)

func Fetch(icsURL string) ([3]string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(icsURL)
	if err != nil {
		return [3]string{}, err
	}
	defer resp.Body.Close()

	cal, err := ical.Parse(resp.Body)
	if err != nil {
		return [3]string{}, err
	}

	now := time.Now()
	cutoff := now.Add(7 * 24 * time.Hour)

	type event struct {
		start   time.Time
		summary string
	}
	var upcoming []event
	for _, e := range cal.Events {
		if e.Start.After(now) && e.Start.Before(cutoff) {
			upcoming = append(upcoming, event{start: e.Start, summary: e.Summary})
		}
	}

	if len(upcoming) == 0 {
		return [3]string{
			layout.Center("NO EVENTS", layout.Cols),
			layout.Center("NEXT 7 DAYS", layout.Cols),
			layout.Center("", layout.Cols),
		}, nil
	}

	sort.Slice(upcoming, func(i, j int) bool {
		return upcoming[i].start.Before(upcoming[j].start)
	})

	next := upcoming[0]
	return [3]string{
		layout.Center(next.start.Format("MON 1/2"), layout.Cols),
		layout.Center(next.start.Format("3:04 PM"), layout.Cols),
		layout.Center(layout.Truncate(next.summary, layout.Cols), layout.Cols),
	}, nil
}

func FetchMulti(icsURL string, n int) ([3]string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(icsURL)
	if err != nil {
		return [3]string{}, err
	}
	defer resp.Body.Close()

	cal, err := ical.Parse(resp.Body)
	if err != nil {
		return [3]string{}, err
	}

	now := time.Now()
	type event struct {
		start   time.Time
		summary string
	}
	var upcoming []event
	for _, e := range cal.Events {
		if e.Start.After(now) {
			upcoming = append(upcoming, event{start: e.Start, summary: e.Summary})
		}
	}
	sort.Slice(upcoming, func(i, j int) bool {
		return upcoming[i].start.Before(upcoming[j].start)
	})
	if len(upcoming) > n {
		upcoming = upcoming[:n]
	}

	var lines [3]string
	for i := 0; i < 3; i++ {
		if i < len(upcoming) {
			e := upcoming[i]
			lines[i] = layout.PadRight(
				fmt.Sprintf("%s %s", e.start.Format("1/2"), layout.Truncate(e.summary, layout.Cols-4)),
				layout.Cols,
			)
		} else {
			lines[i] = layout.PadRight("", layout.Cols)
		}
	}
	return lines, nil
}
