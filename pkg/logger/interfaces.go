package logger

// Logger define a interface para logging no sistema
type Logger interface {
	Info() Event
	Error() Event
	Debug() Event
	Warn() Event
	Fatal() Event
	With() Context
}

// Event representa um evento de log em construção
type Event interface {
	Msg(msg string)
	Err(err error) Event
	Str(key, val string) Event
	Int(key string, val int) Event
	Uint(key string, val uint) Event
	Bool(key string, val bool) Event
}

// Context permite adicionar campos de contexto ao logger
type Context interface {
	Str(key, val string) Context
	Int(key string, val int) Context
	Uint(key string, val uint) Context
	Logger() Logger
}
