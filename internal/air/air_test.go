package air

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClassify(t *testing.T) {
	tests := []struct {
		aqi       int
		wantLabel string
		wantColor int
	}{
		{0, "GOOD", 66},
		{50, "GOOD", 66},
		{51, "MODERATE", 65},
		{100, "MODERATE", 65},
		{101, "UNHEALTHY+", 64},
		{150, "UNHEALTHY+", 64},
		{151, "UNHEALTHY", 63},
		{200, "UNHEALTHY", 63},
		{201, "VERY UNHLTHY", 68},
		{300, "VERY UNHLTHY", 68},
		{301, "HAZARDOUS", 70},
		{500, "HAZARDOUS", 70},
	}
	for _, tt := range tests {
		label, color := classify(tt.aqi)
		if label != tt.wantLabel || color != tt.wantColor {
			t.Errorf("classify(%d) = (%q, %d), want (%q, %d)",
				tt.aqi, label, color, tt.wantLabel, tt.wantColor)
		}
	}
}

func TestClassifyTrivial(t *testing.T) {
	trivialCases := []int{0, 25, 50}
	for _, v := range trivialCases {
		label, _ := classify(v)
		if label != "GOOD" {
			t.Errorf("classify(%d) = %q, expected GOOD (trivial)", v, label)
		}
	}
	nonTrivialCases := []int{51, 100, 200}
	for _, v := range nonTrivialCases {
		label, _ := classify(v)
		if label == "GOOD" {
			t.Errorf("classify(%d) = GOOD, expected non-trivial", v)
		}
	}
}

func TestFetch_success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"current": map[string]interface{}{
				"us_aqi": 75,
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
		t.Error("expected non-trivial for AQI 75")
	}
	// row 0 is a color tile row ({NN} × 15), rows 1–2 are text
	if lines[0] == "" {
		t.Error("row 0 (color row) should not be empty")
	}
	for i := 1; i < 3; i++ {
		if len([]rune(lines[i])) != 15 {
			t.Errorf("row %d len = %d, want 15: %q", i, len([]rune(lines[i])), lines[i])
		}
	}
}

func TestFetch_trivialGood(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"current": map[string]interface{}{
				"us_aqi": 10,
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
		t.Error("expected trivial for AQI 10 (GOOD)")
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


