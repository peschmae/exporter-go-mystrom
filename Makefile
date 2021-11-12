.PHONY: clean run

linux:
	GOOS=linux GOARCH=amd64 go build -o output/mystrom-exporter_linux-amd64
mac:
	GOOS=darwin GOARCH=amd64 go build -o output/mystrom-exporter_mac-amd64
arm64:
	GOOS=linux GOARCH=arm64 go build -o output/mystrom-exporter_linux-arm64
arm:
	GOOS=linux GOARCH=arm go build -o output/mystrom-exporter_linux-arm


all: linux mac arm64 arm