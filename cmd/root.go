package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/slavakurilyak/ctx/internal/app"
	"github.com/slavakurilyak/ctx/internal/config"
	"github.com/slavakurilyak/ctx/internal/history"
	"github.com/slavakurilyak/ctx/internal/telemetry"
	"github.com/slavakurilyak/ctx/internal/tokenizer"
	"github.com/slavakurilyak/ctx/internal/updater"
	"github.com/spf13/cobra"
)

// NewRootCmdWithDI creates a new root command with dependency injection
func NewRootCmdWithDI(version, commit, date string) *cobra.Command {
	// Build the complete help text
	helpText := StyledDescription("Run context (ctx) command to wrap any tool (CLI, shell script, etc) with token awareness/intelligence (example: token count) to empower token-based decisions before and during command execution.") + "\n\n" +
		StyledHeader("USAGE:") + "\n\n" +
		"  " + StyledCommand("ctx [flags] -- <command> [args...]     # POSIX standard separator") + "\n" +
		"  " + StyledCommand("ctx [flags] run <command> [args...]    # Explicit subcommand") + "\n" +
		"  " + StyledCommand("ctx [flags] \"<command with args>\"      # Quoted command (simple cases)") + "\n\n" +
		StyledHeader("EXAMPLES:") + "\n\n" +
		"  " + StyledDescription("# All three methods work:") + "\n" +
		"  " + StyledCommand("ctx -- ls -la") + "\n" +
		"  " + StyledCommand("ctx run ls -la") + "\n" +
		"  " + StyledCommand("ctx \"ls -la\"") + "\n" +
		"  " + StyledDescription("# With ctx flags:") + "\n" +
		"  " + StyledCommand("ctx --max-tokens 5000 -- psql -c \"SELECT * FROM users\"") + "\n" +
		"  " + StyledCommand("ctx --max-tokens 5000 run psql -c \"SELECT * FROM users\"") + "\n" +
		"  " + StyledCommand("ctx --no-tokens \"echo Hello World\"") + "\n\n" +
		config.GenerateEnvSection() + "\n\n" +
		StyledHeader("RECOMMENDED USAGE:") + "\n\n" +
		"  " + StyledDescription("Three workflows to prevent token explosions:") + "\n\n" +
		"  " + StyledSubHeader("PROBE") + "\n" +
		"    " + StyledCommand("ctx psql -c \"SELECT * FROM users\" | jq '.tokens'") + "  " + StyledDescription("# See: 25000 tokens!") + "\n\n" +
		"  " + StyledSubHeader("PROBE-ACT") + "\n" +
		"    " + StyledCommand("ctx --max-lines 1 docker logs app | jq '.tokens'") + "  " + StyledDescription("# Check: 45 tokens") + "\n" +
		"    " + StyledCommand("ctx docker logs app --tail 100") + "  " + StyledDescription("# Act: safe to run") + "\n\n" +
		"  " + StyledSubHeader("PROBE-FILTER-ACT") + "\n" +
		"    " + StyledCommand("ctx psql -c \"SELECT COUNT(*) FROM events\" | jq '.tokens'") + "  " + StyledDescription("# Probe size") + "\n" +
		"    " + StyledCommand("ctx --max-tokens 5000 psql -c \"SELECT id FROM events WHERE status='error'\"") + "  " + StyledDescription("# Filter & act") + "\n\n" +
		StyledHeader("FILTER TYPES:") + "\n\n" +
        "    • " + StyledDescription("ctx filters (safety): --max-tokens, --max-lines, --max-output-bytes, --timeout") + "\n" +
		"    • " + StyledDescription("Command filters (efficiency): LIMIT, WHERE, --tail, --since, head, grep") + "\n" +
		"    • " + StyledDescription("Best: ctx --max-tokens 1000 psql -c \"SELECT id FROM users LIMIT 100\"") + "\n\n" +
		StyledSeparator("════════════════════════════════════════════════════════════════════════════════")

	longDescription := helpText

	var rootCmd = &cobra.Command{
		Use:   "ctx [flags] <command> [args...]",
		Short: "A universal CLI context engine for AI agents.",
		Long:  longDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.ArbitraryArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// 1. Build config, respecting precedence (Flags > Env > Default)
			cfg := config.NewFromFlagsAndEnv(cmd)

			// 2. Initialize Telemetry
			var tel *telemetry.Manager
			var err error
			if !cfg.NoTelemetry {
				tel, err = telemetry.Initialize(cmd.Context())
				if err != nil {
					// Log warning but continue - telemetry is optional
					fmt.Fprintf(os.Stderr, "Warning: failed to initialize telemetry: %v\n", err)
				}
			}

			// 3. Initialize Tokenizer (if not disabled)
			var tok tokenizer.Tokenizer
			if !cfg.NoTokens && cfg.TokenModel != "" {
				factory := &tokenizer.DefaultTokenizerFactory{}
				tokenizerCache := tokenizer.NewTokenizerCache(factory)
				tok, err = tokenizerCache.GetOrCreate(cfg.TokenModel)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Warning: could not initialize tokenizer for provider %q: %v\n", cfg.TokenModel, err)
				}
			}

			// 4. Initialize History Manager
			var hm *history.HistoryManager
			if !cfg.NoHistory {
				hm = history.NewHistoryManager()
			}

			// 5. Build the final AppContext
			appCtx := app.NewAppContext(
				app.WithConfig(cfg),
				app.WithHistory(hm),
				app.WithTelemetry(tel),
				app.WithTokenizer(tok),
			)

			// 6. Inject AppContext into the command's context for RunE to use.
			newCmdCtx := context.WithValue(cmd.Context(), app.AppContextKey, appCtx)
			cmd.SetContext(newCmdCtx)

			// 7. Check for updates if enabled (non-blocking)
			go checkForUpdatesIfNeeded(cfg, version)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Retrieve the AppContext that was prepared in PersistentPreRunE.
			appCtx, ok := cmd.Context().Value(app.AppContextKey).(*app.AppContext)
			if !ok || appCtx == nil {
				return fmt.Errorf("application context not initialized")
			}

			// Check if we need to show help
			if len(args) == 0 {
				// If no command is provided, show help.
				return cmd.Help()
			}

			// Check if this is a quoted command (single arg with spaces)
			// The shell has already removed the quotes for us
			if len(args) == 1 && strings.Contains(args[0], " ") {
				// Split the quoted command into parts
				commandParts := parseQuotedCommand(args[0])
				if len(commandParts) > 0 {
					args = commandParts
				}
			}

			executor := NewCommandExecutor(appCtx)
			
			// Check if streaming is enabled
			isStream, _ := cmd.Flags().GetBool("stream")
			if isStream {
				return executor.ExecuteStreamCommand(cmd.Context(), args)
			}
			
			// Separate ctx flags from the command to be executed.
			// Cobra does this automatically; `args` contains only non-flag arguments.
			return executor.ExecuteCommand(cmd.Context(), args)
		},
	}

	// Add persistent flags that will be available to all subcommands (if any)
	rootCmd.PersistentFlags().String("token-model", "", "Token provider (anthropic, openai, gemini). Overrides CTX_TOKEN_MODEL.")
	rootCmd.PersistentFlags().Bool("no-tokens", false, "Disable token counting. Overrides CTX_NO_TOKENS.")
	rootCmd.PersistentFlags().Int64("max-tokens", 0, "Maximum tokens allowed in output (0 for no limit). Overrides CTX_MAX_TOKENS.")
	rootCmd.PersistentFlags().Int64("max-output-bytes", 0, "Maximum bytes allowed in output (0 for no limit). Overrides CTX_MAX_OUTPUT_BYTES.")
	rootCmd.PersistentFlags().Int64("max-lines", 0, "Maximum lines allowed in output (0 for no limit). Overrides CTX_MAX_LINES.")
	rootCmd.PersistentFlags().Int("max-pipeline-stages", 0, "Maximum pipeline stages allowed (0 for no limit). Overrides CTX_MAX_PIPELINE_STAGES.")
	rootCmd.PersistentFlags().Bool("no-history", false, "Disable saving command history. Overrides CTX_NO_HISTORY.")
	rootCmd.PersistentFlags().Bool("no-telemetry", false, "Disable OpenTelemetry tracing. Overrides CTX_NO_TELEMETRY.")
	rootCmd.PersistentFlags().Bool("private", false, "Enable privacy mode (disables history and telemetry). Overrides CTX_PRIVATE.")
	rootCmd.PersistentFlags().Duration("timeout", 0, "Command execution timeout (e.g., '5s', '1m'). Overrides CTX_TIMEOUT.")
	rootCmd.PersistentFlags().String("output", "json", "Output format ('json').")
	rootCmd.PersistentFlags().Bool("stream", false, "Stream command output line by line for long-running tasks.")
	rootCmd.PersistentFlags().Bool("pretty", false, "Output in pretty format instead of JSON.")

	// Version flag
	rootCmd.Version = fmt.Sprintf("%s, commit %s, built at %s", version, commit, date)
	rootCmd.Flags().Bool("version", false, "Show ctx version")

	// Add telemetry subcommand
	rootCmd.AddCommand(NewTelemetryCmd())
	
	// Add config subcommand
	rootCmd.AddCommand(NewConfigCmd())
	
	// Add version subcommand
	rootCmd.AddCommand(NewVersionCmd(version, commit, date))
	
	// Add run subcommand for explicit command execution
	rootCmd.AddCommand(NewRunCmd())
	
	// Add setup subcommand for setting up coding agents
	rootCmd.AddCommand(NewSetupCmd())

	// Add authentication commands
	rootCmd.AddCommand(NewLoginCmd())
	rootCmd.AddCommand(NewLogoutCmd())
	rootCmd.AddCommand(NewAccountCmd())
	
	// Add update command
	rootCmd.AddCommand(NewUpdateCmd(version))

	return rootCmd
}

