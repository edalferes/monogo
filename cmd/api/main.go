package main

import (
	_ "github.com/edalferes/monetics/docs/openapi"
	"github.com/edalferes/monetics/internal/applications/api"
	"github.com/edalferes/monetics/internal/config"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Create and run application
	app := api.NewApp(cfg)
	app.Run()
}
