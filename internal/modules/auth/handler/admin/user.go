package admin

import (
	"net/http"

	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"github.com/edalferes/monogo/internal/modules/auth/errors"
	"github.com/edalferes/monogo/internal/modules/auth/repository"
	"github.com/edalferes/monogo/internal/modules/auth/service"
	"github.com/labstack/echo/v4"
)

type AdminUserHandler struct {
	UserRepo        repository.UserRepository
	RoleRepo        repository.RoleRepository
	PasswordService service.PasswordService
}

// CreateUser godoc
// @Summary Cria um novo usuário
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user body RegisterDTO true "Dados do usuário"
// @Success 201 {object} map[string]string "created"
// @Failure 400 {object} map[string]string "invalid data or user exists"
// @Failure 500 {object} map[string]string "internal error"
// @Router /v1/admin/users [post]
func (h *AdminUserHandler) CreateUser(c echo.Context) error {
	var input RegisterDTO
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errors.ErrInvalidData.Error()})
	}
	if input.Username == "" || input.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errors.ErrMissingCredentials.Error()})
	}
	if user, _ := h.UserRepo.FindByUsername(input.Username); user != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errors.ErrUserAlreadyExists.Error()})
	}
	hash, err := h.PasswordService.Hash(input.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}
	role, err := h.RoleRepo.FindByName("user")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}
	user := &domain.User{
		Username: input.Username,
		Password: hash,
		Roles:    []domain.Role{*role},
	}
	if err := h.UserRepo.Create(user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}
	return c.JSON(http.StatusCreated, map[string]string{"message": "user created"})
}

type RegisterDTO struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
