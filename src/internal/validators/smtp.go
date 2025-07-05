package validators

import (
	"context"
	"fmt"
	"net"
	"net/smtp"
	"net/textproto"
	"sort"
	"time"
)

// SMTPConfig holds SMTP validator configuration
type SMTPConfig struct {
	Timeout    time.Duration
	FromDomain string
	FromEmail  string
	EnableVRFY bool
	EnableRCPT bool
	DialFunc   DialFunc
}

// SMTPValidator validates email addresses via SMTP
type SMTPValidator struct {
	config  *SMTPConfig
	enabled bool
}

// NewSMTPValidator creates a new SMTP validator
func NewSMTPValidator(config *SMTPConfig) Validator {
	if config.DialFunc == nil {
		config.DialFunc = (&net.Dialer{Timeout: config.Timeout}).DialContext
	}

	return &SMTPValidator{
		config:  config,
		enabled: true,
	}
}

func (v *SMTPValidator) Name() string { return "smtp" }

func (v *SMTPValidator) IsEnabled() bool { return v.enabled }

func (v *SMTPValidator) SetEnabled(enabled bool) { v.enabled = enabled }

func (v *SMTPValidator) Validate(ctx context.Context, email string) *ValidationResult {
	start := time.Now()

	domain := extractDomain(email)

	result := &ValidationResult{
		Details: make(map[string]interface{}),
	}

	// Get MX records first
	mxRecords, err := v.getMXRecords(ctx, domain)
	if err != nil {
		result.Valid = false
		result.Message = "Could not resolve MX records"
		result.Error = err.Error()
		result.Duration = time.Since(start)
		return result
	}

	if len(mxRecords) == 0 {
		result.Valid = false
		result.Message = "No MX records found"
		result.Duration = time.Since(start)
		return result
	}

	// Try SMTP validation with the primary MX server
	smtpResult := v.validateViaSMTP(ctx, email, mxRecords[0])

	result.Valid = smtpResult.CanReceive
	result.Message = smtpResult.Message
	result.Details["smtp_response"] = smtpResult
	result.Details["mx_host"] = mxRecords[0].Host
	result.Duration = time.Since(start)

	return result
}

// getMXRecords retrieves and sorts MX records
func (v *SMTPValidator) getMXRecords(ctx context.Context, domain string) ([]*net.MX, error) {
	resolver := &net.Resolver{
		PreferGo: true,
		Dial:     v.config.DialFunc,
	}

	mxRecords, err := resolver.LookupMX(ctx, domain)
	if err != nil {
		return nil, err
	}

	// Sort by priority
	sort.Slice(mxRecords, func(i, j int) bool {
		return mxRecords[i].Pref < mxRecords[j].Pref
	})

	return mxRecords, nil
}

// validateViaSMTP performs SMTP validation
func (v *SMTPValidator) validateViaSMTP(ctx context.Context, email string, mx *net.MX) *SMTPResponse {
	response := &SMTPResponse{}

	// Connect to MX server
	conn, err := v.config.DialFunc(ctx, "tcp", fmt.Sprintf("%s:25", mx.Host))
	if err != nil {
		response.Code = 0
		response.Message = fmt.Sprintf("Connection failed: %v", err)
		return response
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, mx.Host)
	if err != nil {
		response.Code = 0
		response.Message = fmt.Sprintf("SMTP client creation failed: %v", err)
		return response
	}
	defer client.Quit()

	// HELO/EHLO
	if err := client.Hello(v.config.FromDomain); err != nil {
		response.Code = 0
		response.Message = fmt.Sprintf("HELO failed: %v", err)
		return response
	}

	// MAIL FROM
	if err := client.Mail(v.config.FromEmail); err != nil {
		response.Code = 0
		response.Message = fmt.Sprintf("MAIL FROM failed: %v", err)
		return response
	}

	// Try VRFY command if enabled
	if v.config.EnableVRFY {
		if vrfyResult := v.tryVRFY(client, email); vrfyResult != nil {
			return vrfyResult
		}
	}

	// Try RCPT TO if enabled
	if v.config.EnableRCPT {
		return v.tryRCPT(client, email)
	}

	response.Code = 250
	response.Message = "SMTP connection successful"
	response.CanReceive = true
	return response
}

// tryVRFY attempts VRFY command
func (v *SMTPValidator) tryVRFY(_ *smtp.Client, _ string) *SMTPResponse {
	// VRFY is not directly supported by Go's smtp package
	// We would need to implement raw SMTP commands
	return nil
}

// tryRCPT attempts RCPT TO command
func (v *SMTPValidator) tryRCPT(client *smtp.Client, email string) *SMTPResponse {
	response := &SMTPResponse{}

	err := client.Rcpt(email)
	if err != nil {
		if smtpErr, ok := err.(*textproto.Error); ok {
			response.Code = smtpErr.Code
			response.Message = smtpErr.Msg

			// Analyze SMTP response codes
			switch {
			case smtpErr.Code >= 200 && smtpErr.Code < 300:
				response.CanReceive = true
				response.IsMailbox = true
			case smtpErr.Code == 550:
				response.CanReceive = false
				response.Message = "Mailbox does not exist"
			case smtpErr.Code == 551:
				response.CanReceive = false
				response.Message = "User not local"
			case smtpErr.Code >= 400 && smtpErr.Code < 500:
				response.CanReceive = false
				response.Message = "Temporary failure"
			default:
				response.CanReceive = false
			}
		} else {
			response.Code = 0
			response.Message = err.Error()
			response.CanReceive = false
		}
	} else {
		response.Code = 250
		response.Message = "Recipient accepted"
		response.CanReceive = true
		response.IsMailbox = true
	}

	return response
}
