# COSMOS

A web application for astronomy. Enter your location and discover what is worth
observing in the sky today or over the next few days.

Not just a data page: a real software project with a backend, a database,
external API integration, in-house calculations, and a considered interface.

## What it does

Enter a location, and COSMOS shows:

- Current date
- The selected location
- Upcoming astronomical events
- Highlighted events of the day
- What can be observed, and when
- The ability to save favorite events

Example:

> Today - September 27
> Tandil, Buenos Aires
>
> The Moon rises at 19:42
> ISS visible at 20:15
> Jupiter visible during the early morning hours
> A meteor shower coming up
>
> [View details] [Save]

## Planned modules

| Module | Purpose |
| --- | --- |
| Astronomical events | Moon phases, sun and moon rise/set, visible planets, eclipses, meteor showers, conjunctions, oppositions, constellations, daily highlights. Queryable for today, tomorrow, the next 7 days, or a given date. |
| Satellites | ISS and other known satellites: current position, next visible passes over your location, start time, peak time, end time, direction, max elevation, duration. Built on TLE data and orbital propagation. |
| What can I see today? | Calculates the most interesting observable objects for a given location, ordered by time. |
| Personal agenda | Save, list, and remove events. Mark important ones. Later: reminders, notifications, calendar export. |
| Earth Monitor | Explore satellite imagery of a region: current image, compare across dates, temporal timeline, NDVI, change detection. |
| Exoplanet detector | Load a light curve, visualize it, detect periodic brightness dips, identify transit candidates, estimate orbital period, duration, and depth. |

## Stack

- **Backend**: Go, REST API, JSON
- **Database**: SQLite
- **Frontend**: HTML, CSS, and vanilla JavaScript
- **Libraries (as needed)**: Skyfield for orbital calculations, Leaflet for maps, Plotly or Chart.js for charts, TLE and orbital data sources

## Architecture

```text
                    ┌──────────────────────┐
                    │       FRONTEND       │
                    │   HTML / CSS / JS    │
                    └──────────┬───────────┘
                               │
                            HTTP/JSON
                               │
                    ┌──────────▼───────────┐
                    │       BACKEND        │
                    │          Go          │
                    │      REST API        │
                    └──────┬───────┬───────┘
                           │       │
                 ┌─────────▼─┐   ┌▼──────────────┐
                 │  SQLite   │   │ External APIs │
                 │           │   │ Astronomy     │
                 │  Events   │   │ Satellites    │
                 │  Saved    │   │ Earth data    │
                 └───────────┘   └───────────────┘
```

## Status

Early. The repository is initialized and the foundation is being built up
module by module.

## License

To be decided.
