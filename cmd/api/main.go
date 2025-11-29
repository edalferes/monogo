package main

import (
	"flag"

	_ "github.com/edalferes/monetics/docs/openapi"
	"github.com/edalferes/monetics/internal/applications/api"
	"github.com/edalferes/monetics/internal/config"
)

func main() {
	// Parse command-line flags
	moduleFlag := flag.String("module", "all", "Modules to run (comma-separated): auth, budget, or all")
	flag.Parse()

	// Load configuration
	cfg := config.LoadConfig()

	// Create and run application with selected modules
	app := api.NewApp(*moduleFlag, cfg)
	app.Run()
}
