package cencori

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestChat_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"id": "chat-123",
			"choices": []map[string]any{
				{
					"index": 0,
					"message": map[string]string{
						"role":    "assistant",
						"content": "Hello!",
					},
					"finish_reason": "stop",
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	resp, err := client.Chat.Create(context.Background(), &ChatParams{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{Role: "user", Content: "Hi"},
		},
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.ID != "chat-123" {
		t.Errorf("expected ID chat-123, got %s", resp.ID)
	}
	if resp.Choices[0].Message.Content != "Hello!" {
		t.Errorf("expected content Hello!, got %s", resp.Choices[0].Message.Content)
	}
}

func TestChat_503Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(503)
		w.Write([]byte("Service Unavailable"))
	}))
	defer server.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	_, err := client.Chat.Create(context.Background(), &ChatParams{})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != 503 {
		t.Errorf("expected status 503, got %d", apiErr.StatusCode)
	}
}

func TestChat_Stream_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, "data: {\"choices\": [{\"delta\": {\"content\": \"Hello\"}}]}\n\n")
		fmt.Fprint(w, "data: {\"choices\": [{\"delta\": {\"content\": \"!\"}}]}\n\n")
		fmt.Fprint(w, "data: [DONE]\n\n")
	}))
	defer server.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	stream, err := client.Chat.Stream(context.Background(), &ChatParams{})
	if err != nil {
		t.Fatalf("failed to start stream: %v", err)
	}

	var content string
	for chunk := range stream {
		if chunk.Err != nil {
			t.Fatalf("stream error: %v", chunk.Err)
		}
		if len(chunk.Choices) > 0 {
			content += chunk.Choices[0].Delta.Content
		}
	}

	if content != "Hello!" {
		t.Errorf("expected Hello!, got %s", content)
	}
}

func TestChat_Stream_MidErrorFrame(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, "data: {\"choices\": [{\"delta\": {\"content\": \"Hello\"}}]}\n\n")
		fmt.Fprint(w, "data: {\"error\": \"something went wrong\", \"code\": \"PROVIDER_ERROR\"}\n\n")
	}))
	defer server.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	stream, err := client.Chat.Stream(context.Background(), &ChatParams{})
	if err != nil {
		t.Fatalf("failed to start stream: %v", err)
	}

	var lastErr error
	for chunk := range stream {
		if chunk.Err != nil {
			lastErr = chunk.Err
		}
	}

	if lastErr == nil {
		t.Fatal("expected error in stream, got nil")
	}
	if !errors.Is(lastErr, ErrProvider) {
		t.Errorf("expected ErrProvider, got %v", lastErr)
	}
}

func TestChat_Stream_ContextCancel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")

		fmt.Fprint(w, "data: {\"choices\": [{\"delta\": {\"content\": \"Hello\"}}]}\n\n")
		w.(http.Flusher).Flush()

		<-r.Context().Done()
	}))
	defer server.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	ctx, cancel := context.WithCancel(context.Background())
	stream, err := client.Chat.Stream(ctx, &ChatParams{})
	if err != nil {
		t.Fatalf("failed to start stream: %v", err)
	}

	chunk, ok := <-stream
	if !ok {
		t.Fatal("stream closed early")
	}
	if chunk.Err != nil {
		t.Fatalf("unexpected error: %v", chunk.Err)
	}

	cancel()

	_, ok = <-stream
	if ok {
		t.Fatal("expected stream to close on cancel")
	}
}

