package sunscene

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	red    = 63
	orange = 64
	yellow = 65
	blue   = 67
	violet = 68
	black  = 70
)

var sunriseScene = [3][15]int{
	{black, black, black, violet, violet, blue, blue, blue, violet, violet, black, black, black, black, black},
	{violet, blue, yellow, yellow, yellow, yellow, yellow, yellow, yellow, yellow, blue, violet, violet, black, black},
	{orange, orange, red, red, orange, yellow, yellow, yellow, orange, red, red, orange, orange, violet, black},
}

var sunsetScene = [3][15]int{
	{blue, violet, violet, red, red, orange, orange, orange, red, red, violet, violet, blue, blue, black},
	{orange, red, red, yellow, yellow, yellow, yellow, yellow, yellow, yellow, red, red, orange, violet, blue},
	{red, red, red, orange, red, red, orange, orange, red, red, orange, red, red, red, violet},
}

func toLines(scene [3][15]int) [3]string {
	var lines [3]string
	for r := 0; r < 3; r++ {
		var b strings.Builder
		for c := 0; c < 15; c++ {
			fmt.Fprintf(&b, "{%d}", scene[r][c])
		}
		lines[r] = b.String()
	}
	return lines
}

type apiResponse struct {
	Results struct {
		Sunrise string `json:"sunrise"`
		Sunset  string `json:"sunset"`
	} `json:"results"`
	Status string `json:"status"`
}

func fetchTimes(lat, lon float64, date time.Time) (rise, set time.Time, err error) {
	url := fmt.Sprintf("https://api.sunrise-sunset.org/json?lat=%f&lng=%f&date=%s&formatted=0",
		lat, lon, date.Format("2006-01-02"))
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	defer resp.Body.Close()
	var data apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return time.Time{}, time.Time{}, err
	}
	if data.Status != "OK" {
		return time.Time{}, time.Time{}, fmt.Errorf("sunrise API: %s", data.Status)
	}
	rise, err = time.Parse(time.RFC3339, data.Results.Sunrise)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("parse sunrise: %w", err)
	}
	set, err = time.Parse(time.RFC3339, data.Results.Sunset)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("parse sunset: %w", err)
	}
	return rise, set, nil
}

func Fetch(lat, lon float64) ([3]string, error) {
	now := time.Now()
	rise, set, err := fetchTimes(lat, lon, now)
	if err != nil {
		return [3]string{}, err
	}
	if now.After(set) {
		tomorrow := now.AddDate(0, 0, 1)
		rise, set, err = fetchTimes(lat, lon, tomorrow)
		if err != nil {
			return [3]string{}, err
		}
	}
	if now.Before(rise) {
		return toLines(sunriseScene), nil
	}
	return toLines(sunsetScene), nil
}
