package suntime

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/nicoleyson/vestaboard-note/internal/sunapi"
)

var tileRe = regexp.MustCompile(`\{\d+\}|.`)

func countTiles(s string) int {
	return len(tileRe.FindAllString(s, -1))
}

type testResponse struct {
	Results struct {
		Sunrise string `json:"sunrise"`
		Sunset  string `json:"sunset"`
	} `json:"results"`
	Status string `json:"status"`
}

func makeSunServer(t *testing.T, rise, set string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var resp testResponse
		resp.Status = "OK"
		resp.Results.Sunrise = rise
		resp.Results.Sunset = set
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("encode response: %v", err)
		}
	}))
}

func patchURL(srv *httptest.Server) func() {
	orig := sunapi.APIURL
	sunapi.APIURL = srv.URL
	return func() { sunapi.APIURL = orig }
}

func TestFetchTimes_parsesRFC3339(t *testing.T) {
	riseStr := "2026-08-02T06:00:00+00:00"
	setStr := "2026-08-02T20:00:00+00:00"
	srv := makeSunServer(t, riseStr, setStr)
	defer srv.Close()
	restore := patchURL(srv)
	defer restore()

	rise, set, err := sunapi.FetchTimes(37.7, -122.4, time.Now())
	if err != nil {
		t.Fatalf("fetchTimes: %v", err)
	}
	wantRise, _ := time.Parse(time.RFC3339, riseStr)
	wantSet, _ := time.Parse(time.RFC3339, setStr)
	if !rise.Equal(wantRise) {
		t.Errorf("rise = %v, want %v", rise, wantRise)
	}
	if !set.Equal(wantSet) {
		t.Errorf("set = %v, want %v", set, wantSet)
	}
}

func TestFetchTimes_apiError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var resp testResponse
		resp.Status = "ERROR"
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()
	restore := patchURL(srv)
	defer restore()

	_, _, err := sunapi.FetchTimes(37.7, -122.4, time.Now())
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestFetch_rowLengths(t *testing.T) {
	now := time.Now()
	riseStr := now.Add(1 * time.Hour).UTC().Format(time.RFC3339)
	setStr := now.Add(10 * time.Hour).UTC().Format(time.RFC3339)
	srv := makeSunServer(t, riseStr, setStr)
	defer srv.Close()
	restore := patchURL(srv)
	defer restore()

	lines, err := Fetch(37.7, -122.4)
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	for i, row := range lines {
		if countTiles(row) != 15 {
			t.Errorf("row%d tile count = %d, want 15: %q", i+1, countTiles(row), row)
		}
	}
}

func TestFetch_beforeRise_labelIsSunrise(t *testing.T) {
	now := time.Now()
	riseStr := now.Add(2 * time.Hour).UTC().Format(time.RFC3339)
	setStr := now.Add(14 * time.Hour).UTC().Format(time.RFC3339)
	srv := makeSunServer(t, riseStr, setStr)
	defer srv.Close()
	restore := patchURL(srv)
	defer restore()

	lines, err := Fetch(37.7, -122.4)
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if !strings.Contains(lines[1], "SUNRISE") {
		t.Errorf("row2 = %q, expected SUNRISE", lines[1])
	}
}

func TestFetch_betweenRiseAndSet_labelIsSunset(t *testing.T) {
	now := time.Now()
	riseStr := now.Add(-2 * time.Hour).UTC().Format(time.RFC3339)
	setStr := now.Add(4 * time.Hour).UTC().Format(time.RFC3339)
	srv := makeSunServer(t, riseStr, setStr)
	defer srv.Close()
	restore := patchURL(srv)
	defer restore()

	lines, err := Fetch(37.7, -122.4)
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if !strings.Contains(lines[1], "SUNSET") {
		t.Errorf("row2 = %q, expected SUNSET", lines[1])
	}
}

func TestFetch_afterSet_fetchesTomorrowSunrise(t *testing.T) {
	now := time.Now()
	riseStr := now.Add(-10 * time.Hour).UTC().Format(time.RFC3339)
	setStr := now.Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	tomorrowRiseStr := now.Add(20 * time.Hour).UTC().Format(time.RFC3339)
	tomorrowSetStr := now.Add(34 * time.Hour).UTC().Format(time.RFC3339)

	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		callCount++
		var resp testResponse
		resp.Status = "OK"
		if callCount == 1 {
			resp.Results.Sunrise = riseStr
			resp.Results.Sunset = setStr
		} else {
			resp.Results.Sunrise = tomorrowRiseStr
			resp.Results.Sunset = tomorrowSetStr
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()
	restore := patchURL(srv)
	defer restore()

	lines, err := Fetch(37.7, -122.4)
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if callCount != 2 {
		t.Errorf("expected 2 API calls (today + tomorrow), got %d", callCount)
	}
	if !strings.Contains(lines[1], "SUNRISE") {
		t.Errorf("row2 = %q, expected SUNRISE after past sunset", lines[1])
	}
}
