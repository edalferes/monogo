package testmodule

import (
	"github.com/edalferes/monogo/internal/modules/auth"
	"github.com/edalferes/monogo/internal/modules/testmodule/handler"
	"github.com/labstack/echo/v4"
)

// WireUp register the testmodule routes, protecting with JWT and role check
func WireUp(group *echo.Group, jwtSecret string) {
	h := &handler.Handler{}
	protected := group.Group("/test")
	protected.Use(auth.JWTMiddleware(jwtSecret))
	protected.GET("/protected", h.Protected, auth.RequireRoles("user"))
}
