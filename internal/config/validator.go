package config

import (
	"fmt"
	"strings"
)

// Validator is responsible for validating configurations
type Validator struct{}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{}
}

// Validate checks if the configuration has the required fields
func (v *Validator) Validate(config *Config) error {
	if config.Database.Host == "" {
		return fmt.Errorf("database.host is required")
	}
	if config.Database.User == "" {
		return fmt.Errorf("database.user is required")
	}
	if config.Database.Name == "" {
		return fmt.Errorf("database.name is required")
	}
	if config.JWT.Secret == "" {
		return fmt.Errorf("jwt.secret is required")
	}
	if config.App.Port <= 0 || config.App.Port > 65535 {
		return fmt.Errorf("app.port must be between 1 and 65535")
	}

	// Validate environment
	validEnvs := []string{"development", "staging", "production"}
	found := false
	for _, env := range validEnvs {
		if config.App.Environment == env {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("app.environment must be one of: %s", strings.Join(validEnvs, ", "))
	}

	return nil
}
