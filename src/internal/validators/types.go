package validators

import (
	"context"
	"net"
	"time"
)

// ValidationResult represents the result of a single validator
type ValidationResult struct {
	Valid    bool                   `json:"valid"`
	Message  string                 `json:"message"`
	Details  map[string]interface{} `json:"details,omitempty"`
	Duration time.Duration          `json:"duration"`
	Error    string                 `json:"error,omitempty"`
}

// Validator interface for all validation types
type Validator interface {
	Name() string
	Validate(ctx context.Context, email string) *ValidationResult
	IsEnabled() bool
	SetEnabled(enabled bool)
}

// MXRecord represents an MX record
type MXRecord struct {
	Host     string `json:"host"`
	Priority uint16 `json:"priority"`
	IP       string `json:"ip,omitempty"`
}

// SMTPResponse represents SMTP validation response
type SMTPResponse struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	CanReceive bool   `json:"can_receive"`
	IsMailbox  bool   `json:"is_mailbox"`
	IsCatchAll bool   `json:"is_catch_all"`
}

// DialFunc represents a custom dial function for proxy support
type DialFunc func(ctx context.Context, network, address string) (net.Conn, error)
