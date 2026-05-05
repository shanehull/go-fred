package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// APIError represents a FRED API error response.
type APIError struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("fred: API error %d: %s", e.Code, e.Message)
}

// DoRequest sends a GET request to the FRED API and returns the response body.
// It handles base URL construction, API key injection, JSON file type, and error responses.
// On API errors (non-200), it returns *APIError.
func DoRequest(ctx context.Context, client *http.Client, baseURL, apiKey, path string, params url.Values) ([]byte, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("api_key", apiKey)
	params.Set("file_type", "json")

	u, err := url.Parse(baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("fred: invalid URL %q: %w", baseURL+path, err)
	}
	u.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("fred: request creation failed: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fred: request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("fred: reading response failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr APIError
		if err := json.Unmarshal(body, &apiErr); err != nil || apiErr.Message == "" {
			return nil, fmt.Errorf("fred: HTTP %d: %s", resp.StatusCode, string(body))
		}
		return nil, &apiErr
	}

	return body, nil
}
