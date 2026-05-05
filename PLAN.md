# go-fred — Plan

FRED API client for Go. Covers all 34 FRED API endpoints. Stdlib only. No dependencies.

## Data Format

FRED supports JSON via `file_type=json`. All endpoints are verified with curl. No XML.

Example JSON response for `/series/observations`:

```json
{
  "realtime_start": "2024-01-01",
  "realtime_end": "2024-12-31",
  "observation_start": "2014-01-01",
  "observation_end": "2014-12-01",
  "units": "lin",
  "output_type": 1,
  "file_type": "json",
  "order_by": "observation_date",
  "sort_order": "asc",
  "count": 119,
  "offset": 0,
  "limit": 100000,
  "observations": [
    { "date": "2014-01-01", "value": "233.916" },
    { "date": "2014-02-01", "value": "." }
  ]
}
```

FRED uses `"."` (a single dot) for missing values. The JSON value is the string `"."`.

---

## Package Layout

```
go-fred/
├── fred.go              # Client struct, New(), With* options, package doc
├── fred_test.go
├── series.go            # GetSeriesObservations, GetSeriesInfo, vintage methods
├── series_test.go
├── search.go            # SearchSeries, GetReleaseSeries, GetCategorySeries, pagination
├── search_test.go
├── types.go             # Observation, Series, SearchResult, Category, Release, Source, Tag
├── options.go           # ObservationOption, SearchOption, TagOption, and all option types
├── errors.go            # APIError type
├── internal/
│   └── http.go          # doRequest, buildURL, response handling
├── example_test.go      # Runnable examples (godoc)
├── testdata/            # JSON fixtures for unit tests
├── PLAN.md
├── README.md
├── LICENSE
└── go.mod              # module github.com/shanehull/go-fred
```

Zero dependencies beyond Go stdlib.

---

## 1. Client Initialization

```go
package fred

// Client is a FRED API client. Must be constructed via New.
// Safe for concurrent use by multiple goroutines.
//
// API key resolution order:
//   1. WithAPIKey("...")
//   2. $FRED_API_KEY environment variable
type Client struct {
    apiKey     string
    httpClient *http.Client
    baseURL    string       // https://api.stlouisfed.org
}

func New(opts ...ClientOption) (*Client, error)

type ClientOption func(*Client) error

func WithAPIKey(key string) ClientOption
func WithHTTPClient(hc *http.Client) ClientOption    // for testing
func WithBaseURL(url string) ClientOption             // for testing
```

### Proxy Support

`http.DefaultTransport` picks up `HTTP_PROXY`/`HTTPS_PROXY` automatically via
`http.ProxyFromEnvironment`. No extra code needed.

### Go Module

```
module github.com/shanehull/go-fred
go 1.25.4
```

Import: `import "github.com/shanehull/go-fred"`
Package: `fred`
Usage: `client, err := fred.New()` → `client.GetSeriesObservations(...)`

---

## 2. Core Types

### 2.1 Observation

```go
type Observation struct {
    Date          time.Time
    Value         float64
    RealtimeStart time.Time
    RealtimeEnd   time.Time
    IsNA          bool
}

// rawObservation is the JSON wire format (all fields are strings from FRED).
type rawObservation struct {
    Date          string `json:"date"`
    Value         string `json:"value"`
    RealtimeStart string `json:"realtime_start"`
    RealtimeEnd   string `json:"realtime_end"`
}
```

Missing values (`"."`) → `IsNA = true, Value = 0.0`.

### 2.2 ObservationResponse

```go
type observationResponse struct {
    RealtimeStart    string           `json:"realtime_start"`
    RealtimeEnd      string           `json:"realtime_end"`
    ObservationStart string           `json:"observation_start"`
    ObservationEnd   string           `json:"observation_end"`
    Units            string           `json:"units"`
    OutputType       int              `json:"output_type"`
    FileType         string           `json:"file_type"`
    OrderBy          string           `json:"order_by"`
    SortOrder        string           `json:"sort_order"`
    Count            int              `json:"count"`
    Offset           int              `json:"offset"`
    Limit            int              `json:"limit"`
    Observations     []rawObservation `json:"observations"`
}
```

