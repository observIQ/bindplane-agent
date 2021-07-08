package version

// these will be replaced at link time by make.
var (
	Version = "latest"  // Semantic version, or "latest" by default
	GitHash = "unknown" // Commit hash from which this build was generated
	Date    = "unknown" // Date the build was generated
)
