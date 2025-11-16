package main

import (
	_ "github.com/edalferes/monetics/docs"
	"github.com/edalferes/monetics/internal/applications/api"
	"github.com/edalferes/monetics/internal/config"
)

func main() {
	cfg := config.LoadConfig()
	app := api.NewApp()
	app.Run(cfg)
}
