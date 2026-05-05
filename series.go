package fred

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/shanehull/go-fred/internal"
)

// parseObservation converts a raw wire-format observation into a clean Observation.
// Missing values (".") result in IsNA=true, Value=0.0.
func parseObservation(raw rawObservation) (Observation, error) {
	obs := Observation{}
	if raw.Date != "" {
		t, err := time.Parse("2006-01-02", raw.Date)
		if err != nil {
			return obs, fmt.Errorf("fred: invalid date %q: %w", raw.Date, err)
		}
		obs.Date = t
	}
	if raw.RealtimeStart != "" {
		t, err := time.Parse("2006-01-02", raw.RealtimeStart)
		if err != nil {
			return obs, fmt.Errorf("fred: invalid realtime_start %q: %w", raw.RealtimeStart, err)
		}
		obs.RealtimeStart = t
	}
	if raw.RealtimeEnd != "" {
		t, err := time.Parse("2006-01-02", raw.RealtimeEnd)
		if err != nil {
			return obs, fmt.Errorf("fred: invalid realtime_end %q: %w", raw.RealtimeEnd, err)
		}
		obs.RealtimeEnd = t
	}
	if raw.Value != "" && raw.Value != "." {
		v, err := strconv.ParseFloat(raw.Value, 64)
		if err != nil {
			return obs, fmt.Errorf("fred: invalid value %q: %w", raw.Value, err)
		}
		obs.Value = v
	} else {
		obs.IsNA = true
	}
	return obs, nil
}

// parseObservations converts a slice of raw observations into clean Observations.
func parseObservations(raw []rawObservation) ([]Observation, error) {
	out := make([]Observation, len(raw))
	for i, r := range raw {
		o, err := parseObservation(r)
		if err != nil {
			return nil, fmt.Errorf("fred: observation %d: %w", i, err)
		}
		out[i] = o
	}
	return out, nil
}

// applyObsOptions populates url.Values from ObservationOption slice.
func applyObsOptions(p *obsParams, q url.Values) {
	if p.observationStart != "" {
		q.Set("observation_start", p.observationStart)
	}
	if p.observationEnd != "" {
		q.Set("observation_end", p.observationEnd)
	}
	if p.realtimeStart != "" {
		q.Set("realtime_start", p.realtimeStart)
	}
	if p.realtimeEnd != "" {
		q.Set("realtime_end", p.realtimeEnd)
	}
	if p.units != "" {
		q.Set("units", p.units)
	}
	if p.frequency != "" {
		q.Set("frequency", p.frequency)
	}
	if p.aggregationMethod != "" {
		q.Set("aggregation_method", p.aggregationMethod)
	}
	if p.outputType != 0 {
		q.Set("output_type", strconv.Itoa(p.outputType))
	}
	if len(p.vintageDates) > 0 {
		for _, d := range p.vintageDates {
			q.Add("vintage_dates", d)
		}
	}
	if p.sortOrder != "" {
		q.Set("sort_order", string(p.sortOrder))
	}
	if p.limit != 0 {
		q.Set("limit", strconv.Itoa(p.limit))
	}
	if p.offset != 0 {
		q.Set("offset", strconv.Itoa(p.offset))
	}
}

// GetSeriesInfo returns metadata for a FRED series.
// Endpoint: GET /fred/series?series_id={id}
func (c *Client) GetSeriesInfo(ctx context.Context, seriesID string) (*Series, error) {
	params := url.Values{}
	params.Set("series_id", seriesID)

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/fred/series", params)
	if err != nil {
		return nil, err
	}

	var resp seriesInfoResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding series info: %w", err)
	}

	if len(resp.Series) == 0 {
		return nil, fmt.Errorf("fred: series %q not found", seriesID)
	}

	s := resp.Series[0]

	// Merge tags if present in a separate field
	return &s, nil
}

