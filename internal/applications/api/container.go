package api

import (
	"sync"

	"github.com/edalferes/monetics/internal/config"
	"github.com/edalferes/monetics/pkg/logger"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// ModuleContainer manages module dependencies and services
type ModuleContainer struct {
	mu       sync.RWMutex
	services map[string]interface{}
	config   *config.Config
	db       *gorm.DB
	logger   logger.Logger
	echo     *echo.Group
	modules  []string
}

// NewModuleContainer creates a new dependency injection container
func NewModuleContainer(cfg *config.Config, db *gorm.DB, log logger.Logger, group *echo.Group, modules []string) *ModuleContainer {
	return &ModuleContainer{
		services: make(map[string]interface{}),
		config:   cfg,
		db:       db,
		logger:   log,
		echo:     group,
		modules:  modules,
	}
}

// Register registers a service in the container
func (c *ModuleContainer) Register(name string, service interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.services[name] = service
	c.logger.Debug().Str("service", name).Msg("Service registered in container")
}

// Get retrieves a service from the container
func (c *ModuleContainer) Get(name string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	service, exists := c.services[name]
	return service, exists
}

// Has checks if a service exists in the container
func (c *ModuleContainer) Has(name string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.services[name]
	return exists
}

// IsModuleEnabled checks if a module is enabled in this instance
func (c *ModuleContainer) IsModuleEnabled(moduleName string) bool {
	for _, m := range c.modules {
		if m == moduleName {
			return true
		}
	}
	return false
}

// GetConfig returns the application configuration
func (c *ModuleContainer) GetConfig() *config.Config {
	return c.config
}

// GetDB returns the database connection
func (c *ModuleContainer) GetDB() *gorm.DB {
	return c.db
}

// GetLogger returns the logger
func (c *ModuleContainer) GetLogger() logger.Logger {
	return c.logger
}

// GetEchoGroup returns the Echo group for route registration
func (c *ModuleContainer) GetEchoGroup() *echo.Group {
	return c.echo
}

// ServiceFactory defines a function that creates a service
type ServiceFactory func(container *ModuleContainer) interface{}

// ModuleInitializer defines a function that initializes a module
type ModuleInitializer func(container *ModuleContainer) error
