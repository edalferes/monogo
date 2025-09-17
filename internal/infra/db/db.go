package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/edalferes/monogo/config"

	_ "github.com/lib/pq"
)

// NewGormDB Initialize a new Gorm DB connection
func NewGormDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
