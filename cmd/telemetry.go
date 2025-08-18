package cmd

import (
	"fmt"

	"github.com/slavakurilyak/ctx/internal/config"
	"github.com/spf13/cobra"
)

func NewTelemetryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "telemetry",
		Short: "Manage and view telemetry settings",
		Long:  "Provides tools to check the status of OpenTelemetry data collection.",
	}
	cmd.AddCommand(newTelemetryStatusCmd())
	return cmd
}

func newTelemetryStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show the current telemetry status",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config.NewFromFlagsAndEnv(cmd.Root())

			status := "Enabled"
			if cfg.NoTelemetry {
				status = "Disabled"
			}
			fmt.Printf("Telemetry: %s (Source: %s)\n", status, cfg.NoTelemetrySource)
		},
	}
}
