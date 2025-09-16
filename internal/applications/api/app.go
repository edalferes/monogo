package api

import (
	"database/sql"
	"net/http"

	"github.com/edalferes/monogo/config"
	"github.com/edalferes/monogo/internal/infra/db"
	"github.com/edalferes/monogo/internal/infra/logger"
	"github.com/edalferes/monogo/internal/modules/user"
)

type App struct {
	mux *http.ServeMux
	db  *sql.DB
}

func NewApp() *App {
	logger.Init()
	cfg := config.LoadConfig()
	database, err := db.NewPostgresDB(cfg)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to connect to database")
	}
	return &App{
		mux: http.NewServeMux(),
		db:  database,
	}
}

func (a *App) RegisterModules() {
	user.WireUp(a.mux, a.db)
	// auth.WireUp(a.mux, a.db)
	// billing.WireUp(a.mux, a.db)
	RegisterGlobalRoutes(a.mux)
}

func (a *App) Run() {
	logger.Log.Info().Msg("API running on :8080")
	err := http.ListenAndServe(":8080", a.mux)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("server exited with error")
	}
}
