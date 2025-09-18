package api

import (
	"github.com/edalferes/monogo/config"
	"github.com/edalferes/monogo/internal/infra/db"
	"github.com/edalferes/monogo/internal/infra/logger"
	"github.com/edalferes/monogo/internal/modules/auth"
	"github.com/edalferes/monogo/internal/modules/testmodule"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
)

type App struct {
	echo *echo.Echo
	db   *gorm.DB
}

func NewApp() *App {
	logger.Init()
	cfg := config.LoadConfig()
	database, err := db.NewGormDB(cfg)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to connect to database")
	}

	var entities []interface{}
	entities = append(entities, auth.Entities()...)
	if err := database.AutoMigrate(entities...); err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to migrate database")
	}
	// Seed roles default
	if err := auth.Seed(database); err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to seed roles")
	}
	e := echo.New()
	return &App{
		echo: e,
		db:   database,
	}
}

// RegisterModules register all modules
func (a *App) RegisterModules() {
	v1 := a.echo.Group("/v1")

	// Auth module
	cfg := config.LoadConfig()
	auth.WireUp(v1, a.db, cfg.JWTSecret)

	// Test module
	testmodule.WireUp(v1, cfg.JWTSecret)
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

func (a *App) Run() {
	a.RegisterGlobalRoutes()
	a.RegisterModules()
	logger.Log.Info().Msg("API running on :8080")
	a.echo.Logger.Fatal(a.echo.Start(":8080"))
}
