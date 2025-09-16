package main

import "github.com/edalferes/monogo/internal/applications/api"

func main() {
	app := api.NewApp()
	app.RegisterModules()
	app.Run()
}
