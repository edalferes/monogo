package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	Port       int
}

func LoadConfig() *Config {
	viper.AutomaticEnv()

	viper.SetDefault("PORT", 8080)

	cfg := &Config{
		DBHost:     viper.GetString("DB_HOST"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
		Port:       viper.GetInt("PORT"),
	}

	if cfg.DBHost == "" || cfg.DBUser == "" || cfg.DBName == "" {
		log.Fatal("Missing required environment variables: DB_HOST, DB_USER, DB_NAME")
	}

	return cfg
}
