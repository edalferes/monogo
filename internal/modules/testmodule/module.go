package testmodule

import (
	"github.com/edalferes/monogo/internal/modules/auth"
	"github.com/edalferes/monogo/internal/modules/testmodule/handler"
	"github.com/labstack/echo/v4"
)

func WireUp(group *echo.Group) {
	h := &handler.Handler{}
	group.GET("/test/protected", h.Protected, auth.RequireRoles("user"))
}
