package layout

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	Rows = 3
	Cols = 15
)

type Grid [Rows][Cols]string

func (g Grid) Lines() [Rows]string {
	var lines [Rows]string
	for r := 0; r < Rows; r++ {
		var sb strings.Builder
		for c := 0; c < Cols; c++ {
			sb.WriteString(g[r][c])
		}
		lines[r] = sb.String()
	}
	return lines
}

func FromLines(lines ...string) Grid {
	var g Grid
	for r := 0; r < Rows && r < len(lines); r++ {
		runes := []rune(strings.ToUpper(lines[r]))
		for c := 0; c < Cols; c++ {
			if c < len(runes) {
				g[r][c] = string(runes[c])
			} else {
				g[r][c] = " "
			}
		}
	}
	for r := len(lines); r < Rows; r++ {
		for c := 0; c < Cols; c++ {
			g[r][c] = " "
		}
	}
	return g
}

func Center(s string, width int) string {
	runes := []rune(strings.ToUpper(s))
	if len(runes) >= width {
		return string(runes[:width])
	}
	total := width - len(runes)
	left := total / 2
	right := total - left
	return strings.Repeat(" ", left) + string(runes) + strings.Repeat(" ", right)
}

func PadRight(s string, width int) string {
	runes := []rune(strings.ToUpper(s))
	if len(runes) >= width {
		return string(runes[:width])
	}
	return string(runes) + strings.Repeat(" ", width-len(runes))
}

func Wrap(text string, width int) []string {
	words := strings.Fields(strings.ToUpper(text))
	if len(words) == 0 {
		return nil
	}
	var lines []string
	current := ""
	for _, word := range words {
		wordRunes := []rune(word)
		for utf8.RuneCountInString(word) > width {
			space := width - utf8.RuneCountInString(current)
			if current != "" {
				space--
			}
			if space <= 0 {
				lines = append(lines, PadRight(current, width))
				current = ""
				space = width
			}
			lines = append(lines, PadRight(string(wordRunes[:space]), width))
			wordRunes = wordRunes[space:]
			word = string(wordRunes)
		}
		if word == "" {
			continue
		}
		if current == "" {
			current = word
		} else if utf8.RuneCountInString(current)+1+utf8.RuneCountInString(word) <= width {
			current += " " + word
		} else {
			lines = append(lines, PadRight(current, width))
			current = word
		}
	}
	if current != "" {
		lines = append(lines, PadRight(current, width))
	}
	return lines
}

func Truncate(s string, n int) string {
	runes := []rune(strings.ToUpper(s))
	if len(runes) <= n {
		return string(runes)
	}
	if n <= 1 {
		return string(runes[:n])
	}
	return string(runes[:n-1]) + "~"
}

var shortcodeRe = regexp.MustCompile(`:[a-zA-Z0-9_+-]+:`)

func StripEmoji(s string) string {
	s = strings.ReplaceAll(s, "❤️", "{62}")
	s = strings.ReplaceAll(s, "❤", "{62}")

	s = shortcodeRe.ReplaceAllString(s, "")

	var b strings.Builder
	for _, r := range s {
		if unicode.Is(unicode.So, r) || unicode.Is(unicode.Sm, r) ||
			(r >= 0xFE00 && r <= 0xFE1F) ||
			(r >= 0x1F000 && r <= 0x1FFFF) ||
			(r >= 0x2600 && r <= 0x27BF) {
			continue
		}
		b.WriteRune(r)
	}
	return strings.TrimSpace(b.String())
}

func ColorRow(code int) string {
	tile := fmt.Sprintf("{%d}", code)
	return strings.Repeat(tile, Cols)
}
