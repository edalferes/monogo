// Package cli provides command-line interface functionality for the modular application.
//
// This package implements a clean CLI architecture using Cobra that allows selective
// module execution, following the Single Responsibility Principle and maintaining
// clear separation of concerns.
package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/edalferes/monogo/internal/applications/api"
	"github.com/edalferes/monogo/internal/config"
	"github.com/spf13/cobra"
)

// ModuleRunner defines the contract for running the application with specific modules.
// This interface follows the Dependency Inversion Principle, allowing easy testing
// and different implementations.
type ModuleRunner interface {
	RunWithModules(modules []string, cfg *config.Config) error
}

// AppRunner implements ModuleRunner using the actual application.
type AppRunner struct{}

// RunWithModules starts the application with the specified modules.
func (r *AppRunner) RunWithModules(modules []string, cfg *config.Config) error {
	app := api.NewAppWithModules(modules)
	app.Run(cfg)
	return nil
}

// CLI encapsulates the command-line interface functionality.
type CLI struct {
	runner ModuleRunner
	config *config.Config
}

// NewCLI creates a new CLI instance with clean dependencies.
func NewCLI() *CLI {
	return &CLI{
		runner: &AppRunner{},
		config: config.LoadConfig(),
	}
}

// NewCLIWithRunner creates a CLI instance with a custom runner (useful for testing).
func NewCLIWithRunner(runner ModuleRunner, cfg *config.Config) *CLI {
	return &CLI{
		runner: runner,
		config: cfg,
	}
}

// Execute runs the CLI application.
func (c *CLI) Execute() error {
	return c.createRootCommand().Execute()
}

// createRootCommand builds the root cobra command with clean configuration.
func (c *CLI) createRootCommand() *cobra.Command {
	var modules []string

	rootCmd := &cobra.Command{
		Use:     "monogo",
		Short:   "Modular monolith application with selective module execution",
		Long:    c.getLongDescription(),
		Example: c.getExamples(),
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.runApplication(modules)
		},
	}

	// Configure flags with clean defaults
	rootCmd.Flags().StringSliceVarP(
		&modules,
		"modules",
		"m",
		[]string{"auth", "testmodule"},
		"Modules to run (auth, testmodule)",
	)

	// Add subcommands
	rootCmd.AddCommand(c.createListCommand())
	rootCmd.AddCommand(c.createHealthCommand())

	return rootCmd
}

// runApplication executes the main application logic.
func (c *CLI) runApplication(modules []string) error {
	validModules := c.validateAndCleanModules(modules)

	fmt.Printf("🚀 Starting Monogo with modules: %s\n",
		strings.Join(validModules, ", "))

	return c.runner.RunWithModules(validModules, c.config)
}

// validateAndCleanModules ensures only valid modules are processed.
func (c *CLI) validateAndCleanModules(modules []string) []string {
	validModules := []string{"auth", "testmodule"}
	var result []string

	for _, module := range modules {
		module = strings.TrimSpace(strings.ToLower(module))
		if c.isValidModule(module, validModules) {
			result = append(result, module)
		} else {
			fmt.Printf("⚠️  Unknown module '%s', skipping\n", module)
		}
	}

	if len(result) == 0 {
		fmt.Println("⚠️  No valid modules specified, using default: auth")
		return []string{"auth"}
	}

	return result
}

// isValidModule checks if a module is in the list of valid modules.
func (c *CLI) isValidModule(module string, validModules []string) bool {
	for _, valid := range validModules {
		if module == valid {
			return true
		}
	}
	return false
}

// createListCommand creates a command to list available modules.
func (c *CLI) createListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available modules",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("📦 Available modules:")
			fmt.Println("  • auth      - Authentication and authorization")
			fmt.Println("  • testmodule - Test and development module")
		},
	}
}

// createHealthCommand creates a health check command.
func (c *CLI) createHealthCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "health",
		Short: "Check application health",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("✅ Application is healthy")
			fmt.Printf("📍 Environment: %s\n", c.config.App.Environment)
			fmt.Printf("🔧 Version: %s\n", c.config.App.Name)
		},
	}
}

// getLongDescription returns the detailed description for the CLI.
func (c *CLI) getLongDescription() string {
	return `Monogo is a modular monolith application that supports selective module execution.

This allows you to run specific modules independently, enabling microservice-like 
deployment patterns while maintaining the benefits of a monolithic codebase.

Each module can be deployed as a separate service using the same binary and Docker image,
providing flexibility in scaling and deployment strategies.`
}

// getExamples returns usage examples for the CLI.
func (c *CLI) getExamples() string {
	return `  # Run with default auth module
  monogo

  # Run specific modules
  monogo --modules auth,testmodule
  monogo -m testmodule

  # List available modules
  monogo list

  # Check application health
  monogo health`
}

// ExecuteCLI is the main entry point for the CLI application.
// This function provides a clean interface for main.go.
func ExecuteCLI() {
	cli := NewCLI()
	if err := cli.Execute(); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}
}
