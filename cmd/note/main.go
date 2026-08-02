package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/nicoleyson/vestaboard-note/internal/air"
	"github.com/nicoleyson/vestaboard-note/internal/holiday"
	"github.com/nicoleyson/vestaboard-note/internal/pattern"
	"github.com/nicoleyson/vestaboard-note/internal/tearoff"
	"github.com/nicoleyson/vestaboard-note/internal/calendar"
	"github.com/nicoleyson/vestaboard-note/internal/clock"
	"github.com/nicoleyson/vestaboard-note/internal/countdown"
	"github.com/nicoleyson/vestaboard-note/internal/discogs"
	"github.com/nicoleyson/vestaboard-note/internal/flights"
	"github.com/nicoleyson/vestaboard-note/internal/moonphase"
	"github.com/nicoleyson/vestaboard-note/internal/onthisday"
	"github.com/nicoleyson/vestaboard-note/internal/pollen"
	"github.com/nicoleyson/vestaboard-note/internal/rain"
	"github.com/nicoleyson/vestaboard-note/internal/satellites"
	"github.com/nicoleyson/vestaboard-note/internal/season"
	"github.com/nicoleyson/vestaboard-note/internal/suntime"
	"github.com/nicoleyson/vestaboard-note/internal/sunscene"
	"github.com/nicoleyson/vestaboard-note/internal/uv"
	"github.com/nicoleyson/vestaboard-note/internal/vestaboard"
	"github.com/nicoleyson/vestaboard-note/internal/weather"
)

