package rain

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClassify(t *testing.T) {
	tests := []struct {
		mm        float64
		wantLabel string
		wantColor int
	}{
		{0, "NONE", 65},
		{0.01, "LIGHT", 67},
		{0.5, "LIGHT", 67},
		{1.0, "MODERATE", 67},
		{3.9, "MODERATE", 67},
		{4.0, "HEAVY", 67},
		{10.0, "HEAVY", 67},
	}
	for _, tt := range tests {
		label, color := classify(tt.mm)
		if label != tt.wantLabel || color != tt.wantColor {
			t.Errorf("classify(mm=%.2f) = (%q, %d), want (%q, %d)",
				tt.mm, label, color, tt.wantLabel, tt.wantColor)
		}
	}
}

func TestClassifyTrivial(t *testing.T) {
	label, _ := classify(0)
	if label != "NONE" {
		t.Errorf("classify(0) = %q, want NONE", label)
	}
	label, _ = classify(0.1)
	if label == "NONE" {
		t.Error("classify(0.1) = NONE, want non-trivial")
	}
}

func TestFetch_mmPath(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"current": map[string]interface{}{
				"precipitation": 2.5,
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
		t.Error("Fetch: expected non-trivial for 2.5mm rain")
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

func TestFetch_nonRaining(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"current": map[string]interface{}{
				"precipitation": 0.0,
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
		t.Error("Fetch: expected trivial for 0mm rain")
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
