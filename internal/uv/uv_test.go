package uv

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClassify(t *testing.T) {
	tests := []struct {
		idx       float64
		wantLabel string
		wantColor int
	}{
		{0, "LOW", 66},
		{2.9, "LOW", 66},
		{3, "MODERATE", 65},
		{5.9, "MODERATE", 65},
		{6, "HIGH", 64},
		{7.9, "HIGH", 64},
		{8, "VERY HIGH", 63},
		{10.9, "VERY HIGH", 63},
		{11, "EXTREME", 68},
		{15, "EXTREME", 68},
	}
	for _, tt := range tests {
		label, color := classify(tt.idx)
		if label != tt.wantLabel || color != tt.wantColor {
			t.Errorf("classify(%.1f) = (%q, %d), want (%q, %d)",
				tt.idx, label, color, tt.wantLabel, tt.wantColor)
		}
	}
}

func TestFetch_success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"current": map[string]interface{}{
				"uv_index": 9.0,
			},
		})
	}))
	defer srv.Close()

	old := apiURL
	apiURL = srv.URL
	defer func() { apiURL = old }()

	lines, err := Fetch(37.7, -122.4)
	if err != nil {
		t.Fatalf("Fetch error: %v", err)
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

func TestFetch_httpError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "service unavailable", http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	old := apiURL
	apiURL = srv.URL
	defer func() { apiURL = old }()

	_, err := Fetch(37.7, -122.4)
	if err == nil {
		t.Error("expected error for HTTP 503")
	}
}

func TestFetch_jsonError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("not json {{{"))
	}))
	defer srv.Close()

	old := apiURL
	apiURL = srv.URL
	defer func() { apiURL = old }()

	_, err := Fetch(37.7, -122.4)
	if err == nil {
		t.Error("expected error for invalid JSON response")
	}
}

