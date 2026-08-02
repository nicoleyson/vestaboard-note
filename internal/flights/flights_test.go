package flights

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func patchURL(srv *httptest.Server) func() {
	orig := apiURL
	apiURL = srv.URL
	return func() { apiURL = orig }
}

func serveStates(t *testing.T, states [][]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(apiResponse{States: states}); err != nil {
			t.Errorf("encode: %v", err)
		}
	}))
}

func makeState(callsign, origin string, altitude float64) []interface{} {
	s := make([]interface{}, 17)
	s[1] = callsign
	s[2] = origin
	s[7] = altitude
	return s
}

func TestParseFirst_returnsFirstValidCallsign(t *testing.T) {
	states := [][]interface{}{
		makeState("", "", 5000),
		makeState("UAL123", "US", 9000),
		makeState("DAL456", "US", 8000),
	}
	ac := parseFirst(states)
	if ac == nil {
		t.Fatal("expected aircraft, got nil")
	}
	if ac.callsign != "UAL123" {
		t.Errorf("callsign = %q, want UAL123", ac.callsign)
	}
}

func TestParseFirst_skipsEmptyCallsign(t *testing.T) {
	states := [][]interface{}{
		makeState("   ", "US", 5000),
		makeState("SWA789", "US", 6000),
	}
	ac := parseFirst(states)
	if ac == nil || ac.callsign != "SWA789" {
		t.Errorf("expected SWA789, got %v", ac)
	}
}

func TestParseFirst_nilWhenNoValidStates(t *testing.T) {
	states := [][]interface{}{
		makeState("", "", 0),
		{1, 2},
	}
	if parseFirst(states) != nil {
		t.Error("expected nil for all-empty states")
	}
}

func TestParseFirst_truncatesOriginToFourRunes(t *testing.T) {
	states := [][]interface{}{makeState("AAL1", "UNITED STATES", 5000)}
	ac := parseFirst(states)
	if ac == nil {
		t.Fatal("expected aircraft")
	}
	if len([]rune(ac.origin)) != 4 {
		t.Errorf("origin len = %d, want 4: %q", len([]rune(ac.origin)), ac.origin)
	}
}

func TestParseFirst_altitudeZeroWhenMissing(t *testing.T) {
	s := []interface{}{nil, "AAL1", "US", nil, nil, nil, nil, nil}
	ac := parseFirst([][]interface{}{s})
	if ac == nil {
		t.Fatal("expected aircraft")
	}
	if ac.altitude != 0 {
		t.Errorf("altitude = %v, want 0", ac.altitude)
	}
}

func TestFetch_clearSkies_rowLengths(t *testing.T) {
	srv := serveStates(t, nil)
	defer srv.Close()
	restore := patchURL(srv)
	defer restore()

	lines, trivial, err := Fetch(37.7, -122.4)
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if !trivial {
		t.Error("expected trivial=true for clear skies")
	}
	for i, row := range lines {
		if len([]rune(row)) != 15 {
			t.Errorf("row%d len = %d, want 15: %q", i+1, len([]rune(row)), row)
		}
	}
}

func TestFetch_withAircraft_rowLengths(t *testing.T) {
	states := [][]interface{}{makeState("UAL123", "US", 9000)}
	srv := serveStates(t, states)
	defer srv.Close()
	restore := patchURL(srv)
	defer restore()

	lines, trivial, err := Fetch(37.7, -122.4)
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if trivial {
		t.Error("expected trivial=false for aircraft overhead")
	}
	for i, row := range lines {
		if len([]rune(row)) != 15 {
			t.Errorf("row%d len = %d, want 15: %q", i+1, len([]rune(row)), row)
		}
	}
}

func TestFetch_withAircraft_row2ContainsCallsign(t *testing.T) {
	states := [][]interface{}{makeState("UAL123", "US", 9000)}
	srv := serveStates(t, states)
	defer srv.Close()
	restore := patchURL(srv)
	defer restore()

	lines, _, err := Fetch(37.7, -122.4)
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if !strings.Contains(lines[1], "UAL123") {
		t.Errorf("row2 = %q, expected UAL123", lines[1])
	}
}

func TestFetch_rateLimited_rowsNotEmpty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(429)
	}))
	defer srv.Close()
	restore := patchURL(srv)
	defer restore()

	lines, trivial, err := Fetch(37.7, -122.4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !trivial {
		t.Error("expected trivial=true for rate limit")
	}
	for i, row := range lines {
		if len([]rune(row)) != 15 {
			t.Errorf("row%d len = %d, want 15: %q", i+1, len([]rune(row)), row)
		}
	}
}

func TestFetch_serverError_returnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(503)
		fmt.Fprint(w, "<html>error</html>")
	}))
	defer srv.Close()
	restore := patchURL(srv)
	defer restore()

	_, _, err := Fetch(37.7, -122.4)
	if err == nil {
		t.Fatal("expected error for non-200 response")
	}
}