// parseQuotedCommand splits a command string into arguments
// The shell has already removed the quotes, so we just need to split on spaces
// This handles the case where user types: ctx "ls -la"
func parseQuotedCommand(cmdStr string) []string {
	// Use Fields which handles multiple spaces and gives us clean args
	return strings.Fields(cmdStr)
}

// checkForUpdatesIfNeeded checks for updates in the background if conditions are met
func checkForUpdatesIfNeeded(cfg *config.Config, currentVersion string) {
	// Only check if installation method supports auto-updates
	if cfg.Installation == nil || !cfg.Installation.AutoUpdateCheck {
		return
	}
	
	// Skip if go-install method (can't auto-update)
	if cfg.Installation.Method == "go-install" {
		return
	}
	
	// Skip if we've checked recently
	if time.Since(cfg.Installation.LastUpdateCheck) < cfg.Installation.UpdateCheckInterval {
		return
	}
	
	// Skip if version is unknown/dev (can't compare)
	if currentVersion == "dev" || currentVersion == "" || strings.Contains(currentVersion, "built from source") {
		return
	}
	
	// Perform the update check (with timeout)
	upd := updater.NewUpdater("slavakurilyak", "ctx")
	upd.HTTPClient.Timeout = 5 * time.Second // Quick check
	
	updateInfo, err := upd.CheckForUpdate(currentVersion, false)
	if err != nil {
		// Silently fail - this is non-critical background check
		return
	}
	
	// Update last check time
	cfg.Installation.LastUpdateCheck = time.Now()
	cfg.SaveConfig() // Best effort - ignore errors
	
	// Show update notification if available
	if updateInfo.UpdateNeeded {
		fmt.Fprintf(os.Stderr, "\nUpdate available: %s → %s\n", updateInfo.CurrentVersion, updateInfo.LatestVersion)
		fmt.Fprintf(os.Stderr, "Run 'ctx update' to install the latest version.\n\n")
	}
}