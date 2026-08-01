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
| `status` | Preview all subcommands without sending to the board | — | — |

## Recommended crontab

Run `crontab -e` and add entries like these. Adjust the binary path to wherever you built `note`.

Times are in your system timezone. On macOS this matches your local time automatically. On Linux servers that run in UTC, add `TZ=America/Los_Angeles` (or your timezone) as the first line of your crontab.

```cron
# ── Morning routine ──────────────────────────────────────────────
# Moon phase at 8am — changes slowly, once a day is plenty
0 8 * * * cd ~/repos/vestaboard-note && ./note moon

# On this day at 8:05am — a fun read with your coffee
5 8 * * * cd ~/repos/vestaboard-note && ./note onthisday

# Rain check at 8:10am — is it worth grabbing an umbrella?
10 8 * * * cd ~/repos/vestaboard-note && ./note rain

# Countdown at 8:15am — good for trips, deadlines, anything you're anticipating
15 8 * * * cd ~/repos/vestaboard-note && ./note countdown

# ── Throughout the day ───────────────────────────────────────────
# Weather every 30 minutes
*/30 * * * * cd ~/repos/vestaboard-note && ./note weather

# Calendar at the top of every hour — what's coming up next
0 * * * * cd ~/repos/vestaboard-note && ./note calendar

# Clock at :15 and :45 — fills the gap between weather and calendar
15,45 * * * * cd ~/repos/vestaboard-note && ./note clock

# Sunrise/sunset during waking hours — shows whichever is coming up next
0 8-22 * * * cd ~/repos/vestaboard-note && ./note sunrise

# UV index during peak hours — only meaningful when the sun is up
0 10-15 * * * cd ~/repos/vestaboard-note && ./note uv

# Discogs vibe check at each time-of-day transition
0 9,12,17,20,23 * * * cd ~/repos/vestaboard-note && ./note discogs

# ── Seasonal ─────────────────────────────────────────────────────
# Air quality every hour during fire season (June–October)
0 * * 6-10 * cd ~/repos/vestaboard-note && ./note air

# Pollen levels every morning during spring (March–June)
0 8 * 3-6 * cd ~/repos/vestaboard-note && ./note pollen

# ── Screensaver ──────────────────────────────────────────────────
# Random pattern at midnight — something pretty while you sleep
0 0 * * * cd ~/repos/vestaboard-note && ./note pattern
```

A minimal "set and forget" setup if you just want the basics:

```cron
*/30 * * * * cd ~/repos/vestaboard-note && ./note weather
0 * * * * cd ~/repos/vestaboard-note && ./note calendar
15,45 * * * * cd ~/repos/vestaboard-note && ./note clock
```

Rate limit: 1 message per 15 seconds (enforced client-side).

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
