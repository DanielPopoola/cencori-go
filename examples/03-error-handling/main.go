package main

import (
	"context"
	"errors"
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

	// Attempt chat with comprehensive error handling
	resp, err := client.Chat.Chat(context.Background(), &cencori.ChatParams{
		Model: "gpt-4o",
		Messages: []cencori.Message{
			{Role: "user", Content: "Hello!"},
		},
	})

	if err != nil {
		handleChatError(err)
		return
	}

	fmt.Printf("Success: %s\n", resp.Choices[0].Message.Content)
}

func handleChatError(err error) {
	// Check for specific error types using errors.Is
	switch {
	case errors.Is(err, cencori.ErrInvalidAPIKey):
		log.Fatal("Invalid API key. Check your CENCORI_API_KEY environment variable.")

	case errors.Is(err, cencori.ErrRateLimited):
		log.Fatal("Rate limit exceeded. Wait before retrying.")

	case errors.Is(err, cencori.ErrInsufficientCredits):
		log.Fatal("Insufficient credits. Top up your account.")

	case errors.Is(err, cencori.ErrTierRestricted):
		log.Fatal("This feature requires a paid plan.")

	case errors.Is(err, cencori.ErrProvider):
		// Provider errors are often transient - retry
		log.Println("Provider error. Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
		// In real code, implement exponential backoff

	case errors.Is(err, cencori.ErrContentFiltered):
		log.Fatal("Content filtered by safety policies.")

	default:
		// Unknown error - get details
		var apiErr *cencori.APIError
		if errors.As(err, &apiErr) {
			log.Fatalf("API Error [%d]: %s (code: %s)",
				apiErr.StatusCode, apiErr.Message, apiErr.Code)
		}
		log.Fatalf("Unknown error: %v", err)
	}
}
