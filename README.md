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

## Subcommands

```sh
./note weather     # current conditions (METAR, real observations, global)
./note clock       # date + time
./note calendar    # next event from your iCal feed
./note moon        # current lunar phase
./note air         # air quality index
./note flights     # aircraft currently overhead
./note onthisday   # a historical event from today's date
./note countdown   # days until a configured event
./note discogs     # a record from your Discogs collection matched to weather + time of day
```

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
