package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct{}

// Test godoc
// @Summary Endpoint protegido para teste de roles/permissões
// @Tags testmodule
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /v1/test/protected [get]
func (h *Handler) Protected(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "acesso autorizado!"})
}
