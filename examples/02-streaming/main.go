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

	// Start streaming chat
	stream, err := client.Chat.Stream(context.Background(), &cencori.ChatParams{
		Model: "gpt-4o",
		Messages: []cencori.Message{
			{Role: "user", Content: "Write a haiku about coding"},
		},
	})
	if err != nil {
		log.Fatalf("Failed to start stream: %v", err)
	}

	fmt.Print("Assistant: ")

	// Process chunks as they arrive
	for chunk := range stream {
		// Handle errors in stream
		if chunk.Err != nil {
			log.Fatalf("\nStream error: %v", chunk.Err)
		}

		// Print content progressively
		if len(chunk.Choices) > 0 {
			fmt.Print(chunk.Choices[0].Delta.Content)
		}
	}

	fmt.Println() // Final newline
}
