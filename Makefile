VERSION := $(shell cat VERSION)
GIT_HASH := $(shell git rev-parse HEAD)
DATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

VERSION_INFO_IMPORT_PATH=github.com/observiq/observiq-collector/internal/version

# All source code and documents, used when checking for misspellings
ALLDOC := $(shell find . \( -name "*.md" -o -name "*.yaml" \) \
                                -type f | sort)

GOPATH ?= $(shell go env GOPATH)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

ifeq ($(GOOS), windows)
EXT?=.exe
else
EXT?=
endif

OUTDIR=./build
MODNAME=github.com/observiq/observiq-collector

LINT=$(GOPATH)/bin/golangci-lint
LINT_TIMEOUT?=5m0s
MISSPELL=$(GOPATH)/bin/misspell

LDFLAGS=-ldflags "-s -w -X $(VERSION_INFO_IMPORT_PATH).version=$(VERSION) \
-X $(VERSION_INFO_IMPORT_PATH).gitHash=$(GIT_HASH) \
-X $(VERSION_INFO_IMPORT_PATH).date=$(DATE)"
GOBUILDEXTRAENV=GO111MODULE=on CGO_ENABLED=0
GOBUILD=go build
GOINSTALL=go install
GOTEST=go test
GOTOOL=go tool
GOFORMAT=goimports

# Default build target; making this should build for the current os/arch
.PHONY: collector
collector:
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOBUILDEXTRAENV) \
	$(GOBUILD) $(LDFLAGS) -o $(OUTDIR)/collector_$(GOOS)_$(GOARCH)$(EXT) ./cmd/collector

# Other build targets
.PHONY: amd64_linux
amd64_linux:
	GOOS=linux GOARCH=amd64 $(MAKE) collector

.PHONY: amd64_darwin
amd64_darwin:
	GOOS=darwin GOARCH=amd64 $(MAKE) collector

.PHONY: arm_linux
arm_linux:
	GOOS=linux GOARCH=arm $(MAKE) collector

.PHONY: amd64_windows
amd64_windows:
	GOOS=windows GOARCH=amd64 $(MAKE) collector

.PHONY: build-all
build-all: amd64_linux amd64_darwin amd64_windows arm_linux

# tool-related commands
TOOLS_MOD_DIR := ./internal/tools
.PHONY: install-tools
install-tools:
	cd $(TOOLS_MOD_DIR) && $(GOINSTALL) golang.org/x/tools/cmd/goimports
	cd $(TOOLS_MOD_DIR) && $(GOINSTALL) github.com/golangci/golangci-lint/cmd/golangci-lint@v1.40.1
	cd $(TOOLS_MOD_DIR) && $(GOINSTALL) github.com/client9/misspell/cmd/misspell

.PHONY: lint
lint:
	$(LINT) run --timeout $(LINT_TIMEOUT)

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
	$(MAKE) for-all CMD="rm -fr go.sum"
	$(MAKE) for-all CMD="go mod tidy -go=1.16"
	$(MAKE) for-all CMD="go mod tidy -go=1.17"

# This target performs all checks that CI will do (excluding the build itself)
.PHONY: ci-checks
ci-checks: check-fmt misspell lint test

.PHONY: clean
clean:
	rm -rf $(OUTDIR)

.PHONY: for-all
for-all:
	@echo "running $${CMD} in root"
	@$${CMD}
	@set -e; for dir in $(ALL_MODULES); do \
	  (cd "$${dir}" && \
	  	echo "running $${CMD} in $${dir}" && \
	 	$${CMD} ); \
	done