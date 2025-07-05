// internal/validators/utils.go
package validators

import "strings"

// extractDomain extracts domain from email address
func extractDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ""
	}
	return strings.ToLower(parts[1])
}

// extractLocalPart extracts local part from email address
func extractLocalPart(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ""
	}
	return strings.ToLower(parts[0])
}
