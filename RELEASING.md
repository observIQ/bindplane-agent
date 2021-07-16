# Releasing
Releases are managed through GitHub releases. The steps to create a release are as follows:

1. Update CHANGELOG.md with changes, and bump the VERSION file as appropriate (follow Semantic Versioning)
2. On GitHub, create a release with a tag the same as the version (e.g. v0.0.1 or 0.0.1), and the contents the same as the new entry in CHANGELOG.md. **Mark this release as pre-release**
4. The CD job will run for the tagged commit, run tests, build and attach the binaries, and mark the release as a full release once it is finished.
5. Done! 