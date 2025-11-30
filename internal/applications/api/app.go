package api

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"

	"github.com/edalferes/monetics/internal/config"
	"github.com/edalferes/monetics/internal/infra/db"
	"github.com/edalferes/monetics/internal/infra/validator"
	"github.com/edalferes/monetics/internal/modules/auth"
	"github.com/edalferes/monetics/internal/modules/budget"
	"github.com/edalferes/monetics/pkg/logger"
)

type App struct {
	echo    *echo.Echo
	db      *gorm.DB
	logger  logger.Logger
	modules []string // List of enabled modules
	config  *config.Config
}

func NewApp(enabledModules string, cfg *config.Config) *App {
	// Parse enabled modules
	modules := parseModules(enabledModules)

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

	// Migrate entities from enabled modules
	var entities []interface{}
	if hasModule(modules, "auth") {
		entities = append(entities, auth.Entities()...)
	}
	if hasModule(modules, "budget") {
		entities = append(entities, budget.Entities()...)
	}

	if len(entities) > 0 {
		if err := database.AutoMigrate(entities...); err != nil {
			appLogger.Fatal().Err(err).Msg("failed to migrate database")
		}
	}

	// Seed enabled modules
	if hasModule(modules, "auth") {
		if err := auth.Seed(database); err != nil {
			appLogger.Fatal().Err(err).Msg("failed to seed auth module")
		}
	}

	if hasModule(modules, "budget") {
		var rootUser struct{ ID uint }
		if err := database.Table("users").Select("id").Where("username = ?", "root").First(&rootUser).Error; err != nil {
			appLogger.Fatal().Err(err).Msg("failed to find root user for budget seed")
		}
		if err := budget.Seed(database, rootUser.ID); err != nil {
			appLogger.Fatal().Err(err).Msg("failed to seed budget module")
		}
	}

	e := echo.New()
	e.Validator = validator.NewValidator()

	// Configure CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001", "http://localhost:5173"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	return &App{
		echo:    e,
		db:      database,
		logger:  appLogger,
		modules: modules,
		config:  cfg,
	}
}

// parseModules converts module string to list
func parseModules(input string) []string {
	if input == "" || input == "all" {
		return []string{"auth", "budget"}
	}

	parts := strings.Split(input, ",")
	modules := make([]string, 0, len(parts))
	for _, p := range parts {
		m := strings.TrimSpace(p)
		if m != "" {
			modules = append(modules, m)
		}
	}
	return modules
}

// hasModule checks if module is enabled
func hasModule(modules []string, name string) bool {
	for _, m := range modules {
		if m == name {
			return true
		}
	}
	return false
}

// RegisterModules registers only enabled modules using dependency injection
func (a *App) RegisterModules(cfg *config.Config) {
	v1 := a.echo.Group("/v1")

	// Create dependency injection container
	container := NewModuleContainer(cfg, a.db, a.logger, v1, a.modules)

	// Create module registry
	registry := NewModuleRegistry(container)

	// Register all built-in modules with their dependencies
	registry.RegisterBuiltInModules()

	// Initialize modules in dependency order
	if err := registry.Initialize(); err != nil {
		a.logger.Fatal().Err(err).Msg("Failed to initialize modules")
	}
}

func (a *App) RegisterGlobalRoutes() {
	// Health and readiness probes (Kubernetes/Loki style)
	a.echo.GET("/health", a.HealthHandler)
	a.echo.GET("/ready", a.ReadyHandler)
	a.echo.GET("/live", a.LivenessHandler)

	// Legacy health endpoint for backward compatibility
	a.echo.GET("/healthz", func(c echo.Context) error {
		return c.String(200, "ok")
	})

	a.echo.GET("/metrics", func(c echo.Context) error {
		return c.String(200, "metrics: not implemented")
	})
	a.echo.GET("/swagger/*", echoSwagger.WrapHandler)
}

func (a *App) Run() {
	a.RegisterGlobalRoutes()
	a.RegisterModules(a.config)
	a.logger.Info().
		Int("port", a.config.App.Port).
		Str("env", a.config.App.Environment).
		Msg("Starting API server")
	a.echo.Logger.Fatal(a.echo.Start(fmt.Sprintf(":%d", a.config.App.Port)))
}
