package logger

// New cria um novo logger com configuração
func New(config Config) Logger {
	return NewZerologLogger(config)
}

// NewDefault cria um novo logger com configuração padrão
func NewDefault() Logger {
	return NewZerologLogger(DefaultConfig())
}
