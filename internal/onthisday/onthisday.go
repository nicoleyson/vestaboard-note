package onthisday

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/nicoleyson/vestaboard-note/internal/layout"
)

const apiURL = "https://en.wikipedia.org/api/rest_v1/feed/onthisday/events/%d/%d"

type apiResponse struct {
	Events []struct {
		Year int    `json:"year"`
		Text string `json:"text"`
	} `json:"events"`
}

func Fetch(t time.Time) ([3]string, error) {
	url := fmt.Sprintf(apiURL, t.Month(), t.Day())
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return [3]string{}, err
	}
	defer resp.Body.Close()

	var data apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return [3]string{}, err
	}
	if len(data.Events) == 0 {
		return [3]string{}, fmt.Errorf("no events found for %s %d", t.Month().String()[:3], t.Day())
	}

	event := data.Events[rand.Intn(len(data.Events))]
	yearStr := fmt.Sprintf("IN %d", event.Year)
	wrapped := layout.Wrap(layout.StripEmoji(event.Text), layout.Cols)

	row1 := layout.Center("ON THIS DAY", layout.Cols)
	row2 := layout.Center(yearStr, layout.Cols)
	row3 := ""
	if len(wrapped) > 0 {
		row3 = wrapped[0]
	}

	return [3]string{row1, row2, row3}, nil
}
