# AGENTS.md — go-fred

FRED API client for Go. 34 endpoints. Stdlib only. No dependencies.

## Commands

```bash
go build ./...          # always builds clean (no deps)
go vet ./...            # always clean  
go test ./...           # all 38 tests, needs $FRED_API_KEY
go test ./... -run TestGetSeriesInfo   # run a single test
```

Tests hit the live FRED API. If `$FRED_API_KEY` is absent, they `t.Skip()` cleanly.

## Architecture

```
fred.go         Client, New(), WithAPIKey/WithHTTPClient/WithBaseURL
errors.go       type APIError = internal.APIError  (type alias, use errors.As)
types.go        All response types and raw wire types (rawObservation, etc.)
options.go      Option types: ObservationOption, SearchOption, TagOption, etc.
series.go       Series + vintage + series-tree methods
search.go       SearchSeries, GetReleaseSeries, GetCategorySeries + pagination
categories.go   Category endpoints
releases.go     Release endpoints
sources.go      Source endpoints
tags.go         Tag endpoints
geofred.go      GeoFRED map endpoints
internal/http.go  DoRequest — all HTTP I/O, api_key injection, error parsing
```

Every public method takes `ctx context.Context` as first argument.

## Gotchas

### FRED API quirks

- **`"seriess"` (double s)** — FRED serializes series arrays under the key `"seriess"` not `"series"`.
- **Missing values are `"."`** — a literal dot string, not `null` or `""`. `parseObservation` sets `IsNA=true, Value=0.0`.
- **Vintage dates limit is 2000** — the full realtime range (1776–9999) exceeds the FRED API limit. Tests pass explicit `WithRealtimeStart` to narrow the window.
- **Release tables** — the `elements` key in `/release/tables` is a JSON object (map keyed by element_id), not an array.

### GeoFRED deviations from docs

The live API responses differ from the [FRED docs](https://fred.stlouisfed.org/docs/api/geofred/).

- **`/series/data`** — the `data` field is `map[string][]MapDataEntry` (dates as keys), not a flat array. `GetSeriesData` flattens it.
- **`/regional/data`** — same `{"meta": {...}}` wrapper structure as `/series/data`, not the documented dynamic-key format. Parsed identically.
- **`value` field is numeric** — FRED returns numbers (e.g. `35559`) not strings for GeoFRED. `MapDataEntry` has a custom `UnmarshalJSON` to handle both.

### Pagination

- Search/list endpoints (`SearchSeries`, `GetReleaseSeries`, `GetCategorySeries`) auto-paginate via `paginatedSearch`. Page size is 1000. `limit=0` means unlimited.
- Observation endpoints do NOT paginate (FRED returns up to 100,000 in one response).

### Error handling

- `internal.DoRequest` returns `*internal.APIError` on non-200 responses.
- `fred.APIError` is a type alias (`type APIError = internal.APIError`). Callers can use `errors.As(err, &apiErr)`.
- URL construction errors, network errors, and JSON decode errors use `fmt.Errorf` with `%w`.

### No rate limiting

FRED allows 120 req/min. The client does not enforce this — callers manage their own rate.

### Conventions

- Tests in `fred_test` package (external test package).
- Test functions named `Test<MethodName>`.
- Known test series: `DGS20` (20-year Treasury), `SMU56000000500000001A` (GeoFRED), `WIPCPI` (GeoFRED data), release `53` (GDP), category `1` (Production & Business Activity), source `1`, tag `"gdp"`/`"business"`.
