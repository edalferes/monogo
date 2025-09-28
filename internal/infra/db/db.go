package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/edalferes/monogo/internal/config"

	_ "github.com/lib/pq"
)

// NewGormDB Initialize a new Gorm DB connection
func NewGormDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.GetDSN()
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
