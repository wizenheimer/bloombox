package validators

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/wizenheimer/bloombox/internal/filter"
	"github.com/wizenheimer/bloombox/internal/loader"
)

// BlackListEmailsValidator checks against specific blacklisted email addresses
type BlackListEmailsValidator struct {
	filter  filter.Filter
	enabled bool
}

// NewBlackListEmailsValidator creates a new blacklisted emails validator
func NewBlackListEmailsValidator(filename string, falsePositiveRate float64) (Validator, error) {
	l := loader.NewFileLoader()
	emails, err := l.LoadFromFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load blacklisted emails file: %w", err)
	}

	var f filter.Filter
	if falsePositiveRate == 0 || len(emails) < 1000 {
		f = filter.NewMapFilter()
	} else {
		f = filter.NewCuckooFilter(len(emails), falsePositiveRate)
	}

	for _, email := range emails {
		email = strings.ToLower(strings.TrimSpace(email))
		if email != "" {
			f.Add(email)
		}
	}

	return &BlackListEmailsValidator{
		filter:  f,
		enabled: true,
	}, nil
}

func (v *BlackListEmailsValidator) Name() string { return "blacklist_emails" }

func (v *BlackListEmailsValidator) IsEnabled() bool { return v.enabled }

func (v *BlackListEmailsValidator) SetEnabled(enabled bool) { v.enabled = enabled }

func (v *BlackListEmailsValidator) Validate(ctx context.Context, email string) *ValidationResult {
	start := time.Now()

	emailLower := strings.ToLower(strings.TrimSpace(email))
	isBlacklisted := v.filter.Contains(emailLower)

	result := &ValidationResult{
		Valid: !isBlacklisted,
		Message: func() string {
			if isBlacklisted {
				return "Email is blacklisted"
			} else {
				return "Email not in blacklist"
			}
		}(),
		Details: map[string]interface{}{
			"email":          emailLower,
			"is_blacklisted": isBlacklisted,
		},
		Duration: time.Since(start),
	}

	return result
}
