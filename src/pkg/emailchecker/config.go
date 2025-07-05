package emailchecker

import (
	"net"
	"time"
)

// Config holds configuration for the email checker
type Config struct {
	// File paths - all optional, validators disabled if files not provided
	FreeEmailsFile       string `json:"free_emails_file,omitempty"`
	DisposableEmailsFile string `json:"disposable_emails_file,omitempty"`
	RoleEmailsFile       string `json:"role_emails_file,omitempty"`
	BanWordsFile         string `json:"ban_words_file,omitempty"`
	BlackListEmailsFile  string `json:"blacklist_emails_file,omitempty"`
	BlackListDomainsFile string `json:"blacklist_domains_file,omitempty"`

	// Validator settings
	EnabledValidators []string `json:"enabled_validators,omitempty"`

	// Filter settings
	FalsePositiveRate float64 `json:"false_positive_rate"` // 0 means use map
	CacheSize         int     `json:"cache_size"`

	// Network settings
	ValidationTimeout        time.Duration `json:"validation_timeout"`
	MaxConcurrentValidations int           `json:"max_concurrent_validations"`
	DialFunc                 DialFunc      `json:"-"` // Custom dial function for proxy

	// SMTP settings
	SMTPTimeout    time.Duration `json:"smtp_timeout"`
	SMTPFromDomain string        `json:"smtp_from_domain"`
	SMTPFromEmail  string        `json:"smtp_from_email"`
	EnableSMTPVRFY bool          `json:"enable_smtp_vrfy"`
	EnableSMTPRCPT bool          `json:"enable_smtp_rcpt"`

	// Cache settings
	CacheTimeout time.Duration `json:"cache_timeout"`
}

// DefaultConfig returns a minimal default configuration
func DefaultConfig() *Config {
	return &Config{
		EnabledValidators:        []string{"syntax"}, // Only enable syntax by default
		FalsePositiveRate:        0.01,
		CacheSize:                1000,
		ValidationTimeout:        5 * time.Second,
		MaxConcurrentValidations: 10,
		SMTPTimeout:              5 * time.Second,
		SMTPFromDomain:           "example.com",
		SMTPFromEmail:            "test@example.com",
		EnableSMTPVRFY:           false,
		EnableSMTPRCPT:           true,
		CacheTimeout:             10 * time.Minute,
		DialFunc:                 (&net.Dialer{Timeout: 5 * time.Second}).DialContext,
	}
}
