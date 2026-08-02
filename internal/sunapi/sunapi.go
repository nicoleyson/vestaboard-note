package sunapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var APIURL = "https://api.sunrise-sunset.org/json"

type apiResponse struct {
	Results struct {
		Sunrise string `json:"sunrise"`
		Sunset  string `json:"sunset"`
	} `json:"results"`
	Status string `json:"status"`
}

func FetchTimes(lat, lon float64, date time.Time) (rise, set time.Time, err error) {
	dateStr := date.Local().Format("2006-01-02")
	url := fmt.Sprintf("%s?lat=%f&lng=%f&date=%s&formatted=0", APIURL, lat, lon, dateStr)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	req.Header.Set("User-Agent", "vestaboard-note/1.0 (https://github.com/nicoleyson/vestaboard-note)")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return time.Time{}, time.Time{}, fmt.Errorf("sunrise API HTTP %d", resp.StatusCode)
	}

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
