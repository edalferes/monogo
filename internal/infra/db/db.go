package db

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/edalferes/monogo/internal/config"

	_ "github.com/lib/pq"
)

// NewGormDB Initialize a new Gorm DB connection
func NewGormDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.GetDSN()

	// Configure GORM logger to ignore "record not found" errors
	logConfig := logger.Config{
		SlowThreshold:             time.Second,
		LogLevel:                  logger.Warn, // Only logs warnings and errors
		IgnoreRecordNotFoundError: true,        // Ignores "record not found" logs
		Colorful:                  true,
	}

	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logConfig),
	})
}
