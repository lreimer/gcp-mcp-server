CURRENT_VERSION = $(shell git describe --tags --always --dirty)
VERSION ?= $(CURRENT_VERSION)

.PHONY: default
default: build

:PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags "-s -w -X main.version=$(VERSION)"

clean:
	go clean -i ./...
	rm -rf dist/