// Package fred provides a client for the Federal Reserve Economic Data (FRED) API.
//
// It covers all FRED API endpoints. Zero dependencies beyond Go stdlib.
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
	"time"
)

const (
	defaultBaseURL = "https://api.stlouisfed.org"
	defaultTimeout = 30 * time.Second
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
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		baseURL: defaultBaseURL,
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

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(hc *http.Client) ClientOption {
	return func(c *Client) error {
		if hc == nil {
			return fmt.Errorf("fred: http client cannot be nil")
		}
		c.httpClient = hc
		return nil
	}
}

// WithBaseURL sets a custom base URL.
func WithBaseURL(url string) ClientOption {
	return func(c *Client) error {
		c.baseURL = url
		return nil
	}
}
