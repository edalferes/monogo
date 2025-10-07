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
	echo    *echo.Echo
	db      *gorm.DB
	logger  logger.Logger
	modules []string // modules to load
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
		echo:    e,
		db:      database,
		logger:  appLogger,
		modules: []string{"auth", "testmodule"}, // default all modules
	}
}

// NewAppWithModules creates an app instance with specific modules
func NewAppWithModules(modules []string) *App {
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

	// Only migrate entities for enabled modules
	var entities []interface{}
	for _, module := range modules {
		switch module {
		case "auth":
			entities = append(entities, auth.Entities()...)
		}
	}

	if len(entities) > 0 {
		if err := database.AutoMigrate(entities...); err != nil {
			appLogger.Fatal().Err(err).Msg("failed to migrate database")
		}
	}

	// Seed only if auth module is enabled
	for _, module := range modules {
		if module == "auth" {
			if err := auth.Seed(database); err != nil {
				appLogger.Fatal().Err(err).Msg("failed to seed roles")
			}
			break
		}
	}

	e := echo.New()
	return &App{
		echo:    e,
		db:      database,
		logger:  appLogger,
		modules: modules,
	}
}

// RegisterModules register enabled modules
func (a *App) RegisterModules(cfg *config.Config) {
	v1 := a.echo.Group("/v1")

	// Register only enabled modules
	for _, module := range a.modules {
		switch module {
		case "auth":
			a.logger.Info().Str("module", "auth").Msg("Registering auth module")
			auth.WireUp(v1, a.db, cfg.JWT.Secret, a.logger)
		case "testmodule":
			a.logger.Info().Str("module", "testmodule").Msg("Registering test module")
			testmodule.WireUp(v1, cfg.JWT.Secret)
		default:
			a.logger.Warn().Str("module", module).Msg("Unknown module, skipping")
		}
	}
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
		Str("modules", fmt.Sprintf("%v", a.modules)).
		Msg("Starting API server")
	a.echo.Logger.Fatal(a.echo.Start(fmt.Sprintf(":%d", cfg.App.Port)))
}
