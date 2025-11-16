package auth

import (
	authmiddleware "github.com/edalferes/monetics/internal/modules/auth/adapters/http"
	"github.com/labstack/echo/v4"
)

// JWTMiddleware exposes the JWT middleware for use by other modules
func JWTMiddleware(secret string) echo.MiddlewareFunc {
	return authmiddleware.JWTMiddleware(secret)
}

// RequireRoles exposes the role-based middleware for use by other modules
func RequireRoles(roles ...string) echo.MiddlewareFunc {
	return authmiddleware.RequireRoles(roles...)
}

// RequirePermissions exposes the permission-based middleware for use by other modules
func RequirePermissions(perms ...string) echo.MiddlewareFunc {
	return authmiddleware.RequirePermissions(perms...)
}
