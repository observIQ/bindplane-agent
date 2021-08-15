package version

// these will be replaced at link time by make.
var (
	version = "latest"  // Semantic version, or "latest" by default
	gitHash = "unknown" // Commit hash from which this build was generated
	date    = "unknown" // Date the build was generated
)

// Version returns the version of the collector.
func Version() string {
	return version
}

// GitHash returns the githash associated with the collector's version.
func GitHash() string {
	return gitHash
}

// Date returns the publish date associated with the collector's version.
func Date() string {
	return date
}
