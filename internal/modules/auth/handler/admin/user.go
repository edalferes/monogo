package admin

import (
	"net/http"

	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"github.com/edalferes/monogo/internal/modules/auth/errors"
	"github.com/edalferes/monogo/internal/modules/auth/handler"
	"github.com/edalferes/monogo/internal/modules/auth/repository"
	"github.com/edalferes/monogo/internal/modules/auth/service"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AdminUserHandler struct {
	UserRepo        repository.UserRepository
	RoleRepo        repository.RoleRepository
	PasswordService service.PasswordService
}

// ListUsers godoc
// @Summary List all users
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Success 200 {array} domain.User
// @Router /v1/admin/users [get]
func (h *AdminUserHandler) ListUsers(c echo.Context) error {
	users, err := h.UserRepo.ListAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, users)
}

// CreateUser godoc
// @Summary Create a new user
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user body handler.RegisterDTO true "User data"
// @Success 201 {object} map[string]string "created"
// @Failure 400 {object} map[string]string "invalid data or user exists"
// @Failure 500 {object} map[string]string "internal error"
// @Router /v1/admin/users [post]
func (h *AdminUserHandler) CreateUser(c echo.Context) error {
	var input handler.RegisterDTO
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errors.ErrInvalidData.Error()})
	}
	if err := validator.New().Struct(input); err != nil {
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
