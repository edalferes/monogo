package user

import (
	"net/http"

	"github.com/edalferes/monogo/internal/modules/auth/adapters/http/dto"
	"github.com/edalferes/monogo/internal/modules/auth/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	ChangePasswordUC interface {
		Execute(userID uint, currentPassword, newPassword string) error
	}
}

// ChangePassword godoc
// @Summary Change user password
// @Description Allows authenticated user to change their own password
// @Tags user
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param password body dto.ChangePasswordDTO true "Password data"
// @Success 200 {object} map[string]string "password changed"
// @Failure 400 {object} map[string]string "invalid data"
// @Failure 401 {object} map[string]string "invalid current password"
// @Failure 500 {object} map[string]string "internal error"
// @Router /v1/user/password [put]
func (h *UserHandler) ChangePassword(c echo.Context) error {
	// Get user ID from JWT token
	claims, ok := c.Get("user").(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
	}

	userID, ok := claims["uid"].(float64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token claims"})
	}

	var input dto.ChangePasswordDTO
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errors.ErrInvalidData.Error()})
	}

	if input.CurrentPassword == "" || input.NewPassword == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "current_password and new_password are required"})
	}

	if err := h.ChangePasswordUC.Execute(uint(userID), input.CurrentPassword, input.NewPassword); err != nil {
		if err == errors.ErrInvalidCredentials {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid current password"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "password changed successfully"})
}
