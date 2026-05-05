package fred

import "time"

// SortOrder specifies ascending or descending sort.
type SortOrder string

const (
	SortAsc  SortOrder = "asc"
	SortDesc SortOrder = "desc"
)

// OrderBy specifies the field to sort search results by.
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

// ObservationOption configures a call to GetSeriesObservations or related vintage methods.
type ObservationOption func(*obsParams)

type obsParams struct {
	observationStart  string
	observationEnd    string
	realtimeStart     string
	realtimeEnd       string
	units             string
	frequency         string
	aggregationMethod string
	outputType        int
	vintageDates      []string
	sortOrder         SortOrder
	limit             int
	offset            int
}

// WithObservationStart sets the start of the observation period (FRED: observation_start).
func WithObservationStart(t time.Time) ObservationOption {
	return func(p *obsParams) { p.observationStart = t.Format("2006-01-02") }
}

// WithObservationEnd sets the end of the observation period (FRED: observation_end).
func WithObservationEnd(t time.Time) ObservationOption {
	return func(p *obsParams) { p.observationEnd = t.Format("2006-01-02") }
}

// WithRealtimeStart sets the start of the real-time period (FRED: realtime_start).
func WithRealtimeStart(t time.Time) ObservationOption {
	return func(p *obsParams) { p.realtimeStart = t.Format("2006-01-02") }
}

// WithRealtimeEnd sets the end of the real-time period (FRED: realtime_end).
func WithRealtimeEnd(t time.Time) ObservationOption {
	return func(p *obsParams) { p.realtimeEnd = t.Format("2006-01-02") }
}

// WithUnits sets the data value transformation (FRED: units).
// Values: lin, chg, ch1, pch, pc1, pca, cch, cca, log.
func WithUnits(u string) ObservationOption {
	return func(p *obsParams) { p.units = u }
}

// WithFrequency sets the data frequency aggregation (FRED: frequency).
// Values: d, w, bw, m, q, sa, a, wef, weth, wew, wetu, wem, wesu, wesa, bwew, bwem.
func WithFrequency(f string) ObservationOption {
	return func(p *obsParams) { p.frequency = f }
}

// WithAggregationMethod sets the frequency aggregation method (FRED: aggregation_method).
// Values: avg, sum, eop.
func WithAggregationMethod(m string) ObservationOption {
	return func(p *obsParams) { p.aggregationMethod = m }
}

// WithOutputType sets the output type (FRED: output_type). Values: 1-4.
func WithOutputType(t int) ObservationOption {
	return func(p *obsParams) { p.outputType = t }
}

// WithVintageDates sets vintage dates to request (FRED: vintage_dates).
func WithVintageDates(dates ...string) ObservationOption {
	return func(p *obsParams) { p.vintageDates = dates }
}

// WithObservationSortOrder sets the sort order for observations (FRED: sort_order).
func WithObservationSortOrder(o SortOrder) ObservationOption {
	return func(p *obsParams) { p.sortOrder = o }
}

// WithObservationLimit sets the maximum number of observations to return (FRED: limit).
func WithObservationLimit(n int) ObservationOption {
	return func(p *obsParams) { p.limit = n }
}

// WithObservationOffset sets the offset for paginated observation requests (FRED: offset).
func WithObservationOffset(n int) ObservationOption {
	return func(p *obsParams) { p.offset = n }
}

// SearchOption configures a search or series-list call.
type SearchOption func(*searchParams)

type searchParams struct {
	searchType      string
	orderBy         OrderBy
	sortOrder       SortOrder
	filterVariable  string
	filterValue     string
	limit           int
	tagNames        []string
	excludeTagNames []string
}

// WithSearchType sets the search type (FRED: search_type).
// Values: full_text, series_id.
func WithSearchType(t string) SearchOption {
	return func(p *searchParams) { p.searchType = t }
}

// WithSortOrder sets the result sort order for search/series results.
func WithSortOrder(o SortOrder) SearchOption {
	return func(p *searchParams) { p.sortOrder = o }
}

// WithOrderBy sets the field to sort search results by.
func WithOrderBy(ob OrderBy) SearchOption {
	return func(p *searchParams) { p.orderBy = ob }
}

