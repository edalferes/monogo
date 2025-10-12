package main

import (
	_ "github.com/edalferes/monogo/docs"
	"github.com/edalferes/monogo/internal/applications/api"
	"github.com/edalferes/monogo/internal/config"
)

func main() {
	cfg := config.LoadConfig()
	app := api.NewApp()
	app.Run(cfg)
}
