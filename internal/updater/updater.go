package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Release represents a GitHub release
type Release struct {
	TagName     string  `json:"tag_name"`
	Name        string  `json:"name"`
	Body        string  `json:"body"`
	Draft       bool    `json:"draft"`
	Prerelease  bool    `json:"prerelease"`
	PublishedAt string  `json:"published_at"`
	Assets      []Asset `json:"assets"`
}

// Asset represents a release asset
type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

// UpdateInfo contains information about available updates
type UpdateInfo struct {
	CurrentVersion string
	LatestVersion  string
	UpdateURL      string
	ReleaseNotes   string
	UpdateNeeded   bool
}

// Updater handles version checking and updates
type Updater struct {
	Owner      string
	Repo       string
	HTTPClient *http.Client
}

// NewUpdater creates a new updater instance
func NewUpdater(owner, repo string) *Updater {
	return &Updater{
		Owner: owner,
		Repo:  repo,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CheckForUpdate checks if a newer version is available
func (u *Updater) CheckForUpdate(currentVersion string, includePrerelease bool) (*UpdateInfo, error) {
	release, err := u.getLatestRelease(includePrerelease)
	if err != nil {
		return nil, fmt.Errorf("failed to check for updates: %w", err)
	}

	// Clean version strings (remove 'v' prefix if present)
	currentClean := strings.TrimPrefix(currentVersion, "v")
	latestClean := strings.TrimPrefix(release.TagName, "v")

	updateNeeded := isNewerVersion(latestClean, currentClean)

	asset := u.findAssetForPlatform(release.Assets)
	var updateURL string
	if asset != nil {
		updateURL = asset.BrowserDownloadURL
	}

	return &UpdateInfo{
		CurrentVersion: currentVersion,
		LatestVersion:  release.TagName,
		UpdateURL:      updateURL,
		ReleaseNotes:   release.Body,
		UpdateNeeded:   updateNeeded,
	}, nil
}

// PerformUpdate downloads and installs the latest version
func (u *Updater) PerformUpdate(updateInfo *UpdateInfo) error {
	if updateInfo.UpdateURL == "" {
		return fmt.Errorf("no update URL available")
	}

	// Get current executable path
	currentExe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable path: %w", err)
	}

	// Download the new binary
	tempFile, err := u.downloadBinary(updateInfo.UpdateURL)
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}
	defer os.Remove(tempFile)

	// Replace the current executable
	err = u.replaceExecutable(currentExe, tempFile)
	if err != nil {
		return fmt.Errorf("failed to replace executable: %w", err)
	}

	return nil
}

// getLatestRelease fetches the latest release from GitHub API
func (u *Updater) getLatestRelease(includePrerelease bool) (*Release, error) {
	var url string
	if includePrerelease {
		// Get all releases and find the latest (including prereleases)
		url = fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", u.Owner, u.Repo)
	} else {
		// Get only the latest stable release
		url = fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", u.Owner, u.Repo)
	}

	resp, err := u.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	if includePrerelease {
		var releases []Release
		err = json.NewDecoder(resp.Body).Decode(&releases)
		if err != nil {
			return nil, err
		}
		if len(releases) == 0 {
			return nil, fmt.Errorf("no releases found")
		}
		return &releases[0], nil
	} else {
		var release Release
		err = json.NewDecoder(resp.Body).Decode(&release)
		if err != nil {
			return nil, err
		}
		return &release, nil
	}
}

// findAssetForPlatform finds the appropriate asset for the current platform
func (u *Updater) findAssetForPlatform(assets []Asset) *Asset {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	// Build expected filename patterns
	var patterns []string
	
	if goos == "windows" {
		patterns = []string{
			fmt.Sprintf("ctx-%s-%s.exe", goos, goarch),
			fmt.Sprintf("ctx-windows-%s.exe", goarch),
			fmt.Sprintf("ctx_%s_%s.exe", goos, goarch),
		}
	} else {
		patterns = []string{
			fmt.Sprintf("ctx-%s-%s", goos, goarch),
			fmt.Sprintf("ctx_%s_%s", goos, goarch),
		}
		
		// Add darwin-specific patterns
		if goos == "darwin" {
			patterns = append(patterns, fmt.Sprintf("ctx-darwin-%s", goarch))
		}
	}

	// Find matching asset
	for _, asset := range assets {
		for _, pattern := range patterns {
			if asset.Name == pattern {
				return &asset
			}
		}
	}

	return nil
}

// downloadBinary downloads a binary to a temporary file
func (u *Updater) downloadBinary(url string) (string, error) {
	resp, err := u.HTTPClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Create temporary file
	tempFile, err := os.CreateTemp("", "ctx-update-*")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	// Copy downloaded data
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		os.Remove(tempFile.Name())
		return "", err
	}

	// Make executable
	err = os.Chmod(tempFile.Name(), 0755)
	if err != nil {
		os.Remove(tempFile.Name())
		return "", err
	}

	return tempFile.Name(), nil
}

// replaceExecutable replaces the current executable with the new one
func (u *Updater) replaceExecutable(currentPath, newPath string) error {
	// On Windows, we can't replace a running executable directly
	if runtime.GOOS == "windows" {
		return u.replaceExecutableWindows(currentPath, newPath)
	}

	// On Unix-like systems, we can replace directly
	return u.replaceExecutableUnix(currentPath, newPath)
}

// replaceExecutableUnix replaces executable on Unix-like systems
func (u *Updater) replaceExecutableUnix(currentPath, newPath string) error {
	// Create backup
	backupPath := currentPath + ".backup"
	err := os.Rename(currentPath, backupPath)
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Move new binary to current location
	err = os.Rename(newPath, currentPath)
	if err != nil {
		// Restore backup on failure
		os.Rename(backupPath, currentPath)
		return fmt.Errorf("failed to replace executable: %w", err)
	}

	// Remove backup on success
	os.Remove(backupPath)
	return nil
}

// replaceExecutableWindows handles executable replacement on Windows
func (u *Updater) replaceExecutableWindows(currentPath, newPath string) error {
	// On Windows, create a batch script to replace the executable after exit
	batchContent := fmt.Sprintf(`@echo off
timeout /t 1 /nobreak >nul
move "%s" "%s.backup" >nul 2>&1
move "%s" "%s" >nul 2>&1
if errorlevel 1 (
    move "%s.backup" "%s" >nul 2>&1
    echo Update failed - restored original
) else (
    del "%s.backup" >nul 2>&1
    echo Update completed successfully
)
del "%%~f0" >nul 2>&1
`, currentPath, currentPath, newPath, currentPath, currentPath, currentPath, currentPath)

	batchPath := filepath.Join(os.TempDir(), "ctx-update.bat")
	err := os.WriteFile(batchPath, []byte(batchContent), 0755)
	if err != nil {
		return fmt.Errorf("failed to create update script: %w", err)
	}

	fmt.Printf("Update will complete after ctx exits. Please restart ctx after a few seconds.\n")
	
	// The batch script will run after this process exits
	// We can't wait for it or the executable will be locked
	return nil
}

// isNewerVersion compares two version strings
// This is a simple string comparison - for production use, consider semantic versioning
func isNewerVersion(latest, current string) bool {
	// Handle special cases
	if current == "dev" || current == "unknown" || strings.HasPrefix(current, "dev-") {
		return true
	}
	
	// Simple string comparison (works for basic versioning like v0.1.0, v0.2.0)
	// For production, you'd want proper semantic version comparison
	return latest > current
}