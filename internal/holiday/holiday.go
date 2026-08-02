package holiday

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/nicoleyson/vestaboard-note/internal/layout"
)

const (
	nominatimURL = "https://nominatim.openstreetmap.org/reverse"
	nagerURL     = "https://date.nager.at/api/v3/PublicHolidays"
	userAgent    = "vestaboard-note/1.0 (https://github.com/nicoleyson/vestaboard-note)"
)

type nominatimResponse struct {
	Address struct {
		CountryCode string `json:"country_code"`
	} `json:"address"`
}

type nagerHoliday struct {
	Date string `json:"date"`
	Name string `json:"localName"`
}

func getJSON(url string, target interface{}) error {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
	}
	return json.NewDecoder(resp.Body).Decode(target)
}

func countryCode(lat, lon float64) (string, error) {
	url := fmt.Sprintf("%s?lat=%f&lon=%f&format=json&zoom=3", nominatimURL, lat, lon)
	var resp nominatimResponse
	if err := getJSON(url, &resp); err != nil {
		return "", fmt.Errorf("reverse geocode: %w", err)
	}
	code := strings.ToUpper(resp.Address.CountryCode)
	if code == "" {
		return "", fmt.Errorf("no country code returned for lat=%.4f lon=%.4f", lat, lon)
	}
	return code, nil
}

func fetchHolidays(code string, year int) ([]nagerHoliday, error) {
	url := fmt.Sprintf("%s/%d/%s", nagerURL, year, code)
	var holidays []nagerHoliday
	if err := getJSON(url, &holidays); err != nil {
		return nil, fmt.Errorf("fetch holidays: %w", err)
	}
	return holidays, nil
}

func findToday(holidays []nagerHoliday, today string) string {
	for _, h := range holidays {
		if h.Date == today {
			return h.Name
		}
	}
	return ""
}

func findNext(holidays []nagerHoliday, after string) (name string, days int) {
	afterTime, err := time.Parse("2006-01-02", after)
	if err != nil {
		return "", 0
	}
	for _, h := range holidays {
		if h.Date > after {
			d, err := time.Parse("2006-01-02", h.Date)
			if err != nil {
				continue
			}
			return h.Name, int(d.Sub(afterTime).Hours() / 24)
		}
	}
	return "", 0
}

