# Releasing
Releases are managed through GitHub releases. The steps to create a release are as follows:

1. Update CHANGELOG.md with changes

2. Create a new tag (not release)

3. The CD job will run for the tagged commit. Goreleaser will handle the following without user intervention:
  - Build the binaries
  - Create GitHub release with automatic changelog content (based on commit message)
  - Attach binaries to the GitHub release
  - Mark the release as a full release once it is finished.

4. Done! The collector is released