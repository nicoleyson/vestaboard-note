package dewpoint

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
		DewPoint         float64 `json:"dew_point_2m"`
		RelativeHumidity int     `json:"relative_humidity_2m"`
	} `json:"current"`
}

func Fetch(lat, lon float64) ([3]string, bool, error) {
	url := fmt.Sprintf(
		"%s?latitude=%f&longitude=%f&current=dew_point_2m,relative_humidity_2m",
		apiURL, lat, lon,
	)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return [3]string{}, false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return [3]string{}, false, fmt.Errorf("open-meteo dewpoint status %d", resp.StatusCode)
	}

	var data apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return [3]string{}, false, err
	}

	dpC := data.Current.DewPoint
	dpF := dpC*9/5 + 32
	rh := data.Current.RelativeHumidity

	label, color := classify(dpF)
	trivial := label == "DRY"

	row2 := layout.Center(fmt.Sprintf("DP %dF  RH %d%%", int(dpF), rh), layout.Cols)

	return [3]string{
		layout.ColorRow(color),
		row2,
		layout.Center(label, layout.Cols),
	}, trivial, nil
}

func classify(dpF float64) (label string, color int) {
	switch {
	case dpF <= 55:
		return "DRY", 66 // green
	case dpF <= 60:
		return "COMFY", 65 // yellow
	case dpF <= 65:
		return "HUMID", 64 // orange
	case dpF <= 70:
		return "MUGGY", 63 // red
	default:
		return "OPPRESSIVE", 68 // pink/magenta
	}
}
