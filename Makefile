##
#
.DEFAULT_GOAL := help
.PHONY: generate

version    := $(shell git describe --tags --always)
revision   := $(shell git rev-parse HEAD)
branch     := $(shell git rev-parse --abbrev-ref HEAD)
builduser  := $(shell whoami)
builddate  := $(shell date '+%FT%T_%Z')

versionPkgPrefix := mystrom-exporter/pkg/version

LDFLAGS   := -w -s \
	-X $(versionPkgPrefix).Version=${version} \
	-X $(versionPkgPrefix).Revision=${revision} \
	-X $(versionPkgPrefix).Branch=${branch} \
	-X $(versionPkgPrefix).BuildUser=${builduser} \
	-X $(versionPkgPrefix).BuildDate=${builddate}
GOFLAGS   := -v

linux: generate ## builds the linux version of the exporter
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -ldflags '$(LDFLAGS)'
mac: generate ## builds the macos version of the exporter
	GOOS=darwin GOARCH=amd64 go build $(GOFLAGS) -ldflags '$(LDFLAGS)'
mac-arm: generate ## builds the macos (m1) version of the exporter
	GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) -ldflags '$(LDFLAGS)'
arm64: generate
	GOOS=linux GOARCH=arm64 go build $(GOFLAGS) -ldflags '$(LDFLAGS)'
arm: generate
	GOOS=linux GOARCH=arm go build $(GOFLAGS) -ldflags '$(LDFLAGS)'

# -- see more info on https://pkg.go.dev/golang.org/x/tools/cmd/stringer
generate: $(GOPATH)/bin/stringer
	go generate ./...

$(GOPATH)/bin/stringer:
	go install golang.org/x/tools/cmd/stringer@latest

# --
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST)  | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
