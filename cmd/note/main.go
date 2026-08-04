package main

import (
	"fmt"
	"log"
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
	"github.com/nicoleyson/vestaboard-note/internal/holiday"
	"github.com/nicoleyson/vestaboard-note/internal/moonphase"
	"github.com/nicoleyson/vestaboard-note/internal/pattern"
	"github.com/nicoleyson/vestaboard-note/internal/pollen"
	"github.com/nicoleyson/vestaboard-note/internal/rain"
	"github.com/nicoleyson/vestaboard-note/internal/satellites"
	"github.com/nicoleyson/vestaboard-note/internal/season"
	"github.com/nicoleyson/vestaboard-note/internal/sunscene"
	"github.com/nicoleyson/vestaboard-note/internal/suntime"
	"github.com/nicoleyson/vestaboard-note/internal/tearoff"
	"github.com/nicoleyson/vestaboard-note/internal/uv"
	"github.com/nicoleyson/vestaboard-note/internal/vestaboard"
	"github.com/nicoleyson/vestaboard-note/internal/weather"
)

// version is set at build time via -ldflags "-X main.version=<git-sha>"
var version = "dev"

var subcommands = []string{
	"weather", "clock", "calendar", "moonphase", "air",
	"flights", "countdown", "discogs", "pattern",
	"suntime", "sunscene", "pollen", "uv", "rain", "season", "holiday", "satellites", "tearoff",
	"daemon", "status", "completion",
}

// scheduleEntry defines one scheduled job in config.yaml.
// Fields:
//
//	command        — subcommand name (e.g. "weather")
//	hour           — hour to run (0–23, local timezone)
//	minute         — minute to run (default 0)
//	repeat_minutes — if > 0, repeat every N minutes within the hour window
//	until_hour     — last hour (inclusive) to repeat within (requires repeat_minutes)
//	args           — extra flags passed to the subcommand (e.g. ["--skip-trivial"])
type scheduleEntry struct {
	Command       string   `yaml:"command"`
	Hour          int      `yaml:"hour"`
	Minute        int      `yaml:"minute"`
	RepeatMinutes int      `yaml:"repeat_minutes"`
	UntilHour     int      `yaml:"until_hour"`
	Args          []string `yaml:"args"`
}

