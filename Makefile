BINARY := env-vault

BUILD_DATE := $(shell date -u "+%Y-%m-%dT%H:%M:%SZ")

GIT_COMMIT := $(shell git rev-parse --short HEAD)
GIT_DIRTY := $(if $(shell git status --porcelain),+CHANGES)

GO_MODULE := github.com/romantomjak/$(BINARY)
GO_LDFLAGS := "-X $(GO_MODULE)/version.BuildDate=$(BUILD_DATE) -X $(GO_MODULE)/version.GitCommit=$(GIT_COMMIT)$(GIT_DIRTY) -X $(GO_MODULE)/version.BinaryName=$(BINARY)"

PLATFORM = $(firstword $(subst /, ,$@))
ARCHITECTURE = $(lastword $(subst /, ,$@))
VERSION = $(shell sed -nr -e 's/.+Version = "(.*)"/\1/p' version/version.go)

TARGETS := darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64 windows/arm64

.PHONY: build
build:
	@go build -ldflags $(GO_LDFLAGS) -o $(BINARY)

.PHONY: test
test:
	@go test ./...

.PHONY: version
version:
ifndef v
	@echo "version must be provided, e.g. make version v=1.2.3"
	@exit 1
endif
	@sed -i '' -e 's/Version = .*/Version = "$(v)"/' ./version/version.go
	@git commit -qam "prepare release v$(v)"
	@git tag v$(v)

.PHONY: release
release: $(TARGETS)

.PHONY: $(TARGETS)
$(TARGETS):
	@mkdir -p ./dist/$@
	@GOOS=$(PLATFORM) GOARCH=$(ARCHITECTURE) go build -ldflags $(GO_LDFLAGS) -o ./dist/$@/$(BINARY)
	@zip -j -m -q -X ./dist/$(BINARY)_$(VERSION)_$(PLATFORM)_$(ARCHITECTURE).zip ./dist/$@/env-vault
	@rmdir ./dist/$@ ./dist/$(PLATFORM) 2>/dev/null || true

.PHONY: clean
clean:
	@rm -f ./$(BINARY)
	@rm -rf ./dist
