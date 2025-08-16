package cmd

import (
	"fmt"

	"github.com/slavakurilyak/ctx/internal/auth"
	"github.com/slavakurilyak/ctx/internal/config"
	"github.com/spf13/cobra"
)

// NewAccountCmd creates the account command
func NewAccountCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "account",
		Short: "View your ctx Pro account status",
		Long: `View your current ctx Pro subscription status and account details.

This command displays:
- Authentication status
- Subscription tier
- API key validity
- Expiration date`,
		RunE: runAccount,
	}
}

// NewLogoutCmd creates the logout command
func NewLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Log out of your ctx Pro account",
		Long: `Log out of your ctx Pro account and clear stored credentials.

This command will:
- Remove your API key from the system keychain
- Clear authentication details from the configuration file
- Disable Pro features`,
		RunE: runLogout,
	}
}

func runAccount(cmd *cobra.Command, args []string) error {
	// Load config from file and environment
	cfg := config.NewFromFlagsAndEnv(cmd)
	
	// Check authentication status
	authManager := auth.NewManager(cfg)
	if !authManager.IsAuthenticated() {
		fmt.Println("You are not logged in.")
		fmt.Println("\nRun 'ctx login' to authenticate with your ctx Pro account.")
		return nil
	}

	// Get account information
	userInfo, err := authManager.GetAccountInfo()
	if err != nil {
		fmt.Println("Error retrieving account information:", err)
		fmt.Println("\nYour API key may have expired. Please run 'ctx login' to re-authenticate.")
		return nil
	}

	// Display account information
	fmt.Println("ctx Pro Account Status")
	fmt.Println("======================")
	fmt.Printf("Email:       %s\n", userInfo.Email)
	fmt.Printf("Tier:        %s\n", userInfo.Tier)
	fmt.Printf("Status:      Active\n")
	if userInfo.ValidUntil != "" {
		fmt.Printf("Valid Until: %s\n", userInfo.ValidUntil)
	}

	// Display pricing information if available
	if userInfo.Pricing != nil {
		fmt.Println("\nBilling Information:")
		fmt.Printf("  Plan:         %s\n", userInfo.Pricing.Plan)
		fmt.Printf("  Price:        $%.2f %s/%s\n", userInfo.Pricing.Amount, userInfo.Pricing.Currency, userInfo.Pricing.BillingCycle)
		if userInfo.Pricing.Seats > 0 {
			fmt.Printf("  Seats:        %d\n", userInfo.Pricing.Seats)
			fmt.Printf("  Total:        $%.2f %s/%s\n", userInfo.Pricing.Amount*float64(userInfo.Pricing.Seats), userInfo.Pricing.Currency, userInfo.Pricing.BillingCycle)
		}
		if userInfo.Pricing.NextBilling != "" {
			fmt.Printf("  Next Billing: %s\n", userInfo.Pricing.NextBilling)
		}
	}


	return nil
}

func runLogout(cmd *cobra.Command, args []string) error {
	// Load config from file and environment
	cfg := config.NewFromFlagsAndEnv(cmd)
	
	// Check if logged in
	authManager := auth.NewManager(cfg)
	if !authManager.IsAuthenticated() {
		fmt.Println("You are not logged in.")
		return nil
	}

	// Perform logout
	if err := authManager.Logout(); err != nil {
		return fmt.Errorf("logout failed: %w", err)
	}

	fmt.Println("✓ Successfully logged out.")
	fmt.Println("Pro features have been disabled.")

	return nil
}

func enabledStatus(enabled bool) string {
	if enabled {
		return "Enabled"
	}
	return "Disabled"
}