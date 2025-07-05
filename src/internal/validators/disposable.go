package validators

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/wizenheimer/bloombox/internal/filter"
	"github.com/wizenheimer/bloombox/internal/loader"
)

// DisposableValidator checks against disposable email providers
type DisposableValidator struct {
	filter  filter.Filter
	enabled bool
}

// NewDisposableValidator creates a new disposable email validator
func NewDisposableValidator(filename string, falsePositiveRate float64) (Validator, error) {
	l := loader.NewFileLoader()
	domains, err := l.LoadFromFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load disposable emails file: %w", err)
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

	return &DisposableValidator{
		filter:  f,
		enabled: true,
	}, nil
}

func (v *DisposableValidator) Name() string { return "disposable" }

func (v *DisposableValidator) IsEnabled() bool { return v.enabled }

func (v *DisposableValidator) SetEnabled(enabled bool) { v.enabled = enabled }

func (v *DisposableValidator) Validate(ctx context.Context, email string) *ValidationResult {
	start := time.Now()

	domain := extractDomain(email)
	isDisposable := v.filter.Contains(domain)

	result := &ValidationResult{
		Valid: !isDisposable,
		Message: func() string {
			if isDisposable {
				return "Disposable email detected"
			} else {
				return "Not a disposable email"
			}
		}(),
		Details: map[string]interface{}{
			"domain":        domain,
			"is_disposable": isDisposable,
		},
		Duration: time.Since(start),
	}

	return result
}
