package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// ZerologLogger implementa a interface Logger usando zerolog
type ZerologLogger struct {
	logger zerolog.Logger
}

// NewZerologLogger cria uma nova instância do logger com a configuração especificada
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

	// Configurar nível de log
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

// Info implementa Logger.Info
func (z *ZerologLogger) Info() Event {
	return &ZerologEvent{event: z.logger.Info()}
}

// Error implementa Logger.Error
func (z *ZerologLogger) Error() Event {
	return &ZerologEvent{event: z.logger.Error()}
}

// Debug implementa Logger.Debug
func (z *ZerologLogger) Debug() Event {
	return &ZerologEvent{event: z.logger.Debug()}
}

// Warn implementa Logger.Warn
func (z *ZerologLogger) Warn() Event {
	return &ZerologEvent{event: z.logger.Warn()}
}

// Fatal implementa Logger.Fatal
func (z *ZerologLogger) Fatal() Event {
	return &ZerologEvent{event: z.logger.Fatal()}
}

// With implementa Logger.With
func (z *ZerologLogger) With() Context {
	return &ZerologContext{context: z.logger.With()}
}

// ZerologEvent implementa a interface Event
type ZerologEvent struct {
	event *zerolog.Event
}

// Msg implementa Event.Msg
func (e *ZerologEvent) Msg(msg string) {
	e.event.Msg(msg)
}

// Err implementa Event.Err
func (e *ZerologEvent) Err(err error) Event {
	return &ZerologEvent{event: e.event.Err(err)}
}

// Str implementa Event.Str
func (e *ZerologEvent) Str(key, val string) Event {
	return &ZerologEvent{event: e.event.Str(key, val)}
}

// Int implementa Event.Int
func (e *ZerologEvent) Int(key string, val int) Event {
	return &ZerologEvent{event: e.event.Int(key, val)}
}

// Uint implementa Event.Uint
func (e *ZerologEvent) Uint(key string, val uint) Event {
	return &ZerologEvent{event: e.event.Uint(key, val)}
}

// Bool implementa Event.Bool
func (e *ZerologEvent) Bool(key string, val bool) Event {
	return &ZerologEvent{event: e.event.Bool(key, val)}
}

// ZerologContext implementa a interface Context
type ZerologContext struct {
	context zerolog.Context
}

// Str implementa Context.Str
func (c *ZerologContext) Str(key, val string) Context {
	return &ZerologContext{context: c.context.Str(key, val)}
}

// Int implementa Context.Int
func (c *ZerologContext) Int(key string, val int) Context {
	return &ZerologContext{context: c.context.Int(key, val)}
}

// Uint implementa Context.Uint
func (c *ZerologContext) Uint(key string, val uint) Context {
	return &ZerologContext{context: c.context.Uint(key, val)}
}

// Logger implementa Context.Logger
func (c *ZerologContext) Logger() Logger {
	return &ZerologLogger{logger: c.context.Logger()}
}
