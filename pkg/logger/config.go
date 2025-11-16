package logger

import "io"

// Config defines the logger configuration
type Config struct {
	Level      string // "debug", "info", "warn", "error", "fatal"
	Output     io.Writer
	OutputFile string // Se especificado, escreve em arquivo
	Format     string // "json", "console"
	Service    string // Service name for identification
}

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	return Config{
		Level:   "info",
		Format:  "json",
		Service: "monetics",
	}
}