### 2.3 Series (metadata, from `/series`)

```go
// FRED uses "seriess" (double 's') — not a typo.
type seriesInfoResponse struct {
    Series []Series `json:"seriess"`
}

type Series struct {
    ID                      string `json:"id"`
    RealtimeStart           string `json:"realtime_start"`
    RealtimeEnd             string `json:"realtime_end"`
    Title                   string `json:"title"`
    ObservationStart        string `json:"observation_start"`
    ObservationEnd          string `json:"observation_end"`
    Frequency               string `json:"frequency"`
    FrequencyShort          string `json:"frequency_short"`
    Units                   string `json:"units"`
    UnitsShort              string `json:"units_short"`
    SeasonalAdjustment      string `json:"seasonal_adjustment"`
    SeasonalAdjustmentShort string `json:"seasonal_adjustment_short"`
    LastUpdated             string `json:"last_updated"`
    Popularity              int    `json:"popularity"`
    Notes                   string `json:"notes,omitempty"`
    Tags                    []Tag  `json:"tags,omitempty"`
}
```

### 2.4 SearchResult

```go
// FRED uses "seriess" for the array key in search/release/category series responses.
type searchResponse struct {
    Count  int            `json:"count"`
    Series []SearchResult `json:"seriess"`
}

type SearchResult struct {
    ID                      string `json:"id"`
    RealtimeStart           string `json:"realtime_start"`
    RealtimeEnd             string `json:"realtime_end"`
    Title                   string `json:"title"`
    ObservationStart        string `json:"observation_start"`
    ObservationEnd          string `json:"observation_end"`
    Frequency               string `json:"frequency"`
    FrequencyShort          string `json:"frequency_short"`
    Units                   string `json:"units"`
    UnitsShort              string `json:"units_short"`
    SeasonalAdjustment      string `json:"seasonal_adjustment"`
    SeasonalAdjustmentShort string `json:"seasonal_adjustment_short"`
    LastUpdated             string `json:"last_updated"`
    Popularity              int    `json:"popularity"`
    Notes                   string `json:"notes,omitempty"`
}
```

### 2.5 Category

```go
type Category struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    ParentID int    `json:"parent_id"`
    Notes    string `json:"notes,omitempty"`
}
```

### 2.6 Release

```go
type Release struct {
    ID           int    `json:"id"`
    Name         string `json:"name"`
    PressRelease bool   `json:"press_release"`
    Link         string `json:"link,omitempty"`
    Notes        string `json:"notes,omitempty"`
}

type ReleaseDate struct {
    ReleaseID   int    `json:"release_id"`
    ReleaseName string `json:"release_name,omitempty"`
    Date        string `json:"date"`
}
```

### 2.7 Source

```go
type Source struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Link  string `json:"link,omitempty"`
    Notes string `json:"notes,omitempty"`
}
```

### 2.8 Tag

```go
type Tag struct {
    Name        string `json:"name"`
    GroupID     string `json:"group_id"`
    Notes       string `json:"notes,omitempty"`
    Created     string `json:"created"`
    Popularity  int    `json:"popularity"`
    SeriesCount int    `json:"series_count"`
}
```

### 2.9 VintageDates

```go
type vintageDatesResponse struct {
    Dates []string `json:"vintage_dates"`
}
```

### 2.10 ReleaseTable

```go
type releaseTableResponse struct {
    Elements []ReleaseTableElement `json:"elements"`
}

type ReleaseTableElement struct {
    ElementID int                   `json:"element_id"`
    ReleaseID int                   `json:"release_id"`
    SeriesID  string                `json:"series_id,omitempty"`
    ParentID  int                   `json:"parent_id"`
    Line      string                `json:"line"`
    Type      string                `json:"type"`
    Name      string                `json:"name"`
    Level     int                   `json:"level"`
    Children  []ReleaseTableElement `json:"children,omitempty"`
}
```

### 2.11 APIError

