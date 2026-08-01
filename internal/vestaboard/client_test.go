package vestaboard

import (
	"testing"
)

func TestCharCode(t *testing.T) {
	tests := []struct {
		r    rune
		want int
	}{
		{'A', 1},
		{'Z', 26},
		{'a', 1},
		{'z', 26},
		{' ', 0},
		{'1', 27},
		{'9', 35},
		{'0', 36},
		{'!', 37},
		{'.', 56},
		{'%', 54},
		{'?', 60},
		{'°', 62},
		{'€', 0},
	}
	for _, tt := range tests {
		got := charCode(tt.r)
		if got != tt.want {
			t.Errorf("charCode(%q) = %d, want %d", tt.r, got, tt.want)
		}
	}
}

func TestEncodeLines(t *testing.T) {
	lines := [3]string{"A", "", ""}
	grid := encodeLines(lines)
	if grid[0][0] != 1 {
		t.Errorf("encodeLines: grid[0][0] = %d, want 1 (A)", grid[0][0])
	}
	for c := 1; c < cols; c++ {
		if grid[0][c] != 0 {
			t.Errorf("encodeLines: grid[0][%d] = %d, want 0 (space)", c, grid[0][c])
		}
	}
}

func TestEncodeLinesDimensions(t *testing.T) {
	lines := [3]string{"HELLO", "WORLD", "12345"}
	grid := encodeLines(lines)
	if len(grid) != rows {
		t.Errorf("encodeLines rows = %d, want %d", len(grid), rows)
	}
	if len(grid[0]) != cols {
		t.Errorf("encodeLines cols = %d, want %d", len(grid[0]), cols)
	}
}

func TestEncodeLinesColorPlaceholder(t *testing.T) {
	lines := [3]string{"{67}{67}{67}", "", ""}
	grid := encodeLines(lines)
	for c := 0; c < 3; c++ {
		if grid[0][c] != 67 {
			t.Errorf("encodeLines: color placeholder grid[0][%d] = %d, want 67", c, grid[0][c])
		}
	}
}

func TestEncodeLinesFull(t *testing.T) {
	long := "ABCDEFGHIJKLMNO"
	lines := [3]string{long, "", ""}
	grid := encodeLines(lines)
	for c, r := range []rune(long) {
		want := charCode(r)
		if grid[0][c] != want {
			t.Errorf("encodeLines: grid[0][%d] = %d, want %d (%c)", c, grid[0][c], want, r)
		}
	}
}

func TestEncodeLinesTruncatesAt15(t *testing.T) {
	lines := [3]string{"ABCDEFGHIJKLMNOPQRSTU", "", ""}
	grid := encodeLines(lines)
	if grid[0][14] != charCode('O') {
		t.Errorf("encodeLines: col 14 = %d, want %d (O)", grid[0][14], charCode('O'))
	}
}
