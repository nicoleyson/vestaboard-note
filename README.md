# vestaboard-note

CLI for the [Vestaboard Note](https://www.vestaboard.com/note) (3×15 split-flap display).

## Setup

**1. Get your API token**
Open the Vestaboard mobile app → Settings → Developer → create a token.

**2. Configure**
```sh
cp config.yaml.example config.yaml
```

Fill in your values. At minimum you need `token`, `lat`, and `lon`. Everything else is optional depending on which subcommands you use.

Set `VESTABOARD_CONFIG=/path/to/config.yaml` to use a config file in a non-default location.

**3. Build**
```sh
make build
```

## Install

```sh
make install        # builds and copies note to /usr/local/bin
make uninstall      # removes it
```

`make install` prints suggested crontab entries and shell completion instructions when it finishes.

## Shell completion

```sh
note completion bash >> ~/.bash_profile   # bash
note completion zsh  >> ~/.zshrc          # zsh
note completion fish > ~/.config/fish/completions/note.fish  # fish
```

Reload your shell after adding.

## Subcommands

| Command | What it shows | Timing | Requires |
|---|---|---|---|
| `weather` | Current conditions — temp, sky, color row | Real observed data from nearest airport station, updated every ~30–60 min | `lat`, `lon` |
| `clock` | Date and time | Point in time | — |
| `calendar` | Next upcoming event | Point in time | `ical_urls` |
| `moon` | Lunar phase with illumination color split | Point in time | — |
| `air` | Air quality index with color row | Point in time (hourly model) | `lat`, `lon` |
| `flights` | Aircraft currently overhead | Point in time | `lat`, `lon` |
| `onthisday` | A historical event from today's date | Today's date | — |
| `countdown` | Days until a configured event | Point in time | `countdowns` |
| `discogs` | A record from your collection matched to weather + time of day | Point in time | `discogs_username`, `discogs_token`, `lat`, `lon` |
| `pattern` | Random color art — stripes, checker, hearts, confetti, and more | — | — |
| `sunrise` | Next sunrise or sunset time with color row | Calculated astronomical time | `lat`, `lon` |
| `pollen` | Pollen level and dominant type (grass/tree/weed) with color row | Point in time (hourly model) | `lat`, `lon` |
| `uv` | UV index with category color row (green→yellow→orange→red→violet) | Point in time (hourly model) | `lat`, `lon` |
| `rain` | Precipitation probability and intensity (none/light/moderate/heavy) with color row | Point in time (hourly model) | `lat`, `lon` |
| `season` | Current astronomical season with color row and phase label | Calculated from equinox/solstice math | — |
| `status` | Preview all subcommands without sending to the board | — | — |

## Recommended crontab

Run `crontab -e` and add entries like these. Replace `/Users/yourname` with your actual home directory path.

On macOS, `~` and `TZ=` are not reliably supported in crontab — use full absolute paths and omit `TZ=` (macOS cron inherits your system timezone automatically).

```cron
# Rain, pollen, and air quality skip posting when levels are trivial if --skip-trivial is passed.
# Run manually without the flag to always send (good for testing).

# ── Morning routine ──────────────────────────────────────────────
5 8 * * * cd /Users/yourname/repos/nicoleyson/vestaboard-note && ./note onthisday
10 8 * * * cd /Users/yourname/repos/nicoleyson/vestaboard-note && ./note countdown

# ── Throughout the day ───────────────────────────────────────────
# Weather every 30 minutes
*/30 * * * * cd /Users/yourname/repos/nicoleyson/vestaboard-note && ./note weather

# Calendar at the top of every hour
0 * * * * cd /Users/yourname/repos/nicoleyson/vestaboard-note && ./note calendar

# Sunrise/sunset during waking hours (offset to avoid :00 collision)
2 8-22 * * * cd /Users/yourname/repos/nicoleyson/vestaboard-note && ./note sunrise

# UV index during peak hours (offset to avoid collision)
4 10-15 * * * cd /Users/yourname/repos/nicoleyson/vestaboard-note && ./note uv

# Rain check every morning — skips when chance is negligible
10 8 * * * cd /Users/yourname/repos/nicoleyson/vestaboard-note && ./note rain --skip-trivial

# Air quality every hour — skips when air is good
6 * * * * cd /Users/yourname/repos/nicoleyson/vestaboard-note && ./note air --skip-trivial

# Pollen every morning — skips when levels are low
20 8 * * * cd /Users/yourname/repos/nicoleyson/vestaboard-note && ./note pollen --skip-trivial

# ── Evening ──────────────────────────────────────────────────────
# Moon phase at 9pm — best appreciated at night
0 21 * * * cd /Users/yourname/repos/nicoleyson/vestaboard-note && ./note moon

# ── Screensaver ──────────────────────────────────────────────────
0 0 * * * cd /Users/yourname/repos/nicoleyson/vestaboard-note && ./note pattern
```

If you want the clock cycling constantly (good for a display you glance at all day), use a high-frequency setup instead:

```cron
* * * * * cd /Users/yourname/repos/nicoleyson/vestaboard-note && ./note clock
*/30 * * * * cd /Users/yourname/repos/nicoleyson/vestaboard-note && ./note weather
0 * * * * cd /Users/yourname/repos/nicoleyson/vestaboard-note && ./note calendar
```

The Vestaboard rate limit is 1 message per 15 seconds, so `clock` every minute is the fastest practical cadence.

## Config reference

```yaml
token: YOUR_TOKEN_HERE      # required for all subcommands
lat: 45.5051                # your latitude
lon: -122.6750              # your longitude

# calendar: one or more iCal secret URLs
# Google Calendar → gear → Settings → your calendar → Secret address in iCal format
ical_urls:
  - https://calendar.google.com/calendar/ical/you%40gmail.com/private-abc123/basic.ics

# countdown: named dates to count down to
countdowns:
  - label: VACATION
    date: "2025-07-04T00:00:00Z"

# discogs: your collection username and API token (discogs.com → Settings → Developers)
discogs_username: yourname
discogs_token: YOUR_DISCOGS_TOKEN
```