```go
// FRED returns: {"error_code": 400, "error_message": "Bad Request..."}
type APIError struct {
    Code    int    `json:"error_code"`
    Message string `json:"error_message"`
}

func (e *APIError) Error() string {
    return fmt.Sprintf("fred: API error %d: %s", e.Code, e.Message)
}
```

### 2.12 GeoFRED Types

```go
// SeriesGroup represents metadata for a GeoFRED map series.
type SeriesGroup struct {
    Title      string `json:"title"`
    RegionType string `json:"region_type"`
    SeriesGroup string `json:"series_group"`
    Season     string `json:"season"`
    Units      string `json:"units"`
    Frequency  string `json:"frequency"`
    MinDate    string `json:"min_date"`
    MaxDate    string `json:"max_date"`
}

// MapDataEntry represents a single region's data point.
type MapDataEntry struct {
    Region   string `json:"region"`
    Code     string `json:"code"`
    Value    string `json:"value"`
    SeriesID string `json:"series_id"`
}

// MapData wraps GeoFRED series/regional data responses.
type MapData struct {
    Title       string         `json:"title"`
    Region      string         `json:"region"`
    Seasonality string         `json:"seasonality"`
    Units       string         `json:"units"`
    Frequency   string         `json:"frequency"`
    Date        string         `json:"date"`
    Data        []MapDataEntry `json:"data"`
}
```

---

## 3. Rate Limiting

FRED allows 120 requests per minute. The client does not enforce this — callers
are expected to manage their own rate if needed. A future addition could be an
optional rate limiter via `WithRateLimiter(rate.Limiter)`.

---

### 3.1 GetSeriesInfo

```go
// GetSeriesInfo returns metadata for a FRED series.
// Endpoint: GET /fred/series?series_id={id}
func (c *Client) GetSeriesInfo(ctx context.Context, seriesID string) (*Series, error)
```

### 3.2 GetSeriesObservations

```go
// GetSeriesObservations returns data values for a series.
// Endpoint: GET /fred/series/observations?series_id={id}
func (c *Client) GetSeriesObservations(ctx context.Context, seriesID string, opts ...ObservationOption) ([]Observation, error)
```

### 3.3 GetSeriesAllReleases

```go
// GetSeriesAllReleases returns all observations including every revision.
// Defaults realtime_start="1776-07-04", realtime_end="9999-12-31".
// Endpoint: GET /fred/series/observations (with realtime params)
func (c *Client) GetSeriesAllReleases(ctx context.Context, seriesID string, opts ...ObservationOption) ([]Observation, error)
```

### 3.4 GetSeriesFirstRelease

```go
// GetSeriesFirstRelease returns the first published value for each date.
// Calls GetSeriesAllReleases and groups by Date.
func (c *Client) GetSeriesFirstRelease(ctx context.Context, seriesID string, opts ...ObservationOption) ([]Observation, error)
```

### 3.5 GetSeriesAsOf

```go
// GetSeriesAsOf returns data as known on a specific date.
// Calls GetSeriesAllReleases with realtime_end=asOf.
func (c *Client) GetSeriesAsOf(ctx context.Context, seriesID string, asOf time.Time, opts ...ObservationOption) ([]Observation, error)
```

### 3.6 GetSeriesVintageDates

```go
// GetSeriesVintageDates returns all vintage dates for a series.
// Endpoint: GET /fred/series/vintagedates?series_id={id}
func (c *Client) GetSeriesVintageDates(ctx context.Context, seriesID string, opts ...ObservationOption) ([]time.Time, error)
```

### 3.7 SearchSeries

```go
// SearchSeries searches for series by keywords.
// Auto-paginates. Limit 0 = unlimited. Default limit = 1000.
// Endpoint: GET /fred/series/search?search_text={text}
func (c *Client) SearchSeries(ctx context.Context, text string, opts ...SearchOption) ([]SearchResult, error)
```

### 3.8 GetReleaseSeries

```go
// GetReleaseSeries returns series belonging to a release.
// Auto-paginates. Limit 0 = unlimited. Default limit = 0.
// Endpoint: GET /fred/release/series?release_id={id}
func (c *Client) GetReleaseSeries(ctx context.Context, releaseID int, opts ...SearchOption) ([]SearchResult, error)
```

