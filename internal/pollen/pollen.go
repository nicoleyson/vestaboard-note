package pollen

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
		GrassPollen float64 `json:"grass_pollen"`
		TreePollen  float64 `json:"tree_pollen"`
		WeedPollen  float64 `json:"weed_pollen"`
	} `json:"current"`
}

func Fetch(lat, lon float64) ([3]string, bool, error) {
	url := fmt.Sprintf("%s?latitude=%f&longitude=%f&current=grass_pollen,tree_pollen,weed_pollen", apiURL, lat, lon)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return [3]string{}, false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return [3]string{}, false, fmt.Errorf("open-meteo pollen status %d", resp.StatusCode)
	}

	var data apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return [3]string{}, false, err
	}

	grass := data.Current.GrassPollen
	tree := data.Current.TreePollen
	weed := data.Current.WeedPollen

	dominant, dominantVal := dominantType(grass, tree, weed)
	label, color := classify(dominantVal)

	trivial := label == "LOW"
	return [3]string{
		layout.ColorRow(color),
		layout.Center(fmt.Sprintf("POLLEN %s", label), layout.Cols),
		layout.Center(dominant, layout.Cols),
	}, trivial, nil
}

func dominantType(grass, tree, weed float64) (name string, val float64) {
	if grass >= tree && grass >= weed {
		return "GRASS", grass
	}
	if tree >= weed {
		return "TREE", tree
	}
	return "WEED", weed
}

func classify(val float64) (label string, color int) {
	switch {
	case val < 10:
		return "LOW", 66
	case val < 50:
		return "MODERATE", 65
	case val < 200:
		return "HIGH", 64
	default:
		return "V.HIGH", 63
	}
}
