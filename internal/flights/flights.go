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

func Fetch(lat, lon float64) ([3]string, error) {
	delta := 0.4
	url := fmt.Sprintf("%s?lamin=%f&lomin=%f&lamax=%f&lomax=%f",
		apiURL, lat-delta, lon-delta, lat+delta, lon+delta)

	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return [3]string{}, err
	}
	req.Header.Set("User-Agent", "vestaboard-note/1.0 (https://github.com/nicoleyson/vestaboard-note)")
	resp, err := client.Do(req)
	if err != nil {
		return [3]string{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		return [3]string{
			layout.Center("FLIGHTS", layout.Cols),
			layout.Center("RATE LIMITED", layout.Cols),
			layout.Center("TRY LATER", layout.Cols),
		}, nil
	}

	var data apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return [3]string{}, err
	}

	if len(data.States) == 0 {
		return [3]string{
			layout.Center("FLIGHTS", layout.Cols),
			layout.Center("CLEAR SKIES", layout.Cols),
			layout.Center("NONE OVERHEAD", layout.Cols),
		}, nil
	}

	ac := parseFirst(data.States)

	row1 := layout.Center("OVERHEAD", layout.Cols)
	row2 := layout.Center(ac.callsign, layout.Cols)

	altStr := ""
	if ac.altitude > 0 {
		altStr = fmt.Sprintf("%.0fFT", ac.altitude*3.28084)
	}
	row3 := layout.Center(altStr, layout.Cols)
	if ac.origin != "" && altStr != "" {
		row3 = layout.Center(fmt.Sprintf("%s  %s", ac.origin, altStr), layout.Cols)
	}

	return [3]string{row1, row2, row3}, nil
}

func parseFirst(states [][]interface{}) aircraft {
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
		return aircraft{callsign: callsign, altitude: altitude, origin: origin}
	}
	return aircraft{callsign: "UNKNOWN"}
}
