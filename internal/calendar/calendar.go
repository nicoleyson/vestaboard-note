package calendar

import (
	"net/http"
	"sort"
	"time"

	"github.com/bounoable/ical"
	"github.com/nicoleyson/vestaboard-note/internal/layout"
)

type event struct {
	uid     string
	start   time.Time
	summary string
}

func fetchURL(url string, client *http.Client) ([]event, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	cal, err := ical.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	var events []event
	for _, e := range cal.Events {
		events = append(events, event{uid: e.UID, start: e.Start, summary: e.Summary})
	}
	return events, nil
}

func Fetch(urls []string) ([3]string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	seen := map[string]bool{}
	now := time.Now()
	cutoff := now.Add(7 * 24 * time.Hour)

	var upcoming []event
	for _, url := range urls {
		events, err := fetchURL(url, client)
		if err != nil {
			return [3]string{}, err
		}
		for _, e := range events {
			if seen[e.uid] {
				continue
			}
			seen[e.uid] = true
			if e.start.After(now) && e.start.Before(cutoff) {
				upcoming = append(upcoming, e)
			}
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
		layout.Center(next.start.Format("Mon 1/2"), layout.Cols),
		layout.Center(next.start.Format("3:04 PM"), layout.Cols),
		layout.Center(layout.Truncate(next.summary, layout.Cols), layout.Cols),
	}, nil
}
