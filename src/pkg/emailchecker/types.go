package emailchecker

import (
	"time"
)

// CheckRequest represents a validation request
type CheckRequest struct {
	Email      string   `json:"email"`
	Validators []string `json:"validators,omitempty"` // Specific validators to run
	Timeout    int      `json:"timeout,omitempty"`    // Timeout in seconds
}

// CheckResult represents the result of an email check
type CheckResult struct {
	Email     string                       `json:"email"`
	Timestamp time.Time                    `json:"timestamp"`
	Duration  time.Duration                `json:"duration"`
	Results   map[string]*ValidationResult `json:"results"`
	IsValid   bool                         `json:"is_valid"`
	Summary   *CheckSummary                `json:"summary,omitempty"`
}

// CheckSummary provides a quick summary of validation results
type CheckSummary struct {
	IsDisposable bool `json:"is_disposable"`
	IsFree       bool `json:"is_free"`
	IsRole       bool `json:"is_role"`
}

// ValidationResult represents the result of a single validator
type ValidationResult struct {
	Valid    bool                   `json:"valid"`
	Message  string                 `json:"message"`
	Details  map[string]interface{} `json:"details,omitempty"`
	Duration time.Duration          `json:"duration"`
	Error    string                 `json:"error,omitempty"`
}

// MXRecord represents an MX record
type MXRecord struct {
	Host     string `json:"host"`
	Priority uint16 `json:"priority"`
	IP       string `json:"ip,omitempty"`
}

// SMTPResponse represents SMTP validation response
type SMTPResponse struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	CanReceive bool   `json:"can_receive"`
	IsMailbox  bool   `json:"is_mailbox"`
	IsCatchAll bool   `json:"is_catch_all"`
}

// SMTPConfig holds SMTP validator configuration
type SMTPConfig struct {
	Timeout    time.Duration `json:"timeout"`
	FromDomain string        `json:"from_domain"`
	FromEmail  string        `json:"from_email"`
	EnableVRFY bool          `json:"enable_vrfy"`
	EnableRCPT bool          `json:"enable_rcpt"`
	DialFunc   DialFunc      `json:"-"`
}
