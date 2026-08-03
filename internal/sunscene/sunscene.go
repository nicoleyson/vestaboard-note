package sunscene

import (
	"fmt"
	"strings"
	"time"

	"github.com/nicoleyson/vestaboard-note/internal/sunapi"
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
		if r == 0 {
			return blue
		}
		return orange
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

// nightRowColor maps row index, moon row, and night progress to a sky color.
// The moon row always gets violet so the white disc reads clearly against it.
func nightRowColor(r, mr int, progress float64) int {
	if r == mr {
		return violet
	}
	if r < mr {
		if progress < 0.20 || progress > 0.80 {
			return violet
		}
		return black
	}
	return violet
}

func renderNight(progress float64) [3]string {
	mr := sunRow(progress)

	var g [rows][cols]int
	for r := 0; r < rows; r++ {
		c := nightRowColor(r, mr, progress)
		for col := 0; col < cols; col++ {
			g[r][col] = c
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
	rise, set, err := sunapi.FetchTimes(lat, lon, now)
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
		_, nightStart, err = sunapi.FetchTimes(lat, lon, yesterday)
		if err != nil {
			return renderNight(0.5), nil
		}
		nextRise = rise
	} else {
		nightStart = set
		tomorrow := now.AddDate(0, 0, 1)
		nextRise, _, err = sunapi.FetchTimes(lat, lon, tomorrow)
		if err != nil {
			return renderNight(0.5), nil
		}
	}

	nightLen := nextRise.Sub(nightStart).Seconds()
	if nightLen <= 0 {
		return renderNight(0.5), nil
	}
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
