package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// LivenessHandler - K8s liveness probe
// Returns 200 OK if the process is running
func (app *App) LivenessHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "alive",
	})
}

// ReadinessHandler - K8s readiness probe
// Returns 200 OK if the service can accept traffic
func (app *App) ReadinessHandler(c echo.Context) error {
	sqlDB, err := app.db.DB()
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"status": "not ready",
		})
	}

	if err := sqlDB.Ping(); err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"status": "not ready",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "ready",
	})
}
