package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/nicoleyson/vestaboard-note/internal/air"
	"github.com/nicoleyson/vestaboard-note/internal/calendar"
	"github.com/nicoleyson/vestaboard-note/internal/clock"
	"github.com/nicoleyson/vestaboard-note/internal/countdown"
	"github.com/nicoleyson/vestaboard-note/internal/discogs"
	"github.com/nicoleyson/vestaboard-note/internal/flights"
	"github.com/nicoleyson/vestaboard-note/internal/moon"
	"github.com/nicoleyson/vestaboard-note/internal/onthisday"
	"github.com/nicoleyson/vestaboard-note/internal/vestaboard"
	"github.com/nicoleyson/vestaboard-note/internal/weather"
)

var subcommands = []string{
	"weather", "clock", "calendar", "moon", "air",
	"flights", "onthisday", "countdown", "discogs",
	"status", "completion",
}

type config struct {
	Token           string            `yaml:"token"`
	Lat             float64           `yaml:"lat"`
	Lon             float64           `yaml:"lon"`
	ICalURLs        []string          `yaml:"ical_urls"`
	Countdowns      []countdown.Event `yaml:"countdowns"`
	DiscogsToken    string            `yaml:"discogs_token"`
	DiscogsUsername string            `yaml:"discogs_username"`
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

func printLines(label string, lines [3]string, err error) {
	if err != nil {
		fmt.Printf("  %-12s error: %v\n", label, err)
		return
	}
	fmt.Printf("  %-12s row1: %q\n", label, lines[0])
	fmt.Printf("  %-12s row2: %q\n", "", lines[1])
	fmt.Printf("  %-12s row3: %q\n", "", lines[2])
}

func runStatus(cfg config) {
	now := time.Now()
	fmt.Println("vestaboard-note status")
	fmt.Println(strings.Repeat("─", 50))

	lines, err := weather.Fetch(cfg.Lat, cfg.Lon)
	printLines("weather", lines, err)

	printLines("clock", clock.Format(now), nil)

	if len(cfg.ICalURLs) > 0 {
		lines, err = calendar.Fetch(cfg.ICalURLs)
		printLines("calendar", lines, err)
	} else {
		fmt.Printf("  %-12s skipped (no ical_urls)\n", "calendar")
	}

	printLines("moon", moon.Format(now), nil)

	if cfg.Lat != 0 || cfg.Lon != 0 {
		lines, err = air.Fetch(cfg.Lat, cfg.Lon)
		printLines("air", lines, err)

		lines, err = flights.Fetch(cfg.Lat, cfg.Lon)
		printLines("flights", lines, err)
	} else {
		fmt.Printf("  %-12s skipped (no lat/lon)\n", "air")
		fmt.Printf("  %-12s skipped (no lat/lon)\n", "flights")
	}

	lines, err = onthisday.Fetch(now)
	printLines("onthisday", lines, err)

	printLines("countdown", countdown.Format(cfg.Countdowns, now), nil)

	if cfg.DiscogsToken != "" && cfg.DiscogsUsername != "" && (cfg.Lat != 0 || cfg.Lon != 0) {
		lines, err = discogs.Fetch(cfg.DiscogsUsername, cfg.DiscogsToken, cfg.Lat, cfg.Lon)
		printLines("discogs", lines, err)
	} else {
		fmt.Printf("  %-12s skipped (no discogs_token/discogs_username)\n", "discogs")
	}
}

const bashCompletion = `_note_completions() {
    local cur="${COMP_WORDS[COMP_CWORD]}"
    local cmds="weather clock calendar moon air flights onthisday countdown discogs status completion"
    COMPREPLY=($(compgen -W "${cmds}" -- "${cur}"))
}
complete -F _note_completions note`

const zshCompletion = `#compdef note
_note() {
    local -a cmds
    cmds=(
        'weather:current conditions'
        'clock:date and time'
        'calendar:next calendar event'
        'moon:lunar phase'
        'air:air quality index'
        'flights:aircraft overhead'
        'onthisday:historical event'
        'countdown:days until configured event'
        'discogs:record matched to weather and time'
        'status:preview all subcommands without sending'
        'completion:print shell completion script'
    )
    _describe 'command' cmds
}
_note`

const fishCompletion = `complete -c note -f
complete -c note -n __fish_use_subcommand -a weather    -d 'Current conditions'
complete -c note -n __fish_use_subcommand -a clock      -d 'Date and time'
complete -c note -n __fish_use_subcommand -a calendar   -d 'Next calendar event'
complete -c note -n __fish_use_subcommand -a moon       -d 'Lunar phase'
complete -c note -n __fish_use_subcommand -a air        -d 'Air quality index'
complete -c note -n __fish_use_subcommand -a flights    -d 'Aircraft overhead'
complete -c note -n __fish_use_subcommand -a onthisday  -d 'Historical event'
complete -c note -n __fish_use_subcommand -a countdown  -d 'Days until configured event'
complete -c note -n __fish_use_subcommand -a discogs    -d 'Record matched to weather and time'
complete -c note -n __fish_use_subcommand -a status     -d 'Preview all subcommands without sending'
complete -c note -n __fish_use_subcommand -a completion -d 'Print shell completion script'`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: note <%s>\n", strings.Join(subcommands, "|"))
		os.Exit(1)
	}

	cmd := os.Args[1]

	if cmd == "completion" {
		shell := "bash"
		if len(os.Args) >= 3 {
			shell = os.Args[2]
		}
		switch shell {
		case "bash":
			fmt.Println(bashCompletion)
		case "zsh":
			fmt.Println(zshCompletion)
		case "fish":
			fmt.Println(fishCompletion)
		default:
			fmt.Fprintf(os.Stderr, "unknown shell: %s (use bash, zsh, or fish)\n", shell)
			os.Exit(1)
		}
		return
	}

	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if cmd == "status" {
		runStatus(cfg)
		return
	}

	if cfg.Token == "" {
		fmt.Fprintf(os.Stderr, "error: token is required in config.yaml\n")
		os.Exit(1)
	}

	var lines [3]string

	switch cmd {
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
	case "discogs":
		if cfg.DiscogsToken == "" || cfg.DiscogsUsername == "" {
			fmt.Fprintf(os.Stderr, "error: discogs_token and discogs_username required for discogs\n")
			os.Exit(1)
		}
		if cfg.Lat == 0 || cfg.Lon == 0 {
			fmt.Fprintf(os.Stderr, "error: lat and lon required for discogs\n")
			os.Exit(1)
		}
		lines, err = discogs.Fetch(cfg.DiscogsUsername, cfg.DiscogsToken, cfg.Lat, cfg.Lon)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
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
