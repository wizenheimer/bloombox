// internal/validators/banwords.go
package validators

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/wizenheimer/bloombox/internal/loader"
)

// BanWordsValidator checks for banned words in email username
type BanWordsValidator struct {
	banWords map[string]bool
	enabled  bool
}

// NewBanWordsValidator creates a new ban words validator
func NewBanWordsValidator(filename string) (Validator, error) {
	l := loader.NewFileLoader()
	words, err := l.LoadFromFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load ban words file: %w", err)
	}

	banWordsMap := make(map[string]bool)
	for _, word := range words {
		word = strings.ToLower(strings.TrimSpace(word))
		if word != "" {
			banWordsMap[word] = true
		}
	}

	return &BanWordsValidator{
		banWords: banWordsMap,
		enabled:  true,
	}, nil
}

func (v *BanWordsValidator) Name() string { return "banwords" }

func (v *BanWordsValidator) IsEnabled() bool { return v.enabled }

func (v *BanWordsValidator) SetEnabled(enabled bool) { v.enabled = enabled }

func (v *BanWordsValidator) Validate(ctx context.Context, email string) *ValidationResult {
	start := time.Now()

	localPart := extractLocalPart(email)
	localPartLower := strings.ToLower(localPart)

	result := &ValidationResult{
		Details: make(map[string]interface{}),
	}

	var foundBanWords []string

	// Check for exact matches
	if v.banWords[localPartLower] {
		foundBanWords = append(foundBanWords, localPartLower)
	}

	// Check for partial matches (ban words contained in username)
	for banWord := range v.banWords {
		if len(banWord) >= 3 && strings.Contains(localPartLower, banWord) {
			// Avoid duplicates
			found := false
			for _, existing := range foundBanWords {
				if existing == banWord {
					found = true
					break
				}
			}
			if !found {
				foundBanWords = append(foundBanWords, banWord)
			}
		}
	}

	hasBanWords := len(foundBanWords) > 0

	result.Valid = !hasBanWords
	result.Message = func() string {
		if hasBanWords {
			return fmt.Sprintf("Contains banned words: %s", strings.Join(foundBanWords, ", "))
		}
		return "No banned words found"
	}()

	result.Details["local_part"] = localPart
	result.Details["has_ban_words"] = hasBanWords
	result.Details["ban_words_found"] = foundBanWords
	result.Details["ban_words_count"] = len(foundBanWords)
	result.Duration = time.Since(start)

	return result
}
