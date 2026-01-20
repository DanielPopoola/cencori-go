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

	// Example 1: Request with timeout
	fmt.Println("Request with 5-second timeout...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Chat.Chat(ctx, &cencori.ChatParams{
		Model: "gpt-4o",
		Messages: []cencori.Message{
			{Role: "user", Content: "What is AI?"},
		},
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Println("Request timed out after 5 seconds")
		} else {
			log.Fatalf("Request failed: %v", err)
		}
	} else {
		fmt.Printf("Got response: %s\n", resp.Choices[0].Message.Content[:50]+"...")
	}

	// Example 2: Cancellable streaming
	fmt.Println("\n Cancellable streaming...")
	streamCtx, streamCancel := context.WithCancel(context.Background())

	stream, err := client.Chat.Stream(streamCtx, &cencori.ChatParams{
		Model: "gpt-4o",
		Messages: []cencori.Message{
			{Role: "user", Content: "Count to 100"},
		},
	})
	if err != nil {
		log.Fatalf("Failed to start stream: %v", err)
	}

	// Cancel after 2 seconds
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("\n Cancelling stream...")
		streamCancel()
	}()

	fmt.Print("Assistant: ")
	for chunk := range stream {
		if chunk.Err != nil {
			if errors.Is(chunk.Err, context.Canceled) {
				fmt.Println("\n Stream cancelled successfully")
			} else {
				fmt.Printf("\n Stream error: %v\n", chunk.Err)
			}
			break
		}

		if len(chunk.Choices) > 0 {
			fmt.Print(chunk.Choices[0].Delta.Content)
		}
	}

	// Example 3: Parent context cancellation
	fmt.Println("\n Parent context cancellation...")
	parentCtx, parentCancel := context.WithCancel(context.Background())

	// Cancel parent immediately
	parentCancel()

	_, err = client.Chat.Chat(parentCtx, &cencori.ChatParams{
		Model: "gpt-4o",
		Messages: []cencori.Message{
			{Role: "user", Content: "Hello"},
		},
	})

	if errors.Is(err, context.Canceled) {
		fmt.Println("Request cancelled via parent context")
	}

	// Example 4: Custom deadline
	fmt.Println("\n Custom deadline (expires in 10ms)...")
	deadline := time.Now().Add(10 * time.Millisecond)
	deadlineCtx, deadlineCancel := context.WithDeadline(context.Background(), deadline)
	defer deadlineCancel()

	_, err = client.Chat.Chat(deadlineCtx, &cencori.ChatParams{
		Model: "gpt-4o",
		Messages: []cencori.Message{
			{Role: "user", Content: "Quick question"},
		},
	})

	if errors.Is(err, context.DeadlineExceeded) {
		fmt.Println("Request stopped at deadline")
	}
}
