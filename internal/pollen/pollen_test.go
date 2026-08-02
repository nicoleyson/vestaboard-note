package pollen

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClassify(t *testing.T) {
	tests := []struct {
		val       float64
		wantLabel string
		wantColor int
	}{
		{0, "LOW", 66},
		{9.9, "LOW", 66},
		{10, "MODERATE", 65},
		{49.9, "MODERATE", 65},
		{50, "HIGH", 64},
		{199.9, "HIGH", 64},
		{200, "V.HIGH", 63},
		{500, "V.HIGH", 63},
	}
	for _, tt := range tests {
		label, color := classify(tt.val)
		if label != tt.wantLabel || color != tt.wantColor {
			t.Errorf("classify(%.1f) = (%q, %d), want (%q, %d)",
				tt.val, label, color, tt.wantLabel, tt.wantColor)
		}
	}
}

func TestDominantType(t *testing.T) {
	tests := []struct {
		grass, tree, weed float64
		wantName          string
	}{
		{10, 5, 3, "GRASS"},
		{5, 10, 3, "TREE"},
		{3, 5, 10, "WEED"},
		{10, 10, 5, "GRASS"},
		{5, 10, 10, "TREE"},
		{0, 0, 0, "GRASS"},
	}
	for _, tt := range tests {
		name, _ := dominantType(tt.grass, tt.tree, tt.weed)
		if name != tt.wantName {
			t.Errorf("dominantType(%.0f, %.0f, %.0f) = %q, want %q",
				tt.grass, tt.tree, tt.weed, name, tt.wantName)
		}
	}
}

func TestClassifyTrivial(t *testing.T) {
	trivialCases := []float64{0, 5, 9.9}
	for _, v := range trivialCases {
		label, _ := classify(v)
		if label != "LOW" {
			t.Errorf("classify(%.1f) = %q, expected LOW (trivial)", v, label)
		}
	}
	nonTrivialCases := []float64{10, 50, 200}
	for _, v := range nonTrivialCases {
		label, _ := classify(v)
		if label == "LOW" {
			t.Errorf("classify(%.1f) = LOW, expected non-trivial", v)
		}
	}
}

func TestFetch_success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"current": map[string]interface{}{
				"grass_pollen": 75.0,
				"tree_pollen":  20.0,
				"weed_pollen":  5.0,
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
		t.Error("expected non-trivial for grass pollen 75")
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

func TestFetch_trivialLow(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"current": map[string]interface{}{
				"grass_pollen": 1.0,
				"tree_pollen":  0.5,
				"weed_pollen":  0.0,
			},
		})
	}))
	defer srv.Close()

	old := apiURL
	apiURL = srv.URL
	defer func() { apiURL = old }()

	_, trivial, err := Fetch(37.7, -122.4)
	if err != nil {
		t.Fatalf("Fetch error: %v", err)
	}
	if !trivial {
		t.Error("expected trivial for low pollen")
	}
}

func TestFetch_httpError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	old := apiURL
	apiURL = srv.URL
	defer func() { apiURL = old }()

	_, _, err := Fetch(37.7, -122.4)
	if err == nil {
		t.Error("expected error for HTTP 503")
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
		t.Error("expected error for malformed JSON response")
	}
}

