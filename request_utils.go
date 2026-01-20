package cencori

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func handleError(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body) //nolint:errcheck // reading the response body; error can be ignored here.
	var apiErr APIError
	if err := json.Unmarshal(body, &apiErr); err != nil {
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}
	apiErr.StatusCode = resp.StatusCode
	apiErr.fillSentinel() // Attach the ErrInvalidAPIKey etc.
	return &apiErr
}

// doRequest performs an HTTP request and returns the decoded response.
// It marshals the request body to JSON, sets required headers including the API key,
// executes the HTTP request, and decodes the response body into the specified response type.
//
// Type parameters:
//   - Req: the type of the request body
//   - Resp: the type of the response body
//
// Parameters:
//   - c: the HTTP client configuration containing BaseURL and ApiKey
//   - ctx: the context for the request
//   - method: the HTTP method (GET, POST, etc.)
//   - path: the request path appended to BaseURL
//   - body: the request body to marshal; if nil, no body is sent
//
// Returns:
//   - a pointer to the decoded response of type Resp
//   - an error if marshaling, creating, executing the request, or decoding the response fails
//   - an error if the response status code is not OK (200)
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
	defer resp.Body.Close() //nolint:errcheck // Closing the response body; error can be ignored here.

	if resp.StatusCode != http.StatusOK {
		return nil, handleError(resp)
	}

	var result Resp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &result, nil
}
