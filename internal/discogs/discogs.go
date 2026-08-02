package discogs

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/nicoleyson/vestaboard-note/internal/layout"
)

const apiBase = "https://api.discogs.com"

type collectionPage struct {
	Pagination struct {
		Pages int `json:"pages"`
	} `json:"pagination"`
	Releases []struct {
		BasicInformation struct {
			Title   string `json:"title"`
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
			Styles []string `json:"styles"`
			Genres []string `json:"genres"`
		} `json:"basic_information"`
	} `json:"releases"`
}

type record struct {
	artist string
	title  string
	styles []string
}

func fetchCollection(username, token string) ([]record, error) {
	client := &http.Client{Timeout: 60 * time.Second}
	var records []record
	page := 1
	for {
		url := fmt.Sprintf("%s/users/%s/collection/folders/0/releases?per_page=100&page=%d", apiBase, username, page)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			cancel()
			return nil, err
		}
		req.Header.Set("Authorization", "Discogs token="+token)
		req.Header.Set("User-Agent", "vestaboard-note/1.0 +https://github.com/nicoleyson/vestaboard-note")

		resp, err := client.Do(req)
		cancel()
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("discogs API returned %d", resp.StatusCode)
		}

		var pg collectionPage
		err = json.NewDecoder(resp.Body).Decode(&pg)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}

		for _, rel := range pg.Releases {
			info := rel.BasicInformation
			artist := ""
			if len(info.Artists) > 0 {
				artist = cleanArtistName(info.Artists[0].Name)
			}
			tags := append(info.Styles, info.Genres...)
			records = append(records, record{
				artist: artist,
				title:  info.Title,
				styles: tags,
			})
		}

		if page >= pg.Pagination.Pages {
			break
		}
		page++
	}
	return records, nil
}

func cleanArtistName(name string) string {
	if idx := strings.LastIndex(name, " ("); idx != -1 {
		return name[:idx]
	}
	return name
}

func scoreRecord(r record, weatherStyles, timeStyles []string, slot timeSlot) int {
	score := 1
	for _, want := range weatherStyles {
		for _, have := range r.styles {
			if strings.EqualFold(want, have) {
				score += 2
				goto doneWeather
			}
		}
	}
doneWeather:
	for _, want := range timeStyles {
		for _, have := range r.styles {
			if strings.EqualFold(want, have) {
				score++
				goto doneTime
			}
		}
	}
doneTime:
	if titleMatchesSlot(r.title, slot) {
		score += 3
	}
	return score
}

func weightedPick(records []record, scores []int) record {
	total := 0
	for _, s := range scores {
		total += s
	}
	n := rand.Intn(total)
	for i, s := range scores {
		n -= s
		if n < 0 {
			return records[i]
		}
	}
	return records[len(records)-1]
}

func Fetch(username, token string, lat, lon float64) ([3]string, error) {
	wmoCode, err := fetchWMO(lat, lon)
	if err != nil {
		wmoCode = 0
	}

	now := time.Now()
	s := seasonFor(now)
	slot := timeSlotFor(now)

	weatherStyles := stylesFor(wmoCode, s)
	timeStyles := timeSlotStyles[slot]
	label := vibeLabel(wmoCode, s, slot)

	records, err := fetchCollection(username, token)
	if err != nil {
		return [3]string{}, err
	}
	if len(records) == 0 {
		return [3]string{}, fmt.Errorf("discogs collection is empty")
	}

	scores := make([]int, len(records))
	for i, r := range records {
		scores[i] = scoreRecord(r, weatherStyles, timeStyles, slot)
	}

	chosen := weightedPick(records, scores)

	return [3]string{
		layout.Center(label, layout.Cols),
		layout.PadRight(layout.Truncate(layout.StripEmoji(chosen.artist), layout.Cols), layout.Cols),
		layout.PadRight(layout.Truncate(layout.StripEmoji(chosen.title), layout.Cols), layout.Cols),
	}, nil
}

func fetchWMO(lat, lon float64) (int, error) {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=weather_code", lat, lon)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("open-meteo wmo status %d", resp.StatusCode)
	}

	var data struct {
		Current struct {
			WeatherCode int `json:"weather_code"`
		} `json:"current"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}
	return data.Current.WeatherCode, nil
}
