package pattern

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/nicoleyson/vestaboard-note/internal/season"
)

const currentUserAgent = "vestaboard-note/1.0 (https://github.com/nicoleyson/vestaboard-note)"

type currentNominatimResp struct {
	Address struct {
		CountryCode string `json:"country_code"`
	} `json:"address"`
}

type currentNagerHoliday struct {
	Date string `json:"date"`
	Name string `json:"localName"`
}

func currentGetJSON(url string, target interface{}) error {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", currentUserAgent)
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(target)
}

func currentCountryCode(lat, lon float64) string {
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?lat=%f&lon=%f&format=json&zoom=3", lat, lon)
	var resp currentNominatimResp
	if err := currentGetJSON(url, &resp); err != nil {
		return ""
	}
	return strings.ToUpper(resp.Address.CountryCode)
}

func todaysHoliday(lat, lon float64) string {
	code := currentCountryCode(lat, lon)
	if code == "" {
		return ""
	}
	today := time.Now().Format("2006-01-02")
	year := time.Now().Year()
	url := fmt.Sprintf("https://date.nager.at/api/v3/PublicHolidays/%d/%s", year, code)
	var holidays []currentNagerHoliday
	if err := currentGetJSON(url, &holidays); err != nil {
		return ""
	}
	for _, h := range holidays {
		if h.Date == today {
			return strings.ToUpper(h.Name)
		}
	}
	return ""
}

type theme struct {
	palette  []int
	patterns []string
}

var holidayThemes = []struct {
	keyword string
	t       theme
}{
	{"LUNAR NEW YEAR", theme{[]int{cRed, cYellow}, []string{"stripes", "checker", "diagonal", "bars"}}},
	{"CHINESE NEW YEAR", theme{[]int{cRed, cYellow}, []string{"stripes", "checker", "diagonal", "bars"}}},
	{"NEW YEAR", theme{[]int{cViolet, cYellow, cWhite}, []string{"confetti", "rainbow", "sparkle"}}},
	{"CHRISTMAS", theme{[]int{cRed, cGreen}, []string{"checker", "stripes", "diagonal", "bars"}}},
	{"HALLOWEEN", theme{[]int{cOrange, cBlack}, []string{"checker", "diagonal", "stripes"}}},
	{"VALENTINE", theme{[]int{cRed, cHeart}, []string{"hearts", "stripes", "checker"}}},
	{"PATRICK", theme{[]int{cGreen, cWhite}, []string{"stripes", "checker", "bars"}}},
	{"EASTER", theme{[]int{cRed, cGreen, cYellow}, []string{"checker", "confetti", "fade"}}},
	{"DIWALI", theme{[]int{cYellow, cOrange}, []string{"confetti", "sparkle", "rainbow"}}},
	{"DEEPAVALI", theme{[]int{cYellow, cOrange}, []string{"confetti", "sparkle", "rainbow"}}},
	{"HANUKKAH", theme{[]int{cBlue, cWhite}, []string{"stripes", "bars", "fade"}}},
	{"CHANUKAH", theme{[]int{cBlue, cWhite}, []string{"stripes", "bars", "fade"}}},
	{"EID", theme{[]int{cGreen, cWhite}, []string{"stripes", "checker", "bars"}}},
	{"INDEPENDENCE", theme{[]int{cRed, cWhite, cBlue}, []string{"stripes", "bars", "diagonal"}}},
	{"NATIONAL DAY", theme{[]int{cRed, cWhite, cBlue}, []string{"stripes", "bars", "diagonal"}}},
	{"THANKSGIVING", theme{[]int{cOrange, cYellow}, []string{"fade", "diagonal", "stripes"}}},
	{"MEMORIAL", theme{[]int{cRed, cWhite, cBlue}, []string{"stripes", "bars"}}},
	{"VETERANS", theme{[]int{cRed, cWhite, cBlue}, []string{"stripes", "bars"}}},
	{"BASTILLE", theme{[]int{cBlue, cWhite, cRed}, []string{"stripes", "bars", "diagonal"}}},
	{"ANZAC", theme{[]int{cRed, cWhite}, []string{"stripes", "bars", "checker"}}},
	{"REMEMBRANCE", theme{[]int{cRed, cWhite}, []string{"stripes", "bars", "checker"}}},
	{"GUY FAWKES", theme{[]int{cOrange, cBlack}, []string{"checker", "diagonal", "stripes"}}},
}

