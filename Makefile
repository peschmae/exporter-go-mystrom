##
#
.DEFAULT_GOAL := help

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

linux: ## builds the linux version of the exporter
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -ldflags '$(LDFLAGS)'
mac: ## builds the macos version of the exporter
	GOOS=darwin GOARCH=amd64 go build $(GOFLAGS) -ldflags '$(LDFLAGS)'
arm64:
	GOOS=linux GOARCH=arm64 go build $(GOFLAGS) -ldflags '$(LDFLAGS)'
arm:
	GOOS=linux GOARCH=arm go build $(GOFLAGS) -ldflags '$(LDFLAGS)'

# --
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST)  | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
