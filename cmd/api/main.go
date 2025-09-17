package main

import (
	_ "github.com/edalferes/monogo/docs"
	"github.com/edalferes/monogo/internal/applications/api"
)

func main() {
	app := api.NewApp()
	app.RegisterModules()
	app.Run()
}
