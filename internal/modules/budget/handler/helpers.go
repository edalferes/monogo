package handler

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// GetUserIDFromContext extracts user ID from JWT claims in context
func GetUserIDFromContext(c echo.Context) (uint, error) {
	claims, ok := c.Get("user").(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid user claims")
	}

	uidFloat, ok := claims["uid"].(float64)
	if !ok {
		return 0, fmt.Errorf("user ID not found in claims")
	}

	return uint(uidFloat), nil
}
