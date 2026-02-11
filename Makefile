BINARY_NAME := open-entire
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

.PHONY: build test lint install clean fmt vet

build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/entire

install:
	go install $(LDFLAGS) ./cmd/entire

test:
	go test ./... -v

test-race:
	go test -race ./...

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -rf bin/
	go clean

dev: fmt vet test build
