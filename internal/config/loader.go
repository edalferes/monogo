package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// LoadOptions options for loading configuration
type LoadOptions struct {
	ConfigPath string
	ConfigName string
	ConfigType string
}

// DefaultLoadOptions returns default options
func DefaultLoadOptions() LoadOptions {
	return LoadOptions{
		ConfigPath: ".",
		ConfigName: "config",
		ConfigType: "yaml",
	}
}

// Loader is responsible for loading configurations
type Loader struct {
	viper *viper.Viper
}

// NewLoader creates a new configuration loader
func NewLoader() *Loader {
	return &Loader{
		viper: viper.New(),
	}
}

// Load loads configuration from YAML files and environment variables
func (l *Loader) Load(opts ...LoadOptions) (*Config, error) {
	options := DefaultLoadOptions()
	if len(opts) > 0 {
		options = opts[0]
	}

	// Configure search for configuration files
	l.viper.SetConfigName(options.ConfigName)
	l.viper.SetConfigType(options.ConfigType)
	l.viper.AddConfigPath(options.ConfigPath)
	l.viper.AddConfigPath("./config")
	l.viper.AddConfigPath(".")

	// Set default values
	l.setDefaults()

	// Configure reading of environment variables
	l.viper.AutomaticEnv()
	l.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Try to read configuration file
	if err := l.viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading configuration file: %w", err)
		}
		// If file not found, use only env vars and defaults
	}

	var config Config
	if err := l.viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling configuration: %w", err)
	}

	return &config, nil
}

// setDefaults sets default values
func (l *Loader) setDefaults() {
	// App defaults
	l.viper.SetDefault("app.name", "monogo")
	l.viper.SetDefault("app.version", "1.0.0")
	l.viper.SetDefault("app.environment", "development")
	l.viper.SetDefault("app.port", 8080)

	// Database defaults
	l.viper.SetDefault("database.host", "localhost")
	l.viper.SetDefault("database.port", 5432)
	l.viper.SetDefault("database.ssl_mode", "disable")

	// JWT defaults
	l.viper.SetDefault("jwt.expiry_hour", 24)

	// Logger defaults
	l.viper.SetDefault("logger.level", "info")
	l.viper.SetDefault("logger.format", "json")
}
