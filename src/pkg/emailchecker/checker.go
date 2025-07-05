package emailchecker

import (
	"context"
	"fmt"
	"net/mail"
	"strings"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
)

// EmailChecker is the main email checking service
type EmailChecker struct {
	config     *Config
	validators map[string]Validator
	cache      *lru.Cache[string, *CheckResult]
	semaphore  chan struct{}
	mu         sync.RWMutex
}

// New creates a new EmailChecker instance
func New(config *Config) (*EmailChecker, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// Initialize cache
	cache, err := lru.New[string, *CheckResult](config.CacheSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create cache: %w", err)
	}

	checker := &EmailChecker{
		config:     config,
		validators: make(map[string]Validator),
		cache:      cache,
		semaphore:  make(chan struct{}, config.MaxConcurrentValidations),
	}

	// Initialize validators
	if err := checker.initializeValidators(); err != nil {
		return nil, fmt.Errorf("failed to initialize validators: %w", err)
	}

	return checker, nil
}

// initializeValidators sets up all validators
func (e *EmailChecker) initializeValidators() error {
	// Always create syntax validator
	e.validators["syntax"] = NewSyntaxValidator()

	// File-based validators - only create if files are provided
	if e.config.DisposableEmailsFile != "" {
		validator, err := NewDisposableValidator(
			e.config.DisposableEmailsFile,
			e.config.FalsePositiveRate,
		)
		if err != nil {
			return fmt.Errorf("failed to create disposable validator: %w", err)
		}
		e.validators["disposable"] = validator
	}

	if e.config.FreeEmailsFile != "" {
		validator, err := NewFreeValidator(
			e.config.FreeEmailsFile,
			e.config.FalsePositiveRate,
		)
		if err != nil {
			return fmt.Errorf("failed to create free validator: %w", err)
		}
		e.validators["free"] = validator
	}

	if e.config.RoleEmailsFile != "" {
		validator, err := NewRoleValidator(e.config.RoleEmailsFile)
		if err != nil {
			return fmt.Errorf("failed to create role validator: %w", err)
		}
		e.validators["role"] = validator
	}

	if e.config.BanWordsFile != "" {
		validator, err := NewBanWordsValidator(e.config.BanWordsFile)
		if err != nil {
			return fmt.Errorf("failed to create ban words validator: %w", err)
		}
		e.validators["banwords"] = validator
	}

	if e.config.BlackListEmailsFile != "" {
		validator, err := NewBlackListEmailsValidator(
			e.config.BlackListEmailsFile,
			e.config.FalsePositiveRate,
		)
		if err != nil {
			return fmt.Errorf("failed to create blacklist emails validator: %w", err)
		}
		e.validators["blacklist_emails"] = validator
	}

	if e.config.BlackListDomainsFile != "" {
		validator, err := NewBlackListDomainsValidator(
			e.config.BlackListDomainsFile,
			e.config.FalsePositiveRate,
		)
		if err != nil {
			return fmt.Errorf("failed to create blacklist domains validator: %w", err)
		}
		e.validators["blacklist_domains"] = validator
	}

	// Network validators - always available
	e.validators["mx"] = NewMXValidator(e.config.ValidationTimeout, e.config.DialFunc)

	smtpValidator := NewSMTPValidator(&SMTPConfig{
		Timeout:    e.config.SMTPTimeout,
		FromDomain: e.config.SMTPFromDomain,
		FromEmail:  e.config.SMTPFromEmail,
		EnableVRFY: e.config.EnableSMTPVRFY,
		EnableRCPT: e.config.EnableSMTPRCPT,
		DialFunc:   e.config.DialFunc,
	})
	e.validators["smtp"] = smtpValidator

	e.validators["gravatar"] = NewGravatarValidator(e.config.ValidationTimeout)

	// Set enabled validators
	e.setEnabledValidators()

	return nil
}

// setEnabledValidators enables/disables validators based on config
func (e *EmailChecker) setEnabledValidators() {
	enabledMap := make(map[string]bool)
	for _, name := range e.config.EnabledValidators {
		enabledMap[name] = true
	}

	for name, validator := range e.validators {
		validator.SetEnabled(enabledMap[name])
	}
}

// Check performs email validation
func (e *EmailChecker) Check(email string) *CheckResult {
	return e.CheckWithValidators(email, nil)
}

