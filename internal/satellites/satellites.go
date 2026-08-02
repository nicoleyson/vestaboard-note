package satellites

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/nicoleyson/vestaboard-note/internal/layout"
)

// interesting categories to surface — exclude DEBRIS and STARLINK (hundreds, noisy)
var interestingCategories = map[string]bool{
	"OTHER":   true,
	"GPS":     true,
	"IRIDIUM": true,
}

type satellite struct {
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	ElevationDeg float64 `json:"elevation_deg"`
	Direction   string  `json:"direction"`
}

type overhead struct {
	Count      int         `json:"count"`
	Satellites []satellite `json:"satellites"`
}

// Fetch calls api.satlas.app/api/overhead for the given location.
// Returns lines and trivial=true when nothing interesting is overhead.
// HTTP errors are treated as trivial so --skip-trivial cron jobs never fail noisily.
func Fetch(lat, lon float64) ([3]string, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	url := fmt.Sprintf(
		"https://api.satlas.app/api/overhead?latitude=%g&longitude=%g&min_elevation=10&limit=50",
		lat, lon,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return trivialResult(), true, nil
	}
	req.Header.Set("User-Agent", "vestaboard-note/1.0 (github.com/nicoleyson/vestaboard-note)")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return trivialResult(), true, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return trivialResult(), true, nil
	}

	var result overhead
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return trivialResult(), true, nil
	}

	best := pickBest(result.Satellites)
	if best == nil {
		return trivialResult(), true, nil
	}

	name := cleanName(best.Name)
	row2 := layout.Center(name, layout.Cols)
	row3 := layout.Center(fmt.Sprintf("%d DEG  %s", int(best.ElevationDeg), best.Direction), layout.Cols)

	colorCode := colorForCategory(best.Category)
	lines := [3]string{
		layout.ColorRow(colorCode),
		row2,
		row3,
	}
	return lines, false, nil
}

func pickBest(sats []satellite) *satellite {
	var best *satellite
	for i := range sats {
		s := &sats[i]
		if !interestingCategories[s.Category] {
			continue
		}
		if best == nil || s.ElevationDeg > best.ElevationDeg {
			best = s
		}
	}
	return best
}

// cleanName strips common suffixes like "(ZARYA)" and truncates to 15.
func cleanName(name string) string {
	if i := strings.Index(name, " ("); i > 0 {
		name = name[:i]
	}
	return layout.Truncate(strings.TrimSpace(name), layout.Cols)
}

func colorForCategory(cat string) int {
	switch cat {
	case "GPS":
		return 66 // green
	case "IRIDIUM":
		return 67 // blue
	default:
		return 68 // violet — ISS, weather sats, science sats
	}
}

func trivialResult() [3]string {
	return [3]string{
		layout.ColorRow(70), // black
		layout.Center("NO SATELLITES", layout.Cols),
		layout.Center("OVERHEAD NOW", layout.Cols),
	}
}