var seasonThemes = map[season.Season][]theme{
	season.Spring: {
		{[]int{cGreen, cWhite, cYellow}, []string{"hearts", "fade", "stripes"}},
		{[]int{cGreen, cYellow, cHeart}, []string{"hearts", "confetti", "checker"}},
		{[]int{cGreen, cWhite}, []string{"stripes", "bars", "fade"}},
	},
	season.Summer: {
		{[]int{cYellow, cOrange}, []string{"fade", "rainbow", "bars"}},
		{[]int{cYellow, cOrange, cRed}, []string{"stripes", "diagonal", "fade"}},
		{[]int{cBlue, cYellow}, []string{"checker", "stripes", "diagonal"}},
	},
	season.Fall: {
		{[]int{cOrange, cYellow}, []string{"fade", "diagonal", "stripes"}},
		{[]int{cOrange, cRed, cYellow}, []string{"diagonal", "checker", "bars"}},
		{[]int{cOrange, cBlack}, []string{"diagonal", "checker", "stripes"}},
	},
	season.Winter: {
		{[]int{cWhite, cViolet}, []string{"sparkle", "bars", "stripes"}},
		{[]int{cBlue, cWhite}, []string{"stripes", "checker", "fade"}},
		{[]int{cWhite, cBlue, cViolet}, []string{"sparkle", "fade", "diagonal"}},
	},
}

func themeForProgress(s season.Season, progress float64) theme {
	themes := seasonThemes[s]
	var idx int
	switch {
	case progress < 0.33:
		idx = 0
	case progress < 0.67:
		idx = 1
	default:
		idx = 2
	}
	return themes[idx]
}

func renderWithTheme(t theme) [3]string {
	patName := t.patterns[rand.Intn(len(t.patterns))]
	pal := t.palette
	var g grid
	switch patName {
	case "stripes":
		g = stripesP(pal)
	case "checker":
		g = checkerP(pal)
	case "bars":
		g = barsP(pal)
	case "fade":
		g = fadeP(pal)
	case "diagonal":
		g = diagonalP(pal)
	case "hearts":
		g = heartsP(pal)
	case "confetti":
		g = confettiP(pal)
	case "sparkle":
		g = sparkleP(pal)
	default:
		g = stripes()
	}
	return g.toLines()
}

func Current(lat, lon float64) ([3]string, error) {
	holiday := todaysHoliday(lat, lon)
	if holiday != "" {
		for _, ht := range holidayThemes {
			if strings.Contains(holiday, ht.keyword) {
				return renderWithTheme(ht.t), nil
			}
		}
	}

	s, progress := season.Current(time.Now())
	return renderWithTheme(themeForProgress(s, progress)), nil
}

func stripesP(pal []int) grid {
	var g grid
	width := 1 + rand.Intn(3)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			g[r][c] = pal[(c/width)%len(pal)]
		}
	}
	return g
}

func checkerP(pal []int) grid {
	var g grid
	size := 1 + rand.Intn(2)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			g[r][c] = pal[((r/size)+(c/size))%len(pal)]
		}
	}
	return g
}

func barsP(pal []int) grid {
	var g grid
	for r := 0; r < rows; r++ {
		c := pal[r%len(pal)]
		for col := 0; col < cols; col++ {
			g[r][col] = c
		}
	}
	return g
}

func fadeP(pal []int) grid {
	var g grid
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			idx := c * len(pal) / cols
			g[r][c] = pal[idx]
		}
	}
	return g
}

func diagonalP(pal []int) grid {
	var g grid
	width := 2 + rand.Intn(3)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			g[r][c] = pal[((r+c)/width)%len(pal)]
		}
	}
	return g
}

func heartsP(pal []int) grid {
	var g grid
	bg := pal[0]
	density := 3 + rand.Intn(5)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			g[r][c] = bg
		}
	}
	placed, attempts := 0, 0
	for placed < density && attempts < 100 {
		r := rand.Intn(rows)
		c := rand.Intn(cols)
		if g[r][c] == bg {
			g[r][c] = cHeart
			placed++
		}
		attempts++
	}
	return g
}

func confettiP(pal []int) grid {
	var g grid
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			g[r][c] = pal[rand.Intn(len(pal))]
		}
	}
	return g
}

func sparkleP(pal []int) grid {
	var g grid
	bg := pal[len(pal)-1]
	density := 4 + rand.Intn(6)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			g[r][c] = bg
		}
	}
	placed, attempts := 0, 0
	for placed < density && attempts < 100 {
		r := rand.Intn(rows)
		c := rand.Intn(cols)
		if g[r][c] == bg {
			g[r][c] = pal[rand.Intn(len(pal)-1)]
			placed++
		}
		attempts++
	}
	return g
}
