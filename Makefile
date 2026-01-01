# Variables
BINARY_PATH=bin
BINARY_NAME=ev
GO_FILES=$(shell find . -name '*.go')

# Targets
.PHONY: build test clean help

build: ## Build the binary for the local OS
	@echo " > Building binary..."
	go build -o "./$(BINARY_PATH)/$(BINARY_NAME)" main.go
	@echo " > Done! Binary created at ./$(BINARY_PATH)/$(BINARY_NAME)"

install: build
	@echo " > Installing to $$GOPATH/bin..."
	go install
	@echo " > Done! You can now run 'ev' from anywhere"

test: ## Run all unit and integration tests
	@echo " > Running tests..."
	go test ./... -v

clean: ## Remove build artifacts and temp files
	@echo " > Cleaning up..."
	rm -rf ./$(BINARY_PATH)/
	rm -rf coverage.out


# System specific compilation
build-linux: ## Build for Linux (amd64)
	GOOS=linux GOARCH=amd64 go build -o "./$(BINARY_PATH)/$(BINARY_NAME)-linux-amd64" main.go

build-linux-arm: ## Build for Linux (arm64)
	GOOS=linux GOARCH=arm64 go build -o "./$(BINARY_PATH)/$(BINARY_NAME)-linux-arm64" main.go

build-mac: ## Build for macOS (amd64)
	GOOS=darwin GOARCH=amd64 go build -o "./$(BINARY_PATH)/$(BINARY_NAME)-darwin-amd64" main.go

build-mac-arm: ## Build for macOS (arm64/Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build -o "./$(BINARY_PATH)/$(BINARY_NAME)-darwin-arm64" main.go


# Complete build
build-all: clean build build-linux build-linux-arm build-mac build-mac-arm
	@echo " > Done building for ALL targets"

# Help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
