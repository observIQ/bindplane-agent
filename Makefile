# All source code and documents, used when checking for misspellings
ALLDOC := $(shell find . \( -name "*.md" -o -name "*.yaml" \) \
                                -type f | sort)
ALL_MODULES := $(shell find . -type f -name "go.mod" -exec dirname {} \; | sort )
ALL_MDATAGEN_MODULES := $(shell find . -type f -name "metadata.yaml" -exec dirname {} \; | sort )

# All source code files
ALL_SRC := $(shell find . -name '*.go' -o -name '*.sh' -o -name 'Dockerfile*' -type f | sort)

OUTDIR=./dist
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

INTEGRATION_TEST_ARGS?=-tags integration

TOOLS_MOD_DIR := ./internal/tools

ifeq ($(GOOS), windows)
EXT?=.exe
else
EXT?=
endif

PREVIOUS_TAG := $(shell git tag --sort=v:refname --no-contains HEAD | grep -E "[0-9]+\.[0-9]+\.[0-9]+$$" | tail -n1)
CURRENT_TAG := $(shell git tag --sort=v:refname --points-at HEAD | grep -E "v[0-9]+\.[0-9]+\.[0-9]+$$" | tail -n1)
# Version will be the tag pointing to the current commit, or the previous version tag if there is no such tag
VERSION ?= $(if $(CURRENT_TAG),$(CURRENT_TAG),$(PREVIOUS_TAG))

# Build binaries for current GOOS/GOARCH by default
.DEFAULT_GOAL := build-binaries

# Builds just the agent for current GOOS/GOARCH pair
.PHONY: agent
agent:
	go build -ldflags "-s -w -X github.com/observiq/bindplane-agent/internal/version.version=$(VERSION)" -tags bindplane -o $(OUTDIR)/collector_$(GOOS)_$(GOARCH)$(EXT) ./cmd/collector

# Builds just the updater for current GOOS/GOARCH pair
.PHONY: updater
updater:
	cd ./updater/; go build -ldflags "-s -w -X github.com/observiq/bindplane-agent/internal/version.version=$(VERSION)" -o ../$(OUTDIR)/updater_$(GOOS)_$(GOARCH)$(EXT) ./cmd/updater

# Builds the updater + agent for current GOOS/GOARCH pair
.PHONY: build-binaries
build-binaries: agent updater

.PHONY: build-all
build-all: build-linux build-darwin build-windows

.PHONY: build-linux
build-linux: build-linux-amd64 build-linux-arm64 build-linux-arm build-linux-ppc64 build-linux-ppc64le

.PHONY: build-darwin
build-darwin: build-darwin-amd64 build-darwin-arm64

.PHONY: build-windows
build-windows: build-windows-amd64

.PHONY: build-linux-ppc64
build-linux-ppc64:
	GOOS=linux GOARCH=ppc64 $(MAKE) build-binaries -j2

.PHONY: build-linux-ppc64le
build-linux-ppc64le:
	GOOS=linux GOARCH=ppc64le $(MAKE) build-binaries -j2

.PHONY: build-linux-amd64
build-linux-amd64:
	GOOS=linux GOARCH=amd64 $(MAKE) build-binaries -j2

.PHONY: build-linux-arm64
build-linux-arm64:
	GOOS=linux GOARCH=arm64 $(MAKE) build-binaries -j2

.PHONY: build-linux-arm
build-linux-arm:
	GOOS=linux GOARCH=arm $(MAKE) build-binaries -j2

.PHONY: build-darwin-amd64
build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 $(MAKE) build-binaries -j2

.PHONY: build-darwin-arm64
build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 $(MAKE) build-binaries -j2

.PHONY: build-windows-amd64
build-windows-amd64:
	GOOS=windows GOARCH=amd64 $(MAKE) build-binaries -j2

# tool-related commands
.PHONY: install-tools
install-tools:
	cd $(TOOLS_MOD_DIR) && go install github.com/client9/misspell/cmd/misspell
	cd $(TOOLS_MOD_DIR) && go install github.com/google/addlicense
	cd $(TOOLS_MOD_DIR) && go install github.com/mgechev/revive
	cd $(TOOLS_MOD_DIR) && go install go.opentelemetry.io/collector/cmd/mdatagen
	cd $(TOOLS_MOD_DIR) && go install github.com/securego/gosec/v2/cmd/gosec
# update cosign in release.yml when updating this version
# update cosign in docs/verify-signature.md when updating this version
	go install github.com/sigstore/cosign/cmd/cosign@v1.13.1
	cd $(TOOLS_MOD_DIR) && go install github.com/uw-labs/lichen
	cd $(TOOLS_MOD_DIR) && go install github.com/vektra/mockery/v2
	cd $(TOOLS_MOD_DIR) && go install golang.org/x/tools/cmd/goimports

.PHONY: lint
lint:
	revive -config revive/config.toml -formatter friendly ./...

.PHONY: misspell
misspell:
	misspell -error $(ALLDOC)

.PHONY: misspell-fix
misspell-fix:
	misspell -w $(ALLDOC)

.PHONY: test
test:
	$(MAKE) for-all CMD="go test -race ./..."

.PHONY: test-with-cover
test-with-cover:
	$(MAKE) for-all CMD="go test -coverprofile=cover.out ./..."
	$(MAKE) for-all CMD="go tool cover -html=cover.out -o cover.html"

.PHONY: test-updater-integration
test-updater-integration:
	cd updater; go test $(INTEGRATION_TEST_ARGS) -race ./...

.PHONY: bench
bench:
	$(MAKE) for-all CMD="go test -benchmem -run=^$$ -bench ^* ./..."

.PHONY: check-fmt
check-fmt:
	goimports -d ./ | diff -u /dev/null -

