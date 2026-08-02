# vestaboard-note

A Go CLI for the [Vestaboard Note](https://www.vestaboard.com/note) — the compact 3×15 split-flap display. Point it at your location, add your API token, and a cron job takes care of the rest.

Twenty subcommands. No subscriptions. Almost everything is keyless.

---

## Quick start

**1. Get your API token**
Open the Vestaboard mobile app → Settings → Developer → create a token.

**2. Clone and configure**
```sh
git clone https://github.com/nicoleyson/vestaboard-note.git
cd vestaboard-note
cp config.yaml.example config.yaml
```
Fill in `token`, `lat`, and `lon`. Everything else is optional.

**3. Build and try it**
```sh
make build
./note weather
./note status        # preview all subcommands without sending
```

**4. Install and schedule**
```sh
make install         # copies binary to /usr/local/bin
make cron            # installs the recommended crontab schedule
```

---

## Subcommands

| Command | What it shows | Requires |
|---|---|---|
| `weather` | Temp, conditions, color row — real observed data from nearest airport METAR station | `lat`, `lon` |
| `clock` | Day, date, and time | — |
| `calendar` | Next upcoming event from your calendar | `ical_urls` |
| `moonphase` | Lunar phase name, illumination split across color tiles | — |
| `season` | Current astronomical season with phase label (IS HERE / IN BLOOM / WINDING DOWN…) | — |
| `suntime` | Next sunrise or sunset with color row | `lat`, `lon` |
| `sunscene` | Animated sky scene — sun arc and sky color shift through the day | `lat`, `lon` |
| `air` | US AQI with color row (green → hazardous) | `lat`, `lon` |
| `uv` | UV index with color scale | `lat`, `lon` |
| `rain` | Precipitation level and intensity | `lat`, `lon` |
| `pollen` | Pollen level and dominant type (grass / tree / weed) | `lat`, `lon` |
| `flights` | Aircraft currently overhead — callsign, origin, altitude | `lat`, `lon` |
| `satellites` | Notable satellite overhead — ISS, GPS, Iridium — with elevation and direction | `lat`, `lon` |
| `holiday` | Today's public holiday by location (200+ countries via Nager.Date) | `lat`, `lon` |
| `onthisday` | A historical event from today's date via Wikipedia | — |
| `countdown` | Days until a configured event | `countdowns` in config |
| `discogs` | A record from your collection matched to current weather and time of day | `discogs_username`, `discogs_token`, `lat`, `lon` |
| `tearoff` | Tear-off calendar showing today's date | — |
| `pattern` | Color art — stripes, checker, hearts, confetti, rainbow, and more. `pattern current` picks a seasonal or holiday-aware palette | — |
| `status` | Preview all subcommands without sending anything to the board | — |

### Trivial-skip flag

`rain`, `air`, `pollen`, `flights`, `holiday`, and `satellites` support `--skip-trivial`. When the data is unremarkable (no rain, good air quality, clear skies, no holiday today), the subcommand exits without sending. Useful for cron jobs that should only interrupt when there's something worth showing.

```sh
./note rain --skip-trivial      # only sends if it's actually raining or likely to
./note holiday --skip-trivial   # only sends on actual holidays
```

---

## Scheduling with cron

The repo includes `crontab.template` with a recommended schedule. Install it with:

```sh
make cron            # smart: adds on first run, updates on subsequent runs
make cron-update     # always replace with current template
make cron-uninstall  # remove the vestaboard block from your crontab
```

You can override the install path:
```sh
make cron NOTE_BIN=/usr/local/bin/note
make cron VESTABOARD_DIR=/custom/path
```

**macOS**: entries are merged into your user crontab inside `# BEGIN VESTABOARD` / `# END VESTABOARD` markers. The rest of your crontab is left untouched.

**Linux**: entries are written to `/etc/cron.d/vestaboard` (requires `sudo`).

The default schedule runs `weather` every 30 minutes as a backdrop — the board always returns to current conditions within half an hour after any one-shot subcommand. One-shots fire at the top of each hour throughout the day. Read `crontab.template` directly to see the full schedule and adjust it.

**Clock mode** — if you want the board cycling all day as a live clock:
```cron
* * * * *    cd $VESTABOARD_DIR && ./note clock
*/15 * * * * cd $VESTABOARD_DIR && ./note sunscene
*/30 * * * * cd $VESTABOARD_DIR && ./note weather
0 * * * *    cd $VESTABOARD_DIR && ./note calendar
```

The Vestaboard rate limit is 1 message per 15 seconds. `clock` every minute is the fastest practical cadence without hitting it.

---

## Shell completion

```sh
note completion bash >> ~/.bash_profile
note completion zsh  >> ~/.zshrc
note completion fish > ~/.config/fish/completions/note.fish
```

Reload your shell. Tab-completes subcommand names and `pattern` variants.

---

## Config reference

```yaml
token: YOUR_TOKEN_HERE      # required — from Vestaboard app → Settings → Developer
lat: 45.5051                # your latitude
lon: -122.6750              # your longitude

# calendar: one or more iCal secret URLs
# Google Calendar → gear → Settings → your calendar → Secret address in iCal format
ical_urls:
  - https://calendar.google.com/calendar/ical/you%40gmail.com/private-abc123/basic.ics

# countdown: named dates to count down (or back) to
countdowns:
  - label: VACATION
    date: "2027-07-04T00:00:00Z"

# discogs: your collection username and API token
# discogs.com → Settings → Developers → Generate token
discogs_username: yourname
discogs_token: YOUR_DISCOGS_TOKEN
```

Set `VESTABOARD_CONFIG=/path/to/config.yaml` to use a config file outside the default location.

---

## APIs used

Everything is free. Most require no account.

| Subcommand | Source |
|---|---|
| weather | [aviationweather.gov](https://aviationweather.gov/) METAR (no key) |
| air, uv, rain, pollen | [Open-Meteo](https://open-meteo.com/) (no key) |
| suntime, sunscene | [sunrise-sunset.org](https://sunrise-sunset.org/api) (no key) |
| flights | [OpenSky Network](https://opensky-network.org/) (no key) |
| satellites | [Satlas](https://satlas.app/) (no key) |
| holiday | [Nager.Date](https://date.nager.at/) + [Nominatim](https://nominatim.org/) (no key) |
| onthisday | [Wikipedia REST API](https://en.wikipedia.org/api/rest_v1/) (no key) |
| calendar | iCal secret URL (your Google/iCloud/Outlook calendar) |
| discogs | [Discogs API](https://www.discogs.com/developers) (free token) |
| moonphase, season, countdown, tearoff, clock, pattern | local math only |
