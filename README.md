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
| `moonphase` | Lunar phase with illumination color split | Point in time | — |
| `air` | Air quality index with color row | Point in time (hourly model) | `lat`, `lon` |
| `flights` | Aircraft currently overhead | Point in time | `lat`, `lon` |
| `onthisday` | A historical event from today's date | Today's date | — |
| `countdown` | Days until a configured event | Point in time | `countdowns` |
| `discogs` | A record from your collection matched to weather + time of day | Point in time | `discogs_username`, `discogs_token`, `lat`, `lon` |
| `pattern` | Color art — stripes, checker, hearts, confetti, and more. Use `pattern current` for seasonal/holiday-aware palette | — | — |
| `suntime` | Next sunrise or sunset time with color row | Calculated astronomical time | `lat`, `lon` |
| `sunscene` | Visual color art scene — sunrise (cool sky, warm burst) or sunset (warm sky, deep afterglow) | Calculated astronomical time | `lat`, `lon` |
| `pollen` | Pollen level and dominant type (grass/tree/weed) with color row | Point in time (hourly model) | `lat`, `lon` |
| `uv` | UV index with category color row (green→yellow→orange→red→violet) | Point in time (hourly model) | `lat`, `lon` |
| `rain` | Precipitation probability and intensity (none/light/moderate/heavy) with color row | Point in time (hourly model) | `lat`, `lon` |
| `season` | Current astronomical season with color row and phase label | Calculated from equinox/solstice math | — |
| `holiday` | Today's public holiday by location via Nager.Date (200+ countries) | Real data | `lat`, `lon` |
| `satellites` | Notable satellite currently overhead (ISS, GPS, Iridium) with elevation and direction | Real-time from Satlas | `lat`, `lon` |
| `tearoff` | Tear-off calendar showing today's date — red card with yellow tabs and white number | Point in time | — |
| `status` | Preview all subcommands without sending to the board | — | — |

## Crontab

The repo includes a `crontab.template` with the recommended schedule. Install it with:

```sh
make cron          # smart: adds on first run, updates on subsequent runs
make cron-init     # first-time only — aborts if already installed
make cron-update   # always replace with current template
make cron-uninstall
```

By default `VESTABOARD_DIR` is set to the current directory and `NOTE_BIN` to `$(VESTABOARD_DIR)/note`. Override either:

```sh
make cron VESTABOARD_DIR=/custom/path
make cron NOTE_BIN=/usr/local/bin/note   # if you ran make install
```

**macOS**: entries are merged into your user crontab (via `crontab -e`) inside a `# BEGIN VESTABOARD` / `# END VESTABOARD` block. `make cron-update` replaces just that block, leaving the rest of your crontab untouched.

**Linux**: entries are written to `/etc/cron.d/vestaboard` (requires `sudo`). The username column is added automatically.

The template runs weather every 30 minutes as a backbone, with one-shot subcommands sprinkled through the morning and evening. To see the full schedule, read `crontab.template` directly.

If you want the clock cycling constantly (good for a display you glance at all day), edit `crontab.template` and replace the weather backbone with:

```cron
* * * * *    cd $VESTABOARD_DIR && ./note clock
*/15 * * * * cd $VESTABOARD_DIR && ./note sunscene
*/30 * * * * cd $VESTABOARD_DIR && ./note weather
0 * * * *    cd $VESTABOARD_DIR && ./note calendar
```

The Vestaboard rate limit is 1 message per 15 seconds, so `clock` every minute is the fastest practical cadence. `sunscene` every 15 minutes is enough to show the sun visibly moving through rows as the day progresses.

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
