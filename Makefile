SHELL = bash
PROJECT_ROOT := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
VERSION := 0.4.0
GIT_COMMIT := $(shell git rev-parse --short HEAD)

GO_PKGS := $(shell go list ./...)
GO_LDFLAGS := "-X github.com/romantomjak/env-vault/version.Version=$(VERSION) -X github.com/romantomjak/env-vault/version.GitCommit=$(GIT_COMMIT)"

TARGETS := darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64 windows/arm64

.PHONY: build
build:
	@mkdir -p $(PROJECT_ROOT)/bin
	@go build -ldflags $(GO_LDFLAGS) -o $(PROJECT_ROOT)/bin/env-vault

.PHONY: release
release: $(TARGETS)

.PHONY: $(TARGETS)
$(TARGETS): PLATFORM=$(firstword $(subst /, ,$@))
$(TARGETS): ARCHITECTURE=$(lastword $(subst /, ,$@))
$(TARGETS):
	@echo "==> Building $@..."
	@mkdir -p $(PROJECT_ROOT)/dist/$@
	@GOOS=$(PLATFORM) GOARCH=$(ARCHITECTURE) go build -ldflags $(GO_LDFLAGS) -o $(PROJECT_ROOT)/dist/$@/env-vault github.com/romantomjak/env-vault
	@echo "==> Packaging $@..."
	@zip -q -X -j $(PROJECT_ROOT)/dist/env-vault_$(VERSION)_$(PLATFORM)_$(ARCHITECTURE).zip $(PROJECT_ROOT)/dist/$@/env-vault
	@rm -rf $(PROJECT_ROOT)/dist/$(PLATFORM)

.PHONY: clean
clean:
	@rm -rf "$(PROJECT_ROOT)/bin/"
	@rm -rf "$(PROJECT_ROOT)/dist/"
