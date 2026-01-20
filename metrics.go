package cencori

import (
	"context"
	"fmt"
)

// MetricsService provides methods for fetching analytics for projects
// It uses a Client to communicate with the metrics API endpoints.
type MetricsService struct {
	client *Client
}

// Get retrieves metrics for the specified period.
// It sends a GET request to the metrics API endpoint and returns the metrics response.
// If the request fails or the context is cancelled, an error is returned.
func (s *MetricsService) Get(ctx context.Context, period string) (*MetricsResponse, error) {
	path := fmt.Sprintf("/api/v1/metrics/%s", period)
	return doRequest[any, MetricsResponse](s.client, ctx, "GET", path, nil)
}