// WithFilter filters results by variable and value (FRED: filter_variable, filter_value).
func WithFilter(variable, value string) SearchOption {
	return func(p *searchParams) {
		p.filterVariable = variable
		p.filterValue = value
	}
}

// WithLimit sets the maximum number of results to return. 0 means unlimited.
func WithLimit(n int) SearchOption {
	return func(p *searchParams) { p.limit = n }
}

// WithTagNames filters results to series matching the given tags (FRED: tag_names).
func WithTagNames(tags ...string) SearchOption {
	return func(p *searchParams) { p.tagNames = tags }
}

// WithExcludeTags excludes series matching the given tags (FRED: exclude_tag_names).
func WithExcludeTags(tags ...string) SearchOption {
	return func(p *searchParams) { p.excludeTagNames = tags }
}

// TagOption configures a tag-related call.
type TagOption func(*tagParams)

type tagParams struct {
	groupID         string
	searchText      string
	limit           int
	tagNames        []string
	excludeTagNames []string
	sortOrder       SortOrder
	orderBy         OrderBy
}

// WithTagGroupID filters tags by group ID (FRED: tag_group_id).
func WithTagGroupID(g string) TagOption {
	return func(p *tagParams) { p.groupID = g }
}

// WithTagSearchText filters tags by search text (FRED: search_text).
func WithTagSearchText(text string) TagOption {
	return func(p *tagParams) { p.searchText = text }
}

// WithTagLimit sets the maximum number of tags to return (FRED: limit).
func WithTagLimit(n int) TagOption {
	return func(p *tagParams) { p.limit = n }
}

// WithTagSortOrder sets the sort order for tag results (FRED: sort_order).
func WithTagSortOrder(o SortOrder) TagOption {
	return func(p *tagParams) { p.sortOrder = o }
}

// WithTagOrderBy sets the field to sort tags by (FRED: order_by).
func WithTagOrderBy(ob OrderBy) TagOption {
	return func(p *tagParams) { p.orderBy = ob }
}

// WithTagSetNames filters results by tag names (FRED: tag_names).
func WithTagSetNames(tags ...string) TagOption {
	return func(p *tagParams) { p.tagNames = tags }
}

// WithTagSetExclude excludes results by tag names (FRED: exclude_tag_names).
func WithTagSetExclude(tags ...string) TagOption {
	return func(p *tagParams) { p.excludeTagNames = tags }
}

// ReleaseListOption configures release listing calls.
type ReleaseListOption func(*releaseListParams)

type releaseListParams struct {
	limit     int
	sortOrder SortOrder
}

// WithReleaseLimit sets the maximum number of releases to return (FRED: limit).
func WithReleaseLimit(n int) ReleaseListOption {
	return func(p *releaseListParams) { p.limit = n }
}

// WithReleaseSortOrder sets the sort order for release results (FRED: sort_order).
func WithReleaseSortOrder(o SortOrder) ReleaseListOption {
	return func(p *releaseListParams) { p.sortOrder = o }
}

// ReleaseDateOption configures release date calls.
type ReleaseDateOption func(*releaseDateParams)

type releaseDateParams struct {
	limit         int
	sortOrder     SortOrder
	includeNoData bool
}

// WithReleaseDateLimit sets the maximum number of release dates (FRED: limit).
func WithReleaseDateLimit(n int) ReleaseDateOption {
	return func(p *releaseDateParams) { p.limit = n }
}

// WithReleaseDateSortOrder sets the sort order for release dates (FRED: sort_order).
func WithReleaseDateSortOrder(o SortOrder) ReleaseDateOption {
	return func(p *releaseDateParams) { p.sortOrder = o }
}

// WithIncludeNoData includes releases with no data (FRED: include_release_dates_with_no_data).
func WithIncludeNoData(b bool) ReleaseDateOption {
	return func(p *releaseDateParams) { p.includeNoData = b }
}

// TableOption configures release table calls.
type TableOption func(*tableParams)

type tableParams struct {
	elementID                int
	includeObservationValues bool
	observationDate          string
}

// WithTableElementID filters tables by element ID (FRED: element_id).
func WithTableElementID(id int) TableOption {
	return func(p *tableParams) { p.elementID = id }
}

