package emailchecker

import (
	"context"
	"net"
)

// Validator interface for all validation types
type Validator interface {
	Name() string
	Validate(ctx context.Context, email string) *ValidationResult
	IsEnabled() bool
	SetEnabled(enabled bool)
}

// Filter interface for different filtering implementations
type Filter interface {
	Add(item string) error
	Contains(item string) bool
	Size() int
}

// DialFunc represents a custom dial function for proxy support
type DialFunc func(ctx context.Context, network, address string) (net.Conn, error)
