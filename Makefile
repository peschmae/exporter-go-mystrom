##
#
.DEFAULT_GOAL := help
.PHONY: generate go-tools

version    := $(shell git describe --tags --always)
revision   := $(shell git rev-parse HEAD)
branch     := $(shell git rev-parse --abbrev-ref HEAD)
builduser  := $(shell whoami)
builddate  := $(shell date '+%FT%T_%Z')

versionPkgPrefix := mystrom-exporter/pkg/version

BINDIR    := $(CURDIR)/output
GO        ?= go
GOPATH    ?= $(shell $(GO) env GOPATH)
LDFLAGS   := -w -s \
	-X $(versionPkgPrefix).Version=${version} \
	-X $(versionPkgPrefix).Revision=${revision} \
	-X $(versionPkgPrefix).Branch=${branch} \
	-X $(versionPkgPrefix).BuildUser=${builduser} \
	-X $(versionPkgPrefix).BuildDate=${builddate}
GOFLAGS   := -v
GOX_FLAGS          := -mod=vendor
GO_BUILD_FLAGS     := -v
export GO111MODULE := on

build: go-tools generate ## builds the all platform binaries of the exporter
	$(GOPATH)/bin/gox \
			-os="darwin linux" \
			-arch="amd64 arm arm64" \
			-osarch="!darwin/arm" \
			-output "${BINDIR}/{{.Dir}}-{{.OS}}-{{.Arch}}" \
			-gcflags "$(GO_BUILD_FLAGS)" \
			-ldflags '$(LDFLAGS)' \
			-tags '$(TAGS)' \
			./...

run:
	${BINDIR}/mystrom-exporter-$(shell $(GO) env GOOS)-$(shell $(GO) env GOARCH)


generate: go-tools
	$(GO) generate ./...

go-tools: $(GOPATH)/bin/stringer $(GOPATH)/bin/gox

# -- see more info on https://pkg.go.dev/golang.org/x/tools/cmd/stringer
$(GOPATH)/bin/stringer:
	$(GO) install golang.org/x/tools/cmd/stringer@latest

$(GOPATH)/bin/gox:
	$(GO) install github.com/mitchellh/gox@latest
# --
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST)  | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
