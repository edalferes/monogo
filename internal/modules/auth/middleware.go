package auth

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// JWTMiddleware autentica o token JWT e popula as claims no contexto Echo
func JWTMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid claims"})
			}
			c.Set("user", claims)
			return next(c)
		}
	}
}

func RequireRoles(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, ok := c.Get("user").(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden: no user claims"})
			}
			userRoles, ok := claims["roles"].([]interface{})
			if !ok {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden: no roles"})
			}
			for _, required := range roles {
				for _, ur := range userRoles {
					if roleStr, ok := ur.(string); ok && roleStr == required {
						return next(c)
					}
				}
			}
			return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden: insufficient role"})
		}
	}
}

func RequirePermissions(perms ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, ok := c.Get("user").(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden: no user claims"})
			}
			userPerms, ok := claims["permissions"].([]interface{})
			if !ok {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden: no permissions"})
			}
			for _, required := range perms {
				for _, up := range userPerms {
					if permStr, ok := up.(string); ok && permStr == required {
						return next(c)
					}
				}
			}
			return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden: insufficient permission"})
		}
	}
}
