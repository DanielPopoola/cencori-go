package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/DanielPopoola/cencori-go"
)

func main() {
	client, err := cencori.NewClient(
		cencori.WithAPIKey(os.Getenv("CENCORI_API_KEY")),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Test multiple providers with the same prompt
	prompt := "Explain quantum computing in one sentence."

	models := []string{
		"gpt-4o",          // OpenAI
		"gpt-4o-mini",     // OpenAI (cheaper)
		"claude-3-sonnet", // Anthropic
		"gemini-pro",      // Google
	}

	fmt.Printf("Testing prompt across %d models:\n", len(models))
	fmt.Printf("Prompt: %s\n\n", prompt)

	type result struct {
		model    string
		response string
		tokens   int
		latency  time.Duration
		err      error
	}

	results := make([]result, 0, len(models))

	for _, model := range models {
		start := time.Now()

		resp, err := client.Chat.Chat(context.Background(), &cencori.ChatParams{
			Model: model,
			Messages: []cencori.Message{
				{Role: "user", Content: prompt},
			},
		})

		latency := time.Since(start)

		if err != nil {
			results = append(results, result{
				model: model,
				err:   err,
			})
			continue
		}

		results = append(results, result{
			model:    model,
			response: resp.Choices[0].Message.Content,
			tokens:   resp.Usage.TotalTokens,
			latency:  latency,
		})
	}

	// Display results
	fmt.Println("=" + strings.Repeat("=", 70))
	for i, r := range results {
		fmt.Printf("\n%d. %s\n", i+1, r.model)
		fmt.Println(strings.Repeat("-", 72))

		if r.err != nil {
			fmt.Printf("Error: %v\n", r.err)
			continue
		}

		fmt.Printf("Response: %s\n", r.response)
		fmt.Printf("Tokens:   %d\n", r.tokens)
		fmt.Printf("Latency:  %v\n", r.latency)
	}

	// Find fastest and most efficient
	var fastest, cheapest result
	minLatency := time.Hour
	minTokens := 999999

	for _, r := range results {
		if r.err != nil {
			continue
		}
		if r.latency < minLatency {
			minLatency = r.latency
			fastest = r
		}
		if r.tokens < minTokens {
			minTokens = r.tokens
			cheapest = r
		}
	}

	if fastest.model != "" {
		fmt.Printf("\n Fastest: %s (%v)\n", fastest.model, fastest.latency)
	}
	if cheapest.model != "" {
		fmt.Printf("Most efficient: %s (%d tokens)\n", cheapest.model, cheapest.tokens)
	}
}

var strings = struct {
	Repeat func(string, int) string
}{
	Repeat: func(s string, count int) string {
		result := ""
		for i := 0; i < count; i++ {
			result += s
		}
		return result
	},
}
