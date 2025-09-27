package logger

import "io"

// Config define as configurações do logger
type Config struct {
	Level      string // "debug", "info", "warn", "error", "fatal"
	Output     io.Writer
	OutputFile string // Se especificado, escreve em arquivo
	Format     string // "json", "console"
	Service    string // Nome do serviço para identificação
}

// DefaultConfig retorna uma configuração padrão
func DefaultConfig() Config {
	return Config{
		Level:   "info",
		Format:  "json",
		Service: "monogo",
	}
}
