package cencori

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ChatService provides methods for managing chat-related operations.
// It uses a Client to communicate with the chat API endpoints.
type ChatService struct {
	client *Client
}

// Create sends a chat request to the AI service and returns the response.
// It disables streaming and makes a synchronous request to the /api/ai/chat endpoint.
// The context can be used to cancel the request or set a timeout.
// It returns a ChatResponse on success or an error if the request fails.
func (s *ChatService) Create(ctx context.Context, params *ChatParams) (*ChatResponse, error) {
	params.Stream = false
	return doRequest[ChatParams, ChatResponse](s.client, ctx, "POST", "/api/ai/chat", params)
}

// Completions is a convenience method that wraps Create for simple text completions.
// It takes a single prompt string and converts it to a chat message internally.
// This is useful for quick, single-turn completions without managing conversation history.
func (s *ChatService) Completions(ctx context.Context, params CompletionParams) (*ChatResponse, error) {
	chatParams := &ChatParams{
		Model:       params.Model,
		Temperature: params.Temperature,
		MaxTokens:   params.MaxTokens,
		Messages: []Message{
			{Role: "user", Content: params.Prompt},
		},
	}
	return s.Create(ctx, chatParams)
}

// Embeddings generates vector embeddings for the given input text(s).
// Input can be a single string or a slice of strings.
// Returns an EmbeddingResponse containing the embeddings and token usage.
func (s *ChatService) Embeddings(ctx context.Context, params EmbeddingParams) (*EmbeddingResponse, error) {
	return doRequest[EmbeddingParams, EmbeddingResponse](s.client, ctx, "POST", "/api/v1/embeddings", &params)
}

// Stream sends a chat request with streaming enabled and returns a channel that receives
// chat response chunks as they arrive from the server. The stream continues until the server
// sends a "[DONE]" message or an error occurs. The context can be used to cancel the stream.
// If the context is cancelled, it simply closes.
// The returned channel will be closed when the stream ends or an error occurs.
func (s *ChatService) Stream(ctx context.Context, params *ChatParams) (<-chan StreamChunk, error) {
	params.Stream = true

	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		s.client.BaseURL+"/api/ai/chat",
		bytes.NewReader(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("CENCORI_API_KEY", s.client.APIKey)
	req.Header.Set("Accept", "text/event-stream")

	resp, err := s.client.httpClient.Do(req) //nolint:bodyclose // Body is closed by the streaming goroutine
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		err := handleError(resp)
		resp.Body.Close() //nolint:errcheck // Closing the response body; error can be ignored here.

		return nil, err
	}

	chunks := make(chan StreamChunk)

	go func() {
		defer close(chunks)
		defer resp.Body.Close() //nolint:errcheck // Closing the response body; error can be ignored here.

		go func() {
			<-ctx.Done()
			resp.Body.Close() //nolint:errcheck // Closing the response body; error can be ignored here.
		}()

		reader := bufio.NewReader(resp.Body)

		for {
			line, err := reader.ReadString('\n')

			if err != nil {
				if ctx.Err() != nil {
					return
				}

				if err == io.EOF {
					return
				}

				chunks <- StreamChunk{Err: fmt.Errorf("stream read: %w", err)}
				return
			}

			line = strings.TrimSpace(line)

			// Ignore comments / empty lines / non-data frames
			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")

			if data == "[DONE]" {
				return
			}

			if strings.Contains(data, "\"error\":") {
				var apiErr APIError
				if err := json.Unmarshal([]byte(data), &apiErr); err == nil {
					apiErr.fillSentinel()
					chunks <- StreamChunk{Err: &apiErr}
					return
				}
			}

			var chunk StreamChunk
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				chunks <- StreamChunk{Err: fmt.Errorf("unmarshal chunk: %w", err)}
				return
			}

			chunks <- chunk
		}
	}()

	return chunks, nil
}