func TestCompletions_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify it's calling chat endpoint with correct format
		var req ChatParams
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		// Verify prompt was converted to message
		if len(req.Messages) != 1 {
			t.Fatalf("expected 1 message, got %d", len(req.Messages))
		}
		if req.Messages[0].Role != "user" {
			t.Errorf("expected role 'user', got %s", req.Messages[0].Role)
		}
		if req.Messages[0].Content != "Test prompt" {
			t.Errorf("expected content 'Test prompt', got %s", req.Messages[0].Content)
		}

		// Return chat response
		resp := ChatResponse{
			ID: "cmpl-123",
			Choices: []struct {
				Index        int     `json:"index"`
				Message      Message `json:"message"`
				FinishReason string  `json:"finish_reason"`
			}{
				{
					Index:        0,
					Message:      Message{Role: "assistant", Content: "Completion response"},
					FinishReason: "stop",
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	temp := 0.7
	resp, err := client.Chat.Completions(context.Background(), CompletionParams{
		Prompt:      "Test prompt",
		Model:       "gpt-4o",
		Temperature: &temp,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Choices[0].Message.Content != "Completion response" {
		t.Errorf("unexpected response content: %s", resp.Choices[0].Message.Content)
	}
}

func TestEmbeddings_SingleString(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/embeddings" {
			t.Errorf("expected path /api/v1/embeddings, got %s", r.URL.Path)
		}

		var req EmbeddingParams
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		// Verify input is correct
		if req.Input != "Hello world" {
			t.Errorf("expected input 'Hello world', got %v", req.Input)
		}

		// Return embedding response
		resp := EmbeddingResponse{
			Model:  "text-embedding-3-small",
			Object: "list",
			Data: []EmbeddingData{
				{
					Embedding: []float64{0.1, 0.2, 0.3},
					Index:     0,
				},
			},
			Usage: EmbeddingUsage{
				TotalTokens: 2,
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	resp, err := client.Chat.Embeddings(context.Background(), EmbeddingParams{
		Input: "Hello world",
		Model: "text-embedding-3-small",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 embedding, got %d", len(resp.Data))
	}
	if len(resp.Data[0].Embedding) != 3 {
		t.Errorf("expected embedding dimension 3, got %d", len(resp.Data[0].Embedding))
	}
	if resp.Usage.TotalTokens != 2 {
		t.Errorf("expected 2 tokens, got %d", resp.Usage.TotalTokens)
	}
}

func TestEmbeddings_MultipleStrings(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req EmbeddingParams
		json.NewDecoder(r.Body).Decode(&req)

		// Verify input is array
		inputs, ok := req.Input.([]interface{})
		if !ok {
			t.Fatalf("expected input to be array, got %T", req.Input)
		}
		if len(inputs) != 2 {
			t.Errorf("expected 2 inputs, got %d", len(inputs))
		}

		// Return multiple embeddings
		resp := EmbeddingResponse{
			Model: "text-embedding-3-small",
			Data: []EmbeddingData{
				{Embedding: []float64{0.1, 0.2}, Index: 0},
				{Embedding: []float64{0.3, 0.4}, Index: 1},
			},
			Usage: EmbeddingUsage{TotalTokens: 4},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	resp, err := client.Chat.Embeddings(context.Background(), EmbeddingParams{
		Input: []string{"First text", "Second text"},
		Model: "text-embedding-3-small",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(resp.Data) != 2 {
		t.Fatalf("expected 2 embeddings, got %d", len(resp.Data))
	}
}

func TestEmbeddings_ErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid model",
			"code":  "INVALID_MODEL",
		})
	}))
	defer server.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	_, err := client.Chat.Embeddings(context.Background(), EmbeddingParams{
		Input: "test",
		Model: "invalid-model",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if !errors.Is(err, ErrInvalidModel) {
		t.Error("expected ErrInvalidModel sentinel")
	}
}

func TestPathBuilding_OrgProject(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/organizations/my-org/projects/my-project"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}
		json.NewEncoder(w).Encode(Project{ID: "p1", Name: "My Project"})
	}))
	defer server.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	proj, err := client.Projects.Get(context.Background(), "my-org", "my-project")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if proj.ID != "p1" {
		t.Errorf("expected project ID p1, got %s", proj.ID)
	}
}

func TestErrorDecoding_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Write([]byte("not a json"))
	}))
	defer server.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	_, err := client.Chat.Create(context.Background(), &ChatParams{})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.Message != "not a json" {
		t.Errorf("expected message 'not a json', got %s", apiErr.Message)
	}
}

func TestConcurrency_Race(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(ChatResponse{ID: "race-test"})
	}))
	defer server.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	var wg sync.WaitGroup
	numRequests := 20
	for range numRequests {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := client.Chat.Create(context.Background(), &ChatParams{})
			if err != nil {
				t.Errorf("concurrent request failed: %v", err)
			}
		}()
	}
	wg.Wait()
}

func TestHeaderInjection(t *testing.T) {
	apiKey := "test-key-123"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("CENCORI_API_KEY") != apiKey {
			t.Errorf("expected API key %s, got %s", apiKey, r.Header.Get("CENCORI_API_KEY"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}
		json.NewEncoder(w).Encode(ChatResponse{})
	}))
	defer server.Close()

	client, _ := NewClient(WithAPIKey(apiKey), WithBaseURL(server.URL))

	_, err := client.Chat.Create(context.Background(), &ChatParams{})
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
}
