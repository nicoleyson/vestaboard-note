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
	0: "CLEAR", 1: "MOSTLY CLR", 2: "PARTLY CLDY", 3: "OVERCAST",
	45: "FOG", 48: "ICING FOG",
	51: "LT DRIZZLE", 53: "DRIZZLE", 55: "HVY DRIZZLE",
	61: "LT RAIN", 63: "RAIN", 65: "HVY RAIN",
	71: "LT SNOW", 73: "SNOW", 75: "HVY SNOW", 77: "SNOW GRAINS",
	80: "SHOWERS", 81: "SHOWERS", 82: "HVY SHOWERS",
	85: "SNOW SHWRS", 86: "HVY SNOW SHWRS",
	95: "TSTORM", 96: "TSTORM+HAIL", 99: "TSTORM+HAIL",
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
	line1 := layout.Center(now.Format("Mon Jan 2"), layout.Cols)
	line2 := layout.Center(fmt.Sprintf("%.0fF", data.Current.Temperature), layout.Cols)
	line3 := layout.Center(desc, layout.Cols)

	return [3]string{line1, line2, line3}, nil
}
