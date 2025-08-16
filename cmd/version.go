package cmd

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/slavakurilyak/ctx/internal/models"
	"github.com/spf13/cobra"
)

// VersionInfo contains all version-related information
type VersionInfo struct {
	CTXVersion    string `json:"ctx_version"`
	SchemaVersion string `json:"schema_version"`
	Commit        string `json:"commit"`
	BuildDate     string `json:"build_date"`
	GoVersion     string `json:"go_version"`
	OS            string `json:"os"`
	Arch          string `json:"arch"`
}

// NewVersionCmd creates the version subcommand
func NewVersionCmd(version, commit, date string) *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show ctx version information",
		Long:  `Display detailed version information about ctx, including software version, schema version, and build details.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			info := VersionInfo{
				CTXVersion:    version,
				SchemaVersion: models.CurrentSchemaVersion,
				Commit:        commit,
				BuildDate:     date,
				GoVersion:     runtime.Version(),
				OS:            runtime.GOOS,
				Arch:          runtime.GOARCH,
			}

			// Handle development builds
			if info.CTXVersion == "dev" || info.CTXVersion == "" {
				info.CTXVersion = "dev (built from source)"
			}
			if info.Commit == "none" || info.Commit == "" {
				info.Commit = "unknown"
			}
			if info.BuildDate == "unknown" || info.BuildDate == "" {
				info.BuildDate = "unknown"
			}

			if jsonOutput {
				// JSON output for programmatic use
				output, err := json.MarshalIndent(info, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal version info: %w", err)
				}
				fmt.Println(string(output))
			} else {
				// Human-readable output
				fmt.Printf("ctx version information:\n")
				fmt.Printf("  Software Version: %s\n", info.CTXVersion)
				fmt.Printf("  Schema Version:   %s\n", info.SchemaVersion)
				fmt.Printf("  Commit:           %s\n", info.Commit)
				fmt.Printf("  Build Date:       %s\n", info.BuildDate)
				fmt.Printf("  Go Version:       %s\n", info.GoVersion)
				fmt.Printf("  Platform:         %s/%s\n", info.OS, info.Arch)
				
				// Add compatibility note
				fmt.Printf("\nCompatibility:\n")
				fmt.Printf("  This version outputs JSON with schema version %s\n", info.SchemaVersion)
				fmt.Printf("  Use 'ctx <command>' to see the schema_version in output\n")
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output version information as JSON")

	return cmd
}