package dewpoint

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClassify(t *testing.T) {
	tests := []struct {
		dpF       float64
		wantLabel string
		wantColor int
	}{
		{0, "DRY", 66},
		{55, "DRY", 66},
		{55.1, "COMFY", 65},
		{60, "COMFY", 65},
		{60.1, "HUMID", 64},
		{65, "HUMID", 64},
		{65.1, "MUGGY", 63},
		{70, "MUGGY", 63},
		{70.1, "OPPRESSIVE", 68},
		{85, "OPPRESSIVE", 68},
	}
	for _, tt := range tests {
		label, color := classify(tt.dpF)
		if label != tt.wantLabel || color != tt.wantColor {
			t.Errorf("classify(%.1f) = (%q, %d), want (%q, %d)",
				tt.dpF, label, color, tt.wantLabel, tt.wantColor)
		}
	}
}

func TestClassifyTrivial(t *testing.T) {
	label, _ := classify(55)
	if label != "DRY" {
		t.Errorf("classify(55) = %q, want DRY", label)
	}
	label, _ = classify(56)
	if label == "DRY" {
		t.Error("classify(56) = DRY, want non-trivial")
	}
}

func TestFetch_muggy(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"current": map[string]interface{}{
				"dew_point_2m":         20.0,
				"relative_humidity_2m": 85,
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
		t.Error("Fetch: expected non-trivial for muggy dew point")
	}
	if lines[0] == "" {
		t.Error("row 0 (color row) should not be empty")
	}
	for i := 1; i < 3; i++ {
		if len([]rune(lines[i])) != 15 {
			t.Errorf("row %d len = %d, want 15: %q", i, len([]rune(lines[i])), lines[i])
		}
	}
}

func TestFetch_dry(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"current": map[string]interface{}{
				"dew_point_2m":         5.0,
				"relative_humidity_2m": 30,
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
		t.Error("Fetch: expected trivial for dry dew point")
	}
	if !strings.Contains(lines[2], "DRY") {
		t.Errorf("row 2 should contain DRY, got %q", lines[2])
	}
}

func TestFetch_row2ContainsDPandRH(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"current": map[string]interface{}{
				"dew_point_2m":         18.0,
				"relative_humidity_2m": 72,
			},
		})
	}))
	defer srv.Close()

	old := apiURL
	apiURL = srv.URL
	defer func() { apiURL = old }()

	lines, _, err := Fetch(37.7, -122.4)
	if err != nil {
		t.Fatalf("Fetch error: %v", err)
	}
	if !strings.Contains(lines[1], "DP") {
		t.Errorf("row 1 should contain DP, got %q", lines[1])
	}
	if !strings.Contains(lines[1], "RH") {
		t.Errorf("row 1 should contain RH, got %q", lines[1])
	}
}

func TestFetch_httpError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	old := apiURL
	apiURL = srv.URL
	defer func() { apiURL = old }()

	_, _, err := Fetch(37.7, -122.4)
	if err == nil {
		t.Error("Fetch: expected error for HTTP 503")
	}
}

func TestFetch_jsonDecodeError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{invalid json`))
	}))
	defer srv.Close()

	old := apiURL
	apiURL = srv.URL
	defer func() { apiURL = old }()

	_, _, err := Fetch(37.7, -122.4)
	if err == nil {
		t.Error("Fetch: expected error for malformed JSON response")
	}
}
