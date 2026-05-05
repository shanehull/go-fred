package fred

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/shanehull/go-fred/internal"
)

// applySourceOptions populates url.Values from SourceOption slice.
func applySourceOptions(p *sourceParams, q url.Values) {
	if p.limit != 0 {
		q.Set("limit", strconv.Itoa(p.limit))
	}
	if p.sortOrder != "" {
		q.Set("sort_order", string(p.sortOrder))
	}
}

// GetSources returns all FRED data sources.
// Endpoint: GET /fred/sources
func (c *Client) GetSources(ctx context.Context, opts ...SourceOption) ([]Source, error) {
	p := &sourceParams{}
	for _, opt := range opts {
		opt(p)
	}

	params := url.Values{}
	applySourceOptions(p, params)

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/fred/sources", params)
	if err != nil {
		return nil, err
	}

	var resp sourcesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding sources: %w", err)
	}

	return resp.Sources, nil
}

// GetSource returns information about a FRED data source.
// Endpoint: GET /fred/source?source_id={id}
func (c *Client) GetSource(ctx context.Context, id int) (*Source, error) {
	params := url.Values{}
	params.Set("source_id", strconv.Itoa(id))

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/fred/source", params)
	if err != nil {
		return nil, err
	}

	var resp sourcesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding source: %w", err)
	}

	if len(resp.Sources) == 0 {
		return nil, fmt.Errorf("fred: source %d not found", id)
	}

	return &resp.Sources[0], nil
}

// GetSourceReleases returns the releases for a FRED data source.
// Endpoint: GET /fred/source/releases?source_id={id}
func (c *Client) GetSourceReleases(ctx context.Context, id int, opts ...ReleaseListOption) ([]Release, error) {
	p := &releaseListParams{}
	for _, opt := range opts {
		opt(p)
	}

	params := url.Values{}
	params.Set("source_id", strconv.Itoa(id))
	applyReleaseListOptions(p, params)

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/fred/source/releases", params)
	if err != nil {
		return nil, err
	}

	var resp releasesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding source releases: %w", err)
	}

	return resp.Releases, nil
}
