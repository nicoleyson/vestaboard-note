package discogs

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCleanArtistName(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"The Beatles (2)", "The Beatles"},
		{"Miles Davis", "Miles Davis"},
		{"Prince (3)", "Prince"},
		{"", ""},
		{"Band (no closing paren", "Band"},
	}
	for _, tc := range cases {
		got := cleanArtistName(tc.input)
		if got != tc.want {
			t.Errorf("cleanArtistName(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestScoreRecord_baseScore(t *testing.T) {
	r := record{artist: "X", title: "Y", styles: []string{}}
	got := scoreRecord(r, nil, nil, morning)
	if got != 1 {
		t.Errorf("base score should be 1, got %d", got)
	}
}

func TestScoreRecord_weatherMatch(t *testing.T) {
	r := record{styles: []string{"Jazz"}}
	got := scoreRecord(r, []string{"Jazz"}, nil, morning)
	if got != 3 {
		t.Errorf("weather match should add 2 (total 3), got %d", got)
	}
}

func TestScoreRecord_timeMatch(t *testing.T) {
	r := record{styles: []string{"Ambient"}}
	got := scoreRecord(r, nil, []string{"Ambient"}, morning)
	if got != 2 {
		t.Errorf("time match should add 1 (total 2), got %d", got)
	}
}

func TestScoreRecord_titleMatch(t *testing.T) {
	r := record{title: "Morning Light", styles: []string{}}
	got := scoreRecord(r, nil, nil, morning)
	if got != 4 {
		t.Errorf("title match should add 3 (total 4), got %d", got)
	}
}

func TestScoreRecord_combinedMatch(t *testing.T) {
	r := record{title: "Night Train", styles: []string{"Blues", "Jazz"}}
	got := scoreRecord(r, []string{"Blues"}, []string{"Jazz"}, night)
	if got != 1+2+1+3 {
		t.Errorf("combined score: want 7, got %d", got)
	}
}

func TestWeightedPick_singleRecord(t *testing.T) {
	records := []record{{artist: "Solo", title: "One"}}
	scores := []int{5}
	got := weightedPick(records, scores)
	if got.artist != "Solo" {
		t.Errorf("single record: want Solo, got %q", got.artist)
	}
}

func TestWeightedPick_allSameScore(t *testing.T) {
	records := []record{
		{artist: "A"}, {artist: "B"}, {artist: "C"},
	}
	scores := []int{1, 1, 1}
	seen := map[string]bool{}
	for i := 0; i < 30; i++ {
		got := weightedPick(records, scores)
		seen[got.artist] = true
	}
	if len(seen) < 2 {
		t.Error("expected distribution across records with equal scores")
	}
}

func TestWeightedPick_highScoreDominates(t *testing.T) {
	records := []record{{artist: "Rare"}, {artist: "Common"}}
	scores := []int{1, 99}
	commonCount := 0
	for i := 0; i < 100; i++ {
		if weightedPick(records, scores).artist == "Common" {
			commonCount++
		}
	}
	if commonCount < 80 {
		t.Errorf("high-score record should dominate: got %d/100", commonCount)
	}
}

func makeSinglePageCollection(artists []string) collectionPage {
	var pg collectionPage
	pg.Pagination.Pages = 1
	for _, a := range artists {
		pg.Releases = append(pg.Releases, struct {
			BasicInformation struct {
				Title   string `json:"title"`
				Artists []struct {
					Name string `json:"name"`
				} `json:"artists"`
				Styles []string `json:"styles"`
				Genres []string `json:"genres"`
			} `json:"basic_information"`
		}{BasicInformation: struct {
			Title   string `json:"title"`
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
			Styles []string `json:"styles"`
			Genres []string `json:"genres"`
		}{Title: "Album", Artists: []struct {
			Name string `json:"name"`
		}{{Name: a}}, Styles: []string{"Rock"}}})
	}
	return pg
}

func TestFetchCollection_singlePage(t *testing.T) {
	pg := makeSinglePageCollection([]string{"The Beatles", "Miles Davis"})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pg)
	}))
	defer srv.Close()

	old := apiBase
	apiBase = srv.URL
	defer func() { apiBase = old }()

	records, err := fetchCollection("user", "token")
	if err != nil {
		t.Fatalf("fetchCollection error: %v", err)
	}
	if len(records) != 2 {
		t.Errorf("fetchCollection: got %d records, want 2", len(records))
	}
}

func TestFetchCollection_httpError_returnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "forbidden", http.StatusForbidden)
	}))
	defer srv.Close()

	old := apiBase
	apiBase = srv.URL
	defer func() { apiBase = old }()

	_, err := fetchCollection("user", "token")
	if err == nil {
		t.Error("expected error for HTTP 403")
	}
}

func TestFetchWMO_success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"current": map[string]interface{}{
				"weather_code": 61,
			},
		})
	}))
	defer srv.Close()

	old := wmoURL
	wmoURL = srv.URL
	defer func() { wmoURL = old }()

	code, err := fetchWMO(37.7, -122.4)
	if err != nil {
		t.Fatalf("fetchWMO error: %v", err)
	}
	if code != 61 {
		t.Errorf("fetchWMO: got %d, want 61", code)
	}
}

func TestFetchWMO_httpError_returnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	old := wmoURL
	wmoURL = srv.URL
	defer func() { wmoURL = old }()

	_, err := fetchWMO(37.7, -122.4)
	if err == nil {
		t.Error("expected error for HTTP 503")
	}
}

func TestFetch_rowLengths(t *testing.T) {
	pg := makeSinglePageCollection([]string{"The Beatles"})
	collSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pg)
	}))
	defer collSrv.Close()

	wmoSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"current": map[string]interface{}{"weather_code": 0},
		})
	}))
	defer wmoSrv.Close()

	oldAPI := apiBase
	apiBase = collSrv.URL
	defer func() { apiBase = oldAPI }()

	oldWMO := wmoURL
	wmoURL = wmoSrv.URL
	defer func() { wmoURL = oldWMO }()

	lines, err := Fetch("user", "token", 37.7, -122.4)
	if err != nil {
		t.Fatalf("Fetch error: %v", err)
	}
	for i, line := range lines {
		if len([]rune(line)) != 15 {
			t.Errorf("Fetch row %d len = %d, want 15: %q", i, len([]rune(line)), line)
		}
	}
}

func TestFetch_emptyCollection_returnsError(t *testing.T) {
	var empty collectionPage
	empty.Pagination.Pages = 1
	collSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(empty)
	}))
	defer collSrv.Close()

	wmoSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"current": map[string]interface{}{"weather_code": 0},
		})
	}))
	defer wmoSrv.Close()

	oldAPI := apiBase
	apiBase = collSrv.URL
	defer func() { apiBase = oldAPI }()

	oldWMO := wmoURL
	wmoURL = wmoSrv.URL
	defer func() { wmoURL = oldWMO }()

	_, err := Fetch("user", "token", 37.7, -122.4)
	if err == nil {
		t.Error("expected error for empty collection")
	}
	if !strings.Contains(err.Error(), "empty") {
		t.Errorf("error should mention empty collection: %v", err)
	}
}

