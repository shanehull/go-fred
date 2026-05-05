# go-fred

<p align="center">
  <img src="https://go.dev/blog/go-brand/Go-Logo/PNG/Go-Logo_LightBlue.png" height="80" alt="Go">
  <br>
  <img src="assets/fred-logo.svg" height="48" alt="FRED">
</p>

[![Go Reference](https://pkg.go.dev/badge/github.com/shanehull/go-fred.svg)](https://pkg.go.dev/github.com/shanehull/go-fred)
[![Go Report Card](https://goreportcard.com/badge/github.com/shanehull/go-fred)](https://goreportcard.com/report/github.com/shanehull/go-fred)
[![Go](https://img.shields.io/github/go-mod/go-version/shanehull/go-fred)](https://go.dev/dl/)
[![CI](https://github.com/shanehull/go-fred/actions/workflows/test.yaml/badge.svg)](https://github.com/shanehull/go-fred/actions/workflows/test.yaml)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Go client for the [FRED API](https://fred.stlouisfed.org/docs/api/fred/). Stdlib only — zero dependencies.

## Install

```sh
go get github.com/shanehull/go-fred
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/shanehull/go-fred"
)

func main() {
    client, _ := fred.New() // uses $FRED_API_KEY

    obs, _ := client.GetSeriesObservations(context.Background(), "DGS20",
        fred.WithObservationStart(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
        fred.WithObservationSortOrder(fred.SortDesc),
    )
    for _, o := range obs {
        if !o.IsNA {
            fmt.Printf("%s  %.2f%%\n", o.Date.Format("2006-01-02"), o.Value)
        }
    }
}
```

## API Key

Get a free key at [fred.stlouisfed.org/docs/api/api_key.html](https://fred.stlouisfed.org/docs/api/api_key.html).

```go
// Environment variable (recommended)
client, _ := fred.New()

// Explicit option (overrides env)
client, _ := fred.New(fred.WithAPIKey("your-key"))
```

## Usage

### Observations

```go
obs, err := client.GetSeriesObservations(ctx, "DGS20",
    fred.WithObservationStart(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
    fred.WithObservationSortOrder(fred.SortDesc),
    fred.WithUnits("lin"),
    fred.WithFrequency("m"),
)
```

### Vintage Data

```go
// All releases (every revision)
all, _ := client.GetSeriesAllReleases(ctx, "DGS20",
    fred.WithRealtimeStart(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
)

// First published value for each date
first, _ := client.GetSeriesFirstRelease(ctx, "DGS20")

// As known on a specific date
asof, _ := client.GetSeriesAsOf(ctx, "DGS20", time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC))
```

### Search

```go
// Search auto-paginates. Set limit=0 for all results.
results, _ := client.SearchSeries(ctx, "real GDP",
    fred.WithLimit(50),
    fred.WithOrderBy(fred.OrderByPopularity),
)

// Filter by release or category
series, _ := client.GetReleaseSeries(ctx, 53, fred.WithLimit(100))
series, _ := client.GetCategorySeries(ctx, 1, fred.WithLimit(100))
```

### Categories, Releases, Sources, Tags

```go
cat, _  := client.GetCategory(ctx, 1)
children, _ := client.GetCategoryChildren(ctx, 0)
rel, _  := client.GetRelease(ctx, 53)
src, _  := client.GetSource(ctx, 1)
tags, _ := client.GetTags(ctx, fred.WithTagLimit(20))
```

### GeoFRED

```go
group, _ := client.GetSeriesGroup(ctx, "SMU56000000500000001A")
data, _  := client.GetSeriesData(ctx, "WIPCPI",
    fred.WithMapDate("2012-01-01"),
)
regional, _ := client.GetRegionalData(ctx,
    fred.WithSeriesGroup("882"),
    fred.WithRegionType("state"),
    fred.WithRegionalDate("2013-01-01"),
    fred.WithMapUnits("Dollars"),
    fred.WithSeason("NSA"),
    fred.WithRegionalFrequency("a"),
)
```

### Error Handling

```go
obs, err := client.GetSeriesObservations(ctx, "INVALID")
var apiErr fred.APIError
if errors.As(err, &apiErr) {
    fmt.Printf("FRED API error %d: %s\n", apiErr.Code, apiErr.Message)
}
```

## API Coverage

| Domain     | Methods                                                                                                                                                                                                          |
| ---------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Series     | `GetSeriesInfo`, `GetSeriesObservations`, `GetSeriesCategories`, `GetSeriesRelease`, `GetSeriesTags`, `GetSeriesUpdates`, `GetSeriesVintageDates`, `SearchSeries`, `SearchSeriesTags`, `SearchSeriesRelatedTags` |
| Vintage    | `GetSeriesAllReleases`, `GetSeriesFirstRelease`, `GetSeriesAsOf`                                                                                                                                                 |
| Categories | `GetCategory`, `GetCategoryChildren`, `GetCategoryRelated`, `GetCategorySeries`, `GetCategoryTags`, `GetCategoryRelatedTags`                                                                                     |
| Releases   | `GetReleases`, `GetReleasesDates`, `GetRelease`, `GetReleaseDates`, `GetReleaseSeries`, `GetReleaseSources`, `GetReleaseTags`, `GetReleaseRelatedTags`, `GetReleaseTables`                                       |
| Sources    | `GetSources`, `GetSource`, `GetSourceReleases`                                                                                                                                                                   |
| Tags       | `GetTags`, `GetRelatedTags`, `GetTagsSeries`                                                                                                                                                                     |
| GeoFRED    | `GetSeriesGroup`, `GetSeriesData`, `GetRegionalData`                                                                                                                                                             |

## Examples

See [examples/](examples/) for runnable programs:

- [observations](examples/series/main.go) — fetch and display series data
- [search](examples/search/main.go) — search for series
- [vintage](examples/vintage/main.go) — compare first vs latest release values
- [categories](examples/categories/main.go) — explore the category tree
- [releases](examples/releases/main.go) — release info and dates
- [tags](examples/tags/main.go) — popular tags and related tags
- [geofred](examples/geofred/main.go) — map metadata and regional data
- [sources](examples/sources/main.go) — data sources and their releases
