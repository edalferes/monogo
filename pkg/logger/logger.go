package logger

// New creates a new logger with configuration
func New(config Config) Logger {
	return NewZerologLogger(config)
}

// NewDefault creates a new logger with default configuration
func NewDefault() Logger {
	return NewZerologLogger(DefaultConfig())
}
