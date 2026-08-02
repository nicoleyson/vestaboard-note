package calendar

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"
)

var tileRe = regexp.MustCompile(`\{\d+\}|.`)

func countTiles(s string) int {
	return len(tileRe.FindAllString(s, -1))
}

func makeICal(entries []struct{ uid, summary, dtstart string }) string {
	var sb strings.Builder
	sb.WriteString("BEGIN:VCALENDAR\r\nVERSION:2.0\r\nPRODID:-//test//test//EN\r\n")
	for _, e := range entries {
		t, err := time.Parse(time.RFC3339, e.dtstart)
		if err != nil {
			panic(fmt.Sprintf("bad dtstart %q: %v", e.dtstart, err))
		}
		sb.WriteString("BEGIN:VEVENT\r\n")
		fmt.Fprintf(&sb, "UID:%s\r\n", e.uid)
		fmt.Fprintf(&sb, "SUMMARY:%s\r\n", e.summary)
		fmt.Fprintf(&sb, "DTSTART:%s\r\n", t.UTC().Format("20060102T150405Z"))
		sb.WriteString("END:VEVENT\r\n")
	}
	sb.WriteString("END:VCALENDAR\r\n")
	return sb.String()
}

func serveICal(t *testing.T, body string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/calendar")
		fmt.Fprint(w, body)
	}))
}

func TestFetch_noEvents_returnsNoEventsMessage(t *testing.T) {
	past := time.Now().Add(-48 * time.Hour).UTC().Format(time.RFC3339)
	ical := makeICal([]struct{ uid, summary, dtstart string }{
		{"uid1", "Old Meeting", past},
	})
	srv := serveICal(t, ical)
	defer srv.Close()

	lines, err := Fetch([]string{srv.URL})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	combined := strings.Join(lines[:], " ")
	if !strings.Contains(combined, "NO EVENTS") {
		t.Errorf("expected NO EVENTS message, got: %v", lines)
	}
}

func TestFetch_upcomingEvent_rowLengths(t *testing.T) {
	future := time.Now().Add(2 * time.Hour).UTC().Format(time.RFC3339)
	ical := makeICal([]struct{ uid, summary, dtstart string }{
		{"uid1", "Team Standup", future},
	})
	srv := serveICal(t, ical)
	defer srv.Close()

	lines, err := Fetch([]string{srv.URL})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, row := range lines {
		if countTiles(row) != 15 {
			t.Errorf("row%d tile count = %d, want 15: %q", i+1, countTiles(row), row)
		}
	}
}

func TestFetch_upcomingEvent_row3ContainsSummary(t *testing.T) {
	future := time.Now().Add(2 * time.Hour).UTC().Format(time.RFC3339)
	ical := makeICal([]struct{ uid, summary, dtstart string }{
		{"uid1", "BOARD MEETING", future},
	})
	srv := serveICal(t, ical)
	defer srv.Close()

	lines, err := Fetch([]string{srv.URL})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(lines[2], "BOARD MEETING") {
		t.Errorf("row3 = %q, expected BOARD MEETING", lines[2])
	}
}

func TestFetch_allURLsFail_returnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	srv.Close()

	_, err := Fetch([]string{srv.URL})
	if err == nil {
		t.Fatal("expected error when all URLs fail, got nil")
	}
}

func TestFetch_deduplicatesUID(t *testing.T) {
	future := time.Now().Add(2 * time.Hour).UTC().Format(time.RFC3339)
	ical := makeICal([]struct{ uid, summary, dtstart string }{
		{"uid-dupe", "DUPE EVENT", future},
	})
	srv1 := serveICal(t, ical)
	defer srv1.Close()
	srv2 := serveICal(t, ical)
	defer srv2.Close()

	lines, err := Fetch([]string{srv1.URL, srv2.URL})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(lines[2], "DUPE EVENT") {
		t.Errorf("row3 = %q, expected deduped event to still be shown", lines[2])
	}
}

func TestFetch_picksEarliestEvent(t *testing.T) {
	soon := time.Now().Add(1 * time.Hour).UTC().Format(time.RFC3339)
	later := time.Now().Add(5 * time.Hour).UTC().Format(time.RFC3339)
	ical := makeICal([]struct{ uid, summary, dtstart string }{
		{"uid-later", "FAR EVENT", later},
		{"uid-soon", "NEAR EVENT", soon},
	})
	srv := serveICal(t, ical)
	defer srv.Close()

	lines, err := Fetch([]string{srv.URL})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(lines[2], "NEAR EVENT") {
		t.Errorf("row3 = %q, expected nearest event NEAR EVENT", lines[2])
	}
}

func TestFetch_allDayEvent_todayIncluded(t *testing.T) {
	localNow := time.Now().In(time.Local)
	startOfLocalDay := time.Date(localNow.Year(), localNow.Month(), localNow.Day(), 0, 0, 0, 0, time.Local).UTC()
	ical := makeICal([]struct{ uid, summary, dtstart string }{
		{"uid-allday", "ALL DAY EVENT", startOfLocalDay.Format(time.RFC3339)},
	})
	srv := serveICal(t, ical)
	defer srv.Close()

	lines, err := Fetch([]string{srv.URL})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(lines[2], "ALL DAY EVENT") {
		t.Errorf("row3 = %q, expected all-day event to be included", lines[2])
	}
}

func TestFetch_summaryTruncatedTo15(t *testing.T) {
	future := time.Now().Add(2 * time.Hour).UTC().Format(time.RFC3339)
	longTitle := "THIS IS A VERY LONG CALENDAR EVENT TITLE THAT EXCEEDS FIFTEEN"
	ical := makeICal([]struct{ uid, summary, dtstart string }{
		{"uid1", longTitle, future},
	})
	srv := serveICal(t, ical)
	defer srv.Close()

	lines, err := Fetch([]string{srv.URL})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if countTiles(lines[2]) != 15 {
		t.Errorf("row3 tile count = %d, want 15 (truncated): %q", countTiles(lines[2]), lines[2])
	}
}
