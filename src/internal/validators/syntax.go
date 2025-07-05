package validators

import (
	"context"
	"net/mail"
	"time"
)

// SyntaxValidator validates email syntax using Go's built-in RFC 5322 parser
type SyntaxValidator struct {
	enabled bool
}

// NewSyntaxValidator creates a new syntax validator
func NewSyntaxValidator() Validator {
	return &SyntaxValidator{enabled: true}
}

func (v *SyntaxValidator) Name() string { return "syntax" }

func (v *SyntaxValidator) IsEnabled() bool { return v.enabled }

func (v *SyntaxValidator) SetEnabled(enabled bool) { v.enabled = enabled }

func (v *SyntaxValidator) Validate(ctx context.Context, email string) *ValidationResult {
	start := time.Now()

	result := &ValidationResult{
		Details: make(map[string]interface{}),
	}

	addr, err := mail.ParseAddress(email)
	if err != nil {
		result.Valid = false
		result.Message = "Invalid email syntax"
		result.Error = err.Error()
	} else {
		result.Valid = true
		result.Message = "Valid email syntax"
		result.Details["parsed_address"] = addr.Address
		result.Details["display_name"] = addr.Name
	}

	result.Duration = time.Since(start)
	return result
}
