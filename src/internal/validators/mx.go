// internal/validators/mx.go
package validators

import (
	"context"
	"net"
	"sort"
	"time"
)

// MXValidator validates MX records
type MXValidator struct {
	timeout  time.Duration
	dialFunc DialFunc
	enabled  bool
}

// NewMXValidator creates a new MX validator
func NewMXValidator(timeout time.Duration, dialFunc DialFunc) Validator {
	if dialFunc == nil {
		dialFunc = (&net.Dialer{Timeout: timeout}).DialContext
	}

	return &MXValidator{
		timeout:  timeout,
		dialFunc: dialFunc,
		enabled:  true,
	}
}

func (v *MXValidator) Name() string { return "mx" }

func (v *MXValidator) IsEnabled() bool { return v.enabled }

func (v *MXValidator) SetEnabled(enabled bool) { v.enabled = enabled }

func (v *MXValidator) Validate(ctx context.Context, email string) *ValidationResult {
	start := time.Now()

	domain := extractDomain(email)

	resolver := &net.Resolver{
		PreferGo: true,
		Dial:     v.dialFunc,
	}

	mxRecords, err := resolver.LookupMX(ctx, domain)

	result := &ValidationResult{
		Details:  make(map[string]interface{}),
		Duration: time.Since(start),
	}

	if err != nil {
		// Try A record as fallback
		ips, aErr := resolver.LookupIPAddr(ctx, domain)
		if aErr != nil {
			result.Valid = false
			result.Message = "No MX or A records found"
			result.Error = err.Error()
		} else {
			result.Valid = true
			result.Message = "No MX record, but domain has A record (implicit MX)"
			result.Details["implicit_mx"] = true
			result.Details["a_records"] = len(ips)
		}
	} else {
		result.Valid = true
		result.Message = "Valid MX records found"

		// Sort MX records by priority
		sort.Slice(mxRecords, func(i, j int) bool {
			return mxRecords[i].Pref < mxRecords[j].Pref
		})

		mxDetails := make([]MXRecord, len(mxRecords))
		for i, mx := range mxRecords {
			mxDetails[i] = MXRecord{
				Host:     mx.Host,
				Priority: mx.Pref,
			}

			// Try to resolve IP for the MX host
			if ips, err := resolver.LookupIPAddr(ctx, mx.Host); err == nil && len(ips) > 0 {
				mxDetails[i].IP = ips[0].IP.String()
			}
		}

		result.Details["mx_records"] = mxDetails
		result.Details["mx_count"] = len(mxRecords)
		result.Details["primary_mx"] = mxDetails[0].Host
	}

	return result
}
