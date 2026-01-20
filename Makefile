.PHONY: test lint fmt build examples clean install-tools

# Run all tests with race detection
test:
	go test -v -race -coverprofile=coverage.out ./...

# Run tests with coverage report
test-coverage: test
	go tool cover -html=coverage.out

# Run linter
lint:
	golangci-lint run --timeout 5m

# Format code
fmt:
	gofmt -s -w .
	goimports -w .

# Build the SDK
build:
	go build -v ./...

# Run all examples
examples:
	@for dir in examples/*/; do \
		echo "▶ Running $$dir..."; \
		go run $$dir/main.go || exit 1; \
		echo ""; \
	done

# Run specific example (e.g., make example-01)
example-%:
	@go run examples/$*-*/main.go

# Clean build artifacts
clean:
	rm -f coverage.out
	go clean -cache -testcache

# Install development tools
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

# Verify everything (used in CI)
verify: fmt lint test build

# Help
help:
	@echo "Available targets:"
	@echo "  test           - Run tests with race detection"
	@echo "  test-coverage  - Run tests and show coverage"
	@echo "  lint           - Run golangci-lint"
	@echo "  fmt            - Format code"
	@echo "  build          - Build the SDK"
	@echo "  examples       - Run all examples"
	@echo "  example-01     - Run specific example (01-09)"
	@echo "  clean          - Clean build artifacts"
	@echo "  install-tools  - Install dev dependencies"
	@echo "  verify         - Run all checks (CI)"