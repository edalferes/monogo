package api

import (
	"fmt"

	"github.com/edalferes/monogo/internal/config"
	"github.com/edalferes/monogo/internal/infra/db"
	"github.com/edalferes/monogo/internal/modules/auth"
	"github.com/edalferes/monogo/internal/modules/testmodule"
	"github.com/edalferes/monogo/pkg/logger"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
)

type App struct {
	echo   *echo.Echo
	db     *gorm.DB
	logger logger.Logger
}

func NewApp() *App {
	cfg := config.LoadConfig()

	// Configure logger based on configuration
	loggerConfig := logger.DefaultConfig()
	loggerConfig.Level = cfg.Logger.Level
	loggerConfig.Format = cfg.Logger.Format
	loggerConfig.Service = cfg.App.Name

	appLogger := logger.New(loggerConfig)
	database, err := db.NewGormDB(cfg)
	if err != nil {
		appLogger.Fatal().Err(err).Msg("failed to connect to database")
	}

	var entities []interface{}
	entities = append(entities, auth.Entities()...)
	if err := database.AutoMigrate(entities...); err != nil {
		appLogger.Fatal().Err(err).Msg("failed to migrate database")
	}
	// Seed roles default
	if err := auth.Seed(database); err != nil {
		appLogger.Fatal().Err(err).Msg("failed to seed roles")
	}

	e := echo.New()
	return &App{
		echo:   e,
		db:     database,
		logger: appLogger,
	}
}

// RegisterModules register all modules
func (a *App) RegisterModules(cfg *config.Config) {
	v1 := a.echo.Group("/v1")

	// Auth module
	auth.WireUp(v1, a.db, cfg.JWT.Secret, a.logger)

	// Test module
	testmodule.WireUp(v1, cfg.JWT.Secret)
}

func (a *App) RegisterGlobalRoutes() {
	a.echo.GET("/health", func(c echo.Context) error {
		return c.String(200, "ok")
	})
	a.echo.GET("/metrics", func(c echo.Context) error {
		return c.String(200, "metrics: not implemented")
	})
	a.echo.GET("/swagger/*", echoSwagger.WrapHandler)
}

func (a *App) Run(cfg *config.Config) {
	a.RegisterGlobalRoutes()
	a.RegisterModules(cfg)
	a.logger.Info().Int("port", cfg.App.Port).Str("env", cfg.App.Environment).Msg("Starting API server")
	a.echo.Logger.Fatal(a.echo.Start(fmt.Sprintf(":%d", cfg.App.Port)))
}