// CheckWithValidators performs email validation with specific validators
func (e *EmailChecker) CheckWithValidators(email string, validatorNames []string) *CheckResult {
	email = strings.ToLower(strings.TrimSpace(email))
	start := time.Now()

	// Check cache first
	cacheKey := e.buildCacheKey(email, validatorNames)
	if cached, ok := e.cache.Get(cacheKey); ok {
		if time.Since(cached.Timestamp) < e.config.CacheTimeout {
			return cached
		}
		e.cache.Remove(cacheKey)
	}

	// Create result
	result := &CheckResult{
		Email:     email,
		Timestamp: time.Now(),
		Results:   make(map[string]*ValidationResult),
	}

	// Validate basic email format first
	if _, err := mail.ParseAddress(email); err != nil {
		result.Duration = time.Since(start)
		result.Results["syntax"] = &ValidationResult{
			Valid:   false,
			Message: "Invalid email syntax",
			Error:   err.Error(),
		}
		return result
	}

	// Determine which validators to run
	validatorsToRun := e.getValidatorsToRun(validatorNames)

	// Run validators
	e.runValidators(result, validatorsToRun)

	// Calculate overall validity and summary
	e.calculateSummary(result)

	result.Duration = time.Since(start)

	// Cache result
	e.cache.Add(cacheKey, result)

	return result
}

// runValidators executes the specified validators
func (e *EmailChecker) runValidators(result *CheckResult, validatorNames []string) {
	var wg sync.WaitGroup
	resultsChan := make(chan struct {
		name string
		res  *ValidationResult
	}, len(validatorNames))

	for _, name := range validatorNames {
		validator, exists := e.validators[name]
		if !exists || !validator.IsEnabled() {
			continue
		}

		wg.Add(1)
		go func(vName string, v Validator) {
			defer wg.Done()

			// Acquire semaphore
			select {
			case e.semaphore <- struct{}{}:
				defer func() { <-e.semaphore }()
			default:
				// Skip if too many concurrent validations
				resultsChan <- struct {
					name string
					res  *ValidationResult
				}{vName, &ValidationResult{
					Valid:   false,
					Message: "skipped due to rate limiting",
				}}
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), e.config.ValidationTimeout)
			defer cancel()

			res := v.Validate(ctx, result.Email)

			resultsChan <- struct {
				name string
				res  *ValidationResult
			}{vName, res}
		}(name, validator)
	}

	// Close channel when all goroutines complete
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results
	for res := range resultsChan {
		result.Results[res.name] = res.res
	}
}

// getValidatorsToRun determines which validators to execute
func (e *EmailChecker) getValidatorsToRun(requested []string) []string {
	if len(requested) == 0 {
		// Return all enabled validators
		var enabled []string
		for name, validator := range e.validators {
			if validator.IsEnabled() {
				enabled = append(enabled, name)
			}
		}
		return enabled
	}

	// Filter requested validators to only enabled ones
	var valid []string
	for _, name := range requested {
		if validator, exists := e.validators[name]; exists && validator.IsEnabled() {
			valid = append(valid, name)
		}
	}
	return valid
}

// buildCacheKey creates a cache key for the given email and validators
func (e *EmailChecker) buildCacheKey(email string, validatorNames []string) string {
	if len(validatorNames) == 0 {
		return email
	}
	return fmt.Sprintf("%s:%v", email, validatorNames)
}

// GetValidators returns available validators and their status
func (e *EmailChecker) GetValidators() map[string]bool {
	validators := make(map[string]bool)
	for name, validator := range e.validators {
		validators[name] = validator.IsEnabled()
	}
	return validators
}

// SetValidatorEnabled enables or disables a specific validator
func (e *EmailChecker) SetValidatorEnabled(name string, enabled bool) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	validator, exists := e.validators[name]
	if !exists {
		return fmt.Errorf("validator %s not found", name)
	}

	validator.SetEnabled(enabled)
	return nil
}

// GetValidatorNames returns all available validator names
func (e *EmailChecker) GetValidatorNames() []string {
	var names []string
	for name := range e.validators {
		names = append(names, name)
	}
	return names
}

// calculateSummary calculates overall validity and summary from validation results
func (e *EmailChecker) calculateSummary(result *CheckResult) {
	// Calculate overall validity - email is valid if all enabled validators pass
	result.IsValid = true
	summary := &CheckSummary{}

	for name, validationResult := range result.Results {
		if !validationResult.Valid {
			result.IsValid = false
		}

		// Extract specific information for summary
		switch name {
		case "disposable":
			if !validationResult.Valid {
				summary.IsDisposable = true
			}
		case "free":
			if !validationResult.Valid {
				summary.IsFree = true
			}
		case "role":
			if !validationResult.Valid {
				summary.IsRole = true
			}
		}
	}

	result.Summary = summary
}
