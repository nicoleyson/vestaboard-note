package rain

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nicoleyson/vestaboard-note/internal/layout"
)

const apiURL = "https://api.open-meteo.com/v1/forecast"

type apiResponse struct {
	Current struct {
		Precipitation        float64 `json:"precipitation"`
		PrecipitationProb    int     `json:"precipitation_probability"`
	} `json:"current"`
}

func Fetch(lat, lon float64) ([3]string, bool, error) {
	url := fmt.Sprintf(
		"%s?latitude=%f&longitude=%f&current=precipitation,precipitation_probability",
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

	prob := data.Current.PrecipitationProb
	mm := data.Current.Precipitation
	intensity, color := classify(prob, mm)

	trivial := intensity == "NONE"
	return [3]string{
		layout.ColorRow(color),
		layout.Center(fmt.Sprintf("RAIN  %d%%", prob), layout.Cols),
		layout.Center(intensity, layout.Cols),
	}, trivial, nil
}

// classify returns an intensity label and a Vestaboard color code.
// Color: blue (67) when raining or high chance, green (66) when low chance, yellow (65) when dry.
func classify(prob int, mm float64) (label string, color int) {
	switch {
	case mm >= 4.0:
		return "HEAVY", 67
	case mm >= 1.0:
		return "MODERATE", 67
	case mm > 0:
		return "LIGHT", 67
	case prob >= 60:
		return "LIKELY", 67
	case prob >= 30:
		return "POSSIBLE", 66
	default:
		return "NONE", 65
	}
}
