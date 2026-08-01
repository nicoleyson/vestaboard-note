package pattern

import (
	"strings"
	"testing"

	"github.com/nicoleyson/vestaboard-note/internal/layout"
)

func tileCount(line string) int {
	count := 0
	s := line
	for len(s) > 0 {
		if s[0] == '{' {
			end := strings.IndexByte(s, '}')
			if end > 0 {
				s = s[end+1:]
				count++
				continue
			}
		}
		_, size := runeAt(s)
		s = s[size:]
		count++
	}
	return count
}

func runeAt(s string) (rune, int) {
	for i, r := range s {
		_ = i
		return r, len(string(r))
	}
	return 0, 1
}

func assertLines(t *testing.T, name string, lines [3]string) {
	t.Helper()
	for i, line := range lines {
		count := tileCount(line)
		if count != layout.Cols {
			t.Errorf("%s: row %d has %d tiles, want %d: %q", name, i, count, layout.Cols, line)
		}
	}
}

func TestAllNamedPatterns(t *testing.T) {
	for _, name := range Names {
		lines, err := Generate(name)
		if err != nil {
			t.Errorf("Generate(%q) error: %v", name, err)
			continue
		}
		assertLines(t, name, lines)
	}
}

func TestGenerateUnknown(t *testing.T) {
	_, err := Generate("doesnotexist")
	if err == nil {
		t.Error("Generate(unknown) expected error, got nil")
	}
}

func TestRandom(t *testing.T) {
	for i := 0; i < 10; i++ {
		lines := Random()
		assertLines(t, "random", lines)
	}
}

func TestNamesNotEmpty(t *testing.T) {
	if len(Names) == 0 {
		t.Error("Names is empty")
	}
}
