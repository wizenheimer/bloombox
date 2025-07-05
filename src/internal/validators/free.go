package validators

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/wizenheimer/bloombox/internal/filter"
	"github.com/wizenheimer/bloombox/internal/loader"
)

// FreeValidator checks against free email providers
type FreeValidator struct {
	filter  filter.Filter
	enabled bool
}

// NewFreeValidator creates a new free email validator
func NewFreeValidator(filename string, falsePositiveRate float64) (Validator, error) {
	l := loader.NewFileLoader()
	domains, err := l.LoadFromFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load free emails file: %w", err)
	}

	var f filter.Filter
	if falsePositiveRate == 0 {
		f = filter.NewMapFilter()
	} else {
		f = filter.NewCuckooFilter(len(domains), falsePositiveRate)
	}

	for _, domain := range domains {
		domain = strings.ToLower(strings.TrimSpace(domain))
		if domain != "" {
			f.Add(domain)
		}
	}

	return &FreeValidator{
		filter:  f,
		enabled: true,
	}, nil
}

func (v *FreeValidator) Name() string { return "free" }

func (v *FreeValidator) IsEnabled() bool { return v.enabled }

func (v *FreeValidator) SetEnabled(enabled bool) { v.enabled = enabled }

func (v *FreeValidator) Validate(ctx context.Context, email string) *ValidationResult {
	start := time.Now()

	domain := extractDomain(email)
	isFree := v.filter.Contains(domain)

	result := &ValidationResult{
		Valid: !isFree,
		Message: func() string {
			if isFree {
				return "Free email provider"
			} else {
				return "Not a free email provider"
			}
		}(),
		Details: map[string]interface{}{
			"domain":  domain,
			"is_free": isFree,
		},
		Duration: time.Since(start),
	}

	return result
}
