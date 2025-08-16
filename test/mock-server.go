package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

// Mock API responses
type ValidationResponse struct {
	Valid      bool         `json:"valid"`
	Email      string       `json:"email"`
	Tier       string       `json:"tier"`
	ExpiresAt  string       `json:"expires_at"`
	ValidUntil string       `json:"valid_until"`
	Pricing    *PricingInfo `json:"pricing,omitempty"`
}

type PricingInfo struct {
	Plan         string  `json:"plan"`
	Currency     string  `json:"currency"`
	Amount       float64 `json:"amount"`
	BillingCycle string  `json:"billing_cycle"`
	Seats        int     `json:"seats,omitempty"`
	NextBilling  string  `json:"next_billing,omitempty"`
}

type PreToolUseRequest struct {
	HookType string `json:"hook_type"`
	APIKey   string `json:"api_key"`
	Command  struct {
		Raw    string   `json:"raw"`
		Parsed []string `json:"parsed"`
	} `json:"command"`
}

type PreToolUseResponse struct {
	Action  string `json:"action"`
	Message string `json:"message,omitempty"`
}

type PostToolUseRequest struct {
	HookType string          `json:"hook_type"`
	APIKey   string          `json:"api_key"`
	Result   json.RawMessage `json:"result"`
}

type PostToolUseResponse struct {
	Insights []struct {
		Type        string `json:"type"`
		Title       string `json:"title"`
		Description string `json:"description"`
	} `json:"insights"`
}

func main() {
	// API key validation endpoint
	http.HandleFunc("/v1/validate", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s request to %s", r.Method, r.URL.Path)
		
		// Check authorization header
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		
		apiKey := strings.TrimPrefix(auth, "Bearer ")
		log.Printf("Validating API key: %s", apiKey)
		
		// Mock validation - accept any key starting with "test-"
		if strings.HasPrefix(apiKey, "test-") {
			// Determine pricing based on API key
			var pricing *PricingInfo
			
			if strings.Contains(apiKey, "team") {
				// Team plan pricing
				pricing = &PricingInfo{
					Plan:         "team",
					Currency:     "USD",
					Amount:       20.00, // $20 per seat
					BillingCycle: "monthly",
					Seats:        5, // Example: 5 seats
					NextBilling:  time.Now().Add(30 * 24 * time.Hour).Format("2006-01-02"),
				}
			} else {
				// Individual plan pricing
				pricing = &PricingInfo{
					Plan:         "individual",
					Currency:     "USD",
					Amount:       10.00, // $10 per user
					BillingCycle: "monthly",
					NextBilling:  time.Now().Add(30 * 24 * time.Hour).Format("2006-01-02"),
				}
			}
			
			resp := ValidationResponse{
				Valid:      true,
				Email:      "test@example.com",
				Tier:       "pro",
				ExpiresAt:  time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
				ValidUntil: "2025-02-13",
				Pricing:    pricing,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			log.Printf("API key validated successfully")
		} else {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
		}
	})
	
	// Pre-tool-use webhook endpoint
	http.HandleFunc("/v1/pre-tool-use", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s request to %s", r.Method, r.URL.Path)
		
		var req PreToolUseRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		log.Printf("Pre-tool-use webhook for command: %s", req.Command.Raw)
		
		// Mock webhook logic
		resp := PreToolUseResponse{
			Action: "ALLOW",
		}
		
		// Block dangerous commands
		if strings.Contains(req.Command.Raw, "rm -rf") {
			resp.Action = "BLOCK"
			resp.Message = "Command blocked: potentially dangerous rm -rf detected"
		}
		
		// Warn about sudo commands
		if strings.Contains(req.Command.Raw, "sudo") {
			resp.Action = "WARN"
			resp.Message = "Warning: Command requires elevated privileges"
		}
		
		// Modify ls commands to add -la
		if len(req.Command.Parsed) > 0 && req.Command.Parsed[0] == "ls" {
			resp.Action = "MODIFY"
			resp.Message = "Command modified: added -la flags for detailed listing"
			// Note: In real implementation, would return modified command
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		log.Printf("Pre-tool-use response: %s", resp.Action)
	})
	
	// Post-tool-use webhook endpoint
	http.HandleFunc("/v1/post-tool-use", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s request to %s", r.Method, r.URL.Path)
		
		var req PostToolUseRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		log.Printf("Post-tool-use webhook received")
		
		// Mock insights
		resp := PostToolUseResponse{
			Insights: []struct {
				Type        string `json:"type"`
				Title       string `json:"title"`
				Description string `json:"description"`
			}{
				{
					Type:        "info",
					Title:       "Command Analysis",
					Description: "Command executed successfully with normal token usage",
				},
			},
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		log.Printf("Post-tool-use insights sent")
	})
	
	port := ":8899"
	log.Printf("Starting mock ctx.pro server on %s", port)
	log.Printf("Test with: CTX_API_ENDPOINT=http://localhost:8899")
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}