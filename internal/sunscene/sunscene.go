package sunscene

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	cols = 15
	rows = 3
)

const (
	red    = 63
	orange = 64
	yellow = 65
	blue   = 67
	violet = 68
	white  = 69
	black  = 70
)

type apiResponse struct {
	Results struct {
		Sunrise string `json:"sunrise"`
		Sunset  string `json:"sunset"`
	} `json:"results"`
	Status string `json:"status"`
}

func fetchTimes(lat, lon float64, date time.Time) (rise, set time.Time, err error) {
	url := fmt.Sprintf("https://api.sunrise-sunset.org/json?lat=%f&lng=%f&date=%s&formatted=0",
		lat, lon, date.Format("2006-01-02"))
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	defer resp.Body.Close()
	var data apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return time.Time{}, time.Time{}, err
	}
	if data.Status != "OK" {
		return time.Time{}, time.Time{}, fmt.Errorf("sunrise API: %s", data.Status)
	}
	rise, err = time.Parse(time.RFC3339, data.Results.Sunrise)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("parse sunrise: %w", err)
	}
	set, err = time.Parse(time.RFC3339, data.Results.Sunset)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("parse sunset: %w", err)
	}
	return rise, set, nil
}

// sunRow maps 0–1 daylight progress to a board row (0=top, 2=bottom).
// The sun sits at the bottom near dawn/dusk, rises to mid-sky through the
// morning, and reaches the top row only around solar noon.
func sunRow(progress float64) int {
	switch {
	case progress < 0.15 || progress > 0.85:
		return 2
	case progress < 0.35 || progress > 0.65:
		return 1
	default:
		return 0
	}
}

func rowColor(r, sr int, progress float64) int {
	switch {
	case progress < 0.05 || progress > 0.97:
		if r == sr {
			return red
		}
		return black
	case progress < 0.15 || progress > 0.85:
		switch r {
		case 0:
			return black
		case 1:
			return violet
		default:
			return red
		}
	case progress < 0.25 || progress > 0.75:
		switch r {
		case 0:
			return violet
		case 1:
			return orange
		default:
			return red
		}
	case progress < 0.35 || progress > 0.65:
		switch r {
		case 0:
			return blue
		case 1:
			return orange
		default:
			return orange
		}
	default:
		switch r {
		case 0:
			return blue
		case 1:
			return blue
		default:
			return orange
		}
	}
}

func renderDay(progress float64) [3]string {
	sr := sunRow(progress)

	var g [rows][cols]int
	for r := 0; r < rows; r++ {
		c := rowColor(r, sr, progress)
		for col := 0; col < cols; col++ {
			g[r][col] = c
		}
	}

	g[sr][7] = yellow
	g[sr][6] = orange
	g[sr][8] = orange

	return toLines(g)
}

func renderNight(progress float64) [3]string {
	mr := sunRow(progress)

	var g [rows][cols]int
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if r < mr {
				g[r][c] = black
			} else if r == mr {
				g[r][c] = violet
			} else {
				g[r][c] = black
			}
		}
	}

	g[mr][7] = white
	g[mr][6] = violet
	g[mr][8] = violet

	return toLines(g)
}

func toLines(g [rows][cols]int) [3]string {
	var lines [3]string
	for r := 0; r < rows; r++ {
		var b strings.Builder
		for c := 0; c < cols; c++ {
			fmt.Fprintf(&b, "{%d}", g[r][c])
		}
		lines[r] = b.String()
	}
	return lines
}

func Fetch(lat, lon float64) ([3]string, error) {
	now := time.Now()
	rise, set, err := fetchTimes(lat, lon, now)
	if err != nil {
		return [3]string{}, err
	}

	if now.After(rise) && now.Before(set) {
		dayLen := set.Sub(rise).Seconds()
		elapsed := now.Sub(rise).Seconds()
		return renderDay(elapsed / dayLen), nil
	}

	var nightStart, nextRise time.Time
	if now.Before(rise) {
		yesterday := now.AddDate(0, 0, -1)
		_, nightStart, err = fetchTimes(lat, lon, yesterday)
		if err != nil {
			return renderNight(0.5), nil
		}
		nextRise = rise
	} else {
		nightStart = set
		tomorrow := now.AddDate(0, 0, 1)
		nextRise, _, err = fetchTimes(lat, lon, tomorrow)
		if err != nil {
			return renderNight(0.5), nil
		}
	}

	nightLen := nextRise.Sub(nightStart).Seconds()
	elapsed := now.Sub(nightStart).Seconds()
	progress := elapsed / nightLen
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}
	return renderNight(progress), nil
}
