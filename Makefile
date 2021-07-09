VERSION := $(shell cat VERSION)
GIT_HASH := $(shell git rev-parse HEAD)
DATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

VERSION_INFO_IMPORT_PATH=github.com/observIQ/observIQ-otel-collector/internal/version

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
MODNAME=github.com/observIQ/observIQ-otel-collector

LINT=$(GOPATH)/bin/golangci-lint
MISSPELL=$(GOPATH)/bin/misspell

LDFLAGS=-ldflags "-s -w -X $(VERSION_INFO_IMPORT_PATH).Version=$(VERSION) \
-X $(VERSION_INFO_IMPORT_PATH).GitHash=$(GIT_HASH) \
-X $(VERSION_INFO_IMPORT_PATH).Date=$(DATE)"
GOBUILDEXTRAENV=GO111MODULE=on CGO_ENABLED=0
GOBUILD=go build
GOINSTALL=go install
GOTEST=go test
GOTOOL=go tool
GOFORMAT=goimports
GOTIDY=go mod tidy

# Default build target; making this should build for the current os/arch
.PHONY: observiqcol
observiqcol:
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOBUILDEXTRAENV) \
	$(GOBUILD) $(LDFLAGS) -o $(OUTDIR)/observiqcol_$(GOOS)_$(GOARCH)$(EXT) ./cmd/observiqcol

# Other build targets
.PHONY: amd64_linux
amd64_linux:
	GOOS=linux GOARCH=amd64 $(MAKE) observiqcol

.PHONY: amd64_darwin
amd64_darwin:
	GOOS=darwin GOARCH=amd64 $(MAKE) observiqcol

.PHONY: arm_linux
arm_linux:
	GOOS=linux GOARCH=arm $(MAKE) observiqcol

.PHONY: amd64_windows
amd64_windows:
	GOOS=windows GOARCH=amd64 $(MAKE) observiqcol

.PHONY: build-all
build-all: amd64_linux amd64_darwin amd64_windows arm_linux

# tool-related commands
.PHONY: install-tools
install-tools:
	$(GOINSTALL) golang.org/x/tools/cmd/goimports
	$(GOINSTALL) github.com/golangci/golangci-lint/cmd/golangci-lint@v1.40.1
	$(GOINSTALL) github.com/client9/misspell/cmd/misspell

.PHONY: lint
lint:
	$(LINT) run

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
	$(GOTEST) -vet off -cover cover.out ./...
	$(GOTOOL) cover -html=cover.out -o cover.html

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
	$(GOTIDY)

# This target performs all checks that CI will do
.PHONY: ci-checks
ci-checks: check-fmt misspell lint test

.PHONY: clean
clean:
	rm -rf $(OUTDIR)
