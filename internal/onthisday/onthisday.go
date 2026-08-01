package onthisday

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/nicoleyson/vestaboard-note/internal/layout"
)

const apiURL = "https://en.wikipedia.org/api/rest_v1/feed/onthisday/events/%d/%d"

type wikiEvent struct {
	Year int    `json:"year"`
	Text string `json:"text"`
}

type apiResponse struct {
	Events []wikiEvent `json:"events"`
}

func Fetch(t time.Time) ([3]string, error) {
	url := fmt.Sprintf(apiURL, t.Month(), t.Day())
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return [3]string{}, err
	}
	req.Header.Set("User-Agent", "vestaboard-note/1.0 (https://github.com/nicoleyson/vestaboard-note)")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return [3]string{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var buf bytes.Buffer
		buf.ReadFrom(resp.Body)
		s := buf.String()
		if len(s) > 80 {
			s = s[:80]
		}
		return [3]string{}, fmt.Errorf("wikipedia api %d: %s", resp.StatusCode, s)
	}

	var data apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return [3]string{}, err
	}
	if len(data.Events) == 0 {
		return [3]string{}, fmt.Errorf("no events found for %s %d", t.Month().String()[:3], t.Day())
	}

	e := pickEvent(data.Events)
	wrapped := layout.Wrap(layout.StripEmoji(e.Text), layout.Cols)

	row3 := ""
	if len(wrapped) > 0 {
		row3 = wrapped[0]
	}

	return [3]string{
		layout.Center("ON THIS DAY", layout.Cols),
		layout.Center(fmt.Sprintf("IN %d", e.Year), layout.Cols),
		row3,
	}, nil
}

func pickEvent(events []wikiEvent) wikiEvent {
	var clean []int
	for i, e := range events {
		wrapped := layout.Wrap(layout.StripEmoji(e.Text), layout.Cols)
		if len(wrapped) > 0 && len([]rune(wrapped[0])) == layout.Cols && !endsWithTilde(wrapped[0]) {
			clean = append(clean, i)
		}
	}
	if len(clean) > 0 {
		return events[clean[rand.Intn(len(clean))]]
	}
	return events[rand.Intn(len(events))]
}

func endsWithTilde(s string) bool {
	runes := []rune(s)
	return len(runes) > 0 && runes[len(runes)-1] == '~'
}
