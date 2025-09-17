package handler

import (
	"net/http"

	"github.com/edalferes/monogo/internal/modules/auth/errors"
	"github.com/edalferes/monogo/internal/modules/auth/usecase"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	LoginUseCase    *usecase.LoginUseCase
	RegisterUseCase *usecase.RegisterUseCase
}

// Login godoc
// @Summary User login
// @Description Authenticates user and returns JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginDTO true "User credentials"
// @Success 200 {object} map[string]string "token"
// @Failure 400 {object} map[string]string "invalid data or missing credentials"
// @Failure 401 {object} map[string]string "invalid credentials"
// @Failure 500 {object} map[string]string "internal error"
// @Router /v1/auth/login [post]
func (h *Handler) Login(c echo.Context) error {
	var input LoginDTO
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errors.ErrInvalidData.Error()})
	}
	if input.Username == "" || input.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errors.ErrMissingCredentials.Error()})
	}
	token, err := h.LoginUseCase.Execute(input.Username, input.Password)
	if err != nil {
		if err == errors.ErrInvalidCredentials {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

// Register godoc
// @Summary User registration
// @Description Register a new user with role 'user'
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterDTO true "User registration data"
// @Success 201 {object} map[string]string "created"
// @Failure 400 {object} map[string]string "invalid data or user exists"
// @Failure 500 {object} map[string]string "internal error"
// @Router /v1/auth/register [post]
func (h *Handler) Register(c echo.Context) error {
	var input RegisterDTO
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errors.ErrInvalidData.Error()})
	}
	if input.Username == "" || input.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errors.ErrMissingCredentials.Error()})
	}
	err := h.RegisterUseCase.Execute(input.Username, input.Password)
	if err != nil {
		if err == errors.ErrUserAlreadyExists {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}
	return c.JSON(http.StatusCreated, map[string]string{"message": "user created"})
}
