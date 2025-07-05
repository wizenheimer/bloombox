package emailchecker

import (
	"context"
	"time"

	"github.com/wizenheimer/bloombox/internal/validators"
)

// ValidatorAdapter wraps internal validators to implement the emailchecker.Validator interface
type ValidatorAdapter struct {
	internal validators.Validator
}

func (v *ValidatorAdapter) Name() string {
	return v.internal.Name()
}

func (v *ValidatorAdapter) IsEnabled() bool {
	return v.internal.IsEnabled()
}

func (v *ValidatorAdapter) SetEnabled(enabled bool) {
	v.internal.SetEnabled(enabled)
}

func (v *ValidatorAdapter) Validate(ctx context.Context, email string) *ValidationResult {
	result := v.internal.Validate(ctx, email)
	// Convert validators.ValidationResult to emailchecker.ValidationResult
	return &ValidationResult{
		Valid:    result.Valid,
		Message:  result.Message,
		Details:  result.Details,
		Duration: result.Duration,
		Error:    result.Error,
	}
}

// NewSyntaxValidator creates a new syntax validator
func NewSyntaxValidator() Validator {
	return &ValidatorAdapter{internal: validators.NewSyntaxValidator()}
}

// NewDisposableValidator creates a new disposable email validator
func NewDisposableValidator(filename string, falsePositiveRate float64) (Validator, error) {
	internal, err := validators.NewDisposableValidator(filename, falsePositiveRate)
	if err != nil {
		return nil, err
	}
	return &ValidatorAdapter{internal: internal}, nil
}

// NewFreeValidator creates a new free email validator
func NewFreeValidator(filename string, falsePositiveRate float64) (Validator, error) {
	internal, err := validators.NewFreeValidator(filename, falsePositiveRate)
	if err != nil {
		return nil, err
	}
	return &ValidatorAdapter{internal: internal}, nil
}

// NewRoleValidator creates a new role-based email validator
func NewRoleValidator(filename string) (Validator, error) {
	internal, err := validators.NewRoleValidator(filename)
	if err != nil {
		return nil, err
	}
	return &ValidatorAdapter{internal: internal}, nil
}

// NewBanWordsValidator creates a new ban words validator
func NewBanWordsValidator(filename string) (Validator, error) {
	internal, err := validators.NewBanWordsValidator(filename)
	if err != nil {
		return nil, err
	}
	return &ValidatorAdapter{internal: internal}, nil
}

// NewBlackListEmailsValidator creates a new blacklisted emails validator
func NewBlackListEmailsValidator(filename string, falsePositiveRate float64) (Validator, error) {
	internal, err := validators.NewBlackListEmailsValidator(filename, falsePositiveRate)
	if err != nil {
		return nil, err
	}
	return &ValidatorAdapter{internal: internal}, nil
}

// NewBlackListDomainsValidator creates a new blacklisted domains validator
func NewBlackListDomainsValidator(filename string, falsePositiveRate float64) (Validator, error) {
	internal, err := validators.NewBlackListDomainsValidator(filename, falsePositiveRate)
	if err != nil {
		return nil, err
	}
	return &ValidatorAdapter{internal: internal}, nil
}

// NewMXValidator creates a new MX validator
func NewMXValidator(timeout time.Duration, dialFunc DialFunc) Validator {
	// Convert DialFunc type
	var internalDialFunc validators.DialFunc
	if dialFunc != nil {
		internalDialFunc = validators.DialFunc(dialFunc)
	}
	return &ValidatorAdapter{internal: validators.NewMXValidator(timeout, internalDialFunc)}
}

// NewSMTPValidator creates a new SMTP validator
func NewSMTPValidator(config *SMTPConfig) Validator {
	// Convert emailchecker.SMTPConfig to validators.SMTPConfig
	var internalDialFunc validators.DialFunc
	if config.DialFunc != nil {
		internalDialFunc = validators.DialFunc(config.DialFunc)
	}
	validatorConfig := &validators.SMTPConfig{
		Timeout:    config.Timeout,
		FromDomain: config.FromDomain,
		FromEmail:  config.FromEmail,
		EnableVRFY: config.EnableVRFY,
		EnableRCPT: config.EnableRCPT,
		DialFunc:   internalDialFunc,
	}
	return &ValidatorAdapter{internal: validators.NewSMTPValidator(validatorConfig)}
}

// NewGravatarValidator creates a new Gravatar validator
func NewGravatarValidator(timeout time.Duration) Validator {
	return &ValidatorAdapter{internal: validators.NewGravatarValidator(timeout)}
}
