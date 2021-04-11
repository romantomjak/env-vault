SHELL = bash
PROJECT_ROOT := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
VERSION := 0.1.0
GIT_COMMIT := $(shell git rev-parse --short HEAD)

GO_PKGS := $(shell go list ./...)
GO_LDFLAGS := "-X github.com/romantomjak/env-vault/command.Version=$(VERSION) -X github.com/romantomjak/env-vault/command.GitCommit=$(GIT_COMMIT)"

PLATFORMS := darwin linux windows
os = $(word 1, $@)

.PHONY: build
build:
	@mkdir -p $(PROJECT_ROOT)/bin
	@go build -ldflags $(GO_LDFLAGS) -o $(PROJECT_ROOT)/bin/env-vault

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	@mkdir -p $(PROJECT_ROOT)/dist
	@GOOS=$(os) GOARCH=amd64 go build -ldflags $(GO_LDFLAGS) -o $(PROJECT_ROOT)/dist/$(os)/env-vault github.com/romantomjak/env-vault
	@zip -q -X -j $(PROJECT_ROOT)/dist/env-vault_$(VERSION)_$(os)_amd64.zip $(PROJECT_ROOT)/dist/$(os)/env-vault
	@rm -rf $(PROJECT_ROOT)/dist/$(os)

.PHONY: release
release: windows linux darwin

.PHONY: clean
clean:
	@rm -rf "$(PROJECT_ROOT)/bin/"
	@rm -rf "$(PROJECT_ROOT)/dist/"
