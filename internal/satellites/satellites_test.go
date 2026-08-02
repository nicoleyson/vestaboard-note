package satellites

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPickBest_empty(t *testing.T) {
	if pickBest(nil) != nil {
		t.Error("expected nil for empty slice")
	}
}

func TestPickBest_skipsUninteresting(t *testing.T) {
	sats := []satellite{
		{Name: "STARLINK-1", Category: "STARLINK", ElevationDeg: 80},
		{Name: "DEBRIS-X", Category: "DEBRIS", ElevationDeg: 60},
	}
	if pickBest(sats) != nil {
		t.Error("expected nil when only STARLINK/DEBRIS present")
	}
}

func TestPickBest_picksHighestElevation(t *testing.T) {
	sats := []satellite{
		{Name: "GPS-1", Category: "GPS", ElevationDeg: 40},
		{Name: "GPS-2", Category: "GPS", ElevationDeg: 75},
		{Name: "NOAA-19", Category: "OTHER", ElevationDeg: 50},
	}
	got := pickBest(sats)
	if got == nil {
		t.Fatal("expected non-nil result")
	}
	if got.Name != "GPS-2" {
		t.Errorf("want GPS-2, got %q", got.Name)
	}
}

func TestPickBest_prefersInterestingOverHigherDebris(t *testing.T) {
	sats := []satellite{
		{Name: "STARLINK-999", Category: "STARLINK", ElevationDeg: 89},
		{Name: "ISS", Category: "OTHER", ElevationDeg: 45},
	}
	got := pickBest(sats)
	if got == nil || got.Name != "ISS" {
		t.Errorf("want ISS, got %v", got)
	}
}

func TestPickBest_stationsCategory(t *testing.T) {
	sats := []satellite{
		{Name: "ISS (ZARYA)", Category: "STATIONS", ElevationDeg: 45},
		{Name: "DEBRIS-1", Category: "DEBRIS", ElevationDeg: 80},
	}
	got := pickBest(sats)
	if got == nil || got.Name != "ISS (ZARYA)" {
		t.Errorf("want ISS (ZARYA), got %v", got)
	}
}

func TestCleanName_stripsParenSuffix(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"ISS (ZARYA)", "ISS"},
		{"NOAA 19 (GOES)", "NOAA 19"},
		{"GPS IIF-3", "GPS IIF-3"},
		{"", ""},
	}
	for _, tc := range cases {
		got := cleanName(tc.input)
		if got != tc.want {
			t.Errorf("cleanName(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestCleanName_truncatesLong(t *testing.T) {
	long := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	got := cleanName(long)
	if len(got) > 15 {
		t.Errorf("cleanName did not truncate: len=%d, got=%q", len(got), got)
	}
}

func TestColorForCategory(t *testing.T) {
	cases := []struct {
		cat  string
		want int
	}{
		{"GPS", 66},
		{"IRIDIUM", 67},
		{"STATIONS", 68},
		{"OTHER", 68},
		{"WEATHER", 68},
		{"", 68},
	}
	for _, tc := range cases {
		got := colorForCategory(tc.cat)
		if got != tc.want {
			t.Errorf("colorForCategory(%q) = %d, want %d", tc.cat, got, tc.want)
		}
	}
}

func TestFetch_returnsISS(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(overhead{
			Count: 1,
			Satellites: []satellite{
				{Name: "ISS (ZARYA)", Category: "STATIONS", ElevationDeg: 45, Direction: "NW"},
			},
		})
	}))
	defer srv.Close()

	old := apiURL
	apiURL = srv.URL
	defer func() { apiURL = old }()

	lines, trivial, err := Fetch(37.7, -122.4)
	if err != nil {
		t.Fatalf("Fetch error: %v", err)
	}
	if trivial {
		t.Error("expected non-trivial when ISS is overhead")
	}
	if lines[1] != "      ISS      " {
		t.Errorf("row1 want centered ISS, got %q", lines[1])
	}
}

func TestFetch_noInterestingSatellites_returnsTrivial(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(overhead{
			Count: 2,
			Satellites: []satellite{
				{Name: "STARLINK-1234", Category: "STARLINK", ElevationDeg: 80},
				{Name: "DEBRIS-X", Category: "DEBRIS", ElevationDeg: 60},
			},
		})
	}))
	defer srv.Close()

	old := apiURL
	apiURL = srv.URL
	defer func() { apiURL = old }()

	lines, trivial, err := Fetch(37.7, -122.4)
	if err != nil {
		t.Fatalf("Fetch error: %v", err)
	}
	if !trivial {
		t.Error("expected trivial when only Starlink/Debris present")
	}
	_ = lines
}

func TestFetch_httpError_returnsTrivial(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	old := apiURL
	apiURL = srv.URL
	defer func() { apiURL = old }()

	_, trivial, err := Fetch(37.7, -122.4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !trivial {
		t.Error("expected trivial on HTTP error (satellites treats errors as trivial)")
	}
}

func TestFetch_emptyResponse_returnsTrivial(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(overhead{Count: 0, Satellites: nil})
	}))
	defer srv.Close()

	old := apiURL
	apiURL = srv.URL
	defer func() { apiURL = old }()

	_, trivial, err := Fetch(37.7, -122.4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !trivial {
		t.Error("expected trivial for empty satellite list")
	}
}

func TestFetch_rowLengths(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(overhead{
			Count: 1,
			Satellites: []satellite{
				{Name: "GPS IIR-M (1)", Category: "GPS", ElevationDeg: 72, Direction: "SE"},
			},
		})
	}))
	defer srv.Close()

	old := apiURL
	apiURL = srv.URL
	defer func() { apiURL = old }()

	lines, trivial, err := Fetch(37.7, -122.4)
	if err != nil {
		t.Fatalf("Fetch error: %v", err)
	}
	if trivial {
		t.Error("expected non-trivial for GPS satellite")
	}
	for i := 1; i < 3; i++ {
		if len([]rune(lines[i])) != 15 {
			t.Errorf("row %d len = %d, want 15: %q", i, len([]rune(lines[i])), lines[i])
		}
	}
}
