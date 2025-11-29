package api

import (
	"fmt"

	"github.com/edalferes/monetics/internal/modules/auth"
	"github.com/edalferes/monetics/internal/modules/budget"
	"github.com/edalferes/monetics/internal/modules/testmodule"
	"github.com/edalferes/monetics/pkg/logger"
)

// ModuleRegistry manages module registration and initialization order
type ModuleRegistry struct {
	container    *ModuleContainer
	initializers map[string]ModuleInitializer
	dependencies map[string][]string
	initialized  map[string]bool
	logger       logger.Logger
}

// NewModuleRegistry creates a new module registry
func NewModuleRegistry(container *ModuleContainer) *ModuleRegistry {
	return &ModuleRegistry{
		container:    container,
		initializers: make(map[string]ModuleInitializer),
		dependencies: make(map[string][]string),
		initialized:  make(map[string]bool),
		logger:       container.GetLogger(),
	}
}

// Register registers a module with its initializer and dependencies
func (r *ModuleRegistry) Register(name string, deps []string, initializer ModuleInitializer) {
	r.initializers[name] = initializer
	r.dependencies[name] = deps
	r.logger.Debug().
		Str("module", name).
		Int("dependencies_count", len(deps)).
		Msg("Module registered")
}

// Initialize initializes all registered modules in dependency order
func (r *ModuleRegistry) Initialize() error {
	// Get enabled modules
	enabledModules := r.container.modules

	// Initialize modules in dependency order
	for _, moduleName := range enabledModules {
		if err := r.initializeModule(moduleName); err != nil {
			return fmt.Errorf("failed to initialize module %s: %w", moduleName, err)
		}
	}

	return nil
}

// initializeModule initializes a single module and its dependencies recursively
func (r *ModuleRegistry) initializeModule(name string) error {
	// Check if already initialized
	if r.initialized[name] {
		return nil
	}

	// Check if module is registered
	initializer, exists := r.initializers[name]
	if !exists {
		return fmt.Errorf("module %s not registered", name)
	}

	// Initialize dependencies first
	if deps, hasDeps := r.dependencies[name]; hasDeps {
		for _, dep := range deps {
			// Check if dependency module is enabled
			if !r.container.IsModuleEnabled(dep) {
				r.logger.Warn().
					Str("module", name).
					Str("dependency", dep).
					Msg("Dependency module not enabled, will use remote service if configured")
				continue
			}

			if err := r.initializeModule(dep); err != nil {
				return fmt.Errorf("failed to initialize dependency %s for module %s: %w", dep, name, err)
			}
		}
	}

	// Initialize the module
	r.logger.Info().Str("module", name).Msg("Initializing module")
	if err := initializer(r.container); err != nil {
		return err
	}

	r.initialized[name] = true
	r.logger.Info().Str("module", name).Msg("Module initialized successfully")
	return nil
}

// RegisterBuiltInModules registers all built-in modules
func (r *ModuleRegistry) RegisterBuiltInModules() {
	// Auth module - no dependencies
	r.Register("auth", []string{}, func(c *ModuleContainer) error {
		cfg := c.GetConfig()
		db := c.GetDB()
		log := c.GetLogger()
		group := c.GetEchoGroup()

		// Initialize auth and register UserService in container
		userService := auth.WireUpWithService(group, db, cfg.JWT.Secret, log)
		c.Register("auth.UserService", userService)

		return nil
	})

	// Budget module - depends on Auth
	r.Register("budget", []string{"auth"}, func(c *ModuleContainer) error {
		cfg := c.GetConfig()
		db := c.GetDB()
		log := c.GetLogger()
		group := c.GetEchoGroup()

		// Check if Auth is running locally
		if c.IsModuleEnabled("auth") {
			log.Info().Msg("Using local Auth service for Budget module")
			budget.WireUp(group, db, cfg.JWT.Secret, log)
		} else if cfg.Modules.Auth.URL != "" {
			log.Info().
				Str("auth_url", cfg.Modules.Auth.URL).
				Msg("Using remote Auth service for Budget module")
			budget.WireUpWithHTTP(group, db, cfg.JWT.Secret, log, cfg.Modules.Auth.URL)
		} else {
			return fmt.Errorf("budget module requires auth service: enable auth module or configure auth URL")
		}

		return nil
	})

	// Test module - no dependencies
	r.Register("test", []string{}, func(c *ModuleContainer) error {
		cfg := c.GetConfig()
		group := c.GetEchoGroup()

		testmodule.WireUp(group, cfg.JWT.Secret)
		return nil
	})
}
