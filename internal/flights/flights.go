package flights

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/nicoleyson/vestaboard-note/internal/layout"
)

const apiURL = "https://opensky-network.org/api/states/all"

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
	req, err := http.NewRequest("GET", url, nil)
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
		if isFahrenheitCountry(lat, lon) {
			altStr = fmt.Sprintf("%.0fFT", ac.altitude*3.28084)
		} else {
			altStr = fmt.Sprintf("%.0fM", ac.altitude)
		}
	}
	row3 := layout.Center(altStr, layout.Cols)
	if ac.origin != "" && altStr != "" {
		row3 = layout.Center(fmt.Sprintf("%s  %s", ac.origin, altStr), layout.Cols)
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
			if len(origin) > 4 {
				origin = origin[:4]
			}
		}
		return &aircraft{callsign: callsign, altitude: altitude, origin: origin}
	}
	return nil
}

func isFahrenheitCountry(lat, lon float64) bool {
	if lat >= 24 && lat <= 49.5 && lon >= -125 && lon <= -66 {
		return true
	}
	if lat >= 54 && lat <= 72 && lon >= -168 && lon <= -130 {
		return true
	}
	if lat >= 18 && lat <= 23 && lon >= -161 && lon <= -154 {
		return true
	}
	if lat >= 17 && lat <= 18.5 && lon >= -68 && lon <= -64 {
		return true
	}
	if lat >= 4 && lat <= 8.5 && lon >= -11.5 && lon <= -7.5 {
		return true
	}
	return false
}
