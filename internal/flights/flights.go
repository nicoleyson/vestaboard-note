package flights

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/nicoleyson/vestaboard-note/internal/geo"
	"github.com/nicoleyson/vestaboard-note/internal/layout"
)

var apiURL = "https://opensky-network.org/api/states/all"

type apiResponse struct {
	States [][]interface{} `json:"states"`
}

type aircraft struct {
	callsign string
	altitude float64
	origin   string
}

func Fetch(lat, lon float64) ([3]string, bool, error) {
	delta := 0.4
	url := fmt.Sprintf("%s?lamin=%f&lomin=%f&lamax=%f&lomax=%f",
		apiURL, lat-delta, lon-delta, lat+delta, lon+delta)

	client := &http.Client{Timeout: 15 * time.Second}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return [3]string{}, true, err
	}
	req.Header.Set("User-Agent", "vestaboard-note/1.0 (https://github.com/nicoleyson/vestaboard-note)")
	resp, err := client.Do(req)
	if err != nil {
		return [3]string{}, true, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		return [3]string{
			layout.Center("FLIGHTS", layout.Cols),
			layout.Center("RATE LIMITED", layout.Cols),
			layout.Center("TRY LATER", layout.Cols),
		}, true, nil
	}
	if resp.StatusCode != http.StatusOK {
		return [3]string{}, true, fmt.Errorf("opensky HTTP %d", resp.StatusCode)
	}

	var data apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return [3]string{}, true, err
	}

	if len(data.States) == 0 {
		return [3]string{
			layout.Center("FLIGHTS", layout.Cols),
			layout.Center("CLEAR SKIES", layout.Cols),
			layout.Center("NONE OVERHEAD", layout.Cols),
		}, true, nil
	}

	ac := parseFirst(data.States)
	if ac == nil {
		return [3]string{
			layout.Center("FLIGHTS", layout.Cols),
			layout.Center("CLEAR SKIES", layout.Cols),
			layout.Center("NONE OVERHEAD", layout.Cols),
		}, true, nil
	}

	row1 := layout.Center("OVERHEAD", layout.Cols)
	row2 := layout.Center(ac.callsign, layout.Cols)

	altStr := ""
	if ac.altitude > 0 {
		if geo.IsFahrenheitCountry(lat, lon) {
			altStr = fmt.Sprintf("%.0fFT", ac.altitude*3.28084)
		} else {
			altStr = fmt.Sprintf("%.0fM", ac.altitude)
		}
	}
	row3 := layout.Center(altStr, layout.Cols)
	if ac.origin != "" && altStr != "" {
		row3 = layout.Center(fmt.Sprintf("%s %s", ac.origin, altStr), layout.Cols)
	}

	return [3]string{row1, row2, row3}, false, nil
}

func parseFirst(states [][]interface{}) *aircraft {
	for _, s := range states {
		if len(s) < 8 {
			continue
		}
		callsign := ""
		if cs, ok := s[1].(string); ok {
			callsign = strings.TrimSpace(cs)
		}
		if callsign == "" {
			continue
		}
		altitude := 0.0
		if alt, ok := s[7].(float64); ok {
			altitude = alt
		}
		origin := ""
		if o, ok := s[2].(string); ok {
			origin = strings.ToUpper(strings.TrimSpace(o))
			if runes := []rune(origin); len(runes) > 4 {
				origin = string(runes[:4])
			}
		}
		return &aircraft{callsign: callsign, altitude: altitude, origin: origin}
	}
	return nil
}
