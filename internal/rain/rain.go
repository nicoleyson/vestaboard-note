package rain

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nicoleyson/vestaboard-note/internal/layout"
)

var apiURL = "https://api.open-meteo.com/v1/forecast"

type apiResponse struct {
	Current struct {
		Precipitation float64 `json:"precipitation"`
	} `json:"current"`
}

func Fetch(lat, lon float64) ([3]string, bool, error) {
	url := fmt.Sprintf(
		"%s?latitude=%f&longitude=%f&current=precipitation",
		apiURL, lat, lon,
	)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return [3]string{}, false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return [3]string{}, false, fmt.Errorf("open-meteo status %d", resp.StatusCode)
	}

	var data apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return [3]string{}, false, err
	}

	mm := data.Current.Precipitation
	intensity, color := classify(mm)

	trivial := intensity == "NONE"
	return [3]string{
		layout.ColorRow(color),
		layout.Center(fmt.Sprintf("RAIN %.1fMM", mm), layout.Cols),
		layout.Center(intensity, layout.Cols),
	}, trivial, nil
}

// classify returns an intensity label and a Vestaboard color code based on
// precipitation in mm. open-meteo's precipitation_probability is hourly-only
// and silently returns 0 in the current-weather context, so we use mm only.
func classify(mm float64) (label string, color int) {
	switch {
	case mm >= 4.0:
		return "HEAVY", 67
	case mm >= 1.0:
		return "MODERATE", 67
	case mm > 0:
		return "LIGHT", 67
	default:
		return "NONE", 65
	}
}
