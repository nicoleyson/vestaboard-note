package air

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nicoleyson/vestaboard-note/internal/layout"
)

const apiURL = "https://air-quality-api.open-meteo.com/v1/air-quality"

type apiResponse struct {
	Current struct {
		USAQI int `json:"us_aqi"`
	} `json:"current"`
}

func Fetch(lat, lon float64) ([3]string, bool, error) {
	url := fmt.Sprintf("%s?latitude=%f&longitude=%f&current=us_aqi", apiURL, lat, lon)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return [3]string{}, false, err
	}
	defer resp.Body.Close()

	var data apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return [3]string{}, false, err
	}

	aqi := data.Current.USAQI
	label, color := classify(aqi)

	trivial := label == "GOOD"
	return [3]string{
		layout.ColorRow(color),
		layout.Center(fmt.Sprintf("AQI %d", aqi), layout.Cols),
		layout.Center(label, layout.Cols),
	}, trivial, nil
}

func classify(aqi int) (label string, color int) {
	switch {
	case aqi <= 50:
		return "GOOD", 66
	case aqi <= 100:
		return "MODERATE", 65
	case aqi <= 150:
		return "UNHEALTHY+", 64
	case aqi <= 200:
		return "UNHEALTHY", 63
	case aqi <= 300:
		return "VERY UNHLTHY", 68
	default:
		return "HAZARDOUS", 70
	}
}