### 3.9 GetCategorySeries

```go
// GetCategorySeries returns series belonging to a category.
// Auto-paginates. Limit 0 = unlimited. Default limit = 0.
// Endpoint: GET /fred/category/series?category_id={id}
func (c *Client) GetCategorySeries(ctx context.Context, categoryID int, opts ...SearchOption) ([]SearchResult, error)
```

---

## 4. Option Types

### 4.1 ObservationOption

```go
type ObservationOption func(*obsParams)

type obsParams struct {
    observationStart  string
    observationEnd    string
    realtimeStart     string
    realtimeEnd       string
    units             string    // lin, chg, ch1, pch, pc1, pca, cch, cca, log
    frequency         string    // d, w, bw, m, q, sa, a, wef, weth, wew, wetu, wem, wesu, wesa, bwew, bwem
    aggregationMethod string    // avg, sum, eop
    outputType        int       // 1-4
    vintageDates      []string
    sortOrder         SortOrder
    limit             int
    offset            int
}

func WithObservationStart(t time.Time) ObservationOption
func WithObservationEnd(t time.Time) ObservationOption
func WithRealtimeStart(t time.Time) ObservationOption
func WithRealtimeEnd(t time.Time) ObservationOption
func WithUnits(u string) ObservationOption
func WithFrequency(f string) ObservationOption
func WithAggregationMethod(m string) ObservationOption
func WithOutputType(t int) ObservationOption
func WithVintageDates(dates ...string) ObservationOption
func WithObservationSortOrder(o SortOrder) ObservationOption
func WithObservationLimit(n int) ObservationOption
func WithObservationOffset(n int) ObservationOption
```

### 4.2 SearchOption

```go
type SearchOption func(*searchParams)

type searchParams struct {
    searchType      string    // full_text, series_id
    orderBy         OrderBy
    sortOrder       SortOrder
    filterVariable  string    // frequency, units, seasonal_adjustment
    filterValue     string
    limit           int
    tagNames        []string
    excludeTagNames []string
}

func WithSearchType(t string) SearchOption
func WithSortOrder(o SortOrder) SearchOption
func WithOrderBy(ob OrderBy) SearchOption
func WithFilter(variable, value string) SearchOption
func WithLimit(n int) SearchOption
func WithTagNames(tags ...string) SearchOption
func WithExcludeTags(tags ...string) SearchOption
```

### 4.3 Enums

```go
type SortOrder string
const (
    SortAsc  SortOrder = "asc"
    SortDesc SortOrder = "desc"
)

type OrderBy string
const (
    OrderBySearchRank         OrderBy = "search_rank"
    OrderBySeriesID           OrderBy = "series_id"
    OrderByTitle              OrderBy = "title"
    OrderByUnits              OrderBy = "units"
    OrderByFrequency          OrderBy = "frequency"
    OrderBySeasonalAdjustment OrderBy = "seasonal_adjustment"
    OrderByRealtimeStart      OrderBy = "realtime_start"
    OrderByRealtimeEnd        OrderBy = "realtime_end"
    OrderByLastUpdated        OrderBy = "last_updated"
    OrderByObservationStart   OrderBy = "observation_start"
    OrderByObservationEnd     OrderBy = "observation_end"
    OrderByPopularity         OrderBy = "popularity"
)
```

---

## 5. Pagination

FRED returns max 1000 results per request. Search/list endpoints auto-paginate.

```
1. First request without offset → get page 0 and total count
2. Determine max results:
   - limit == 0 → max = total_count (get EVERYTHING)
   - otherwise → max = limit
3. If max > 1000:
   for offset := 1000; offset < max; offset += 1000:
       fetch URL + "&offset={offset}"
       append results
4. Truncate to max if needed
```

