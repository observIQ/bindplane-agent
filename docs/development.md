# Development

## Initial Setup

Clone this repository, and run `make install-tools`

## Building

To create a build for your current machine, run `make collector`

To build for a specific architecture, see the [Makefile](./Makefile)

To build all targets, run `make build-all`

Build files will show up in the `./dist` directory

## Running Tests

Tests can be run with `make test`.

## Running CI checks locally

The CI runs the `ci-checks` make target, which includes linting, testing, and checking documentation for misspelling.
CI also does a build of all targets (`make build-all`)

## Releasing
To release the collector, see [releasing documentation](releasing.md).