// GetSeriesObservations returns data values for a series.
// Endpoint: GET /fred/series/observations?series_id={id}
// Does not auto-paginate — FRED returns up to 100,000 observations per request.
func (c *Client) GetSeriesObservations(ctx context.Context, seriesID string, opts ...ObservationOption) ([]Observation, error) {
	p := &obsParams{}
	for _, opt := range opts {
		opt(p)
	}

	params := url.Values{}
	params.Set("series_id", seriesID)
	applyObsOptions(p, params)

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/fred/series/observations", params)
	if err != nil {
		return nil, err
	}

	var resp observationResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding observations: %w", err)
	}

	return parseObservations(resp.Observations)
}

// GetSeriesAllReleases returns all observations including every revision.
// Defaults realtime_start="1776-07-04", realtime_end="9999-12-31" unless overridden
// by ObservationOption. Uses the same endpoint as GetSeriesObservations.
func (c *Client) GetSeriesAllReleases(ctx context.Context, seriesID string, opts ...ObservationOption) ([]Observation, error) {
	p := &obsParams{}
	for _, opt := range opts {
		opt(p)
	}

	defaultStart := time.Date(1776, 7, 4, 0, 0, 0, 0, time.UTC)
	defaultEnd := time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC)

	params := url.Values{}
	params.Set("series_id", seriesID)

	if p.realtimeStart == "" {
		params.Set("realtime_start", defaultStart.Format("2006-01-02"))
	}
	if p.realtimeEnd == "" {
		params.Set("realtime_end", defaultEnd.Format("2006-01-02"))
	}
	applyObsOptions(p, params)

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/fred/series/observations", params)
	if err != nil {
		return nil, err
	}

	var resp observationResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding observations: %w", err)
	}

	return parseObservations(resp.Observations)
}

// GetSeriesFirstRelease returns the first published value for each date.
// Calls GetSeriesAllReleases and groups by Date, picking the earliest
// RealtimeStart for each Date.
func (c *Client) GetSeriesFirstRelease(ctx context.Context, seriesID string, opts ...ObservationOption) ([]Observation, error) {
	all, err := c.GetSeriesAllReleases(ctx, seriesID, opts...)
	if err != nil {
		return nil, err
	}

	byDate := make(map[time.Time]Observation)
	for _, o := range all {
		existing, ok := byDate[o.Date]
		if !ok || o.RealtimeStart.Before(existing.RealtimeStart) {
			byDate[o.Date] = o
		}
	}

	out := make([]Observation, 0, len(byDate))
	for _, o := range byDate {
		out = append(out, o)
	}

	// Sort by date ascending
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[i].Date.After(out[j].Date) {
				out[i], out[j] = out[j], out[i]
			}
		}
	}

	return out, nil
}

// GetSeriesAsOf returns data as known on a specific date.
// Calls GetSeriesObservations with realtime_end=asOf.
func (c *Client) GetSeriesAsOf(ctx context.Context, seriesID string, asOf time.Time, opts ...ObservationOption) ([]Observation, error) {
	allOpts := []ObservationOption{WithRealtimeEnd(asOf)}
	allOpts = append(allOpts, opts...)
	return c.GetSeriesObservations(ctx, seriesID, allOpts...)
}

// GetSeriesVintageDates returns all vintage dates for a series.
// Endpoint: GET /fred/series/vintagedates?series_id={id}
func (c *Client) GetSeriesVintageDates(ctx context.Context, seriesID string, opts ...ObservationOption) ([]time.Time, error) {
	p := &obsParams{}
	for _, opt := range opts {
		opt(p)
	}

	params := url.Values{}
	params.Set("series_id", seriesID)
	applyObsOptions(p, params)

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/fred/series/vintagedates", params)
	if err != nil {
		return nil, err
	}

	var resp vintageDatesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding vintage dates: %w", err)
	}

	dates := make([]time.Time, 0, len(resp.Dates))
	for _, d := range resp.Dates {
		t, err := time.Parse("2006-01-02", d)
		if err != nil {
			return nil, fmt.Errorf("fred: invalid vintage date %q: %w", d, err)
		}
		dates = append(dates, t)
	}

	return dates, nil
}

