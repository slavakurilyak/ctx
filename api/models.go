package api

import "time"

// ValidationRequest represents a request to validate an API key
type ValidationRequest struct {
	APIKey string `json:"api_key"`
}

// ValidationResponse represents the response from API key validation
type ValidationResponse struct {
	Valid      bool         `json:"valid"`
	Email      string       `json:"email"`
	Tier       string       `json:"tier"`
	ExpiresAt  time.Time    `json:"expires_at"`
	ValidUntil string       `json:"valid_until"`
	Pricing    *PricingInfo `json:"pricing,omitempty"`
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// SubscriptionInfo represents subscription information
type SubscriptionInfo struct {
	Tier       string       `json:"tier"`
	Status     string       `json:"status"`
	StartedAt  time.Time    `json:"started_at"`
	ExpiresAt  time.Time    `json:"expires_at"`
	Features   []string     `json:"features"`
	Limits     Limits       `json:"limits"`
	Pricing    *PricingInfo `json:"pricing,omitempty"`
}

// PricingInfo represents pricing information for a subscription
type PricingInfo struct {
	Plan         string  `json:"plan"`           // "individual" or "team"
	Currency     string  `json:"currency"`       // USD, EUR, etc.
	Amount       float64 `json:"amount"`         // Current price
	BillingCycle string  `json:"billing_cycle"`  // monthly, annual
	Seats        int     `json:"seats,omitempty"` // Number of seats for team plans
	NextBilling  string  `json:"next_billing,omitempty"`
}

// Limits represents subscription limits
type Limits struct {
	MaxRequests       int `json:"max_requests"`
	MaxTokensPerMonth int `json:"max_tokens_per_month"`
	MaxWebhookTimeout int `json:"max_webhook_timeout"`
}

// WebhookStats represents webhook usage statistics
type WebhookStats struct {
	PreToolUse  WebhookStat `json:"pre_tool_use"`
	PostToolUse WebhookStat `json:"post_tool_use"`
}

// WebhookStat represents statistics for a single webhook type
type WebhookStat struct {
	TotalCalls      int     `json:"total_calls"`
	SuccessfulCalls int     `json:"successful_calls"`
	FailedCalls     int     `json:"failed_calls"`
	AverageLatency  float64 `json:"average_latency_ms"`
	LastCalled      string  `json:"last_called,omitempty"`
}