package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// ZerologLogger implements the Logger interface using zerolog
type ZerologLogger struct {
	logger zerolog.Logger
}

// NewZerologLogger creates a new logger instance with the specified configuration
func NewZerologLogger(config Config) Logger {
	var output = config.Output
	if output == nil {
		if config.OutputFile != "" {
			file, err := os.OpenFile(config.OutputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				output = os.Stdout // fallback
			} else {
				output = file
			}
		} else {
			output = os.Stdout
		}
	}

	logger := zerolog.New(output).With().Timestamp()

	// Adicionar service name se especificado
	if config.Service != "" {
		logger = logger.Str("service", config.Service)
	}

	zeroLogger := logger.Logger()

	// Configure log level
	if config.Level != "" {
		if level, err := zerolog.ParseLevel(config.Level); err == nil {
			zeroLogger = zeroLogger.Level(level)
		}
	}

	// Configurar formato se for console
	if config.Format == "console" {
		output = zerolog.ConsoleWriter{Out: output}
		zeroLogger = zerolog.New(output).With().Timestamp().Logger()
	}

	return &ZerologLogger{logger: zeroLogger}
}

// Info implements Logger.Info
func (z *ZerologLogger) Info() Event {
	return &ZerologEvent{event: z.logger.Info()}
}

// Error implements Logger.Error
func (z *ZerologLogger) Error() Event {
	return &ZerologEvent{event: z.logger.Error()}
}

// Debug implements Logger.Debug
func (z *ZerologLogger) Debug() Event {
	return &ZerologEvent{event: z.logger.Debug()}
}

// Warn implements Logger.Warn
func (z *ZerologLogger) Warn() Event {
	return &ZerologEvent{event: z.logger.Warn()}
}

// Fatal implements Logger.Fatal
func (z *ZerologLogger) Fatal() Event {
	return &ZerologEvent{event: z.logger.Fatal()}
}

// With implements Logger.With
func (z *ZerologLogger) With() Context {
	return &ZerologContext{context: z.logger.With()}
}

// ZerologEvent implements the Event interface
type ZerologEvent struct {
	event *zerolog.Event
}

// Msg implements Event.Msg
func (e *ZerologEvent) Msg(msg string) {
	e.event.Msg(msg)
}

// Err implements Event.Err
func (e *ZerologEvent) Err(err error) Event {
	return &ZerologEvent{event: e.event.Err(err)}
}

// Str implements Event.Str
func (e *ZerologEvent) Str(key, val string) Event {
	return &ZerologEvent{event: e.event.Str(key, val)}
}

// Int implements Event.Int
func (e *ZerologEvent) Int(key string, val int) Event {
	return &ZerologEvent{event: e.event.Int(key, val)}
}

// Uint implements Event.Uint
func (e *ZerologEvent) Uint(key string, val uint) Event {
	return &ZerologEvent{event: e.event.Uint(key, val)}
}

// Bool implements Event.Bool
func (e *ZerologEvent) Bool(key string, val bool) Event {
	return &ZerologEvent{event: e.event.Bool(key, val)}
}

// ZerologContext implements the Context interface
type ZerologContext struct {
	context zerolog.Context
}

// Str implements Context.Str
func (c *ZerologContext) Str(key, val string) Context {
	return &ZerologContext{context: c.context.Str(key, val)}
}

// Int implements Context.Int
func (c *ZerologContext) Int(key string, val int) Context {
	return &ZerologContext{context: c.context.Int(key, val)}
}

// Uint implements Context.Uint
func (c *ZerologContext) Uint(key string, val uint) Context {
	return &ZerologContext{context: c.context.Uint(key, val)}
}

// Logger implements Context.Logger
func (c *ZerologContext) Logger() Logger {
	return &ZerologLogger{logger: c.context.Logger()}
}
