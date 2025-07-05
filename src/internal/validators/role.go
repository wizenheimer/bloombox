package validators

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/wizenheimer/bloombox/internal/loader"
)

// RoleValidator checks for role-based email addresses
type RoleValidator struct {
	roleEmails map[string]bool
	enabled    bool
}

// NewRoleValidator creates a new role-based email validator
func NewRoleValidator(filename string) (Validator, error) {
	l := loader.NewFileLoader()
	roles, err := l.LoadFromFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load role emails file: %w", err)
	}

	roleMap := make(map[string]bool)
	for _, role := range roles {
		role = strings.ToLower(strings.TrimSpace(role))
		if role != "" {
			roleMap[role] = true
		}
	}

	return &RoleValidator{
		roleEmails: roleMap,
		enabled:    true,
	}, nil
}

func (v *RoleValidator) Name() string { return "role" }

func (v *RoleValidator) IsEnabled() bool { return v.enabled }

func (v *RoleValidator) SetEnabled(enabled bool) { v.enabled = enabled }

func (v *RoleValidator) Validate(ctx context.Context, email string) *ValidationResult {
	start := time.Now()

	localPart := extractLocalPart(email)
	isRole := v.roleEmails[localPart]

	result := &ValidationResult{
		Valid: !isRole,
		Message: func() string {
			if isRole {
				return "Role-based email address"
			} else {
				return "Not a role-based email"
			}
		}(),
		Details: map[string]interface{}{
			"local_part": localPart,
			"is_role":    isRole,
		},
		Duration: time.Since(start),
	}

	return result
}
