# Windows MSI building

This directory contains sources for building the Windows MSI using [go-msi](https://github.com/observIQ/go-msi/) and the [Wix toolset](https://wixtoolset.org/).

## Building Locally with Vagrant

A local build may be performed with [vagrant](https://www.vagrantup.com/). 

The following make targets are available for local development:
* `vagrant-prep`: Starts up the vagrant box and prepares it for building and testing. The vagrant box must be up in order for building or testing to work. **PLEASE NOTE** that valid Windows licensing is your responsibility.
* `fetch-dependencies`: Fetches dependencies for building the MSI.
* `build-msi`: Builds the MSI. Depends on the `fetch-dependencies` target (`fetch-dependencies` will be run every time this is run).
* `test-install-msi`: Test installing the MSI. `build-msi` should be run before this is run.
* `test-uninstall-msi`: Test uninstalling the MSI. `build-msi` should be run before this is run, and the MSI should be installed (e.g. by running `test-install-msi`)
