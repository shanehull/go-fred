// Package fred provides a client for the Federal Reserve Economic Data (FRED) API.
//
// It covers all 34 FRED API endpoints. Zero dependencies beyond Go stdlib.
//
// Usage:
//
//	client, err := fred.New(fred.WithAPIKey("your-key"))
//	if err != nil {
//	    log.Fatal(err)
//	}
//	obs, err := client.GetSeriesObservations(ctx, "DGS20")
package fred

import (
	"fmt"
	"net/http"
	"os"
)

// Client is a FRED API client. Must be constructed via New.
// Safe for concurrent use by multiple goroutines.
//
// API key resolution order:
//  1. WithAPIKey("...")
//  2. $FRED_API_KEY environment variable
type Client struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

// New creates a new FRED API client.
func New(opts ...ClientOption) (*Client, error) {
	c := &Client{
		httpClient: http.DefaultClient,
		baseURL:    "https://api.stlouisfed.org",
	}

	if key := os.Getenv("FRED_API_KEY"); key != "" {
		c.apiKey = key
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	if c.apiKey == "" {
		return nil, fmt.Errorf("fred: API key required; set $FRED_API_KEY or use WithAPIKey()")
	}

	return c, nil
}

// ClientOption configures a Client.
type ClientOption func(*Client) error

// WithAPIKey sets the FRED API key.
func WithAPIKey(key string) ClientOption {
	return func(c *Client) error {
		c.apiKey = key
		return nil
	}
}

// WithHTTPClient sets a custom HTTP client (useful for testing with httptest).
func WithHTTPClient(hc *http.Client) ClientOption {
	return func(c *Client) error {
		c.httpClient = hc
		return nil
	}
}

// WithBaseURL sets a custom base URL (useful for testing).
func WithBaseURL(url string) ClientOption {
	return func(c *Client) error {
		c.baseURL = url
		return nil
	}
}
