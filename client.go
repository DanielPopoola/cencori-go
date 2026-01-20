package cencori

import (
	"errors"
	"net/http"
	"time"
)

type ClientOptions struct {
	APIKey  string
	BaseURL string
	Timeout time.Duration
}

func WithAPIKey(apiKey string) Option {
	return func(c *ClientOptions) { c.APIKey = apiKey }
}

func WithBaseURL(baseURL string) Option {
	return func(c *ClientOptions) { c.BaseURL = baseURL }
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *ClientOptions) { c.Timeout = timeout }
}

type Client struct {
	APIKey     string
	BaseURL    string
	httpClient *http.Client

	Chat     *ChatService
	Projects *ProjectsService
	APIKeys  *APIKeysService
	Metrics  *MetricsService
}

type Option func(*ClientOptions)

func NewClient(opts ...Option) (*Client, error) {
	config := &ClientOptions{
		BaseURL: "https://cencori.com",
		Timeout: 30 * time.Second,
	}
	for _, opt := range opts {
		opt(config)
	}

	if config.APIKey == "" {
		return nil, errors.New("you need a valid API Key to use this client")
	}

	c := &Client{
		APIKey:  config.APIKey,
		BaseURL: config.BaseURL,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}

	c.Chat = &ChatService{client: c}
	c.Projects = &ProjectsService{client: c}
	c.APIKeys = &APIKeysService{client: c}
	c.Metrics = &MetricsService{client: c}

	return c, nil
}
