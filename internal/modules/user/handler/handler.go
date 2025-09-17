package handler

import (
	"github.com/edalferes/monogo/internal/modules/user/service"
	"github.com/edalferes/monogo/pkg/responses"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Service *service.Service
}

// RegisterRoutesEcho registra as rotas do módulo user usando Echo
func (h *Handler) RegisterRoutesEcho(g *echo.Group) {
	g.POST("/users", h.CreateUserEcho)
}

// CreateUserEcho godoc
// @Summary Cria um novo usuário
// @Description Cria um usuário com nome e email
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserDTO true "Dados do usuário"
// @Success 201 {object} responses.CreatedResponse "created"
// @Failure 400 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /v1/users [post]
func (h *Handler) CreateUserEcho(c echo.Context) error {
	var input CreateUserDTO
	if err := c.Bind(&input); err != nil {
		return responses.BadRequest(c, err)
	}
	user, err := h.Service.Register(input.Name, input.Email)
	if err != nil {
		return responses.InternalServerError(c, err)
	}
	return responses.Created(c, user.ID.String())
}