// GetSeriesCategories returns categories for a FRED series.
// Endpoint: GET /fred/series/categories?series_id={id}
func (c *Client) GetSeriesCategories(ctx context.Context, seriesID string) ([]Category, error) {
	params := url.Values{}
	params.Set("series_id", seriesID)

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/fred/series/categories", params)
	if err != nil {
		return nil, err
	}

	var resp categoriesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding series categories: %w", err)
	}

	return resp.Categories, nil
}

// GetSeriesRelease returns the release for a FRED series.
// Endpoint: GET /fred/series/release?series_id={id}
func (c *Client) GetSeriesRelease(ctx context.Context, seriesID string) (*Release, error) {
	params := url.Values{}
	params.Set("series_id", seriesID)

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/fred/series/release", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Releases []Release `json:"releases"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding series release: %w", err)
	}

	if len(resp.Releases) == 0 {
		return nil, fmt.Errorf("fred: no release found for series %q", seriesID)
	}

	return &resp.Releases[0], nil
}

// GetSeriesTags returns the tags for a FRED series.
// Endpoint: GET /fred/series/tags?series_id={id}
func (c *Client) GetSeriesTags(ctx context.Context, seriesID string, opts ...TagOption) ([]Tag, error) {
	p := &tagParams{}
	for _, opt := range opts {
		opt(p)
	}

	params := url.Values{}
	params.Set("series_id", seriesID)
	applyTagOptions(p, params)

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/fred/series/tags", params)
	if err != nil {
		return nil, err
	}

	var resp tagsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding series tags: %w", err)
	}

	return resp.Tags, nil
}

// SearchSeriesTags searches for tags by series search text.
// Endpoint: GET /fred/series/search/tags?series_search_text={text}
func (c *Client) SearchSeriesTags(ctx context.Context, text string, opts ...TagOption) ([]Tag, error) {
	p := &tagParams{}
	for _, opt := range opts {
		opt(p)
	}

	params := url.Values{}
	params.Set("series_search_text", text)
	applyTagOptions(p, params)

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/fred/series/search/tags", params)
	if err != nil {
		return nil, err
	}

	var resp tagsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding search tags: %w", err)
	}

	return resp.Tags, nil
}

// SearchSeriesRelatedTags searches for related tags by series search text and tag names.
// Endpoint: GET /fred/series/search/related_tags?series_search_text={text}&tag_names={tags}
func (c *Client) SearchSeriesRelatedTags(ctx context.Context, text string, tagNames []string, opts ...TagOption) ([]Tag, error) {
	p := &tagParams{}
	for _, opt := range opts {
		opt(p)
	}

	params := url.Values{}
	params.Set("series_search_text", text)
	params.Set("tag_names", strings.Join(tagNames, ";"))
	applyTagOptions(p, params)

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/fred/series/search/related_tags", params)
	if err != nil {
		return nil, err
	}

	var resp tagsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding search related tags: %w", err)
	}

	return resp.Tags, nil
}

// applyUpdateOptions populates url.Values from UpdateOption slice.
func applyUpdateOptions(p *updateParams, q url.Values) {
	if p.startTime != "" {
		q.Set("start_time", p.startTime)
	}
	if p.endTime != "" {
		q.Set("end_time", p.endTime)
	}
	if p.filterValue != "" {
		q.Set("filter_value", p.filterValue)
	}
	if p.limit != 0 {
		q.Set("limit", strconv.Itoa(p.limit))
	}
}

// GetSeriesUpdates returns series that were recently updated.
// Endpoint: GET /fred/series/updates
func (c *Client) GetSeriesUpdates(ctx context.Context, opts ...UpdateOption) ([]SearchResult, error) {
	p := &updateParams{}
	for _, opt := range opts {
		opt(p)
	}

	params := url.Values{}
	applyUpdateOptions(p, params)

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/fred/series/updates", params)
	if err != nil {
		return nil, err
	}

	var resp searchResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding series updates: %w", err)
	}

	return resp.Series, nil
}
