package version

var (
	// Set by ldflags during build
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

func GetVersion() string {
	if Version == "dev" || Version == "" {
		return "dev (built from source)"
	}
	return Version
}

func GetUserAgent() string {
	return "ctx-cli/" + GetVersion()
}