.PHONY: fmt
fmt:
	goimports -w .

.PHONY: tidy
tidy:
	$(MAKE) for-all CMD="go mod tidy -compat=1.22"

.PHONY: gosec
gosec:
	gosec \
	  -exclude-dir=updater \
	  -exclude-dir=receiver/sapnetweaverreceiver \
	  -exclude-dir=extension/bindplaneextension \
	  -exclude-dir=processor/snapshotprocessor \
	  -exclude-dir=internal/tools \
	  ./...
# exclude the testdata dir; it contains a go program for testing.
	cd updater; gosec -exclude-dir internal/service/testdata ./...
	cd extension/bindplaneextension; gosec ./...
	cd processor/snapshotprocessor; gosec ./...
	cd receiver/sapnetweaverreceiver; gosec ./...

# This target performs all checks that CI will do (excluding the build itself)
.PHONY: ci-checks
ci-checks: check-fmt check-license misspell lint gosec test

# This target checks that license copyright header is on every source file
.PHONY: check-license
check-license:
	@ADDLICENSEOUT=`addlicense -check $(ALL_SRC) 2>&1`; \
		if [ "$$ADDLICENSEOUT" ]; then \
			echo "addlicense FAILED => add License errors:\n"; \
			echo "$$ADDLICENSEOUT\n"; \
			echo "Use 'make add-license' to fix this."; \
			exit 1; \
		else \
			echo "Check License finished successfully"; \
		fi

# This target adds a license copyright header is on every source file that is missing one
.PHONY: add-license
add-license:
	@ADDLICENSEOUT=`addlicense -y "" -c "observIQ, Inc." $(ALL_SRC) 2>&1`; \
		if [ "$$ADDLICENSEOUT" ]; then \
			echo "addlicense FAILED => add License errors:\n"; \
			echo "$$ADDLICENSEOUT\n"; \
			exit 1; \
		else \
			echo "Add License finished successfully"; \
		fi

# update-otel attempts to update otel dependencies in go.mods,
# and update the otel versions in the docs.
# Usage: make update-otel OTEL_VERSION=vx.x.x CONTRIB_VERSION=vx.x.x PDATA_VERSION=vx.x.x-rcx
.PHONY: update-otel
update-otel:
	./scripts/update-otel.sh "$(OTEL_VERSION)" "$(CONTRIB_VERSION)" "$(PDATA_VERSION)"
	./scripts/update-docs.sh "$(OTEL_VERSION)" "$(CONTRIB_VERSION)"
	$(MAKE) tidy
# Double make tidy - this unfortunately is needed due to the order in which modules are tidied.
# The modules this seems to effect are plugindocgen and bindplaneextension
	$(MAKE) tidy

# update-modules updates all submodules to be the new version.
# Usage: make update-modules NEW_VERSION=vx.x.x
.PHONY: update-modules
update-modules:
	./scripts/update-module-version.sh "$(NEW_VERSION)"
	$(MAKE) tidy

# Downloads and setups dependencies that are packaged with binary
.PHONY: release-prep
release-prep:
	@rm -rf release_deps
	@mkdir release_deps
	@echo 'v$(CURR_VERSION)' > release_deps/VERSION.txt
	./buildscripts/download-dependencies.sh release_deps
	@cp -r ./plugins release_deps/
	@cp config/example.yaml release_deps/config.yaml
	@cp config/logging.yaml release_deps/logging.yaml
	@cp service/com.observiq.collector.plist release_deps/com.observiq.collector.plist
	@jq ".files[] | select(.service != null)" windows/wix.json >> release_deps/windows_service.json
	@cp service/observiq-otel-collector.service release_deps/observiq-otel-collector.service
	@cp service/observiq-otel-collector release_deps/observiq-otel-collector
	@cp -r ./service/sysconfig release_deps/

# Build and sign, skip release and ignore dirty git tree
.PHONY: release-test
release-test:
# If there is no MSI in the root dir, we'll create a dummy one so that goreleaser can complete successfully
	if [ ! -e "./observiq-otel-collector.msi" ]; then touch ./observiq-otel-collector.msi; fi
	GORELEASER_CURRENT_TAG=$(VERSION) goreleaser release --parallelism 4 --skip=publish --skip=validate --skip=sign --clean --snapshot

.PHONY: for-all
for-all:
	@echo "running $${CMD} in root"
	@$${CMD}
	@set -e; for dir in $(ALL_MODULES); do \
	  (cd "$${dir}" && \
	  	echo "running $${CMD} in $${dir}" && \
	 	$${CMD} ); \
	done

# Release a new version of the agent. This will also tag all submodules
.PHONY: release
release:
	@if [ -z "$(version)" ]; then \
		echo "version was not set"; \
		exit 1; \
	fi

	@if ! [[ "$(version)" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$$ ]]; then \
		echo "version $(version) is invalid semver"; \
		exit 1; \
	fi

	@git tag $(version)
	@git push --tags

	@set -e; for dir in $(ALL_MODULES); do \
	  if [ $${dir} == \. ]; then \
	  	continue; \
	  else \
	    echo "$${dir}" | sed -e "s+^./++" -e 's+$$+/$(version)+' | awk '{print $1}' | git tag $$(cat)  ; \
	  fi; \
	done

	@git push --tags

.PHONY: clean
clean:
	rm -rf $(OUTDIR)

.PHONY: scan-licenses
scan-licenses:
	lichen --config=./license.yaml $$(find dist/collector_* dist/updater_*)

.PHONY: generate
generate:
	$(MAKE) for-all CMD="go generate ./..."
	$(MAKE) fmt

.PHONY: create-plugin-docs
create-plugin-docs:
	cd cmd/plugindocgen; go run .