// WithIncludeObservationValues includes observation values in table rows (FRED: include_observation_values).
func WithIncludeObservationValues(b bool) TableOption {
	return func(p *tableParams) { p.includeObservationValues = b }
}

// WithObservationDate sets the observation date for table values (FRED: observation_date).
func WithObservationDate(d string) TableOption {
	return func(p *tableParams) { p.observationDate = d }
}

// UpdateOption configures series updates calls.
type UpdateOption func(*updateParams)

type updateParams struct {
	startTime   string
	endTime     string
	filterValue string
	limit       int
}

// WithStartTime filters updates by start time (FRED: start_time).
func WithStartTime(s string) UpdateOption {
	return func(p *updateParams) { p.startTime = s }
}

// WithEndTime filters updates by end time (FRED: end_time).
func WithEndTime(s string) UpdateOption {
	return func(p *updateParams) { p.endTime = s }
}

// WithFilterValue filters updates by category (FRED: filter_value).
// Values: macro, regional, all.
func WithFilterValue(v string) UpdateOption {
	return func(p *updateParams) { p.filterValue = v }
}

// WithUpdateLimit sets the maximum number of updated series (FRED: limit).
func WithUpdateLimit(n int) UpdateOption {
	return func(p *updateParams) { p.limit = n }
}

// SourceOption configures source calls.
type SourceOption func(*sourceParams)

type sourceParams struct {
	limit     int
	sortOrder SortOrder
}

// WithSourceLimit sets the maximum number of sources (FRED: limit).
func WithSourceLimit(n int) SourceOption {
	return func(p *sourceParams) { p.limit = n }
}

// WithSourceSortOrder sets the sort order for source results (FRED: sort_order).
func WithSourceSortOrder(o SortOrder) SourceOption {
	return func(p *sourceParams) { p.sortOrder = o }
}

// MapDataOption configures GeoFRED series data calls.
type MapDataOption func(*mapDataParams)

type mapDataParams struct {
	date      string
	startDate string
}

// WithMapDate sets the date for map data (FRED: date).
func WithMapDate(d string) MapDataOption {
	return func(p *mapDataParams) { p.date = d }
}

// WithMapStartDate sets the start date for map data range (FRED: start_date).
func WithMapStartDate(d string) MapDataOption {
	return func(p *mapDataParams) { p.startDate = d }
}

// RegionalDataOption configures GeoFRED regional data calls.
type RegionalDataOption func(*regionalDataParams)

type regionalDataParams struct {
	seriesGroup    string
	regionType     string
	date           string
	season         string
	units          string
	transformation string
	frequency      string
}

// WithSeriesGroup sets the series group ID (FRED: series_group). Required.
func WithSeriesGroup(g string) RegionalDataOption {
	return func(p *regionalDataParams) { p.seriesGroup = g }
}

// WithRegionType sets the region type (FRED: region_type).
// Values: bea, msa, frb, necta, state, country, county, censusregion.
func WithRegionType(t string) RegionalDataOption {
	return func(p *regionalDataParams) { p.regionType = t }
}

// WithRegionalDate sets the date for regional data (FRED: date).
func WithRegionalDate(d string) RegionalDataOption {
	return func(p *regionalDataParams) { p.date = d }
}

// WithSeason sets the seasonality adjustment (FRED: season).
// Values: SA, NSA, SSA, SAAR, NSAAR.
func WithSeason(s string) RegionalDataOption {
	return func(p *regionalDataParams) { p.season = s }
}

// WithMapUnits sets the units for map data (FRED: units).
func WithMapUnits(u string) RegionalDataOption {
	return func(p *regionalDataParams) { p.units = u }
}

// WithTransformation sets the data transformation (FRED: transformation).
// Values: lin, chg, ch1, pch, pc1, pca, cch, cca, log.
func WithTransformation(t string) RegionalDataOption {
	return func(p *regionalDataParams) { p.transformation = t }
}

// WithRegionalFrequency sets the frequency aggregation (FRED: frequency).
func WithRegionalFrequency(f string) RegionalDataOption {
	return func(p *regionalDataParams) { p.frequency = f }
}
