package validators

import (
	"context"
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// GravatarValidator checks if email has a Gravatar account
type GravatarValidator struct {
	timeout time.Duration
	client  *http.Client
	enabled bool
}

// NewGravatarValidator creates a new Gravatar validator
func NewGravatarValidator(timeout time.Duration) Validator {
	return &GravatarValidator{
		timeout: timeout,
		client: &http.Client{
			Timeout: timeout,
		},
		enabled: true,
	}
}

func (v *GravatarValidator) Name() string { return "gravatar" }

func (v *GravatarValidator) IsEnabled() bool { return v.enabled }

func (v *GravatarValidator) SetEnabled(enabled bool) { v.enabled = enabled }

func (v *GravatarValidator) Validate(ctx context.Context, email string) *ValidationResult {
	start := time.Now()

	result := &ValidationResult{
		Details: make(map[string]interface{}),
	}

	// Generate MD5 hash of email
	email = strings.ToLower(strings.TrimSpace(email))
	hash := fmt.Sprintf("%x", md5.Sum([]byte(email)))

	// Check Gravatar URL
	gravatarURL := fmt.Sprintf("https://www.gravatar.com/avatar/%s?d=404", hash)

	req, err := http.NewRequestWithContext(ctx, "GET", gravatarURL, nil)
	if err != nil {
		result.Valid = false
		result.Message = "Failed to create request"
		result.Error = err.Error()
		result.Duration = time.Since(start)
		return result
	}

	resp, err := v.client.Do(req)
	if err != nil {
		result.Valid = false
		result.Message = "Failed to check Gravatar"
		result.Error = err.Error()
		result.Duration = time.Since(start)
		return result
	}
	defer resp.Body.Close()

	hasGravatar := resp.StatusCode == 200

	result.Valid = hasGravatar
	result.Message = func() string {
		if hasGravatar {
			return "Gravatar account exists"
		}
		return "No Gravatar account found"
	}()
	result.Details["gravatar_url"] = gravatarURL
	result.Details["gravatar_hash"] = hash
	result.Details["has_gravatar"] = hasGravatar
	result.Duration = time.Since(start)

	return result
}
