package fred

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/shanehull/go-fred/internal"
)

// GetTags returns all FRED tags.
// Endpoint: GET /fred/tags
func (c *Client) GetTags(ctx context.Context, opts ...TagOption) ([]Tag, error) {
	p := &tagParams{}
	for _, opt := range opts {
		opt(p)
	}

	params := url.Values{}
	applyTagOptions(p, params)

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/fred/tags", params)
	if err != nil {
		return nil, err
	}

	var resp tagsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding tags: %w", err)
	}

	return resp.Tags, nil
}

// GetRelatedTags returns related tags for given tag names.
// Endpoint: GET /fred/related_tags?tag_names={tags}
func (c *Client) GetRelatedTags(ctx context.Context, tagNames []string, opts ...TagOption) ([]Tag, error) {
	p := &tagParams{}
	for _, opt := range opts {
		opt(p)
	}

	params := url.Values{}
	params.Set("tag_names", strings.Join(tagNames, ";"))
	applyTagOptions(p, params)

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/fred/related_tags", params)
	if err != nil {
		return nil, err
	}

	var resp tagsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding related tags: %w", err)
	}

	return resp.Tags, nil
}

// GetTagsSeries returns series matching given tag names.
// Endpoint: GET /fred/tags/series?tag_names={tags}
func (c *Client) GetTagsSeries(ctx context.Context, tagNames []string, opts ...SearchOption) ([]SearchResult, error) {
	p := &searchParams{}
	for _, opt := range opts {
		opt(p)
	}

	params := url.Values{}
	params.Set("tag_names", strings.Join(tagNames, ";"))
	applySearchOptions(p, params)

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/fred/tags/series", params)
	if err != nil {
		return nil, err
	}

	var resp searchResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding tags series: %w", err)
	}

	return resp.Series, nil
}
