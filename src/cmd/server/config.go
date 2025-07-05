package main

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/wizenheimer/bloombox/pkg/emailchecker"
)

// loadConfigFromEnv loads configuration from environment variables
func loadConfigFromEnv() *emailchecker.Config {
	config := emailchecker.DefaultConfig()

	if val := os.Getenv("FREE_EMAILS_FILE"); val != "" {
		config.FreeEmailsFile = val
	}
	if val := os.Getenv("DISPOSABLE_EMAILS_FILE"); val != "" {
		config.DisposableEmailsFile = val
	}
	if val := os.Getenv("ROLE_EMAILS_FILE"); val != "" {
		config.RoleEmailsFile = val
	}
	if val := os.Getenv("BAN_WORDS_FILE"); val != "" {
		config.BanWordsFile = val
	}
	if val := os.Getenv("BLACKLIST_EMAILS_FILE"); val != "" {
		config.BlackListEmailsFile = val
	}
	if val := os.Getenv("BLACKLIST_DOMAINS_FILE"); val != "" {
		config.BlackListDomainsFile = val
	}
	if val := os.Getenv("FALSE_POSITIVE_RATE"); val != "" {
		if rate, err := strconv.ParseFloat(val, 64); err == nil {
			config.FalsePositiveRate = rate
		}
	}
	if val := os.Getenv("CACHE_SIZE"); val != "" {
		if size, err := strconv.Atoi(val); err == nil {
			config.CacheSize = size
		}
	}
	if val := os.Getenv("VALIDATION_TIMEOUT"); val != "" {
		if timeout, err := time.ParseDuration(val); err == nil {
			config.ValidationTimeout = timeout
		}
	}
	if val := os.Getenv("SMTP_TIMEOUT"); val != "" {
		if timeout, err := time.ParseDuration(val); err == nil {
			config.SMTPTimeout = timeout
		}
	}
	if val := os.Getenv("SMTP_FROM_DOMAIN"); val != "" {
		config.SMTPFromDomain = val
	}
	if val := os.Getenv("SMTP_FROM_EMAIL"); val != "" {
		config.SMTPFromEmail = val
	}
	if val := os.Getenv("ENABLED_VALIDATORS"); val != "" {
		config.EnabledValidators = strings.Split(val, ",")
	}

	return config
}
