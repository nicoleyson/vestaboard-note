package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/nicoleyson/vestaboard-note/internal/layout"
)

const openMeteoURL = "https://api.open-meteo.com/v1/forecast"

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
	case wmoCode <= 1:
		return 65
	case wmoCode <= 3:
		return 69
	case wmoCode <= 48:
		return 69
	case wmoCode <= 55:
		return 67
	case wmoCode <= 65:
		return 67
	case wmoCode <= 77:
		return 69
	case wmoCode <= 82:
		return 67
	case wmoCode <= 86:
		return 69
	case wmoCode >= 95:
		return 68
	default:
		return 64
	}
}

func colorForDesc(desc string) int {
	d := strings.ToUpper(desc)
	switch {
	case strings.Contains(d, "THUNDER") || strings.Contains(d, "STORM"):
		return 68
	case strings.Contains(d, "SNOW") || strings.Contains(d, "ICE") || strings.Contains(d, "BLIZZARD") || strings.Contains(d, "FLURR"):
		return 69
	case strings.Contains(d, "RAIN") || strings.Contains(d, "DRIZZLE") || strings.Contains(d, "SHOWER"):
		return 67
	case strings.Contains(d, "FOG") || strings.Contains(d, "MIST") || strings.Contains(d, "HAZE"):
		return 69
	case strings.Contains(d, "CLOUD") || strings.Contains(d, "OVERCAST"):
		return 69
	case strings.Contains(d, "CLEAR") || strings.Contains(d, "SUNNY") || strings.Contains(d, "FAIR"):
		return 65
	default:
		return 64
	}
}

func descFromNWS(raw string) string {
	d := strings.ToUpper(raw)
	switch {
	case strings.Contains(d, "THUNDER"):
		return "THUNDERSTORM"
	case strings.Contains(d, "BLIZZARD"):
		return "BLIZZARD"
	case strings.Contains(d, "SNOW SHOWER") || strings.Contains(d, "SNOW SHOWERS"):
		return "SNOW SHOWERS"
	case strings.Contains(d, "SNOW"):
		return "SNOW"
	case strings.Contains(d, "FLURR"):
		return "FLURRIES"
	case strings.Contains(d, "ICE") || strings.Contains(d, "SLEET") || strings.Contains(d, "FREEZING"):
		return "ICE"
	case strings.Contains(d, "SHOWER"):
		return "SHOWERS"
	case strings.Contains(d, "DRIZZLE"):
		return "DRIZZLE"
	case strings.Contains(d, "RAIN"):
		return "RAIN"
	case strings.Contains(d, "FOG"):
		return "FOG"
	case strings.Contains(d, "MIST"):
		return "MIST"
	case strings.Contains(d, "HAZE"):
		return "HAZE"
	case strings.Contains(d, "OVERCAST"):
		return "OVERCAST"
	case strings.Contains(d, "CLOUD") || strings.Contains(d, "MOSTLY CLOUDY") || strings.Contains(d, "PARTLY CLOUDY"):
		return "CLOUDY"
	case strings.Contains(d, "CLEAR") || strings.Contains(d, "SUNNY") || strings.Contains(d, "FAIR"):
		return "CLEAR"
	default:
		if raw == "" {
			return "UNKNOWN"
		}
		words := strings.Fields(d)
		if len(words) > 0 {
			return words[0]
		}
		return "UNKNOWN"
	}
}

func isUS(lat, lon float64) bool {
	return lat >= 24 && lat <= 50 && lon >= -125 && lon <= -66
}

func getJSON(ctx context.Context, url string, target interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "vestaboard-note/1.0 +https://github.com/nicoleyson/vestaboard-note")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
	}
	return json.NewDecoder(resp.Body).Decode(target)
}

func fetchNWS(lat, lon float64) (tempF float64, desc string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var pointsResp struct {
		Properties struct {
			ObservationStations string `json:"observationStations"`
		} `json:"properties"`
	}
	pointsURL := fmt.Sprintf("https://api.weather.gov/points/%.4f,%.4f", lat, lon)
	if err = getJSON(ctx, pointsURL, &pointsResp); err != nil {
		return 0, "", fmt.Errorf("nws points: %w", err)
	}

	var stationsResp struct {
		Features []struct {
			Properties struct {
				StationIdentifier string `json:"stationIdentifier"`
			} `json:"properties"`
		} `json:"features"`
	}
	if err = getJSON(ctx, pointsResp.Properties.ObservationStations, &stationsResp); err != nil {
		return 0, "", fmt.Errorf("nws stations: %w", err)
	}
	if len(stationsResp.Features) == 0 {
		return 0, "", fmt.Errorf("nws: no stations found")
	}
	stationID := stationsResp.Features[0].Properties.StationIdentifier

	var obsResp struct {
		Properties struct {
			Temperature struct {
				Value *float64 `json:"value"`
			} `json:"temperature"`
			TextDescription string `json:"textDescription"`
		} `json:"properties"`
	}
	obsURL := fmt.Sprintf("https://api.weather.gov/stations/%s/observations/latest", stationID)
	if err = getJSON(ctx, obsURL, &obsResp); err != nil {
		return 0, "", fmt.Errorf("nws observation: %w", err)
	}

	if obsResp.Properties.Temperature.Value == nil {
		return 0, "", fmt.Errorf("nws: no temperature in observation")
	}
	tempC := *obsResp.Properties.Temperature.Value
	tempF = math.Round(tempC*9/5 + 32)
	desc = descFromNWS(obsResp.Properties.TextDescription)
	return tempF, desc, nil
}

func fetchOpenMeteo(lat, lon float64) (tempF float64, desc string, color int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	url := fmt.Sprintf("%s?latitude=%f&longitude=%f&current=temperature_2m,weathercode&temperature_unit=fahrenheit",
		openMeteoURL, lat, lon)

	var data struct {
		Current struct {
			Temperature float64 `json:"temperature_2m"`
			WeatherCode int     `json:"weathercode"`
		} `json:"current"`
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, "", 0, err
	}
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, "", 0, err
	}

	d, ok := weatherDescriptions[data.Current.WeatherCode]
	if !ok {
		d = fmt.Sprintf("CODE %d", data.Current.WeatherCode)
	}
	return data.Current.Temperature, d, colorForCode(data.Current.WeatherCode), nil
}

func Fetch(lat, lon float64) ([3]string, error) {
	now := time.Now()
	var tempF float64
	var desc string
	var color int

	if isUS(lat, lon) {
		t, d, err := fetchNWS(lat, lon)
		if err != nil {
			var t2 float64
			var d2 string
			var c2 int
			t2, d2, c2, err = fetchOpenMeteo(lat, lon)
			if err != nil {
				return [3]string{}, err
			}
			tempF, desc, color = t2, d2, c2
		} else {
			tempF = t
			desc = d
			color = colorForDesc(d)
		}
	} else {
		t, d, c, err := fetchOpenMeteo(lat, lon)
		if err != nil {
			return [3]string{}, err
		}
		tempF, desc, color = t, d, c
	}

	return [3]string{
		layout.ColorRow(color),
		layout.Center(fmt.Sprintf("%.0fF  %s", tempF, now.Format("Mon Jan 2")), layout.Cols),
		layout.Center(desc, layout.Cols),
	}, nil
}
