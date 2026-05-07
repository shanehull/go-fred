package fred

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/shanehull/go-fred/internal"
)

// applyMapDataOptions populates url.Values from MapDataOption slice.
func applyMapDataOptions(p *mapDataParams, q url.Values) {
	if p.date != "" {
		q.Set("date", p.date)
	}
	if p.startDate != "" {
		q.Set("start_date", p.startDate)
	}
}

// applyRegionalDataOptions populates url.Values from RegionalDataOption slice.
func applyRegionalDataOptions(p *regionalDataParams, q url.Values) {
	if p.seriesGroup != "" {
		q.Set("series_group", p.seriesGroup)
	}
	if p.regionType != "" {
		q.Set("region_type", p.regionType)
	}
	if p.date != "" {
		q.Set("date", p.date)
	}
	if p.season != "" {
		q.Set("season", p.season)
	}
	if p.units != "" {
		q.Set("units", p.units)
	}
	if p.transformation != "" {
		q.Set("transformation", p.transformation)
	}
	if p.frequency != "" {
		q.Set("frequency", p.frequency)
	}
}

// GetSeriesGroup returns the GeoFRED series group metadata for a series.
// Endpoint: GET /geofred/series/group?series_id={id}
func (c *Client) GetSeriesGroup(ctx context.Context, seriesID string) (*SeriesGroup, error) {
	params := url.Values{}
	params.Set("series_id", seriesID)

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/geofred/series/group", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		SeriesGroup SeriesGroup `json:"series_group"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding series group: %w", err)
	}

	return &resp.SeriesGroup, nil
}

// GetSeriesData returns GeoFRED series data for regional maps.
// Endpoint: GET /geofred/series/data?series_id={id}
func (c *Client) GetSeriesData(ctx context.Context, seriesID string, opts ...MapDataOption) (*MapData, error) {
	p := &mapDataParams{}
	for _, opt := range opts {
		opt(p)
	}

	params := url.Values{}
	params.Set("series_id", seriesID)
	applyMapDataOptions(p, params)

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/geofred/series/data", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Meta mapDataResponse `json:"meta"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding series data: %w", err)
	}

	// Flatten map[string][]MapDataEntry into []MapDataEntry
	result := &MapData{
		Title:       resp.Meta.Title,
		Region:      resp.Meta.Region,
		Seasonality: resp.Meta.Seasonality,
		Units:       resp.Meta.Units,
		Frequency:   resp.Meta.Frequency,
		Date:        resp.Meta.Date,
	}
	for _, entries := range resp.Meta.Data {
		result.Data = append(result.Data, entries...)
	}

	return result, nil
}

// GetRegionalData returns GeoFRED regional data for a series group and region type.
// Endpoint: GET /geofred/regional/data
func (c *Client) GetRegionalData(ctx context.Context, opts ...RegionalDataOption) (*MapData, error) {
	p := &regionalDataParams{}
	for _, opt := range opts {
		opt(p)
	}

	if p.seriesGroup == "" {
		return nil, fmt.Errorf("fred: WithSeriesGroup is required for GetRegionalData")
	}
	if p.regionType == "" {
		return nil, fmt.Errorf("fred: WithRegionType is required for GetRegionalData")
	}
	if p.date == "" {
		return nil, fmt.Errorf("fred: WithRegionalDate is required for GetRegionalData")
	}
	if p.season == "" {
		return nil, fmt.Errorf("fred: WithSeason is required for GetRegionalData")
	}
	if p.units == "" {
		return nil, fmt.Errorf("fred: WithMapUnits is required for GetRegionalData")
	}

	params := url.Values{}
	applyRegionalDataOptions(p, params)

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/geofred/regional/data", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Meta mapDataResponse `json:"meta"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding regional data: %w", err)
	}

	result := &MapData{
		Title:       resp.Meta.Title,
		Region:      resp.Meta.Region,
		Seasonality: resp.Meta.Seasonality,
		Units:       resp.Meta.Units,
		Frequency:   resp.Meta.Frequency,
		Date:        resp.Meta.Date,
	}
	for _, entries := range resp.Meta.Data {
		result.Data = append(result.Data, entries...)
	}

	return result, nil
}
