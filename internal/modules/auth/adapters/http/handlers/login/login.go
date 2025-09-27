package login

import (
	"net/http"

	"github.com/edalferes/monogo/internal/infra/logger"
	"github.com/edalferes/monogo/internal/modules/auth/adapters/http/dto"
	"github.com/edalferes/monogo/internal/modules/auth/errors"
	userUC "github.com/edalferes/monogo/internal/modules/auth/usecase/user"
	"github.com/edalferes/monogo/pkg/utils"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	LoginUseCase *userUC.LoginWithAuditUseCase
	Logger       logger.Logger
}

// Login godoc
// @Summary User login
// @Description Authenticates user and returns JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body dto.LoginDTO true "User credentials"
// @Success 200 {object} map[string]string "token"
// @Failure 400 {object} map[string]string "invalid data or missing credentials"
// @Failure 401 {object} map[string]string "invalid credentials"
// @Failure 500 {object} map[string]string "internal error"
// @Router /v1/auth/login [post]
func (h *Handler) Login(c echo.Context) error {
	var input dto.LoginDTO
	if err := c.Bind(&input); err != nil {
		h.Logger.Error().Err(err).Msg("failed to bind login request")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errors.ErrInvalidData.Error()})
	}
	if input.Username == "" || input.Password == "" {
		h.Logger.Warn().Str("username", input.Username).Msg("login attempt with missing credentials")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errors.ErrMissingCredentials.Error()})
	}
	ip := utils.ToIPv4(c.RealIP())
	h.Logger.Info().Str("username", input.Username).Str("ip", ip).Msg("login attempt")

	token, err := h.LoginUseCase.Execute(input.Username, input.Password, ip)
	if err != nil {
		if err == errors.ErrInvalidCredentials {
			h.Logger.Warn().Str("username", input.Username).Str("ip", ip).Msg("invalid credentials")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}
		h.Logger.Error().Err(err).Str("username", input.Username).Msg("internal error during login")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}

	h.Logger.Info().Str("username", input.Username).Str("ip", ip).Msg("successful login")
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
