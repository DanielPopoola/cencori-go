package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/DanielPopoola/cencori-go"
)

func main() {
	// Initialize client with API key from environment
	client, err := cencori.NewClient(
		cencori.WithAPIKey(os.Getenv("CENCORI_API_KEY")),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Simple chat completion
	resp, err := client.Chat.Chat(context.Background(), &cencori.ChatParams{
		Model: "gpt-4o",
		Messages: []cencori.Message{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: "What is the capital of France?"},
		},
	})
	if err != nil {
		log.Fatalf("Chat failed: %v", err)
	}

	// Extract and print response
	if len(resp.Choices) > 0 {
		fmt.Printf("Assistant: %s\n", resp.Choices[0].Message.Content)
		fmt.Printf("\nTokens used: %d\n", resp.Usage.TotalTokens)
	}
}
