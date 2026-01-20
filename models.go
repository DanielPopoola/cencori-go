package cencori

import "time"

// Shared Components.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Stats struct {
	TotalRequests int        `json:"total_requests"`
	TotalCostUSD  float64    `json:"total_cost_usd"`
	LastUsedAt    *time.Time `json:"last_used_at,omitempty"`
}

// Chat Models.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatParams struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature *float64  `json:"temperature,omitempty"`
	MaxTokens   *int      `json:"max_tokens,omitempty"`
	TopP        *float64  `json:"top_p,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
	User        *string   `json:"user,omitempty"`
}

type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int     `json:"index"`
		Message      Message `json:"message"`
		FinishReason string  `json:"finish_reason"`
	} `json:"choices"`
	Usage Usage `json:"usage"`
}

// Completions Models.
type CompletionParams struct {
	Prompt      string   `json:"prompt"`
	Model       string   `json:"model"`
	Temperature *float64 `json:"temperature,omitempty"`
	MaxTokens   *int     `json:"max_tokens,omitempty"`
}

// Embedding Models.
type EmbeddingParams struct {
	Input any    `json:"input"` // string or []string
	Model string `json:"model"`
}

type EmbeddingUsage struct {
	TotalTokens int `json:"total_tokens"`
}

type EmbeddingData struct {
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

type EmbeddingResponse struct {
	Model  string          `json:"model"`
	Data   []EmbeddingData `json:"data"`
	Usage  EmbeddingUsage  `json:"usage"`
	Object string          `json:"object"`
}

// Stream Response.
type StreamChunk struct {
	ID      string         `json:"id,omitempty"`
	Object  string         `json:"object,omitempty"`
	Created int64          `json:"created,omitempty"`
	Model   string         `json:"model,omitempty"`
	Choices []StreamChoice `json:"choices,omitempty"`
	Err     error          `json:"-"`
}

type StreamChoice struct {
	Index        int         `json:"index"`
	Delta        StreamDelta `json:"delta"`
	FinishReason *string     `json:"finish_reason,omitempty"`
}

type StreamDelta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

// --- Project Models ---.
type Project struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Visibility  string    `json:"visibility"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Stats       *Stats    `json:"stats,omitempty"` // Included in "Get" only
}

type CreateProjectParams struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Visibility  string `json:"visibility,omitempty"` // "public" | "private"
}

// --- API Key Models ---.

type APIKey struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Prefix      string     `json:"prefix,omitempty"`
	Key         string     `json:"key,omitempty"` // Only present on Create
	Environment string     `json:"environment"`
	CreatedAt   time.Time  `json:"created_at"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
	UsageCount  int        `json:"usage_count,omitempty"`
}

type CreateAPIKeyParams struct {
	Name        string `json:"name"`
	Environment string `json:"environment"` // "production" | "test"
}

type KeyUsageStats struct {
	KeyID           string         `json:"key_id"`
	TotalRequests   int            `json:"total_requests"`
	TotalCostUSD    float64        `json:"total_cost_usd"`
	LastUsedAt      time.Time      `json:"last_used_at"`
	RequestsByDay   []DailyStat    `json:"requests_by_day"`
	RequestsByModel map[string]int `json:"requests_by_model"`
}

type DailyStat struct {
	Date    string  `json:"date"`
	Count   int     `json:"count"`
	CostUSD float64 `json:"cost_usd"`
}

// Metric Models.
type MetricsResponse struct {
	Period    string               `json:"period"`
	StartDate time.Time            `json:"start_date"`
	EndDate   time.Time            `json:"end_date"`
	Requests  RequestMetrics       `json:"requests"`
	Cost      CostMetrics          `json:"cost"`
	Tokens    TokenMetrics         `json:"tokens"`
	Latency   LatencyMetrics       `json:"latency"`
	Providers map[string]Breakdown `json:"providers"`
	Models    map[string]Breakdown `json:"models"`
}

type RequestMetrics struct {
	Total       int     `json:"total"`
	Success     int     `json:"success"`
	Error       int     `json:"error"`
	Filtered    int     `json:"filtered"`
	SuccessRate float64 `json:"success_rate"`
}

type CostMetrics struct {
	TotalUSD             float64 `json:"total_usd"`
	AveragePerRequestUSD float64 `json:"average_per_request_usd"`
}

type TokenMetrics struct {
	Prompt     int `json:"prompt"`
	Completion int `json:"completion"`
	Total      int `json:"total"`
}

type LatencyMetrics struct {
	AvgMS int `json:"avg_ms"`
	P50MS int `json:"p50_ms"`
	P90MS int `json:"p90_ms"`
	P99MS int `json:"p99_ms"`
}

type Breakdown struct {
	Requests int     `json:"requests"`
	CostUSD  float64 `json:"cost_usd"`
}
