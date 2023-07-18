# Releasing
Releases are managed through GitHub releases. The steps to create a release are as follows:

1. Run `make version={VERSION} release` where `{VERSION}` is the version to release. This will create a tag and push it to GitHub.

2. The CD job will run for the tagged commit. Goreleaser will handle the following without user intervention:
  - Build the binaries
  - Create GitHub release with automatic changelog content (based on commit message)
  - Attach binaries to the GitHub release
  - Mark the release as a full release once it is finished.
  - Create a CHANGELOG

3. Done! The agent is released

# Testing Release locally

In order to run the `make release-test` you need to setup the following in your environment:

1. Run `make install-tools`
2. Setup a github token per [goreleaser instructions](https://goreleaser.com/scm/github/#api-token)
3. Run `make release-test`
