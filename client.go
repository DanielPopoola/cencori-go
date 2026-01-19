package cencori

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ClientOptions struct {
	ApiKey  string
	BaseURL string
	Timeout int
}

func WithApiKey(apiKey string) Option {
	return func(c *ClientOptions) { c.ApiKey = apiKey }
}

func WithBaseURL(baseURL string) Option {
	return func(c *ClientOptions) { c.BaseURL = baseURL }
}

type Client struct {
	ApiKey     string
	BaseURL    string
	httpClient *http.Client
}

type Option func(*ClientOptions)

func NewClient(opts ...Option) (*Client, error) {
	clientOpts := &ClientOptions{
		BaseURL: "https://cencori.com",
		Timeout: 30,
	}
	for _, opt := range opts {
		opt(clientOpts)
	}
	httpClient := &http.Client{
		Timeout: time.Duration(clientOpts.Timeout) * time.Second,
	}

	if clientOpts.ApiKey == "" {
		return nil, errors.New("You need a valid API Key to use this client")
	}

	return &Client{
		ApiKey:     clientOpts.ApiKey,
		BaseURL:    clientOpts.BaseURL,
		httpClient: httpClient,
	}, nil
}

func doRequest[Req any, Resp any](
	c *Client,
	ctx context.Context,
	method, path string,
	body *Req,
) (*Resp, error) {
	url := c.BaseURL + path

	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("CENCORI_API_KEY", c.ApiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, &APIError{
				StatusCode: resp.StatusCode,
				Code:       "READ_ERROR",
				Message:    fmt.Sprintf("failed to read response body: %v", err),
			}
		}
		var apiErr APIError
		if err := json.Unmarshal(body, &apiErr); err != nil {
			return nil, &APIError{
				StatusCode: resp.StatusCode,
				Code:       "UNKNOWN",
				Message:    string(body),
			}
		}
		apiErr.StatusCode = resp.StatusCode
		return nil, &apiErr
	}

	var result Resp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &result, nil
}
