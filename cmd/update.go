package cmd

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/slavakurilyak/ctx/internal/updater"
	"github.com/spf13/cobra"
)

// NewUpdateCmd creates the update command
func NewUpdateCmd(version string) *cobra.Command {
	var (
		checkOnly      bool
		includePrerelease bool
		force         bool
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update ctx to the latest version",
		Long: `Check for and install the latest version of ctx.

This command connects to GitHub to check for newer releases and can automatically 
download and install updates. Use --check to only check for updates without installing.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create updater instance
			upd := updater.NewUpdater("slavakurilyak", "ctx")
			
			// Check for updates
			fmt.Println("Checking for updates...")
			updateInfo, err := upd.CheckForUpdate(version, includePrerelease)
			if err != nil {
				return fmt.Errorf("failed to check for updates: %w", err)
			}

			// Display current and latest versions
			fmt.Printf("Current version: %s\n", updateInfo.CurrentVersion)
			fmt.Printf("Latest version:  %s\n", updateInfo.LatestVersion)

			if !updateInfo.UpdateNeeded && !force {
				fmt.Println("✓ You are running the latest version!")
				return nil
			}

			if updateInfo.UpdateNeeded {
				fmt.Printf("🔄 Update available: %s → %s\n", updateInfo.CurrentVersion, updateInfo.LatestVersion)
			}

			// If only checking, stop here
			if checkOnly {
				if updateInfo.UpdateNeeded {
					fmt.Println("\nTo install the update, run: ctx update")
				}
				return nil
			}

			// Check if we can update (has download URL)
			if updateInfo.UpdateURL == "" {
				fmt.Println("❌ No binary available for your platform")
				fmt.Printf("   Platform: %s\n", getPlatformString())
				fmt.Println("   Please download manually from: https://github.com/slavakurilyak/ctx/releases")
				return nil
			}

			// Show what we're about to do
			if updateInfo.UpdateNeeded {
				fmt.Printf("\nDownloading %s...\n", updateInfo.LatestVersion)
			} else if force {
				fmt.Printf("\nForce reinstalling %s...\n", updateInfo.LatestVersion)
			}

			// Perform the update
			err = upd.PerformUpdate(updateInfo)
			if err != nil {
				return fmt.Errorf("update failed: %w", err)
			}

			fmt.Printf("✅ Successfully updated to %s!\n", updateInfo.LatestVersion)
			
			// Show release notes if available
			if updateInfo.ReleaseNotes != "" && updateInfo.UpdateNeeded {
				fmt.Println("\n📋 Release Notes:")
				fmt.Println(strings.Repeat("-", 50))
				fmt.Println(updateInfo.ReleaseNotes)
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&checkOnly, "check", false, "Only check for updates, don't install")
	cmd.Flags().BoolVar(&includePrerelease, "pre-release", false, "Include pre-release versions")
	cmd.Flags().BoolVar(&force, "force", false, "Force reinstall current version")

	return cmd
}

// getPlatformString returns a human-readable platform string
func getPlatformString() string {
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	
	osName := map[string]string{
		"darwin":  "macOS",
		"linux":   "Linux", 
		"windows": "Windows",
	}[goos]
	if osName == "" {
		osName = goos
	}
	
	archName := map[string]string{
		"amd64": "Intel/AMD64",
		"arm64": "ARM64",
		"386":   "32-bit",
	}[goarch]
	if archName == "" {
		archName = goarch
	}
	
	return fmt.Sprintf("%s/%s", osName, archName)
}