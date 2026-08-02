package uv

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
		UVIndex float64 `json:"uv_index"`
	} `json:"current"`
}

func Fetch(lat, lon float64) ([3]string, error) {
	url := fmt.Sprintf("%s?latitude=%f&longitude=%f&current=uv_index", apiURL, lat, lon)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return [3]string{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return [3]string{}, fmt.Errorf("open-meteo uv status %d", resp.StatusCode)
	}

	var data apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return [3]string{}, err
	}

	idx := data.Current.UVIndex
	label, color := classify(idx)

	return [3]string{
		layout.ColorRow(color),
		layout.Center(fmt.Sprintf("UV INDEX  %.0f", idx), layout.Cols),
		layout.Center(label, layout.Cols),
	}, nil
}

func classify(idx float64) (label string, color int) {
	switch {
	case idx < 3:
		return "LOW", 66
	case idx < 6:
		return "MODERATE", 65
	case idx < 8:
		return "HIGH", 64
	case idx < 11:
		return "VERY HIGH", 63
	default:
		return "EXTREME", 68
	}
}
