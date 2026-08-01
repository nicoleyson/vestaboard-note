package main

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/nicoleyson/vestaboard-note/internal/air"
	"github.com/nicoleyson/vestaboard-note/internal/calendar"
	"github.com/nicoleyson/vestaboard-note/internal/clock"
	"github.com/nicoleyson/vestaboard-note/internal/countdown"
	"github.com/nicoleyson/vestaboard-note/internal/flights"
	"github.com/nicoleyson/vestaboard-note/internal/moon"
	"github.com/nicoleyson/vestaboard-note/internal/onthisday"
	"github.com/nicoleyson/vestaboard-note/internal/vestaboard"
	"github.com/nicoleyson/vestaboard-note/internal/weather"
)

type config struct {
	Token      string              `yaml:"token"`
	Lat        float64             `yaml:"lat"`
	Lon        float64             `yaml:"lon"`
	ICalURLs   []string            `yaml:"ical_urls"`
	Countdowns []countdown.Event   `yaml:"countdowns"`
}

func loadConfig() (config, error) {
	path := os.Getenv("VESTABOARD_CONFIG")
	if path == "" {
		path = "config.yaml"
	}
	f, err := os.Open(path)
	if err != nil {
		return config{}, fmt.Errorf("open config: %w", err)
	}
	defer f.Close()
	var cfg config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return config{}, fmt.Errorf("parse config: %w", err)
	}
	return cfg, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: note <weather|clock|calendar|moon|air|flights|onthisday|countdown>\n")
		os.Exit(1)
	}

	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if cfg.Token == "" {
		fmt.Fprintf(os.Stderr, "error: token is required in config.yaml\n")
		os.Exit(1)
	}

	var lines [3]string

	switch os.Args[1] {
	case "weather":
		if cfg.Lat == 0 || cfg.Lon == 0 {
			fmt.Fprintf(os.Stderr, "error: lat and lon required for weather\n")
			os.Exit(1)
		}
		lines, err = weather.Fetch(cfg.Lat, cfg.Lon)
	case "clock":
		lines = clock.Format(time.Now())
	case "calendar":
		if len(cfg.ICalURLs) == 0 {
			fmt.Fprintf(os.Stderr, "error: ical_urls required for calendar\n")
			os.Exit(1)
		}
		lines, err = calendar.Fetch(cfg.ICalURLs)
	case "moon":
		lines = moon.Format(time.Now())
	case "air":
		if cfg.Lat == 0 || cfg.Lon == 0 {
			fmt.Fprintf(os.Stderr, "error: lat and lon required for air\n")
			os.Exit(1)
		}
		lines, err = air.Fetch(cfg.Lat, cfg.Lon)
	case "flights":
		if cfg.Lat == 0 || cfg.Lon == 0 {
			fmt.Fprintf(os.Stderr, "error: lat and lon required for flights\n")
			os.Exit(1)
		}
		lines, err = flights.Fetch(cfg.Lat, cfg.Lon)
	case "onthisday":
		lines, err = onthisday.Fetch(time.Now())
	case "countdown":
		lines = countdown.Format(cfg.Countdowns, time.Now())
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	client := vestaboard.New(cfg.Token)
	if err := client.SendLines(lines); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("[%s]\n[%s]\n[%s]\n", lines[0], lines[1], lines[2])
}
