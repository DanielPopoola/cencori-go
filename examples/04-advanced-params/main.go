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
		log.Fatalf("Failed to create client: %v", err)
	}

	// Advanced parameters demonstration
	temp := 0.7
	maxTokens := 150
	topP := 0.9
	userID := "user-123"

	resp, err := client.Chat.Chat(context.Background(), &cencori.ChatParams{
		Model: "gpt-4o",
		Messages: []cencori.Message{
			{
				Role:    "system",
				Content: "You are a concise technical writer.",
			},
			{
				Role:    "user",
				Content: "Explain microservices in 2 sentences.",
			},
		},
		Temperature: &temp,      // Control randomness (0.0-2.0)
		MaxTokens:   &maxTokens, // Limit response length
		TopP:        &topP,      // Nucleus sampling
		User:        &userID,    // Track per-user usage
	})
	if err != nil {
		log.Fatalf("Chat failed: %v", err)
	}

	// Display results with metadata
	fmt.Printf("Model: %s\n", resp.Model)
	fmt.Printf("Response: %s\n\n", resp.Choices[0].Message.Content)
	fmt.Printf("Tokens:\n")
	fmt.Printf("  Prompt: %d\n", resp.Usage.PromptTokens)
	fmt.Printf("  Completion: %d\n", resp.Usage.CompletionTokens)
	fmt.Printf("  Total: %d\n", resp.Usage.TotalTokens)
	fmt.Printf("Finish Reason: %s\n", resp.Choices[0].FinishReason)
}