type config struct {
	Token           string          `yaml:"token"`
	Lat             float64         `yaml:"lat"`
	Lon             float64         `yaml:"lon"`
	ICalURLs        []string        `yaml:"ical_urls"`
	Countdowns      []countdown.Event `yaml:"countdowns"`
	DiscogsToken    string          `yaml:"discogs_token"`
	DiscogsUsername string          `yaml:"discogs_username"`
	Timezone        string          `yaml:"timezone"`
	Schedule        []scheduleEntry `yaml:"schedule"`
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
	fmt.Printf("vestaboard-note status (build %s)\n", version)
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

		lines, _, err = flights.Fetch(cfg.Lat, cfg.Lon)
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
    local cmds="weather clock calendar moonphase air flights countdown discogs pattern suntime sunscene pollen uv rain season holiday satellites tearoff daemon status completion"
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
        'countdown:days until configured event'
        'discogs:record matched to weather and time'
        'pattern:random art and color patterns'
        'suntime:next sunrise or sunset time'
        'sunscene:visual sunrise or sunset color art'
        'pollen:pollen levels by type'
        'uv:uv index with color scale'
        'rain:precipitation level and intensity'
        'season:current astronomical season with color'
        'holiday:today'"'"'s public holiday by location'
        'satellites:notable satellites currently overhead'
        'tearoff:tear-off calendar showing today'"'"'s date'
        'daemon:run scheduled jobs continuously'
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
complete -c note -n __fish_use_subcommand -a countdown  -d 'Days until configured event'
complete -c note -n __fish_use_subcommand -a discogs    -d 'Record matched to weather and time'
complete -c note -n __fish_use_subcommand -a pattern    -d 'Random art and color patterns'
complete -c note -n __fish_use_subcommand -a suntime     -d 'Next sunrise or sunset time'
complete -c note -n __fish_use_subcommand -a sunscene   -d 'Next sunrise or sunset color art'
complete -c note -n __fish_use_subcommand -a pollen     -d 'Pollen levels by type'
complete -c note -n __fish_use_subcommand -a uv         -d 'UV index with color scale'
complete -c note -n __fish_use_subcommand -a rain       -d 'Precipitation level and intensity'
complete -c note -n __fish_use_subcommand -a season     -d 'Current astronomical season with color'
complete -c note -n __fish_use_subcommand -a holiday    -d 'Today'"'"'s public holiday by location'
complete -c note -n __fish_use_subcommand -a satellites -d 'Notable satellites currently overhead'
complete -c note -n __fish_use_subcommand -a tearoff    -d 'Tear-off calendar showing today'"'"'s date'
complete -c note -n __fish_use_subcommand -a daemon     -d 'Run scheduled jobs continuously'
complete -c note -n __fish_use_subcommand -a status     -d 'Preview all subcommands without sending'
complete -c note -n __fish_use_subcommand -a completion -d 'Print shell completion script'
complete -c note -n '__fish_seen_subcommand_from pattern' -a 'current random stripes checker bars fade diagonal hearts confetti sparkle pulse rainbow'
complete -c note -n '__fish_seen_subcommand_from completion' -a 'bash zsh fish'`

func logInvocation(cmd string, cfg config) {
	source := "manual"
	if fi, err := os.Stdin.Stat(); err == nil && (fi.Mode()&os.ModeCharDevice) == 0 {
		source = "cron"
	}
	host, _ := os.Hostname()
	fmt.Fprintf(os.Stderr, "note %s %s %s@%s lat=%.2f lon=%.2f\n",
		version, cmd, source, host, cfg.Lat, cfg.Lon)
}

var defaultSchedule = []scheduleEntry{
	{Command: "weather", Hour: 8, Minute: 0, RepeatMinutes: 30, UntilHour: 19},
	{Command: "uv", Hour: 9, Minute: 0, Args: []string{"--skip-trivial"}},
	{Command: "pollen", Hour: 11, Minute: 0, Args: []string{"--skip-trivial"}},
	{Command: "rain", Hour: 12, Minute: 0, Args: []string{"--skip-trivial"}},
	{Command: "calendar", Hour: 13, Minute: 0},
	{Command: "season", Hour: 15, Minute: 0},
	{Command: "holiday", Hour: 16, Minute: 0},
	{Command: "air", Hour: 17, Minute: 0, Args: []string{"--skip-trivial"}},
	{Command: "suntime", Hour: 18, Minute: 0},
	{Command: "satellites", Hour: 19, Minute: 0, Args: []string{"--skip-trivial"}},
	{Command: "holiday", Hour: 20, Minute: 0},
	{Command: "moonphase", Hour: 21, Minute: 0},
	{Command: "season", Hour: 22, Minute: 0},
	{Command: "tearoff", Hour: 23, Minute: 0},
	{Command: "tearoff", Hour: 0, Minute: 0},
}

func isDue(entry scheduleEntry, now time.Time) bool {
	h, m := now.Hour(), now.Minute()
	if entry.RepeatMinutes > 0 {
		if h < entry.Hour || h > entry.UntilHour {
			return false
		}
		return m%entry.RepeatMinutes == 0
	}
	return h == entry.Hour && m == entry.Minute
}

func runDaemon(cfg config) {
	loc := time.Local
	if cfg.Timezone != "" {
		if l, err := time.LoadLocation(cfg.Timezone); err == nil {
			loc = l
		} else {
			log.Printf("daemon: unknown timezone %q, falling back to local: %v", cfg.Timezone, err)
		}
	}

	schedule := cfg.Schedule
	if len(schedule) == 0 {
		schedule = defaultSchedule
	}

	log.Printf("daemon: starting (timezone=%s, jobs=%d)", loc, len(schedule))

	fired := make(map[string]time.Time)

	for {
		now := time.Now().In(loc)
		key := now.Format("2006-01-02 15:04")

		for _, entry := range schedule {
			if !isDue(entry, now) {
				continue
			}
			jobKey := key + " " + entry.Command + " " + strings.Join(entry.Args, " ")
			if _, done := fired[jobKey]; done {
				continue
			}
			fired[jobKey] = now

			go func(e scheduleEntry) {
				log.Printf("daemon: running %s %s", e.Command, strings.Join(e.Args, " "))
				runCommand(e.Command, e.Args, cfg)
			}(entry)
		}

		cleanupFired(fired, now)

		next := now.Truncate(time.Minute).Add(time.Minute)
		time.Sleep(time.Until(next))
	}
}

func cleanupFired(fired map[string]time.Time, now time.Time) {
	cutoff := now.Add(-2 * time.Minute)
	for k, t := range fired {
		if t.Before(cutoff) {
			delete(fired, k)
		}
	}
}

func runCommand(cmd string, args []string, cfg config) {
	skipTrivial := false
	for _, a := range args {
		if a == "--skip-trivial" {
			skipTrivial = true
		}
	}

	var lines [3]string
	var trivial bool
	var err error

	switch cmd {
	case "weather":
		lines, err = weather.Fetch(cfg.Lat, cfg.Lon)
	case "clock":
		lines = clock.Format(time.Now())
	case "calendar":
		lines, err = calendar.Fetch(cfg.ICalURLs)
	case "moonphase":
		lines = moonphase.Format(time.Now())
	case "air":
		lines, trivial, err = air.Fetch(cfg.Lat, cfg.Lon)
	case "flights":
		lines, trivial, err = flights.Fetch(cfg.Lat, cfg.Lon)
	case "countdown":
		lines = countdown.Format(cfg.Countdowns, time.Now())
	case "discogs":
		lines, err = discogs.Fetch(cfg.DiscogsUsername, cfg.DiscogsToken, cfg.Lat, cfg.Lon)
	case "pattern":
		lines, err = pattern.Current(cfg.Lat, cfg.Lon)
	case "suntime":
		lines, err = suntime.Fetch(cfg.Lat, cfg.Lon)
	case "sunscene":
		lines, err = sunscene.Fetch(cfg.Lat, cfg.Lon)
	case "pollen":
		lines, trivial, err = pollen.Fetch(cfg.Lat, cfg.Lon)
	case "uv":
		lines, err = uv.Fetch(cfg.Lat, cfg.Lon)
	case "rain":
		lines, trivial, err = rain.Fetch(cfg.Lat, cfg.Lon)
	case "season":
		lines = season.Format(time.Now())
	case "tearoff":
		lines = tearoff.Format(time.Now())
	case "holiday":
		lines, trivial, err = holiday.Fetch(cfg.Lat, cfg.Lon)
	case "satellites":
		lines, trivial, err = satellites.Fetch(cfg.Lat, cfg.Lon)
	default:
		log.Printf("daemon: unknown command %q", cmd)
		return
	}

	if err != nil {
		log.Printf("daemon: %s error: %v", cmd, err)
		return
	}
	if trivial && skipTrivial {
		log.Printf("daemon: %s skipped (trivial)", cmd)
		return
	}

	client := vestaboard.New(cfg.Token)
	if err := client.SendLines(lines); err != nil {
		log.Printf("daemon: %s send error: %v", cmd, err)
		return
	}
	log.Printf("daemon: %s sent", cmd)
}

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

	logInvocation(cmd, cfg)

	if cmd == "status" {
		runStatus(cfg)
		return
	}

	if cmd == "daemon" {
		runDaemon(cfg)
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
		lines, trivial, err = flights.Fetch(cfg.Lat, cfg.Lon)
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
