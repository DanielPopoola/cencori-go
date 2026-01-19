package cencori

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDoRequest_Sucess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("CENCORI_API_KEY") != "test-key" {
			t.Fatal("API key not sent")
		}

		json.NewEncoder(w).Encode(map[string]string{"id": "test-id"})
	}))
	defer server.Close()

	client := &Client{
		ApiKey:     "test-key",
		BaseURL:    server.URL,
		httpClient: &http.Client{},
	}

	resp, err := doRequest[any, map[string]string](client, context.Background(), "GET", "/test", nil)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil || (*resp)["id"] != "test-id" {
		t.Fatal("response not decoded correctly")
	}
}

func TestDoRequest_401MapsToTypedError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid API key",
			"code":  "INVALID_API_KEY",
		})
	}))
	defer server.Close()

	client := &Client{
		ApiKey:     "wrong-key",
		BaseURL:    server.URL,
		httpClient: &http.Client{},
	}

	_, err := doRequest[any, any](client, context.Background(), "GET", "/test", nil)

	if !errors.Is(err, ErrInvalidApiKey) {
		t.Fatalf("expected ErrInvalidApiKey, got %v", err)
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatal("error not an APIError")
	}
	if apiErr.StatusCode != 401 {
		t.Fatalf("wrong status code: %d", apiErr.StatusCode)
	}
}

func TestDoRequest_429MapsToErrRateLimited(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(429)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Too many requests",
			"code":  "RATE_LIMIT_EXCEEDED",
		})
	}))
	defer server.Close()

	client, _ := NewClient(WithApiKey("test-key"), WithBaseURL(server.URL))

	_, err := doRequest[any, any](client, context.Background(), "GET", "/test", nil)

	if !errors.Is(err, ErrRateLimited) {
		t.Fatalf("expected ErrRateLimited, got %v", err)
	}
}
