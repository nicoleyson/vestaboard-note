# vestaboard-note

CLI for the [Vestaboard Note](https://www.vestaboard.com/note) (3×15 split-flap display).

## Setup

**1. Get your API token**
Open the Vestaboard mobile app → Settings → Developer → create a token.

**2. Get your Google Calendar ICS URL**
Google Calendar → gear icon → Settings → your calendar → "Secret address in iCal format". Copy the `https://` URL.

**3. Configure**
Copy `config.yaml` and fill in your values:

```yaml
token: YOUR_TOKEN_HERE
lat: 37.7749      # your latitude
lon: -122.4194    # your longitude
ical_url: https://calendar.google.com/calendar/ical/...
```

Set `VESTABOARD_CONFIG=/path/to/config.yaml` if you keep it elsewhere.

**4. Build**
```sh
make build
```

## Usage

```sh
./note weather    # current conditions via Open-Meteo (free, no key)
./note clock      # date + time
./note calendar   # next event from your Google Calendar
```

## Cron example

```cron
# weather every 30 minutes
*/30 * * * * /usr/local/bin/note weather

# clock on the hour
0 * * * * /usr/local/bin/note clock
```

Rate limit: 1 message per 15 seconds (enforced client-side).
