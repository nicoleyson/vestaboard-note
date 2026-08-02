package suntime

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nicoleyson/vestaboard-note/internal/layout"
)

const apiURL = "https://api.sunrise-sunset.org/json"

type apiResponse struct {
	Results struct {
		Sunrise string `json:"sunrise"`
		Sunset  string `json:"sunset"`
	} `json:"results"`
	Status string `json:"status"`
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

	timeStr = eventTime.Local().Format("3:04 PM")
	dateStr := eventTime.Local().Format("Mon Jan 2")

	row2 := layout.Center(fmt.Sprintf("%s  %s", label, timeStr), layout.Cols)
	row3 := layout.Center(dateStr, layout.Cols)

	return [3]string{
		layout.ColorRow(color),
		row2,
		row3,
	}, nil
}

func fetchTimes(lat, lon float64, date time.Time) (rise, set time.Time, err error) {
	dateStr := date.Local().Format("2006-01-02")
	url := fmt.Sprintf("%s?lat=%f&lng=%f&date=%s&formatted=0", apiURL, lat, lon, dateStr)

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
		return time.Time{}, time.Time{}, fmt.Errorf("sunrise API status: %s", data.Status)
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
