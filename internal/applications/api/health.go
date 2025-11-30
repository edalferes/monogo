package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status  string            `json:"status"`
	Service string            `json:"service"`
	Details map[string]string `json:"details,omitempty"`
}

// HealthHandler handles health check requests
// Similar to Loki's /ready and /health endpoints
func (app *App) HealthHandler(c echo.Context) error {
	response := HealthResponse{
		Status:  "healthy",
		Service: "monetics",
		Details: map[string]string{
			"database": "connected",
		},
	}

	return c.JSON(http.StatusOK, response)
}

// ReadyHandler handles readiness check requests
// Returns 200 OK when the service is ready to accept traffic
func (app *App) ReadyHandler(c echo.Context) error {
	// Check if database is connected
	sqlDB, err := app.db.DB()
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"status": "not ready",
			"reason": "database connection failed",
		})
	}

	if err := sqlDB.Ping(); err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"status": "not ready",
			"reason": "database ping failed",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "ready",
	})
}

// LivenessHandler handles liveness probe requests
// Returns 200 OK as long as the process is running
func (app *App) LivenessHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "alive",
	})
}
