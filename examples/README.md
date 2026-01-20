# Cencori Go SDK Examples

Complete, production-ready examples for the Cencori Go SDK.

## Prerequisites

```bash
export CENCORI_API_KEY="your-key-here"
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
    client, _ := cencori.NewClient(
        cencori.WithAPIKey(os.Getenv("CENCORI_API_KEY")),
    )
    
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

## Examples by Category

### Chat API
- **[basic.go](chat/basic.go)** - Simple chat completion
- **[streaming.go](chat/streaming.go)** - Real-time streaming responses
- **[error_handling.go](chat/error_handling.go)** - Comprehensive error handling
- **[advanced.go](chat/advanced.go)** - All parameters (temperature, max_tokens, etc.)

### Projects Management
- **[projects.go](projects/projects.go)** - CRUD operations for projects
- **[lifecycle.go](projects/lifecycle.go)** - Complete project lifecycle

### API Keys
- **[rotation.go](api_keys/rotation.go)** - Production key rotation pattern
- **[monitoring.go](api_keys/monitoring.go)** - Track key usage and stats

### Metrics & Analytics
- **[dashboard.go](metrics/dashboard.go)** - Custom analytics dashboard
- **[alerts.go](metrics/alerts.go)** - Cost and performance alerting

### Advanced Patterns
- **[context.go](advanced/context.go)** - Timeouts and cancellation
- **[multi_provider.go](advanced/multi_provider.go)** - Compare providers/models
- **[retry.go](advanced/retry.go)** - Robust retry logic
- **[concurrent.go](advanced/concurrent.go)** - Parallel requests

## Running Examples

```bash
# Run specific example
go run examples/chat/basic.go

# Run with race detector
go run -race examples/chat/streaming.go

# Run all examples
for f in examples/**/*.go; do go run "$f"; done
```

## Common Patterns

### Error Handling

```go
import "errors"

_, err := client.Chat.Chat(ctx, params)
if err != nil {
    switch {
    case errors.Is(err, cencori.ErrInvalidAPIKey):
        // Handle auth error
    case errors.Is(err, cencori.ErrRateLimited):
        // Implement backoff
    case errors.Is(err, context.DeadlineExceeded):
        // Handle timeout
    default:
        // Unknown error
    }
}
```

### Streaming

```go
stream, err := client.Chat.Stream(ctx, params)
if err != nil {
    log.Fatal(err)
}

for chunk := range stream {
    if chunk.Err != nil {
        log.Fatal(chunk.Err)
    }
    fmt.Print(chunk.Choices[0].Delta.Content)
}
```

### Context Timeout

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

resp, err := client.Chat.Chat(ctx, params)
// Request auto-cancelled after 10s
```

## Testing

Each example includes inline comments explaining:
- What it demonstrates
- Expected output
- Common gotchas
- Production considerations

## Need Help?

- 📚 [Full Documentation](https://cencori.com/docs)
