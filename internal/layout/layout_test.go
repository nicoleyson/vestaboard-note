package layout

import (
	"strings"
	"testing"
)

func TestCenter(t *testing.T) {
	tests := []struct {
		input string
		width int
		want  string
	}{
		{"hello", 15, "     HELLO     "},
		{"A", 15, "       A       "},
		{"ABCDEFGHIJKLMNO", 15, "ABCDEFGHIJKLMNO"},
		{"ABCDEFGHIJKLMNOP", 15, "ABCDEFGHIJKLMNO"},
		{"AB", 6, "  AB  "},
		{"ABC", 6, " ABC  "},
	}
	for _, tt := range tests {
		got := Center(tt.input, tt.width)
		if got != tt.want {
			t.Errorf("Center(%q, %d) = %q, want %q", tt.input, tt.width, got, tt.want)
		}
		if len([]rune(got)) != tt.width && len([]rune(tt.input)) <= tt.width {
			t.Errorf("Center(%q, %d) length = %d, want %d", tt.input, tt.width, len(got), tt.width)
		}
	}
}

func TestCenterExtraPaddingTrails(t *testing.T) {
	got := Center("A", 4)
	if got != " A  " {
		t.Errorf("Center odd padding: got %q, want %q", got, " A  ")
	}
}

func TestPadRight(t *testing.T) {
	tests := []struct {
		input string
		width int
		want  string
	}{
		{"HI", 5, "HI   "},
		{"HELLO", 5, "HELLO"},
		{"TOOLONG", 4, "TOOL"},
		{"hi", 5, "HI   "},
	}
	for _, tt := range tests {
		got := PadRight(tt.input, tt.width)
		if got != tt.want {
			t.Errorf("PadRight(%q, %d) = %q, want %q", tt.input, tt.width, got, tt.want)
		}
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input string
		n     int
		want  string
	}{
		{"hello", 10, "HELLO"},
		{"hello world", 5, "HELL~"},
		{"hi", 2, "HI"},
		{"hi", 1, "H"},
		{"abcde", 5, "ABCDE"},
	}
	for _, tt := range tests {
		got := Truncate(tt.input, tt.n)
		if got != tt.want {
			t.Errorf("Truncate(%q, %d) = %q, want %q", tt.input, tt.n, got, tt.want)
		}
	}
}

func TestWrap(t *testing.T) {
	lines := Wrap("hello world", 7)
	if len(lines) != 2 {
		t.Fatalf("Wrap: got %d lines, want 2: %v", len(lines), lines)
	}
	if lines[0] != "HELLO  " {
		t.Errorf("Wrap line 0 = %q, want %q", lines[0], "HELLO  ")
	}
	if lines[1] != "WORLD  " {
		t.Errorf("Wrap line 1 = %q, want %q", lines[1], "WORLD  ")
	}
	for _, l := range lines {
		if len([]rune(l)) != 7 {
			t.Errorf("Wrap: line %q has length %d, want 7", l, len([]rune(l)))
		}
	}
}

func TestWrapEmpty(t *testing.T) {
	if lines := Wrap("", 10); lines != nil {
		t.Errorf("Wrap empty: got %v, want nil", lines)
	}
}

func TestStripEmoji(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello 😀 world", "hello  world"},
		{"❤️ love", "{62} love"},
		{"❤ love", "{62} love"},
		{":heart: hello", "hello"},
		{"plain text", "plain text"},
		{"☀️ sunny", "sunny"},
	}
	for _, tt := range tests {
		got := StripEmoji(tt.input)
		if got != tt.want {
			t.Errorf("StripEmoji(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestColorRow(t *testing.T) {
	row := ColorRow(67)
	if !strings.HasPrefix(row, "{67}") {
		t.Errorf("ColorRow(67) = %q, want prefix {67}", row)
	}
	tiles := strings.Count(row, "{67}")
	if tiles != Cols {
		t.Errorf("ColorRow(67) has %d tiles, want %d", tiles, Cols)
	}
}

func TestFromLines(t *testing.T) {
	g := FromLines("hello", "world")
	lines := g.Lines()
	if lines[0] != "HELLO          " {
		t.Errorf("FromLines row 0 = %q", lines[0])
	}
	if lines[1] != "WORLD          " {
		t.Errorf("FromLines row 1 = %q", lines[1])
	}
	if lines[2] != "               " {
		t.Errorf("FromLines row 2 (empty) = %q", lines[2])
	}
	for r := 0; r < Rows; r++ {
		if len([]rune(lines[r])) != Cols {
			t.Errorf("FromLines row %d length = %d, want %d", r, len(lines[r]), Cols)
		}
	}
}

func TestWrap_longWordSplit(t *testing.T) {
	word := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lines := Wrap(word, 15)
	if len(lines) == 0 {
		t.Fatal("Wrap of long word returned no lines")
	}
	for i, l := range lines {
		if len([]rune(l)) != 15 {
			t.Errorf("Wrap long word: line %d len = %d, want 15: %q", i, len([]rune(l)), l)
		}
	}
	combined := ""
	for _, l := range lines {
		combined += strings.TrimRight(l, " ")
	}
	if combined != word {
		t.Errorf("Wrap long word: reassembled = %q, want %q", combined, word)
	}
}

