package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/slavakurilyak/ctx/internal/auth"
	"github.com/slavakurilyak/ctx/internal/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// NewLoginCmd creates the login command
func NewLoginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Authenticate with your ctx Pro account",
		Long: `Authenticate with your ctx Pro account to enable premium features.

You will be prompted to enter your API key, which will be securely stored
in your system's keychain. Your subscription details will be saved to the
local configuration file.`,
		RunE: runLogin,
	}
}

func runLogin(cmd *cobra.Command, args []string) error {
	// Load config from file and environment
	cfg := config.NewFromFlagsAndEnv(cmd)

	// Check if already logged in
	authManager := auth.NewManager(cfg)
	if authManager.IsAuthenticated() {
		fmt.Println("You are already logged in. Use 'ctx logout' to log out first.")
		return nil
	}

	// Prompt for API key
	fmt.Print("Enter your ctx Pro API key: ")
	apiKey, err := readAPIKey()
	if err != nil {
		return fmt.Errorf("failed to read API key: %w", err)
	}

	// Trim whitespace
	apiKey = strings.TrimSpace(apiKey)

	// Validate and store credentials
	fmt.Println("\nAuthenticating...")
	if err := authManager.Login(apiKey); err != nil {
		return err
	}

	// Get account info to display
	userInfo, err := authManager.GetAccountInfo()
	if err != nil {
		fmt.Println("✓ Authentication successful!")
		return nil
	}

	// Display success message with account details
	fmt.Println("\n✓ Successfully logged in!")
	fmt.Printf("  Account: %s\n", userInfo.Email)
	fmt.Printf("  Tier: %s\n", userInfo.Tier)
	if userInfo.ValidUntil != "" {
		fmt.Printf("  Valid until: %s\n", userInfo.ValidUntil)
	}
	fmt.Println("\nPro features are now enabled. Run 'ctx account' to view your account status.")

	return nil
}

// readAPIKey reads the API key from stdin, hiding the input if possible
func readAPIKey() (string, error) {
	// Try to read password without echo
	if term.IsTerminal(int(syscall.Stdin)) {
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			// Fall back to regular input
			return readAPIKeyPlain()
		}
		return string(bytePassword), nil
	}

	// Not a terminal, read normally
	return readAPIKeyPlain()
}

// readAPIKeyPlain reads the API key in plain text
func readAPIKeyPlain() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	apiKey, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(apiKey), nil
}
