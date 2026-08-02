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
	Date   string `json:"date"`
	Name   string `json:"localName"`
	Global bool   `json:"global"`
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

func todayHoliday(code string, now time.Time) (string, error) {
	url := fmt.Sprintf("%s/%d/%s", nagerURL, now.Year(), code)
	var holidays []nagerHoliday
	if err := getJSON(url, &holidays); err != nil {
		return "", fmt.Errorf("fetch holidays: %w", err)
	}
	today := now.Format("2006-01-02")
	for _, h := range holidays {
		if h.Date == today {
			return h.Name, nil
		}
	}
	return "", nil
}

func Fetch(lat, lon float64) ([3]string, bool, error) {
	now := time.Now()

	code, err := countryCode(lat, lon)
	if err != nil {
		return [3]string{}, false, err
	}

	name, err := todayHoliday(code, now)
	if err != nil {
		return [3]string{}, false, err
	}

	if name == "" {
		return [3]string{
			layout.ColorRow(70),
			layout.Center("NO HOLIDAY", layout.Cols),
			layout.Center("TODAY", layout.Cols),
		}, true, nil
	}

	upper := strings.ToUpper(name)
	lines := layout.Wrap(upper, layout.Cols)

	row2 := ""
	row3 := ""
	if len(lines) > 0 {
		row2 = layout.Center(lines[0], layout.Cols)
	}
	if len(lines) > 1 {
		row3 = layout.Center(lines[1], layout.Cols)
	}

	return [3]string{
		layout.ColorRow(65),
		row2,
		row3,
	}, false, nil
}