func artStrip(name string) string {
	n := strings.ToUpper(name)

	type rule struct {
		keyword string
		strip   string
	}

	rules := []rule{
		// Christmas: red/green alternating
		{"CHRISTMAS", "{63}{66}{63}{66}{63}{66}{63}{66}{63}{66}{63}{66}{63}{66}{63}"},
		// Lunar New Year / Chinese New Year: red/gold(yellow) — must precede NEW YEAR
		{"LUNAR NEW YEAR", "{63}{65}{63}{65}{63}{65}{63}{65}{63}{65}{63}{65}{63}{65}{63}"},
		{"CHINESE NEW YEAR", "{63}{65}{63}{65}{63}{65}{63}{65}{63}{65}{63}{65}{63}{65}{63}"},
		// New Year: violet/yellow sparkle
		{"NEW YEAR", "{68}{65}{68}{65}{68}{65}{68}{65}{68}{65}{68}{65}{68}{65}{68}"},
		// Halloween: orange/black
		{"HALLOWEEN", "{64}{70}{64}{70}{64}{70}{64}{70}{64}{70}{64}{70}{64}{70}{64}"},
		// Independence / 4th of July / National Day (generic): red/white/blue
		{"INDEPENDENCE", "{63}{69}{67}{63}{69}{67}{63}{69}{67}{63}{69}{67}{63}{69}{67}"},
		{"NATIONAL DAY", "{63}{69}{67}{63}{69}{67}{63}{69}{67}{63}{69}{67}{63}{69}{67}"},
		// Thanksgiving: orange/yellow harvest
		{"THANKSGIVING", "{64}{65}{64}{65}{64}{65}{64}{65}{64}{65}{64}{65}{64}{65}{64}"},
		// Easter: red/green/yellow (bright pastels)
		{"EASTER", "{63}{66}{65}{63}{66}{65}{63}{66}{65}{63}{66}{65}{63}{66}{65}"},
		// Valentine's Day: red with hearts
		{"VALENTINE", "{63}{62}{63}{62}{63}{62}{63}{62}{63}{62}{63}{62}{63}{62}{63}"},
		// St. Patrick's Day: green/white
		{"PATRICK", "{66}{69}{66}{69}{66}{69}{66}{69}{66}{69}{66}{69}{66}{69}{66}"},
		// Diwali / Deepavali: yellow/orange festival of lights
		{"DIWALI", "{65}{64}{65}{64}{65}{64}{65}{64}{65}{64}{65}{64}{65}{64}{65}"},
		{"DEEPAVALI", "{65}{64}{65}{64}{65}{64}{65}{64}{65}{64}{65}{64}{65}{64}{65}"},
		// Hanukkah: blue/white
		{"HANUKKAH", "{67}{69}{67}{69}{67}{69}{67}{69}{67}{69}{67}{69}{67}{69}{67}"},
		{"CHANUKAH", "{67}{69}{67}{69}{67}{69}{67}{69}{67}{69}{67}{69}{67}{69}{67}"},
		// Eid: green/white crescent palette
		{"EID", "{66}{69}{66}{69}{66}{69}{66}{69}{66}{69}{66}{69}{66}{69}{66}"},
		// Labor / Workers: blue/white
		{"LABOR", "{67}{69}{67}{69}{67}{69}{67}{69}{67}{69}{67}{69}{67}{69}{67}"},
		{"WORKERS", "{67}{69}{67}{69}{67}{69}{67}{69}{67}{69}{67}{69}{67}{69}{67}"},
		// Memorial: red/white/blue
		{"MEMORIAL", "{63}{69}{67}{63}{69}{67}{63}{69}{67}{63}{69}{67}{63}{69}{67}"},
		// Veterans: red/white/blue
		{"VETERANS", "{63}{69}{67}{63}{69}{67}{63}{69}{67}{63}{69}{67}{63}{69}{67}"},
		// Martin Luther King: blue/white
		{"MARTIN LUTHER", "{67}{69}{67}{69}{67}{69}{67}{69}{67}{69}{67}{69}{67}{69}{67}"},
		// Presidents / Washington / Lincoln: red/white/blue
		{"PRESIDENTS", "{63}{69}{67}{63}{69}{67}{63}{69}{67}{63}{69}{67}{63}{69}{67}"},
		// Guy Fawkes: orange/black bonfire
		{"GUY FAWKES", "{64}{70}{64}{70}{64}{70}{64}{70}{64}{70}{64}{70}{64}{70}{64}"},
		// Bastille: blue/white/red (France)
		{"BASTILLE", "{67}{69}{63}{67}{69}{63}{67}{69}{63}{67}{69}{63}{67}{69}{63}"},
		// Anzac / Remembrance / Armistice: red/white (poppy)
		{"ANZAC", "{63}{69}{63}{69}{63}{69}{63}{69}{63}{69}{63}{69}{63}{69}{63}"},
		{"REMEMBRANCE", "{63}{69}{63}{69}{63}{69}{63}{69}{63}{69}{63}{69}{63}{69}{63}"},
		{"ARMISTICE", "{63}{69}{63}{69}{63}{69}{63}{69}{63}{69}{63}{69}{63}{69}{63}"},
	}

	for _, r := range rules {
		if strings.Contains(n, r.keyword) {
			return r.strip
		}
	}

	return layout.ColorRow(65)
}

func TodaysName(lat, lon float64) string {
	code, err := countryCode(lat, lon)
	if err != nil {
		return ""
	}
	holidays, err := fetchHolidays(code, time.Now().Year())
	if err != nil {
		return ""
	}
	return strings.ToUpper(findToday(holidays, time.Now().Format("2006-01-02")))
}

func Fetch(lat, lon float64) ([3]string, bool, error) {
	now := time.Now()
	today := now.Format("2006-01-02")

	code, err := countryCode(lat, lon)
	if err != nil {
		return [3]string{}, false, err
	}

	holidays, err := fetchHolidays(code, now.Year())
	if err != nil {
		return [3]string{}, false, err
	}

	if name := findToday(holidays, today); name != "" {
		upper := strings.ToUpper(name)
		lines := layout.Wrap(upper, layout.Cols)
		row2, row3 := "", ""
		if len(lines) > 0 {
			row2 = layout.Center(lines[0], layout.Cols)
		}
		if len(lines) > 1 {
			row3 = layout.Center(lines[1], layout.Cols)
		}
		return [3]string{artStrip(name), row2, row3}, false, nil
	}

	nextName, days := findNext(holidays, today)
	if nextName == "" {
		nextYear, err := fetchHolidays(code, now.Year()+1)
		if err == nil {
			nextName, days = findNext(nextYear, today)
		}
	}

	if nextName == "" {
		return [3]string{
			layout.ColorRow(70),
			layout.Center("NO HOLIDAYS", layout.Cols),
			layout.Center("FOUND", layout.Cols),
		}, true, nil
	}

	upper := strings.ToUpper(nextName)
	lines := layout.Wrap(upper, layout.Cols)
	row2 := ""
	if len(lines) > 0 {
		row2 = layout.Center(lines[0], layout.Cols)
	}

	var row3 string
	switch days {
	case 1:
		row3 = layout.Center("TOMORROW", layout.Cols)
	default:
		row3 = layout.Center(fmt.Sprintf("IN %d DAYS", days), layout.Cols)
	}

	return [3]string{artStrip(nextName), row2, row3}, true, nil
}
