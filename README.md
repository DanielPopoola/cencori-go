# Cencori Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/cencori/cencori-go.svg)](https://pkg.go.dev/github.com/cencori/cencori-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/cencori/cencori-go)](https://goreportcard.com/report/github.com/cencori/cencori-go)
[![CI](https://github.com/cencori/cencori-go/workflows/Go/badge.svg)](https://github.com/cencori/cencori-go/actions)

Official Go SDK for the Cencori - unified access to OpenAI, Anthropic, and Google models with built-in security, logging, and cost tracking.

## Features

- **Multi-provider support** - OpenAI, Anthropic, Google with one API
- **Chat & Completions** - Conversational AI and text generation
- **Embeddings** - Vector embeddings for semantic search
- **Built-in security** - Automatic PII detection and content filtering
- **Complete analytics** - Token usage and cost tracking
- **Streaming support** - Real-time response streaming
- **Type-safe** - Full Go type safety with generics

## Installation

```bash
go get github.com/cencori/cencori-go
```

## Quick Start

### Chat Completion

```go
client, _ := cencori.NewClient(
    cencori.WithApiKey(os.Getenv("CENCORI_API_KEY")),
)

resp, err := client.Chat.Create(context.Background(), cencori.ChatParams{
    Model: "gpt-4o",
    Messages: []cencori.Message{
        {Role: "user", Content: "Hello!"},
    },
})
```

### Text Completion

```go
resp, err := client.Chat.Completions(context.Background(), cencori.CompletionParams{
    Prompt: "Write a haiku about coding",
    Model:  "gpt-4o",
})
```

### Embeddings

```go
resp, err := client.Chat.Embeddings(context.Background(), cencori.EmbeddingParams{
    Input: "Hello, world!",
    Model: "text-embedding-3-small",
})

// Batch embeddings
resp, err := client.Chat.Embeddings(context.Background(), cencori.EmbeddingParams{
    Input: []string{"First text", "Second text"},
    Model: "text-embedding-3-small",
})
```

## Examples

See [examples/](examples/) for complete, runnable examples:

### Chat & Completions
- **[01-basic-chat](examples/01-basic-chat/)** - Simple chat completion
- **[02-streaming](examples/02-streaming/)** - Real-time streaming
- **[03-error-handling](examples/03-error-handling/)** - Robust error handling
- **[04-advanced-params](examples/04-advanced-params/)** - All parameters
- **[10-completions](examples/10-completions/)** - Text completions

### Embeddings
- **[11-embeddings](examples/11-embeddings/)** - Vector embeddings & semantic search

### Management
- **[05-projects](examples/05-projects/)** - Project management
- **[06-key-rotation](examples/06-key-rotation/)** - API key rotation
- **[07-metrics](examples/07-metrics/)** - Analytics dashboard

### Advanced
- **[08-context-timeout](examples/08-context-timeout/)** - Context patterns
- **[09-multi-provider](examples/09-multi-provider/)** - Provider comparison

Run any example:
```bash
go run examples/01-basic-chat/main.go
```

## API Reference

### Chat API

```go
// Conversational chat
resp, err := client.Chat.Create(ctx, cencori.ChatParams{
    Model: "gpt-4o",
    Messages: []cencori.Message{
        {Role: "system", Content: "You are helpful"},
        {Role: "user", Content: "Hello"},
    },
})

// Simple text completion
resp, err := client.Chat.Completions(ctx, cencori.CompletionParams{
    Prompt: "Write a story about...",
    Model:  "gpt-4o",
})

// Streaming
stream, err := client.Chat.Stream(ctx, cencori.ChatParams{...})
for chunk := range stream {
    if chunk.Err != nil {
        log.Fatal(chunk.Err)
    }
    fmt.Print(chunk.Choices[0].Delta.Content)
}
```

### Embeddings API

```go
// Single text
resp, err := client.Chat.Embeddings(ctx, cencori.EmbeddingParams{
    Input: "Text to embed",
    Model: "text-embedding-3-small",
})

// Multiple texts (batch)
resp, err := client.Chat.Embeddings(ctx, cencori.EmbeddingParams{
    Input: []string{"Text 1", "Text 2", "Text 3"},
    Model: "text-embedding-3-small",
})

// Access embeddings
for i, data := range resp.Data {
    fmt.Printf("Text %d: %d dimensions\n", i, len(data.Embedding))
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
})
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
```

### Metrics API

```go
metrics, err := client.Metrics.Get(ctx, "24h")

fmt.Printf("Requests: %d\n", metrics.Requests.Total)
fmt.Printf("Cost: $%.4f\n", metrics.Cost.TotalUSD)
```

## Error Handling

All errors are typed with sentinel values:

```go
import "errors"

_, err := client.Chat.Chat(ctx, params)

switch {
case errors.Is(err, cencori.ErrInvalidApiKey):
    // Invalid API key
case errors.Is(err, cencori.ErrRateLimited):
    // Rate limit exceeded
case errors.Is(err, cencori.ErrInsufficientCredits):
    // Out of credits
case errors.Is(err, cencori.ErrInvalidModel):
    // Invalid model name
}
```

## Development

```bash
# Install dev tools
make install-tools

# Run tests
make test

# Run linter
make lint

# Run all checks
make verify
```

## Related Tools

**[@cencori/scan](https://www.npmjs.com/package/@cencori/scan)** — Security scanner for AI apps. Detects hardcoded secrets, PII, and vulnerabilities with AI-powered auto-fix.

```bash
npx @cencori/scan
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines.

## License

MIT License - see [LICENSE](LICENSE) for details.