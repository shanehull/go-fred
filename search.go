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

const pageSize = 1000

// applySearchOptions populates url.Values from SearchOption slice.
func applySearchOptions(p *searchParams, q url.Values) {
	if p.searchType != "" {
		q.Set("search_type", p.searchType)
	}
	if p.orderBy != "" {
		q.Set("order_by", string(p.orderBy))
	}
	if p.sortOrder != "" {
		q.Set("sort_order", string(p.sortOrder))
	}
	if p.filterVariable != "" {
		q.Set("filter_variable", p.filterVariable)
	}
	if p.filterValue != "" {
		q.Set("filter_value", p.filterValue)
	}
	if p.limit != 0 {
		q.Set("limit", strconv.Itoa(p.limit))
	}
	if len(p.tagNames) > 0 {
		q.Set("tag_names", strings.Join(p.tagNames, ";"))
	}
	if len(p.excludeTagNames) > 0 {
		q.Set("exclude_tag_names", strings.Join(p.excludeTagNames, ";"))
	}
}

// paginatedSearch fetches all results from a search endpoint with auto-pagination.
// path is the API path (e.g. "/fred/series/search").
// queryParams are the base query parameters (without api_key or file_type).
// limit 0 means unlimited (fetch all).
func (c *Client) paginatedSearch(ctx context.Context, path string, queryParams url.Values, limit int) ([]SearchResult, error) {
	// First page — get count and first batch
	firstParams := cloneValues(queryParams)
	firstParams.Set("offset", "0")
	firstParams.Set("limit", strconv.Itoa(pageSize))

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, path, firstParams)
	if err != nil {
		return nil, err
	}

	var resp searchResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding search results: %w", err)
	}

	results := resp.Series

	totalCount := resp.Count

	// Determine max results to fetch
	max := totalCount
	if limit > 0 && limit < max {
		max = limit
	}

	// Fetch remaining pages
	for offset := pageSize; offset < max && offset < totalCount; offset += pageSize {
		p := cloneValues(queryParams)
		p.Set("offset", strconv.Itoa(offset))
		p.Set("limit", strconv.Itoa(pageSize))

		body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, path, p)
		if err != nil {
			return nil, fmt.Errorf("fred: fetching page at offset %d: %w", offset, err)
		}

		var r searchResponse
		if err := json.Unmarshal(body, &r); err != nil {
			return nil, fmt.Errorf("fred: decoding search results at offset %d: %w", offset, err)
		}

		results = append(results, r.Series...)
	}

	// Truncate to max
	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// cloneValues returns a copy of url.Values.
func cloneValues(v url.Values) url.Values {
	c := make(url.Values, len(v))
	for k, vs := range v {
		c[k] = append([]string(nil), vs...)
	}
	return c
}

// SearchSeries searches for series by keywords.
// Auto-paginates. Limit 0 = unlimited. Default limit = 1000.
// Endpoint: GET /fred/series/search?search_text={text}
func (c *Client) SearchSeries(ctx context.Context, text string, opts ...SearchOption) ([]SearchResult, error) {
	p := &searchParams{}
	for _, opt := range opts {
		opt(p)
	}

	params := url.Values{}
	params.Set("search_text", text)
	applySearchOptions(p, params)

	return c.paginatedSearch(ctx, "/fred/series/search", params, p.limit)
}

// GetReleaseSeries returns series belonging to a release.
// Auto-paginates. Limit 0 = unlimited. Default limit = 0.
// Endpoint: GET /fred/release/series?release_id={id}
func (c *Client) GetReleaseSeries(ctx context.Context, releaseID int, opts ...SearchOption) ([]SearchResult, error) {
	p := &searchParams{}
	for _, opt := range opts {
		opt(p)
	}

	params := url.Values{}
	params.Set("release_id", strconv.Itoa(releaseID))
	applySearchOptions(p, params)

	return c.paginatedSearch(ctx, "/fred/release/series", params, p.limit)
}

// GetCategorySeries returns series belonging to a category.
// Auto-paginates. Limit 0 = unlimited. Default limit = 0.
// Endpoint: GET /fred/category/series?category_id={id}
func (c *Client) GetCategorySeries(ctx context.Context, categoryID int, opts ...SearchOption) ([]SearchResult, error) {
	p := &searchParams{}
	for _, opt := range opts {
		opt(p)
	}

	params := url.Values{}
	params.Set("category_id", strconv.Itoa(categoryID))
	applySearchOptions(p, params)

	return c.paginatedSearch(ctx, "/fred/category/series", params, p.limit)
}
