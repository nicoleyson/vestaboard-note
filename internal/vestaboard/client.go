package vestaboard

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var baseURL = "https://cloud.vestaboard.com"

const (
	minInterval = 15 * time.Second
	rows        = 3
	cols        = 15
)

type Client struct {
	token    string
	http     *http.Client
	lastSent time.Time
}

func New(token string) *Client {
	return &Client{token: token, http: &http.Client{Timeout: 10 * time.Second}}
}

func charCode(r rune) int {
	switch {
	case r >= 'A' && r <= 'Z':
		return int(r-'A') + 1
	case r >= 'a' && r <= 'z':
		return int(r-'a') + 1
	case r == ' ':
		return 0
	case r >= '1' && r <= '9':
		return int(r-'1') + 27
	case r == '0':
		return 36
	case r == '!':
		return 37
	case r == '@':
		return 38
	case r == '#':
		return 39
	case r == '$':
		return 40
	case r == '(':
		return 41
	case r == ')':
		return 42
	case r == '-':
		return 44
	case r == '+':
		return 46
	case r == '&':
		return 47
	case r == '=':
		return 48
	case r == ';':
		return 49
	case r == ':':
		return 50
	case r == '\'':
		return 52
	case r == '"':
		return 53
	case r == '%':
		return 54
	case r == ',':
		return 55
	case r == '.':
		return 56
	case r == '/':
		return 59
	case r == '?':
		return 60
	case r == '°':
		return 0 // no degree glyph on Vestaboard; render as space
	default:
		return 0
	}
}

func encodeLines(lines [rows]string) [rows][cols]int {
	var grid [rows][cols]int
	for r, line := range lines {
		var codes []int
		s := strings.ToUpper(line)
		for len(s) > 0 {
			if s[0] == '{' {
				end := strings.IndexByte(s, '}')
				if end > 1 {
					if n, err := strconv.Atoi(s[1:end]); err == nil {
						codes = append(codes, n)
						s = s[end+1:]
						continue
					}
				}
			}
			ru := []rune(s)[0]
			codes = append(codes, charCode(ru))
			s = s[len(string(ru)):]
		}
		for c := 0; c < cols && c < len(codes); c++ {
			grid[r][c] = codes[c]
		}
	}
	return grid
}

func (c *Client) SendLines(lines [rows]string) error {
	if wait := minInterval - time.Since(c.lastSent); wait > 0 {
		return fmt.Errorf("rate limit: wait %v", wait.Round(time.Second))
	}

	grid := encodeLines(lines)

	body, err := json.Marshal(map[string][rows][cols]int{"characters": grid})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, baseURL+"/", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("X-Vestaboard-Token", c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var buf bytes.Buffer
		buf.ReadFrom(resp.Body)
		return fmt.Errorf("api %d: %s", resp.StatusCode, buf.String())
	}
	c.lastSent = time.Now()
	return nil
}
