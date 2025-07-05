package validators

import (
	"context"
	"strings"
	"time"

	"github.com/wizenheimer/bloombox/internal/filter"
	"github.com/wizenheimer/bloombox/internal/loader"
)

// BlackListDomainsValidator checks against blacklisted domains
type BlackListDomainsValidator struct {
	filter  filter.Filter
	enabled bool
}

// NewBlackListDomainsValidator creates a new blacklisted domains validator
func NewBlackListDomainsValidator(filename string, falsePositiveRate float64) (Validator, error) {
	l := loader.NewFileLoader()
	domains, err := l.LoadFromFile(filename)
	if err != nil {
		// If file doesn't exist, start with empty list
		domains = []string{}
	}

	var f filter.Filter
	if falsePositiveRate == 0 || len(domains) < 1000 {
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

	return &BlackListDomainsValidator{
		filter:  f,
		enabled: true,
	}, nil
}

func (v *BlackListDomainsValidator) Name() string { return "blacklist_domains" }

func (v *BlackListDomainsValidator) IsEnabled() bool { return v.enabled }

func (v *BlackListDomainsValidator) SetEnabled(enabled bool) { v.enabled = enabled }

func (v *BlackListDomainsValidator) Validate(ctx context.Context, email string) *ValidationResult {
	start := time.Now()

	domain := extractDomain(email)
	isBlacklisted := v.filter.Contains(domain)

	result := &ValidationResult{
		Valid: !isBlacklisted,
		Message: func() string {
			if isBlacklisted {
				return "Domain is blacklisted"
			} else {
				return "Domain not in blacklist"
			}
		}(),
		Details: map[string]interface{}{
			"domain":         domain,
			"is_blacklisted": isBlacklisted,
		},
		Duration: time.Since(start),
	}

	return result
}
