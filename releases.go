package fred

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/shanehull/go-fred/internal"
)

// applyReleaseListOptions populates url.Values from ReleaseListOption slice.
func applyReleaseListOptions(p *releaseListParams, q url.Values) {
	if p.limit != 0 {
		q.Set("limit", strconv.Itoa(p.limit))
	}
	if p.sortOrder != "" {
		q.Set("sort_order", string(p.sortOrder))
	}
}

// applyReleaseDateOptions populates url.Values from ReleaseDateOption slice.
func applyReleaseDateOptions(p *releaseDateParams, q url.Values) {
	if p.limit != 0 {
		q.Set("limit", strconv.Itoa(p.limit))
	}
	if p.sortOrder != "" {
		q.Set("sort_order", string(p.sortOrder))
	}
	if p.includeNoData {
		q.Set("include_release_dates_with_no_data", "true")
	}
}

// applyTableOptions populates url.Values from TableOption slice.
func applyTableOptions(p *tableParams, q url.Values) {
	if p.elementID != 0 {
		q.Set("element_id", strconv.Itoa(p.elementID))
	}
	if p.includeObservationValues {
		q.Set("include_observation_values", "true")
	}
	if p.observationDate != "" {
		q.Set("observation_date", p.observationDate)
	}
}

// GetReleases returns all FRED releases.
// Endpoint: GET /fred/releases
func (c *Client) GetReleases(ctx context.Context, opts ...ReleaseListOption) ([]Release, error) {
	p := &releaseListParams{}
	for _, opt := range opts {
		opt(p)
	}

	params := url.Values{}
	applyReleaseListOptions(p, params)

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/fred/releases", params)
	if err != nil {
		return nil, err
	}

	var resp releasesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding releases: %w", err)
	}

	return resp.Releases, nil
}

// GetReleasesDates returns release dates for all FRED releases.
// Endpoint: GET /fred/releases/dates
func (c *Client) GetReleasesDates(ctx context.Context, opts ...ReleaseDateOption) ([]ReleaseDate, error) {
	p := &releaseDateParams{}
	for _, opt := range opts {
		opt(p)
	}

	params := url.Values{}
	applyReleaseDateOptions(p, params)

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/fred/releases/dates", params)
	if err != nil {
		return nil, err
	}

	var resp releasesDatesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding releases dates: %w", err)
	}

	return resp.Dates, nil
}

// GetRelease returns information about a FRED release.
// Endpoint: GET /fred/release?release_id={id}
func (c *Client) GetRelease(ctx context.Context, id int) (*Release, error) {
	params := url.Values{}
	params.Set("release_id", strconv.Itoa(id))

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/fred/release", params)
	if err != nil {
		return nil, err
	}

	var resp releasesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding release: %w", err)
	}

	if len(resp.Releases) == 0 {
		return nil, fmt.Errorf("fred: release %d not found", id)
	}

	return &resp.Releases[0], nil
}

// GetReleaseDates returns release dates for a specific FRED release.
// Endpoint: GET /fred/release/dates?release_id={id}
func (c *Client) GetReleaseDates(ctx context.Context, id int, opts ...ReleaseDateOption) ([]ReleaseDate, error) {
	p := &releaseDateParams{}
	for _, opt := range opts {
		opt(p)
	}

	params := url.Values{}
	params.Set("release_id", strconv.Itoa(id))
	applyReleaseDateOptions(p, params)

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/fred/release/dates", params)
	if err != nil {
		return nil, err
	}

	var resp releasesDatesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding release dates: %w", err)
	}

	return resp.Dates, nil
}

// GetReleaseSources returns the sources for a FRED release.
// Endpoint: GET /fred/release/sources?release_id={id}
func (c *Client) GetReleaseSources(ctx context.Context, id int) ([]Source, error) {
	params := url.Values{}
	params.Set("release_id", strconv.Itoa(id))

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/fred/release/sources", params)
	if err != nil {
		return nil, err
	}

	var resp sourcesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding release sources: %w", err)
	}

	return resp.Sources, nil
}

// GetReleaseTags returns the tags for a FRED release.
// Endpoint: GET /fred/release/tags?release_id={id}
func (c *Client) GetReleaseTags(ctx context.Context, id int, opts ...TagOption) ([]Tag, error) {
	p := &tagParams{}
	for _, opt := range opts {
		opt(p)
	}

	params := url.Values{}
	params.Set("release_id", strconv.Itoa(id))
	applyTagOptions(p, params)

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/fred/release/tags", params)
	if err != nil {
		return nil, err
	}

	var resp tagsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding release tags: %w", err)
	}

	return resp.Tags, nil
}

// GetReleaseRelatedTags returns related tags for a FRED release.
// Endpoint: GET /fred/release/related_tags?release_id={id}&tag_names={tags}
func (c *Client) GetReleaseRelatedTags(ctx context.Context, id int, tagNames []string, opts ...TagOption) ([]Tag, error) {
	p := &tagParams{}
	for _, opt := range opts {
		opt(p)
	}

	params := url.Values{}
	params.Set("release_id", strconv.Itoa(id))
	params.Set("tag_names", strings.Join(tagNames, ";"))
	applyTagOptions(p, params)

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/fred/release/related_tags", params)
	if err != nil {
		return nil, err
	}

	var resp tagsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding release related tags: %w", err)
	}

	return resp.Tags, nil
}

// GetReleaseTables returns the release table for a FRED release.
// Endpoint: GET /fred/release/tables?release_id={id}
func (c *Client) GetReleaseTables(ctx context.Context, id int, opts ...TableOption) ([]ReleaseTableElement, error) {
	p := &tableParams{}
	for _, opt := range opts {
		opt(p)
	}

	params := url.Values{}
	params.Set("release_id", strconv.Itoa(id))
	applyTableOptions(p, params)

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/fred/release/tables", params)
	if err != nil {
		return nil, err
	}

	var resp releaseTableResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding release tables: %w", err)
	}

	elements := make([]ReleaseTableElement, 0, len(resp.Elements))
	for _, el := range resp.Elements {
		if el.LevelRaw != "" {
			l, err := strconv.Atoi(el.LevelRaw)
			if err == nil {
				el.Level = l
			}
		}
		elements = append(elements, el)
	}

	return elements, nil
}
