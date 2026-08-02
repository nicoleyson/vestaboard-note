package sunscene

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/nicoleyson/vestaboard-note/internal/sunapi"
)

// tileCount counts rendered tiles: each {N} is one tile, each plain rune is one tile.
var tileRe = regexp.MustCompile(`\{[0-9]+\}|.`)

func tileCount(s string) int {
	return len(tileRe.FindAllString(s, -1))
}

func TestSunRow(t *testing.T) {
	cases := []struct {
		progress float64
		wantRow  int
	}{
		{0.0, 2},   // deep dawn → bottom
		{0.10, 2},  // early morning → bottom
		{0.20, 1},  // morning → mid
		{0.50, 0},  // solar noon → top
		{0.80, 1},  // afternoon → mid
		{0.90, 2},  // dusk → bottom
		{1.0, 2},   // deep dusk → bottom
	}
	for _, tc := range cases {
		got := sunRow(tc.progress)
		if got != tc.wantRow {
			t.Errorf("sunRow(%v) = %d, want %d", tc.progress, got, tc.wantRow)
		}
	}
}

func TestRowColor_midday(t *testing.T) {
	// At solar noon (p=0.5), row 0 and 1 should be blue, row 2 orange
	if c := rowColor(0, 0, 0.5); c != blue {
		t.Errorf("rowColor(r=0, sr=0, p=0.5) = %d, want blue(%d)", c, blue)
	}
	if c := rowColor(1, 0, 0.5); c != blue {
		t.Errorf("rowColor(r=1, sr=0, p=0.5) = %d, want blue(%d)", c, blue)
	}
	if c := rowColor(2, 0, 0.5); c != orange {
		t.Errorf("rowColor(r=2, sr=0, p=0.5) = %d, want orange(%d)", c, orange)
	}
}

func TestRowColor_goldenHour(t *testing.T) {
	// At golden hour (p=0.87), row 0 should be violet, row 2 should be red
	if c := rowColor(0, 2, 0.87); c != violet {
		t.Errorf("rowColor(r=0, sr=2, p=0.87) = %d, want violet(%d)", c, violet)
	}
	if c := rowColor(2, 2, 0.87); c != red {
		t.Errorf("rowColor(r=2, sr=2, p=0.87) = %d, want red(%d)", c, red)
	}
}

func TestRowColor_deepDusk_sunRow(t *testing.T) {
	// At deep dusk (p=0.98), sun row → red, other rows → black
	sr := sunRow(0.98) // = 2
	if c := rowColor(sr, sr, 0.98); c != red {
		t.Errorf("rowColor(sunRow, sunRow, 0.98) = %d, want red(%d)", c, red)
	}
	if c := rowColor(0, sr, 0.98); c != black {
		t.Errorf("rowColor(0, sunRow, 0.98) = %d, want black(%d)", c, black)
	}
}

func TestNightRowColor(t *testing.T) {
	// Moon row always violet
	if c := nightRowColor(1, 1, 0.5); c != violet {
		t.Errorf("moon row: got %d, want violet(%d)", c, violet)
	}
	// Row below moon always violet
	if c := nightRowColor(2, 1, 0.5); c != violet {
		t.Errorf("row below moon: got %d, want violet(%d)", c, violet)
	}
	// Row above moon at deep night → black
	if c := nightRowColor(0, 1, 0.5); c != black {
		t.Errorf("row above moon at deep night: got %d, want black(%d)", c, black)
	}
	// Row above moon near dusk → violet
	if c := nightRowColor(0, 1, 0.10); c != violet {
		t.Errorf("row above moon near dusk (p=0.10): got %d, want violet(%d)", c, violet)
	}
}

func TestRenderDay_tileCounts(t *testing.T) {
	for _, p := range []float64{0.0, 0.1, 0.25, 0.5, 0.75, 0.9, 1.0} {
		lines := renderDay(p)
		for i, line := range lines {
			if n := tileCount(line); n != cols {
				t.Errorf("renderDay(%.2f) row %d: %d tiles, want %d", p, i, n, cols)
			}
		}
	}
}

func TestRenderNight_tileCounts(t *testing.T) {
	for _, p := range []float64{0.0, 0.1, 0.5, 0.9, 1.0} {
		lines := renderNight(p)
		for i, line := range lines {
			if n := tileCount(line); n != cols {
				t.Errorf("renderNight(%.2f) row %d: %d tiles, want %d", p, i, n, cols)
			}
		}
	}
}

func TestRenderDay_sunDiscPresent(t *testing.T) {
	// At solar noon the sun disc (yellow center tile) should appear in row 0
	lines := renderDay(0.5)
	sunRowIdx := sunRow(0.5) // = 0
	if !strings.Contains(lines[sunRowIdx], "{65}") { // yellow
		t.Errorf("renderDay(0.5): expected yellow tile in sun row %d, got %q", sunRowIdx, lines[sunRowIdx])
	}
}

func TestRenderNight_moonDiscPresent(t *testing.T) {
	// At mid-night the moon disc (white) should appear in the moon row
	lines := renderNight(0.5)
	moonRowIdx := sunRow(0.5) // same logic, = 0
	if !strings.Contains(lines[moonRowIdx], "{69}") { // white
		t.Errorf("renderNight(0.5): expected white tile in moon row %d, got %q", moonRowIdx, lines[moonRowIdx])
	}
}

func makeSunServer(t *testing.T, rise, set string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "OK",
			"results": map[string]string{
				"sunrise": rise,
				"sunset":  set,
			},
		})
	}))
}

func patchURL(srv *httptest.Server) func() {
	orig := sunapi.APIURL
	sunapi.APIURL = srv.URL
	return func() { sunapi.APIURL = orig }
}

func TestFetch_duringDay_returnsDayScene(t *testing.T) {
	rise := "2026-08-02T06:00:00+00:00"
	set := "2026-08-02T22:00:00+00:00"
	srv := makeSunServer(t, rise, set)
	defer srv.Close()
	restore := patchURL(srv)
	defer restore()

	lines, err := Fetch(37.7, -122.4)
	if err != nil {
		t.Fatalf("Fetch error: %v", err)
	}
	for i, line := range lines {
		if n := tileCount(line); n != cols {
			t.Errorf("row %d: %d tiles, want %d", i, n, cols)
		}
	}
}

func TestFetch_apiError_returnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}))
	defer srv.Close()
	restore := patchURL(srv)
	defer restore()

	_, err := Fetch(37.7, -122.4)
	if err == nil {
		t.Error("expected error when sunapi returns HTTP 500")
	}
}
