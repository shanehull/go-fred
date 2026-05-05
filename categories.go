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

// applyTagOptions populates url.Values from TagOption slice.
func applyTagOptions(p *tagParams, q url.Values) {
	if p.groupID != "" {
		q.Set("tag_group_id", p.groupID)
	}
	if p.searchText != "" {
		q.Set("search_text", p.searchText)
	}
	if p.limit != 0 {
		q.Set("limit", strconv.Itoa(p.limit))
	}
	if p.sortOrder != "" {
		q.Set("sort_order", string(p.sortOrder))
	}
	if p.orderBy != "" {
		q.Set("order_by", string(p.orderBy))
	}
	if len(p.tagNames) > 0 {
		q.Set("tag_names", strings.Join(p.tagNames, ";"))
	}
	if len(p.excludeTagNames) > 0 {
		q.Set("exclude_tag_names", strings.Join(p.excludeTagNames, ";"))
	}
}

// GetCategory returns information about a FRED category.
// Endpoint: GET /fred/category?category_id={id}
func (c *Client) GetCategory(ctx context.Context, id int) (*Category, error) {
	params := url.Values{}
	params.Set("category_id", strconv.Itoa(id))

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/fred/category", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Categories []Category `json:"categories"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding category: %w", err)
	}

	if len(resp.Categories) == 0 {
		return nil, fmt.Errorf("fred: category %d not found", id)
	}

	return &resp.Categories[0], nil
}

// GetCategoryChildren returns the child categories for a given FRED category.
// Endpoint: GET /fred/category/children?category_id={id}
func (c *Client) GetCategoryChildren(ctx context.Context, id int) ([]Category, error) {
	params := url.Values{}
	params.Set("category_id", strconv.Itoa(id))

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/fred/category/children", params)
	if err != nil {
		return nil, err
	}

	var resp categoriesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding category children: %w", err)
	}

	return resp.Categories, nil
}

// GetCategoryRelated returns related categories for a given FRED category.
// Endpoint: GET /fred/category/related?category_id={id}
func (c *Client) GetCategoryRelated(ctx context.Context, id int) ([]Category, error) {
	params := url.Values{}
	params.Set("category_id", strconv.Itoa(id))

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/fred/category/related", params)
	if err != nil {
		return nil, err
	}

	var resp categoriesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding related categories: %w", err)
	}

	return resp.Categories, nil
}

// GetCategoryTags returns the tags for a FRED category.
// Endpoint: GET /fred/category/tags?category_id={id}
func (c *Client) GetCategoryTags(ctx context.Context, id int, opts ...TagOption) ([]Tag, error) {
	p := &tagParams{}
	for _, opt := range opts {
		opt(p)
	}

	params := url.Values{}
	params.Set("category_id", strconv.Itoa(id))
	applyTagOptions(p, params)

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/fred/category/tags", params)
	if err != nil {
		return nil, err
	}

	var resp tagsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding category tags: %w", err)
	}

	return resp.Tags, nil
}

// GetCategoryRelatedTags returns related tags for a FRED category.
// Endpoint: GET /fred/category/related_tags?category_id={id}&tag_names={tags}
func (c *Client) GetCategoryRelatedTags(ctx context.Context, id int, tagNames []string, opts ...TagOption) ([]Tag, error) {
	p := &tagParams{}
	for _, opt := range opts {
		opt(p)
	}

	params := url.Values{}
	params.Set("category_id", strconv.Itoa(id))
	params.Set("tag_names", strings.Join(tagNames, ";"))
	applyTagOptions(p, params)

	body, err := internal.DoRequest(ctx, c.httpClient, c.baseURL, c.apiKey, "/fred/category/related_tags", params)
	if err != nil {
		return nil, err
	}

	var resp tagsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fred: decoding category related tags: %w", err)
	}

	return resp.Tags, nil
}