var subcommands = []string{
	"weather", "clock", "calendar", "moonphase", "air",
	"flights", "onthisday", "countdown", "discogs", "pattern",
	"suntime", "sunscene", "pollen", "uv", "rain", "season", "holiday", "satellites", "tearoff",
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

	printLines("moonphase", moonphase.Format(now), nil)
	printLines("season", season.Format(now), nil)

	if cfg.Lat != 0 || cfg.Lon != 0 {
		lines, _, err = air.Fetch(cfg.Lat, cfg.Lon)
		printLines("air", lines, err)

		lines, err = flights.Fetch(cfg.Lat, cfg.Lon)
		printLines("flights", lines, err)

		lines, err = suntime.Fetch(cfg.Lat, cfg.Lon)
		printLines("suntime", lines, err)

		lines, err = sunscene.Fetch(cfg.Lat, cfg.Lon)
		printLines("sunscene", lines, err)

		lines, _, err = pollen.Fetch(cfg.Lat, cfg.Lon)
		printLines("pollen", lines, err)

		lines, err = uv.Fetch(cfg.Lat, cfg.Lon)
		printLines("uv", lines, err)

		lines, _, err = rain.Fetch(cfg.Lat, cfg.Lon)
		printLines("rain", lines, err)

		lines, _, err = holiday.Fetch(cfg.Lat, cfg.Lon)
		printLines("holiday", lines, err)

		lines, _, err = satellites.Fetch(cfg.Lat, cfg.Lon)
		printLines("satellites", lines, err)
	} else {
		fmt.Printf("  %-12s skipped (no lat/lon)\n", "air")
		fmt.Printf("  %-12s skipped (no lat/lon)\n", "flights")
		fmt.Printf("  %-12s skipped (no lat/lon)\n", "suntime")
		fmt.Printf("  %-12s skipped (no lat/lon)\n", "sunscene")
		fmt.Printf("  %-12s skipped (no lat/lon)\n", "pollen")
		fmt.Printf("  %-12s skipped (no lat/lon)\n", "uv")
		fmt.Printf("  %-12s skipped (no lat/lon)\n", "rain")
		fmt.Printf("  %-12s skipped (no lat/lon)\n", "holiday")
		fmt.Printf("  %-12s skipped (no lat/lon)\n", "satellites")
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

	printLines("pattern", pattern.Random(), nil)
	printLines("tearoff", tearoff.Format(now), nil)
}

const bashCompletion = `_note_completions() {
    local cur="${COMP_WORDS[COMP_CWORD]}"
    local prev="${COMP_WORDS[COMP_CWORD-1]}"
    local cmds="weather clock calendar moonphase air flights onthisday countdown discogs pattern suntime sunscene pollen uv rain season holiday satellites tearoff status completion"
    local patterns="current random stripes checker bars fade diagonal hearts confetti sparkle pulse rainbow"
    if [[ "${prev}" == "pattern" ]]; then
        COMPREPLY=($(compgen -W "${patterns}" -- "${cur}"))
    elif [[ "${prev}" == "completion" ]]; then
        COMPREPLY=($(compgen -W "bash zsh fish" -- "${cur}"))
    else
        COMPREPLY=($(compgen -W "${cmds}" -- "${cur}"))
    fi
}
complete -F _note_completions note`

const zshCompletion = `#compdef note
_note() {
    local -a cmds patterns
    cmds=(
        'weather:current conditions'
        'clock:date and time'
        'calendar:next calendar event'
        'moonphase:lunar phase'
        'air:air quality index'
        'flights:aircraft overhead'
        'onthisday:historical event'
        'countdown:days until configured event'
        'discogs:record matched to weather and time'
        'pattern:random art and color patterns'
        'suntime:next sunrise or sunset time'
        'sunscene:visual sunrise or sunset color art'
        'pollen:pollen levels by type'
        'uv:uv index with color scale'
        'rain:precipitation probability and intensity'
        'season:current astronomical season with color'
        'holiday:today'"'"'s public holiday by location'
        'satellites:notable satellites currently overhead'
        'tearoff:tear-off calendar showing today'"'"'s date'
        'status:preview all subcommands without sending'
        'completion:print shell completion script'
    )
    patterns=(
        'current:seasonal and holiday-aware pattern based on location'
        'random:pick a random pattern'
        'stripes:vertical color bands'
        'checker:checkerboard of two colors'
        'bars:three horizontal color bars'
        'fade:left-to-right color gradient'
        'diagonal:diagonal color bands'
        'hearts:color background with hearts'
        'confetti:random colors and letters'
        'sparkle:dark background with bright spots'
        'pulse:concentric color rings from center'
        'rainbow:spectrum stripes'
    )
    if (( CURRENT == 2 )); then
        _describe 'command' cmds
    elif [[ ${words[2]} == pattern ]]; then
        _describe 'pattern' patterns
    elif [[ ${words[2]} == completion ]]; then
        _values 'shell' bash zsh fish
    fi
}
_note`

const fishCompletion = `complete -c note -f
complete -c note -n __fish_use_subcommand -a weather    -d 'Current conditions'
complete -c note -n __fish_use_subcommand -a clock      -d 'Date and time'
complete -c note -n __fish_use_subcommand -a calendar   -d 'Next calendar event'
complete -c note -n __fish_use_subcommand -a moonphase   -d 'Lunar phase'
complete -c note -n __fish_use_subcommand -a air        -d 'Air quality index'
complete -c note -n __fish_use_subcommand -a flights    -d 'Aircraft overhead'
complete -c note -n __fish_use_subcommand -a onthisday  -d 'Historical event'
complete -c note -n __fish_use_subcommand -a countdown  -d 'Days until configured event'
complete -c note -n __fish_use_subcommand -a discogs    -d 'Record matched to weather and time'
complete -c note -n __fish_use_subcommand -a pattern    -d 'Random art and color patterns'
complete -c note -n __fish_use_subcommand -a suntime     -d 'Next sunrise or sunset time'
complete -c note -n __fish_use_subcommand -a sunscene   -d 'Next sunrise or sunset color art'
complete -c note -n __fish_use_subcommand -a pollen     -d 'Pollen levels by type'
complete -c note -n __fish_use_subcommand -a uv         -d 'UV index with color scale'
complete -c note -n __fish_use_subcommand -a rain       -d 'Precipitation probability and intensity'
complete -c note -n __fish_use_subcommand -a season     -d 'Current astronomical season with color'
complete -c note -n __fish_use_subcommand -a holiday    -d 'Today'"'"'s public holiday by location'
complete -c note -n __fish_use_subcommand -a satellites -d 'Notable satellites currently overhead'
complete -c note -n __fish_use_subcommand -a tearoff    -d 'Tear-off calendar showing today'"'"'s date'
complete -c note -n __fish_use_subcommand -a status     -d 'Preview all subcommands without sending'
complete -c note -n __fish_use_subcommand -a completion -d 'Print shell completion script'
complete -c note -n '__fish_seen_subcommand_from pattern' -a 'current random stripes checker bars fade diagonal hearts confetti sparkle pulse rainbow'
complete -c note -n '__fish_seen_subcommand_from completion' -a 'bash zsh fish'`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: note <%s>\n", strings.Join(subcommands, "|"))
		os.Exit(1)
	}

	cmd := os.Args[1]

	skipTrivial := false
	for _, arg := range os.Args[2:] {
		if arg == "--skip-trivial" {
			skipTrivial = true
		}
	}

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
	var trivial bool

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
	case "moonphase":
		lines = moonphase.Format(time.Now())
	case "air":
		if cfg.Lat == 0 || cfg.Lon == 0 {
			fmt.Fprintf(os.Stderr, "error: lat and lon required for air\n")
			os.Exit(1)
		}
		lines, trivial, err = air.Fetch(cfg.Lat, cfg.Lon)
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
	case "pattern":
		name := "random"
		if len(os.Args) >= 3 {
			name = os.Args[2]
		}
		switch name {
		case "random":
			lines = pattern.Random()
		case "current":
			if cfg.Lat == 0 || cfg.Lon == 0 {
				fmt.Fprintf(os.Stderr, "error: lat and lon required for pattern current\n")
				os.Exit(1)
			}
			lines, err = pattern.Current(cfg.Lat, cfg.Lon)
		default:
			lines, err = pattern.Generate(name)
		}
	case "suntime":
		if cfg.Lat == 0 || cfg.Lon == 0 {
			fmt.Fprintf(os.Stderr, "error: lat and lon required for suntime\n")
			os.Exit(1)
		}
		lines, err = suntime.Fetch(cfg.Lat, cfg.Lon)
	case "sunscene":
		if cfg.Lat == 0 || cfg.Lon == 0 {
			fmt.Fprintf(os.Stderr, "error: lat and lon required for sunscene\n")
			os.Exit(1)
		}
		lines, err = sunscene.Fetch(cfg.Lat, cfg.Lon)
	case "pollen":
		if cfg.Lat == 0 || cfg.Lon == 0 {
			fmt.Fprintf(os.Stderr, "error: lat and lon required for pollen\n")
			os.Exit(1)
		}
		lines, trivial, err = pollen.Fetch(cfg.Lat, cfg.Lon)
	case "uv":
		if cfg.Lat == 0 || cfg.Lon == 0 {
			fmt.Fprintf(os.Stderr, "error: lat and lon required for uv\n")
			os.Exit(1)
		}
		lines, err = uv.Fetch(cfg.Lat, cfg.Lon)
	case "rain":
		if cfg.Lat == 0 || cfg.Lon == 0 {
			fmt.Fprintf(os.Stderr, "error: lat and lon required for rain\n")
			os.Exit(1)
		}
		lines, trivial, err = rain.Fetch(cfg.Lat, cfg.Lon)
	case "season":
		lines = season.Format(time.Now())
	case "tearoff":
		lines = tearoff.Format(time.Now())
	case "holiday":
		if cfg.Lat == 0 || cfg.Lon == 0 {
			fmt.Fprintf(os.Stderr, "error: lat and lon required for holiday\n")
			os.Exit(1)
		}
		lines, trivial, err = holiday.Fetch(cfg.Lat, cfg.Lon)
	case "satellites":
		if cfg.Lat == 0 || cfg.Lon == 0 {
			fmt.Fprintf(os.Stderr, "error: lat and lon required for satellites\n")
			os.Exit(1)
		}
		lines, trivial, err = satellites.Fetch(cfg.Lat, cfg.Lon)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if trivial && skipTrivial {
		return
	}

	client := vestaboard.New(cfg.Token)
	if err := client.SendLines(lines); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("[%s]\n[%s]\n[%s]\n", lines[0], lines[1], lines[2])
}
