package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/nicoleyson/vestaboard-note/internal/geo"
	"github.com/nicoleyson/vestaboard-note/internal/layout"
)

var (
	metarStationsURL = "https://aviationweather.gov/api/data/airport"
	metarURL         = "https://aviationweather.gov/api/data/metar"
)

const metarSearchDelta = 2.0

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

func nearestStation(ctx context.Context, lat, lon float64) (string, error) {
	delta := metarSearchDelta
	for attempts := 0; attempts < 3; attempts++ {
		url := fmt.Sprintf("%s?bbox=%.4f,%.4f,%.4f,%.4f&format=json",
			metarStationsURL, lat-delta, lon-delta, lat+delta, lon+delta)

		var stations []struct {
			IcaoID string  `json:"icaoId"`
			Lat    float64 `json:"lat"`
			Lon    float64 `json:"lon"`
		}
		if err := getJSON(ctx, url, &stations); err != nil {
			return "", err
		}
		if len(stations) > 0 {
			best := stations[0]
			bestDist := dist(lat, lon, best.Lat, best.Lon)
			for _, s := range stations[1:] {
				if d := dist(lat, lon, s.Lat, s.Lon); d < bestDist {
					best, bestDist = s, d
				}
			}
			return best.IcaoID, nil
		}
		delta *= 2
	}
	return "", fmt.Errorf("no metar station found near %.4f,%.4f", lat, lon)
}

func dist(lat1, lon1, lat2, lon2 float64) float64 {
	dlat := lat1 - lat2
	dlon := lon1 - lon2
	return dlat*dlat + dlon*dlon
}

func descFromMetar(cover, rawOb string) string {
	raw := strings.ToUpper(rawOb)
	switch {
	case strings.Contains(raw, " TS"):
		return "THUNDERSTORM"
	case strings.Contains(raw, " FZRA") || strings.Contains(raw, " FZDZ"):
		return "FREEZING RAIN"
	case strings.Contains(raw, " BLSN") || strings.Contains(raw, " BLIZZARD"):
		return "BLIZZARD"
	case strings.Contains(raw, " RASN") || strings.Contains(raw, " SNRA"):
		return "RAIN AND SNOW"
	case strings.Contains(raw, " SN") || strings.Contains(raw, " SG") || strings.Contains(raw, " PL"):
		return "SNOW"
	case strings.Contains(raw, " SHRA"):
		return "SHOWERS"
	case strings.Contains(raw, " RA") || strings.Contains(raw, " DZ"):
		return "RAIN"
	case strings.Contains(raw, " FG") || strings.Contains(raw, " BR") || strings.Contains(raw, " HZ"):
		return "FOG"
	}
	switch cover {
	case "OVC":
		return "OVERCAST"
	case "BKN":
		return "CLOUDY"
	case "SCT":
		return "PARTLY CLOUDY"
	case "FEW":
		return "MOSTLY CLEAR"
	default:
		return "CLEAR"
	}
}

func colorForDesc(desc string) int {
	d := strings.ToUpper(desc)
	switch {
	case strings.Contains(d, "THUNDER"):
		return 68
	case strings.Contains(d, "SNOW") || strings.Contains(d, "ICE") || strings.Contains(d, "BLIZZARD") || strings.Contains(d, "FREEZ"):
		return 69
	case strings.Contains(d, "RAIN") || strings.Contains(d, "DRIZZLE") || strings.Contains(d, "SHOWER"):
		return 67
	case strings.Contains(d, "FOG"):
		return 69
	case strings.Contains(d, "CLOUD") || strings.Contains(d, "OVERCAST"):
		return 69
	case strings.Contains(d, "CLEAR") || strings.Contains(d, "SUNNY") || strings.Contains(d, "FAIR"):
		return 65
	default:
		return 64
	}
}

func Fetch(lat, lon float64) ([3]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	stationID, err := nearestStation(ctx, lat, lon)
	if err != nil {
		return [3]string{}, fmt.Errorf("metar station lookup: %w", err)
	}

	url := fmt.Sprintf("%s?ids=%s&format=json", metarURL, stationID)
	var results []struct {
		Temp  float64 `json:"temp"`
		Cover string  `json:"cover"`
		RawOb string  `json:"rawOb"`
	}
	if err := getJSON(ctx, url, &results); err != nil {
		return [3]string{}, fmt.Errorf("metar fetch: %w", err)
	}
	if len(results) == 0 {
		return [3]string{}, fmt.Errorf("metar: no data for station %s", stationID)
	}

	obs := results[0]
	fahrenheit := geo.IsFahrenheitCountry(lat, lon)
	var temp float64
	var unit string
	if fahrenheit {
		temp = math.Round(obs.Temp*9/5 + 32)
		unit = "F"
	} else {
		temp = math.Round(obs.Temp)
		unit = "C"
	}

	desc := descFromMetar(obs.Cover, obs.RawOb)
	color := colorForDesc(desc)
	now := time.Now()

	return [3]string{
		layout.ColorRow(color),
		layout.Center(fmt.Sprintf("%.0f%s %s", temp, unit, now.Format("1/2")), layout.Cols),
		layout.Center(desc, layout.Cols),
	}, nil
}
