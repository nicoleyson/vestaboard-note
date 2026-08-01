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

| Command | What it shows | Requires |
|---|---|---|
| `weather` | Current conditions — temp, sky, color row | `lat`, `lon` |
| `clock` | Date and time | — |
| `calendar` | Next upcoming event | `ical_urls` |
| `moon` | Lunar phase with illumination color split | — |
| `air` | Air quality index with color row | `lat`, `lon` |
| `flights` | Aircraft currently overhead | `lat`, `lon` |
| `onthisday` | A historical event from today's date | — |
| `countdown` | Days until a configured event | `countdowns` |
| `discogs` | A record from your collection matched to weather + time of day | `discogs_username`, `discogs_token`, `lat`, `lon` |
| `pattern` | Random color art — stripes, checker, hearts, confetti, and more | — |
| `sunrise` | Next sunrise or sunset time with color row | `lat`, `lon` |
| `pollen` | Pollen level and dominant type (grass/tree/weed) with color row | `lat`, `lon` |
| `status` | Preview all subcommands without sending to the board | — |

## Recommended crontab

Run `crontab -e` and add entries like these. Adjust the binary path to wherever you built `note`.

```cron
# Weather every 30 minutes
*/30 * * * * cd ~/repos/vestaboard-note && ./note weather

# Calendar at the top of every hour
0 * * * * cd ~/repos/vestaboard-note && ./note calendar

# Clock every 15 minutes when you just want the time
*/15 * * * * cd ~/repos/vestaboard-note && ./note clock

# Moon phase once a day at 8am
0 8 * * * cd ~/repos/vestaboard-note && ./note moon

# On this day, once a day at 9am
0 9 * * * cd ~/repos/vestaboard-note && ./note onthisday

# Discogs vibe check at the start of each time slot
0 6,9,12,17,20,23 * * * cd ~/repos/vestaboard-note && ./note discogs

# Air quality every hour during fire season (June–October)
0 * * 6-10 * cd ~/repos/vestaboard-note && ./note air

# Pollen levels every morning during spring (March–June)
0 7 * 3-6 * cd ~/repos/vestaboard-note && ./note pollen

# Sunrise/sunset — show what's coming up next, once an hour
0 * * * * cd ~/repos/vestaboard-note && ./note sunrise

# Random pattern — good as a screensaver when nothing else is scheduled
0 2 * * * cd ~/repos/vestaboard-note && ./note pattern
```

A simple "set and forget" setup that covers most of the day:

```cron
# Weather every 30 min, calendar on the hour, clock in between
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
