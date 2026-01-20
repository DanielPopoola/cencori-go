package cencori

import (
	"errors"
	"fmt"
)

type APIError struct {
	StatusCode int
	Code       string         `json:"code"`
	Message    string         `json:"error"`
	Details    map[string]any `json:"details,omitempty"`
	Err        error          `json:"-"`
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("cencori: %s (code: %s, status: %d)", e.Message, e.Code, e.StatusCode)
	}
	return fmt.Sprintf("cencori: %s (status: %d)", e.Message, e.StatusCode)
}

func (e *APIError) Unwrap() error {
	return e.Err
}

func (e *APIError) fillSentinel() {
	switch e.Code {
	case "INVALID_API_KEY":
		e.Err = ErrInvalidAPIKey
	case "RATE_LIMIT_EXCEEDED":
		e.Err = ErrRateLimited
	case "INSUFFICIENT_CREDITS":
		e.Err = ErrInsufficientCredits
	case "TIER_RESTRICTED":
		e.Err = ErrTierRestricted
	case "INVALID_MODEL":
		e.Err = ErrInvalidModel
	case "PROVIDER_ERROR":
		e.Err = ErrProvider
	case "CONTENT_FILTERED":
		e.Err = ErrContentFiltered
	}
}

var (
	ErrInvalidAPIKey       = errors.New("INVALID_API_KEY")
	ErrSecurityViolation   = errors.New("SECURITY_VIOLATION")
	ErrRateLimited         = errors.New("RATE_LIMIT_EXCEEDED")
	ErrInsufficientCredits = errors.New("INSUFFICIENT_CREDITS")
	ErrTierRestricted      = errors.New("TIER_RESTRICTED")
	ErrInvalidModel        = errors.New("INVALID_MODEL")
	ErrProvider            = errors.New("PROVIDER_ERROR")
	ErrContentFiltered     = errors.New("CONTENT_FILTERED")
)
