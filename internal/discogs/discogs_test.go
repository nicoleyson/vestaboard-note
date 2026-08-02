package discogs

import (
	"testing"
)

func TestCleanArtistName(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"The Beatles (2)", "The Beatles"},
		{"Miles Davis", "Miles Davis"},
		{"Prince (3)", "Prince"},
		{"", ""},
		{"Band (no closing paren", "Band"},
	}
	for _, tc := range cases {
		got := cleanArtistName(tc.input)
		if got != tc.want {
			t.Errorf("cleanArtistName(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestScoreRecord_baseScore(t *testing.T) {
	r := record{artist: "X", title: "Y", styles: []string{}}
	got := scoreRecord(r, nil, nil, morning)
	if got != 1 {
		t.Errorf("base score should be 1, got %d", got)
	}
}

func TestScoreRecord_weatherMatch(t *testing.T) {
	r := record{styles: []string{"Jazz"}}
	got := scoreRecord(r, []string{"Jazz"}, nil, morning)
	if got != 3 {
		t.Errorf("weather match should add 2 (total 3), got %d", got)
	}
}

func TestScoreRecord_timeMatch(t *testing.T) {
	r := record{styles: []string{"Ambient"}}
	got := scoreRecord(r, nil, []string{"Ambient"}, morning)
	if got != 2 {
		t.Errorf("time match should add 1 (total 2), got %d", got)
	}
}

func TestScoreRecord_titleMatch(t *testing.T) {
	r := record{title: "Morning Light", styles: []string{}}
	got := scoreRecord(r, nil, nil, morning)
	if got != 4 {
		t.Errorf("title match should add 3 (total 4), got %d", got)
	}
}

func TestScoreRecord_combinedMatch(t *testing.T) {
	r := record{title: "Night Train", styles: []string{"Blues", "Jazz"}}
	got := scoreRecord(r, []string{"Blues"}, []string{"Jazz"}, night)
	if got != 1+2+1+3 {
		t.Errorf("combined score: want 7, got %d", got)
	}
}

func TestWeightedPick_singleRecord(t *testing.T) {
	records := []record{{artist: "Solo", title: "One"}}
	scores := []int{5}
	got := weightedPick(records, scores)
	if got.artist != "Solo" {
		t.Errorf("single record: want Solo, got %q", got.artist)
	}
}

func TestWeightedPick_allSameScore(t *testing.T) {
	records := []record{
		{artist: "A"}, {artist: "B"}, {artist: "C"},
	}
	scores := []int{1, 1, 1}
	seen := map[string]bool{}
	for i := 0; i < 30; i++ {
		got := weightedPick(records, scores)
		seen[got.artist] = true
	}
	if len(seen) < 2 {
		t.Error("expected distribution across records with equal scores")
	}
}

func TestWeightedPick_highScoreDominates(t *testing.T) {
	records := []record{{artist: "Rare"}, {artist: "Common"}}
	scores := []int{1, 99}
	commonCount := 0
	for i := 0; i < 100; i++ {
		if weightedPick(records, scores).artist == "Common" {
			commonCount++
		}
	}
	if commonCount < 80 {
		t.Errorf("high-score record should dominate: got %d/100", commonCount)
	}
}