GetSeriesObservations does NOT paginate — it returns all results in one response (FRED's limit is 100,000).

---

## 6. HTTP Layer (`internal/http.go`)

```go
package internal

// doRequest sends a GET request and returns the response body.
// Callers pass the path and params. The layer handles baseURL, api_key, file_type=json.
func doRequest(ctx context.Context, client *http.Client, baseURL, path string, params url.Values) ([]byte, error)
```

Example call from `GetSeriesObservations`:

```go
params := url.Values{}
params.Set("series_id", seriesID)
params.Set("sort_order", "desc")
return doRequest(ctx, c.httpClient, c.baseURL, "/fred/series/observations", params)
```

The layer:

1. Prepends `baseURL` to `path` to form the full URL
2. Appends `api_key` and `file_type=json` to params
3. Encodes params and sends GET with context support
4. On non-200: reads body, unmarshals `APIError`, returns it
5. On 200: returns body bytes

---

## 7. Value Parsing

```go
func parseObservation(raw rawObservation) (Observation, error) {
    obs := Observation{}
    if raw.Date != "" {
        t, err := time.Parse("2006-01-02", raw.Date)
        if err != nil { return obs, fmt.Errorf("fred: invalid date %q: %w", raw.Date, err) }
        obs.Date = t
    }
    if raw.Value != "" && raw.Value != "." {
        v, err := strconv.ParseFloat(raw.Value, 64)
        if err != nil { return obs, fmt.Errorf("fred: invalid value %q: %w", raw.Value, err) }
        obs.Value = v
    } else {
        obs.IsNA = true
    }
    return obs, nil
}
```

---

## 8. Go Conventions

- Context on every method
- Error wrapping with `%w`
- Never panic
- Thread-safe (no global state, no init)
- Go 1.25

---

## 9. Complete API Coverage — All 34 FRED Endpoints

### 9.1 Series (10 endpoints)

| 1 | `/series` | `GetSeriesInfo(ctx, id string) (*Series, error)` |
| 2 | `/series/observations` | `GetSeriesObservations(ctx, id string, opts ...ObservationOption) ([]Observation, error)` |
| 3 | `/series/categories` | `GetSeriesCategories(ctx, id string) ([]Category, error)` |
| 4 | `/series/release` | `GetSeriesRelease(ctx, id string) (*Release, error)` |
| 5 | `/series/search` | `SearchSeries(ctx, text string, opts ...SearchOption) ([]SearchResult, error)` |
| 6 | `/series/search/tags` | `SearchSeriesTags(ctx, text string, opts ...TagOption) ([]Tag, error)` |
| 7 | `/series/search/related_tags` | `SearchSeriesRelatedTags(ctx, text string, tagNames []string, opts ...TagOption) ([]Tag, error)` |
| 8 | `/series/tags` | `GetSeriesTags(ctx, id string, opts ...TagOption) ([]Tag, error)` |
| 9 | `/series/updates` | `GetSeriesUpdates(ctx, opts ...UpdateOption) ([]SearchResult, error)` |
| 10 | `/series/vintagedates` | `GetSeriesVintageDates(ctx, id string, opts ...ObservationOption) ([]time.Time, error)` |

```go
GetSeriesAllReleases(ctx, id, opts...)  // observations with full realtime range
GetSeriesFirstRelease(ctx, id, opts...) // groups by Date, earliest revision
GetSeriesAsOf(ctx, id, asOf, opts...)   // observations as known on a date
```

### 9.2 Categories (6 endpoints)

| 11 | `/category` | `GetCategory(ctx, id int) (*Category, error)` |
| 12 | `/category/children` | `GetCategoryChildren(ctx, id int) ([]Category, error)` |
| 13 | `/category/related` | `GetCategoryRelated(ctx, id int) ([]Category, error)` |
| 14 | `/category/series` | `GetCategorySeries(ctx, id int, opts ...SearchOption) ([]SearchResult, error)` |
| 15 | `/category/tags` | `GetCategoryTags(ctx, id int, opts ...TagOption) ([]Tag, error)` |
| 16 | `/category/related_tags` | `GetCategoryRelatedTags(ctx, id int, tagNames []string, opts ...TagOption) ([]Tag, error)` |

### 9.3 Releases (9 endpoints)

| 17 | `/releases` | `GetReleases(ctx, opts ...ReleaseListOption) ([]Release, error)` |
| 18 | `/releases/dates` | `GetReleasesDates(ctx, opts ...ReleaseDateOption) ([]ReleaseDate, error)` |
| 19 | `/release` | `GetRelease(ctx, id int) (*Release, error)` |
| 20 | `/release/dates` | `GetReleaseDates(ctx, id int, opts ...ReleaseDateOption) ([]ReleaseDate, error)` |
| 21 | `/release/series` | `GetReleaseSeries(ctx, id int, opts ...SearchOption) ([]SearchResult, error)` |
| 22 | `/release/sources` | `GetReleaseSources(ctx, id int) ([]Source, error)` |
| 23 | `/release/tags` | `GetReleaseTags(ctx, id int, opts ...TagOption) ([]Tag, error)` |
| 24 | `/release/related_tags` | `GetReleaseRelatedTags(ctx, id int, tagNames []string, opts ...TagOption) ([]Tag, error)` |
| 25 | `/release/tables` | `GetReleaseTables(ctx, id int, opts ...TableOption) ([]ReleaseTableElement, error)` |

### 9.4 Sources (3 endpoints)

| 26 | `/sources` | `GetSources(ctx, opts ...SourceOption) ([]Source, error)` |
| 27 | `/source` | `GetSource(ctx, id int) (*Source, error)` |
| 28 | `/source/releases` | `GetSourceReleases(ctx, id int, opts ...ReleaseListOption) ([]Release, error)` |

### 9.5 Tags (3 endpoints)

| 29 | `/tags` | `GetTags(ctx, opts ...TagOption) ([]Tag, error)` |
| 30 | `/related_tags` | `GetRelatedTags(ctx, tagNames []string, opts ...TagOption) ([]Tag, error)` |
| 31 | `/tags/series` | `GetTagsSeries(ctx, tagNames []string, opts ...SearchOption) ([]SearchResult, error)` |

### 9.6 GeoFRED Maps (3 endpoints)

| 32 | `/geofred/series/group` | `GetSeriesGroup(ctx, id string) (*SeriesGroup, error)` |
| 33 | `/geofred/series/data` | `GetSeriesData(ctx, id string, opts ...MapDataOption) (*MapData, error)` |
| 34 | `/geofred/regional/data` | `GetRegionalData(ctx, opts ...RegionalDataOption) (*MapData, error)` |

### 10.1 TagOption

```go
type TagOption func(*tagParams)

func WithTagGroupID(g string) TagOption
func WithTagSearchText(text string) TagOption
func WithTagLimit(n int) TagOption
func WithTagNames(tags ...string) TagOption
func WithExcludeTagNames(tags ...string) TagOption
```

### 10.2 ReleaseListOption

```go
type ReleaseListOption func(*releaseListParams)

func WithReleaseLimit(n int) ReleaseListOption
func WithReleaseSortOrder(o SortOrder) ReleaseListOption
```

### 10.3 ReleaseDateOption

```go
type ReleaseDateOption func(*releaseDateParams)

func WithReleaseDateLimit(n int) ReleaseDateOption
func WithReleaseDateSortOrder(o SortOrder) ReleaseDateOption
func WithIncludeNoData(b bool) ReleaseDateOption
```

### 10.4 TableOption

```go
func WithTableElementID(id int) TableOption
func WithIncludeObservationValues(b bool) TableOption
func WithObservationDate(d string) TableOption
```

### 10.5 UpdateOption

```go
func WithStartTime(s string) UpdateOption
func WithEndTime(s string) UpdateOption
func WithFilterValue(v string) UpdateOption   // macro, regional, all
func WithUpdateLimit(n int) UpdateOption
```

### 10.6 SourceOption

```go
func WithSourceLimit(n int) SourceOption
func WithSourceSortOrder(o SortOrder) SourceOption
```

### 10.7 MapDataOption

```go
func WithMapDate(d string) MapDataOption
func WithMapStartDate(d string) MapDataOption
```

### 10.8 RegionalDataOption

```go
func WithSeriesGroup(g string) RegionalDataOption
func WithRegionType(t string) RegionalDataOption     // bea, msa, frb, necta, state, country, county, censusregion
func WithMapDate(d string) RegionalDataOption
func WithSeason(s string) RegionalDataOption         // SA, NSA, SSA, SAAR, NSAAR
func WithMapUnits(u string) RegionalDataOption
func WithTransformation(t string) RegionalDataOption // lin, chg, ch1, pch, pc1, pca, cch, cca, log
```

---

## 11. Testing Strategy

Tests are written alongside implementation — one endpoint, one test file.
Tests hit the live FRED API using `FRED_API_KEY` from the environment.
If the key is absent, tests `t.Skip()` cleanly rather than failing.

```go
// series_test.go
func TestGetSeriesInfo(t *testing.T) {
    if os.Getenv("FRED_API_KEY") == "" {
        t.Skip("FRED_API_KEY not set")
    }
    client, err := fred.New()
    if err != nil {
        t.Fatal(err)
    }
    ctx := context.Background()
    s, err := client.GetSeriesInfo(ctx, "DGS20")
    if err != nil {
        t.Fatal(err)
    }
    if s.ID != "DGS20" {
        t.Errorf("expected DGS20, got %s", s.ID)
    }
}
```

```go
func ExampleClient_GetSeriesObservations() {
    client, _ := fred.New(fred.WithAPIKey("your-key"))
    obs, _ := client.GetSeriesObservations(context.Background(), "DGS20",
        fred.WithObservationStart(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
    )
    for _, o := range obs {
        fmt.Printf("%s: %.2f\n", o.Date.Format("2006-01-02"), o.Value)
    }
}
```

---

## Implementation Order

| Step              | What                                                                                                                                   | Effort |
| ----------------- | -------------------------------------------------------------------------------------------------------------------------------------- | ------ |
| 1. Scaffold       | go.mod, fred.go (Client, New, options), errors.go, types.go                                                                            | Small  |
| 2. HTTP layer     | internal/http.go                                                                                                                       | Medium |
| 3. Series         | GetSeriesObservations, GetSeriesInfo, observations parsing                                                                             | Medium |
| 4. Vintage        | GetSeriesAllReleases, GetSeriesFirstRelease, GetSeriesAsOf, GetSeriesVintageDates                                                      | Medium |
| 5. Search         | SearchSeries, GetReleaseSeries, GetCategorySeries, pagination                                                                          | Medium |
| 6. Categories     | GetCategory, GetCategoryChildren, GetCategoryRelated, GetCategorySeries, GetCategoryTags, GetCategoryRelatedTags                       | Medium |
| 7. Series tree    | GetSeriesCategories, GetSeriesRelease, GetSeriesTags, SearchSeriesTags, SearchSeriesRelatedTags, GetSeriesUpdates                      | Medium |
| 8. Releases       | GetReleases, GetReleasesDates, GetRelease, GetReleaseDates, GetReleaseSources, GetReleaseTags, GetReleaseRelatedTags, GetReleaseTables | Medium |
| 9. Sources + Tags | GetSources, GetSource, GetSourceReleases, GetTags, GetRelatedTags, GetTagsSeries                                                       | Small  |
| 10. GeoFRED       | GetSeriesGroup, GetSeriesData, GetRegionalData                                                                                         | Small  |

Each step includes its own test file.

## CI

```yaml
# .github/workflows/lint.yaml
name: Lint

on:
  pull_request:
    branches:
      - main

permissions:
  contents: read

jobs:
  Lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
        with:
          go-version: "1.25"
          cache: false
      - name: Lint
        uses: golangci/golangci-lint-action@v9
        with:
          version: v2.6

  Test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
        with:
          go-version: "1.25"
          cache: false
      - run: go test ./...
```

```yaml
# .github/workflows/release.yaml
name: Release

on:
  push:
    branches:
      - main

permissions:
  contents: write
  pull-requests: write

jobs:
  Release:
    runs-on: ubuntu-latest
    steps:
      - name: Release Please
        uses: googleapis/release-please-action@v5
        with:
          token: ${{secrets.RELEASES_GITHUB_TOKEN}}
```

