# Cencori Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/DanielPopoola/cencori-go.svg)](https://pkg.go.dev/github.com/DanielPopoola/cencori-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/DanielPopoola/cencori-go)](https://goreportcard.com/report/github.com/DanielPopoola/cencori-go)
[![CI](https://github.com/DanielPopoola/cencori-go/workflows/Go/badge.svg)](https://github.com/DanielPopoola/cencori-go/actions)
[![codecov](https://codecov.io/gh/DanielPopoola/cencori-go/branch/main/graph/badge.svg)](https://codecov.io/gh/DanielPopoola/cencori-go)

Unofficial Go SDK for the Cencori AI API - unified access to OpenAI, Anthropic, and Google models with built-in security, logging, and cost tracking.

## Features

- 🤖 **Multi-provider support** - OpenAI, Anthropic, Google with one API
- 🔒 **Built-in security** - Automatic PII detection and content filtering
- 📊 **Complete analytics** - Token usage and cost tracking
- ⚡ **Streaming support** - Real-time response streaming
- 🎯 **Type-safe** - Full Go type safety with generics
- 🔄 **Context-aware** - Proper timeout and cancellation support

## Installation

```bash
go get github.com/DanielPopoola/cencori-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    
    "github.com/DanielPopoola/cencori-go"
)

func main() {
    client, err := cencori.NewClient(
        cencori.WithAPIKey(os.Getenv("CENCORI_API_KEY")),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    resp, err := client.Chat.Chat(context.Background(), cencori.*ChatParams{
        Model: "gpt-4o",
        Messages: []cencori.Message{
            {Role: "user", Content: "Hello!"},
        },
    })
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(resp.Choices[0].Message.Content)
}
```

## Examples

See [examples/](examples/) for complete, runnable examples:

- **[01-basic-chat](examples/01-basic-chat/)** - Simple chat completion
- **[02-streaming](examples/02-streaming/)** - Real-time streaming
- **[03-error-handling](examples/03-error-handling/)** - Robust error handling
- **[04-advanced-params](examples/04-advanced-params/)** - All parameters
- **[05-projects](examples/05-projects/)** - Project management
- **[06-key-rotation](examples/06-key-rotation/)** - API key rotation
- **[07-metrics](examples/07-metrics/)** - Analytics dashboard
- **[08-context-timeout](examples/08-context-timeout/)** - Context patterns
- **[09-multi-provider](examples/09-multi-provider/)** - Provider comparison

Run any example:
```bash
go run examples/01-basic-chat/main.go
```

## API Reference

### Chat API

```go
// Non-streaming
resp, err := client.Chat.Chat(ctx, cencori.*ChatParams{
    Model: "gpt-4o",
    Messages: []cencori.Message{
        {Role: "user", Content: "Hello"},
    },
})

// Streaming
stream, err := client.Chat.Stream(ctx, cencori.*ChatParams{
    Model: "gpt-4o",
    Messages: messages,
})

for chunk := range stream {
    if chunk.Err != nil {
        log.Fatal(chunk.Err)
    }
    fmt.Print(chunk.Choices[0].Delta.Content)
}
```

### Projects API

```go
// List projects
projects, err := client.Projects.List(ctx, "org-slug")

// Create project
project, err := client.Projects.Create(ctx, "org-slug", cencori.CreateProjectParams{
    Name:        "My Project",
    Description: "AI assistant",
    Visibility:  "public",
})

// Get project
project, err := client.Projects.Get(ctx, "org-slug", "project-slug")

// Delete project
err := client.Projects.Delete(ctx, "org-slug", "project-slug")
```

### API Keys API

```go
// List keys
keys, err := client.APIKeys.List(ctx, "project-id", "production")

// Create key
key, err := client.APIKeys.Create(ctx, "project-id", cencori.CreateAPIKeyParams{
    Name:        "Production Key",
    Environment: "production",
})

// Get key stats
stats, err := client.APIKeys.GetStats(ctx, "project-id", "key-id")

// Revoke key
err := client.APIKeys.Revoke(ctx, "project-id", "key-id")
```

### Metrics API

```go
metrics, err := client.Metrics.Get(ctx, "24h") // or "7d", "30d", "mtd"

fmt.Printf("Total requests: %d\n", metrics.Requests.Total)
fmt.Printf("Total cost: $%.4f\n", metrics.Cost.TotalUSD)
fmt.Printf("P99 latency: %dms\n", metrics.Latency.P99MS)
```

## Error Handling

All errors are typed with sentinel values:

```go
import "errors"

_, err := client.Chat.Chat(ctx, params)

switch {
case errors.Is(err, cencori.ErrInvalidAPIKey):
    // Invalid or revoked API key
case errors.Is(err, cencori.ErrRateLimited):
    // Rate limit exceeded
case errors.Is(err, cencori.ErrInsufficientCredits):
    // Out of credits
case errors.Is(err, cencori.ErrProvider):
    // Provider error (retry)
case errors.Is(err, context.DeadlineExceeded):
    // Request timeout
default:
    // Unknown error
}
```

## Development

### Prerequisites

- Go 1.25+
- golangci-lint
- goimports

### Setup

```bash
# Install dev tools
make install-tools

# Run tests
make test

# Run linter
make lint

# Format code
make fmt

# Run all checks
make verify
```

### Testing

```bash
# Unit tests with race detection
make test

# Coverage report
make test-coverage

# Run examples
make examples
```

### CI/CD

GitHub Actions automatically:
- ✅ Runs linter (golangci-lint)
- ✅ Runs tests with race detection
- ✅ Checks code coverage
- ✅ Verifies build
- ✅ Validates go.mod

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Support

- 📚 [Documentation](https://docs.cencori.com)
- 💬 [Discord](https://discord.gg/cencori)
- 🐛 [Issues](https://github.com/DanielPopoola/cencori-go/issues)
- 📧 [Email](mailto:support@cencori.com)