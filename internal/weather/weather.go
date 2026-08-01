package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nicoleyson/vestaboard-note/internal/layout"
)

const apiURL = "https://api.open-meteo.com/v1/forecast"

type response struct {
	Current struct {
		Temperature float64 `json:"temperature_2m"`
		WeatherCode int     `json:"weathercode"`
	} `json:"current"`
}

var weatherDescriptions = map[int]string{
	0: "CLEAR", 1: "CLEAR", 2: "CLOUDY", 3: "OVERCAST",
	45: "FOG", 48: "FOG",
	51: "DRIZZLE", 53: "DRIZZLE", 55: "DRIZZLE",
	61: "RAIN", 63: "RAIN", 65: "RAIN",
	71: "SNOW", 73: "SNOW", 75: "SNOW", 77: "SNOW",
	80: "SHOWERS", 81: "SHOWERS", 82: "SHOWERS",
	85: "SNOW SHOWERS", 86: "SNOW SHOWERS",
	95: "THUNDERSTORM", 96: "THUNDERSTORM", 99: "THUNDERSTORM",
}

func colorForCode(wmoCode int) int {
	switch {
	case wmoCode == 0 || wmoCode == 1:
		return 65
	case wmoCode == 2 || wmoCode == 3:
		return 69
	case wmoCode == 45 || wmoCode == 48:
		return 69
	case wmoCode >= 51 && wmoCode <= 55:
		return 67
	case wmoCode >= 61 && wmoCode <= 65:
		return 67
	case wmoCode >= 71 && wmoCode <= 77:
		return 69
	case wmoCode >= 80 && wmoCode <= 82:
		return 67
	case wmoCode >= 85 && wmoCode <= 86:
		return 69
	case wmoCode >= 95:
		return 68
	default:
		return 64
	}
}

func Fetch(lat, lon float64) ([3]string, error) {
	url := fmt.Sprintf("%s?latitude=%f&longitude=%f&current=temperature_2m,weathercode&temperature_unit=fahrenheit",
		apiURL, lat, lon)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return [3]string{}, err
	}
	defer resp.Body.Close()

	var data response
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return [3]string{}, err
	}

	desc, ok := weatherDescriptions[data.Current.WeatherCode]
	if !ok {
		desc = fmt.Sprintf("CODE %d", data.Current.WeatherCode)
	}

	now := time.Now()
	color := colorForCode(data.Current.WeatherCode)

	return [3]string{
		layout.ColorRow(color),
		layout.Center(fmt.Sprintf("%.0fF  %s", data.Current.Temperature, now.Format("Mon Jan 2")), layout.Cols),
		layout.Center(desc, layout.Cols),
	}, nil
}
