package onthisday

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/nicoleyson/vestaboard-note/internal/layout"
)

func TestPickEvent_prefersCleanFit(t *testing.T) {
	short := wikiEvent{Year: 1900, Text: "Hi"}
	exact := wikiEvent{Year: 1961, Text: "U.S. DEFENSE DEPT CREATES DARPA RESEARCH"}
	events := []wikiEvent{short, exact}

	got := pickEvent(events)
	wrapped := layout.Wrap(layout.StripEmoji(got.Text), layout.Cols)
	if len(wrapped) == 0 {
		t.Fatal("pickEvent returned event with no wrapped lines")
	}
	if len([]rune(wrapped[0])) != layout.Cols {
		t.Errorf("preferred event first line len = %d, want %d: %q", len([]rune(wrapped[0])), layout.Cols, wrapped[0])
	}
}

func TestPickEvent_fallsBackToAny(t *testing.T) {
	events := []wikiEvent{
		{Year: 1900, Text: "Hi"},
		{Year: 1901, Text: "OK"},
	}
	got := pickEvent(events)
	if got.Year != 1900 && got.Year != 1901 {
		t.Errorf("pickEvent fallback returned unexpected year %d", got.Year)
	}
}

func TestEndsWithTilde(t *testing.T) {
	if !endsWithTilde("hello~") {
		t.Error("expected true for string ending with tilde")
	}
	if endsWithTilde("hello") {
		t.Error("expected false for string not ending with tilde")
	}
	if endsWithTilde("") {
		t.Error("expected false for empty string")
	}
}

func TestFetch_success(t *testing.T) {
	body, _ := json.Marshal(map[string]interface{}{
		"events": []map[string]interface{}{
			{"year": 1969, "text": "Apollo 11 lands on the Moon"},
		},
	})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}))
	defer srv.Close()

	old := apiURL
	// apiURL format has %d/%d, replace the whole base so fmt.Sprintf still works
	apiURL = srv.URL + "/%d/%d"
	defer func() { apiURL = old }()

	lines, err := Fetch(time.Now())
	if err != nil {
		t.Fatalf("Fetch error: %v", err)
	}
	for i, l := range lines {
		if len([]rune(l)) != 15 {
			t.Errorf("row %d len = %d, want 15: %q", i, len([]rune(l)), l)
		}
	}
	if !strings.Contains(lines[0], "ON THIS DAY") {
		t.Errorf("row 0 should contain ON THIS DAY, got %q", lines[0])
	}
}

func TestFetch_emptyEvents(t *testing.T) {
	body, _ := json.Marshal(map[string]interface{}{
		"events": []map[string]interface{}{},
	})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}))
	defer srv.Close()

	old := apiURL
	apiURL = srv.URL + "/%d/%d"
	defer func() { apiURL = old }()

	_, err := Fetch(time.Now())
	if err == nil {
		t.Error("expected error for empty events response")
	}
}

func TestFetch_httpError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer srv.Close()

	old := apiURL
	apiURL = srv.URL + "/%d/%d"
	defer func() { apiURL = old }()

	_, err := Fetch(time.Now())
	if err == nil {
		t.Error("expected error for HTTP 404")
	}
}

func TestFetch_live(t *testing.T) {
	resp, err := http.Get("https://en.wikipedia.org")
	if err != nil || resp.StatusCode != 200 {
		t.Skip("no network access, skipping live test")
	}
	resp.Body.Close()

	lines, err := Fetch(time.Now())
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}
	for i, l := range lines {
		if len([]rune(l)) != 15 {
			t.Errorf("row %d: got len %d, want 15: %q", i+1, len([]rune(l)), l)
		}
	}
}

func TestFetch_row3IsCentered(t *testing.T) {
	body, _ := json.Marshal(map[string]interface{}{
		"events": []map[string]interface{}{
			{"year": 1776, "text": "Hi"},
		},
	})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}))
	defer srv.Close()

	old := apiURL
	apiURL = srv.URL + "/%d/%d"
	defer func() { apiURL = old }()

	lines, err := Fetch(time.Now())
	if err != nil {
		t.Fatalf("Fetch error: %v", err)
	}
	row3 := lines[2]
	if len([]rune(row3)) != 15 {
		t.Fatalf("row3 len = %d, want 15", len([]rune(row3)))
	}
	if !strings.HasPrefix(row3, " ") {
		t.Errorf("row3 should be centered (leading space for short text), got %q", row3)
	}
	if !strings.HasSuffix(strings.TrimRight(row3, " "), "HI") {
		t.Errorf("row3 should contain HI centered, got %q", row3)
	}
}

