export GIT_HASH = $(shell git rev-parse HEAD)
export DATE = $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
export VERSION_INFO_IMPORT_PATH = github.com/observiq/observiq-collector/internal/version

# All source code and documents, used when checking for misspellings
ALLDOC := $(shell find . \( -name "*.md" -o -name "*.yaml" \) \
                                -type f | sort)

# All source code files
ALL_SRC := $(shell find . -name '*.go' -o -name '*.sh' -o -name 'Dockerfile' -type f | sort)

GOPATH ?= $(shell go env GOPATH)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

ifeq ($(GOOS), windows)
EXT?=.exe
else
EXT?=
endif

LINT=$(GOPATH)/bin/golangci-lint
LINT_TIMEOUT?=5m0s
MISSPELL=$(GOPATH)/bin/misspell

GOINSTALL=go install
GOTEST=go test
GOTOOL=go tool
GOFORMAT=goimports
ALL_MODULES := $(shell find . -type f -name "go.mod" -exec dirname {} \; | sort )

# tool-related commands
TOOLS_MOD_DIR := ./internal/tools
.PHONY: install-tools
install-tools:
	cd $(TOOLS_MOD_DIR) && go install github.com/mgechev/revive@latest 
	cd $(TOOLS_MOD_DIR) && go install github.com/google/addlicense
	cd $(TOOLS_MOD_DIR) && $(GOINSTALL) golang.org/x/tools/cmd/goimports	
	cd $(TOOLS_MOD_DIR) && $(GOINSTALL) github.com/client9/misspell/cmd/misspell
	cd $(TOOLS_MOD_DIR) && $(GOINSTALL) github.com/sigstore/cosign/cmd/cosign
	cd $(TOOLS_MOD_DIR) && $(GOINSTALL) github.com/goreleaser/goreleaser@v1.3.1

.PHONY: lint
lint:
	revive -config revive/config.toml -formatter friendly ./...

.PHONY: misspell
misspell:
	$(MISSPELL) $(ALLDOC)

.PHONY: misspell-fix
misspell-fix:
	$(MISSPELL) -w $(ALLDOC)

.PHONY: test
test:
	$(GOTEST) -vet off -race ./...

.PHONY: test-with-cover
test-with-cover:
	$(GOTEST) -vet off -coverprofile=cover.out ./...
	$(GOTOOL) cover -html=cover.out -o cover.html

.PHONY: bench
bench:
	go test -benchmem -run=^$$ -bench ^* ./...

.PHONY: check-fmt
check-fmt:
	@GOFMTOUT=`$(GOFORMAT) -d .`; \
		if [ "$$GOFMTOUT" ]; then \
			echo "$(GOFORMAT) SUGGESTED CHANGES:"; \
			echo "$$GOFMTOUT\n"; \
			exit 1; \
		else \
			echo "$(GOFORMAT) completed successfully"; \
		fi

.PHONY: fmt
fmt:
	$(GOFORMAT) -w .

.PHONY: tidy
tidy:
	$(MAKE) for-all CMD="go mod tidy -go=1.16"
	$(MAKE) for-all CMD="go mod tidy -go=1.17"

# This target performs all checks that CI will do (excluding the build itself)
.PHONY: ci-checks
ci-checks: check-fmt misspell lint test

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

.PHONY: clean
clean:
	rm -rf $(OUTDIR)

# Default build target; making this should build for the current os/arch
.PHONY: collector
collector:
	goreleaser build --single-target --skip-validate --snapshot --rm-dist

.PHONY: build-all
build-all:
	goreleaser build --skip-validate --snapshot --rm-dist

# Build, sign, and release
.PHONY: release
release:
	goreleaser release --parallelism 4 --rm-dist

# Build and sign, skip release and ignore dirty git tree
.PHONY: release-test
release-test:
	goreleaser release --parallelism 4 --skip-validate --skip-publish --rm-dist

.PHONY: for-all
for-all:
	@echo "running $${CMD} in root"
	@$${CMD}
	@set -e; for dir in $(ALL_MODULES); do \
	  (cd "$${dir}" && \
	  	echo "running $${CMD} in $${dir}" && \
	 	$${CMD} ); \
	done
