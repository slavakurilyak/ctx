package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/slavakurilyak/ctx/internal/version"
)

// Client represents the ctx.pro API client
type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

// NewClient creates a new API client
func NewClient(baseURL, apiKey string) *Client {
	// baseURL must be provided via environment or config
	// No hardcoded default

	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}

// ValidateAPIKey validates an API key
func (c *Client) ValidateAPIKey(ctx context.Context) (*ValidationResponse, error) {
	req, err := c.newRequest(ctx, "GET", "/v1/validate", nil)
	if err != nil {
		return nil, err
	}

	var resp ValidationResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetSubscriptionInfo retrieves subscription information
func (c *Client) GetSubscriptionInfo(ctx context.Context) (*SubscriptionInfo, error) {
	req, err := c.newRequest(ctx, "GET", "/v1/subscription", nil)
	if err != nil {
		return nil, err
	}

	var resp SubscriptionInfo
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetWebhookStats retrieves webhook usage statistics
func (c *Client) GetWebhookStats(ctx context.Context) (*WebhookStats, error) {
	req, err := c.newRequest(ctx, "GET", "/v1/webhook-stats", nil)
	if err != nil {
		return nil, err
	}

	var resp WebhookStats
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// newRequest creates a new HTTP request
func (c *Client) newRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	url := c.baseURL + path

	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("User-Agent", version.GetUserAgent())

	return req, nil
}

// do executes an HTTP request and decodes the response
func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return fmt.Errorf("API error: status %d", resp.StatusCode)
		}
		return fmt.Errorf("API error: %s", errResp.Message)
	}

	// Decode successful response
	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// SetAPIKey updates the API key
func (c *Client) SetAPIKey(apiKey string) {
	c.apiKey = apiKey
}

// SetBaseURL updates the base URL
func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}
