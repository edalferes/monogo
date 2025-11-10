package api

import (
	"fmt"

	"github.com/edalferes/monogo/internal/config"
	"github.com/edalferes/monogo/internal/infra/db"
	"github.com/edalferes/monogo/internal/infra/validator"
	"github.com/edalferes/monogo/internal/modules/auth"
	"github.com/edalferes/monogo/internal/modules/budget"
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
	entities = append(entities, budget.Entities()...)
	if err := database.AutoMigrate(entities...); err != nil {
		appLogger.Fatal().Err(err).Msg("failed to migrate database")
	}

	// Seed auth module (roles, permissions, root user)
	if err := auth.Seed(database); err != nil {
		appLogger.Fatal().Err(err).Msg("failed to seed auth module")
	}

	// Seed budget module with default categories for root user
	// Get root user ID
	var rootUser struct{ ID uint }
	if err := database.Table("users").Select("id").Where("username = ?", "root").First(&rootUser).Error; err != nil {
		appLogger.Fatal().Err(err).Msg("failed to find root user for budget seed")
	}
	if err := budget.Seed(database, rootUser.ID); err != nil {
		appLogger.Fatal().Err(err).Msg("failed to seed budget module")
	}

	e := echo.New()
	e.Validator = validator.NewValidator()

	return &App{
		echo:   e,
		db:     database,
		logger: appLogger,
	}
}

// RegisterModules register all available modules
func (a *App) RegisterModules(cfg *config.Config) {
	v1 := a.echo.Group("/v1")

	// Register all modules
	a.logger.Info().Str("module", "auth").Msg("Registering auth module")
	auth.WireUp(v1, a.db, cfg.JWT.Secret, a.logger)

	a.logger.Info().Str("module", "budget").Msg("Registering budget module")
	budget.WireUp(v1, a.db, cfg.JWT.Secret, a.logger)

	a.logger.Info().Str("module", "testmodule").Msg("Registering test module")
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
	a.logger.Info().
		Int("port", cfg.App.Port).
		Str("env", cfg.App.Environment).
		Msg("Starting API server")
	a.echo.Logger.Fatal(a.echo.Start(fmt.Sprintf(":%d", cfg.App.Port)))
}
