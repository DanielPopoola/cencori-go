package cencori

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  APIError
		want string
	}{
		{
			name: "with code",
			err: APIError{
				StatusCode: 401,
				Code:       "INVALID_API_KEY",
				Message:    "Invalid API key",
			},
			want: "cencori: Invalid API key (code: INVALID_API_KEY, status: 401)",
		},
		{
			name: "without code",
			err: APIError{
				StatusCode: 500,
				Message:    "Server error",
			},
			want: "cencori: Server error (status: 500)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIError_Unwrap(t *testing.T) {
	baseErr := errors.New("base error")
	apiErr := &APIError{
		StatusCode: 500,
		Message:    "wrapped",
		Err:        baseErr,
	}

	if !errors.Is(apiErr, baseErr) {
		t.Error("Unwrap() should allow errors.Is to find wrapped error")
	}
}

func TestAPIError_fillSentinel(t *testing.T) {
	tests := []struct {
		code     string
		sentinel error
	}{
		{"INVALID_API_KEY", ErrInvalidApiKey},
		{"RATE_LIMIT_EXCEEDED", ErrRateLimited},
		{"INSUFFICIENT_CREDITS", ErrInsufficientCredits},
		{"TIER_RESTRICTED", ErrTierRestricted},
		{"INVALID_MODEL", ErrInvalidModel},
		{"PROVIDER_ERROR", ErrProvider},
		{"CONTENT_FILTERED", ErrContentFiltered},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			apiErr := &APIError{Code: tt.code}
			apiErr.fillSentinel()

			if !errors.Is(apiErr, tt.sentinel) {
				t.Errorf("fillSentinel() for %s: errors.Is failed", tt.code)
			}
		})
	}
}

func TestHandleError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantErr    error
		wantCode   string
	}{
		{
			name:       "invalid api key",
			statusCode: 401,
			body:       `{"code":"INVALID_API_KEY","error":"Invalid key"}`,
			wantErr:    ErrInvalidApiKey,
			wantCode:   "INVALID_API_KEY",
		},
		{
			name:       "rate limited",
			statusCode: 429,
			body:       `{"code":"RATE_LIMIT_EXCEEDED","error":"Too many requests"}`,
			wantErr:    ErrRateLimited,
			wantCode:   "RATE_LIMIT_EXCEEDED",
		},
		{
			name:       "malformed json",
			statusCode: 500,
			body:       `not json`,
			wantErr:    nil, // No sentinel for unknown errors
			wantCode:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				StatusCode: tt.statusCode,
				Body:       io.NopCloser(bytes.NewBufferString(tt.body)),
			}

			err := handleError(resp)
			if err == nil {
				t.Fatal("handleError() should return an error")
			}

			apiErr, ok := err.(*APIError)
			if !ok {
				t.Fatalf("handleError() should return *APIError, got %T", err)
			}

			if apiErr.StatusCode != tt.statusCode {
				t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, tt.statusCode)
			}

			if tt.wantCode != "" && apiErr.Code != tt.wantCode {
				t.Errorf("Code = %s, want %s", apiErr.Code, tt.wantCode)
			}

			if tt.wantErr != nil && !errors.Is(apiErr, tt.wantErr) {
				t.Errorf("errors.Is() failed for %v", tt.wantErr)
			}
		})
	}
}
