package vestaboard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	baseURL     = "https://cloud.vestaboard.com"
	minInterval = 15 * time.Second
)

type Client struct {
	token    string
	http     *http.Client
	lastSent time.Time
}

func New(token string) *Client {
	return &Client{token: token, http: &http.Client{Timeout: 10 * time.Second}}
}

func (c *Client) SendLines(lines [3]string) error {
	if wait := minInterval - time.Since(c.lastSent); wait > 0 {
		return fmt.Errorf("rate limit: wait %v", wait.Round(time.Second))
	}
	body, err := json.Marshal(map[string]string{"text": strings.Join(lines[:], "\n")})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, baseURL+"/message", bytes.NewReader(body))
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
