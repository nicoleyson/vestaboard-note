package holiday

import (
	"strings"
	"testing"
)

var testHolidays = []nagerHoliday{
	{Date: "2026-01-01", Name: "New Year's Day"},
	{Date: "2026-07-04", Name: "Independence Day"},
	{Date: "2026-10-31", Name: "Halloween"},
	{Date: "2026-12-25", Name: "Christmas Day"},
}

func TestFindToday_match(t *testing.T) {
	got := findToday(testHolidays, "2026-07-04")
	if got != "Independence Day" {
		t.Errorf("got %q, want %q", got, "Independence Day")
	}
}

func TestFindToday_noMatch(t *testing.T) {
	got := findToday(testHolidays, "2026-06-15")
	if got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

func TestFindToday_empty(t *testing.T) {
	got := findToday(nil, "2026-01-01")
	if got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

func TestFindNext_basic(t *testing.T) {
	name, days := findNext(testHolidays, "2026-06-01")
	if name != "Independence Day" {
		t.Errorf("name: got %q, want %q", name, "Independence Day")
	}
	if days != 33 {
		t.Errorf("days: got %d, want 33", days)
	}
}

func TestFindNext_tomorrow(t *testing.T) {
	_, days := findNext(testHolidays, "2026-01-01")
	if days != 184 {
		t.Errorf("days from Jan 1 to Jul 4 2026: got %d, want 184", days)
	}
}

func TestFindNext_oneDayAway(t *testing.T) {
	_, days := findNext(testHolidays, "2026-07-03")
	if days != 1 {
		t.Errorf("got %d, want 1", days)
	}
}

func TestFindNext_noRemaining(t *testing.T) {
	name, days := findNext(testHolidays, "2026-12-31")
	if name != "" || days != 0 {
		t.Errorf("got (%q, %d), want (\"\", 0)", name, days)
	}
}

func TestArtStrip_knownHolidays(t *testing.T) {
	cases := []struct {
		input   string
		contain string
	}{
		{"Christmas Day", "{63}{66}"},
		{"New Year's Day", "{68}{65}"},
		{"Halloween", "{64}{70}"},
		{"Independence Day", "{63}{69}{67}"},
		{"Thanksgiving Day", "{64}{65}"},
		{"Easter Sunday", "{63}{66}{65}"},
		{"Valentine's Day", "{63}{62}"},
		{"St. Patrick's Day", "{66}{69}"},
		{"Diwali", "{65}{64}"},
		{"Hanukkah", "{67}{69}"},
		{"Lunar New Year", "{63}{65}"},
		{"Eid al-Fitr", "{66}{69}"},
		{"Labor Day", "{67}{69}"},
		{"Memorial Day", "{63}{69}{67}"},
		{"Veterans Day", "{63}{69}{67}"},
		{"Martin Luther King Jr. Day", "{67}{69}"},
		{"Presidents' Day", "{63}{69}{67}"},
		{"Guy Fawkes Night", "{64}{70}"},
		{"Bastille Day", "{67}{69}{63}"},
		{"Anzac Day", "{63}{69}"},
		{"Remembrance Day", "{63}{69}"},
		{"Armistice Day", "{63}{69}"},
	}
	for _, c := range cases {
		got := artStrip(c.input)
		if !strings.Contains(got, c.contain) {
			t.Errorf("artStrip(%q): got %q, want it to contain %q", c.input, got, c.contain)
		}
	}
}

func TestArtStrip_unknownFallback(t *testing.T) {
	got := artStrip("Some Obscure Day")
	if !strings.Contains(got, "{65}") {
		t.Errorf("fallback should be yellow ({65}), got %q", got)
	}
}

func TestArtStrip_caseInsensitive(t *testing.T) {
	lower := artStrip("christmas day")
	upper := artStrip("CHRISTMAS DAY")
	if lower != upper {
		t.Errorf("artStrip should be case-insensitive: %q != %q", lower, upper)
	}
}

func tileCount(strip string) int {
	count := 0
	i := 0
	for i < len(strip) {
		if strip[i] == '{' {
			j := strings.Index(strip[i:], "}")
			if j >= 0 {
				i += j + 1
				count++
				continue
			}
		}
		count++
		i++
	}
	return count
}

func TestArtStrip_allExactly15Tiles(t *testing.T) {
	names := []string{
		"Christmas Day", "New Year's Day", "Halloween", "Independence Day",
		"Thanksgiving Day", "Easter Sunday", "Valentine's Day", "St. Patrick's Day",
		"Diwali", "Deepavali", "Hanukkah", "Chanukah", "Lunar New Year",
		"Chinese New Year", "Eid al-Fitr", "Labor Day", "Workers' Day",
		"Memorial Day", "Veterans Day", "Martin Luther King Jr. Day",
		"Presidents' Day", "Guy Fawkes Night", "Bastille Day",
		"Anzac Day", "Remembrance Day", "Armistice Day",
		"Some Unknown Holiday",
	}
	for _, name := range names {
		strip := artStrip(name)
		n := tileCount(strip)
		if n != 15 {
			t.Errorf("artStrip(%q) has %d tiles, want 15 (strip: %q)", name, n, strip)
		}
	}
}
