package user

import (
	"github.com/edalferes/monetics/internal/modules/auth/errors"
	"github.com/edalferes/monetics/internal/modules/auth/usecase/interfaces"
	"github.com/edalferes/monetics/pkg/logger"
)

type LoginUseCase struct {
	UserRepo        interfaces.User
	PasswordService interfaces.PasswordService
	JWTService      interfaces.JWTService
	logger          logger.Logger
}

func NewLoginUseCase(userRepo interfaces.User, passwordService interfaces.PasswordService, jwtService interfaces.JWTService, log logger.Logger) *LoginUseCase {
	return &LoginUseCase{
		UserRepo:        userRepo,
		PasswordService: passwordService,
		JWTService:      jwtService,
		logger:          log.With().Str("usecase", "auth.login").Logger(),
	}
}

func (u *LoginUseCase) Execute(username, password string) (string, error) {
	u.logger.Debug().Str("username", username).Msg("login attempt")

	user, err := u.UserRepo.FindByUsername(username)
	if err != nil || user == nil {
		u.logger.Warn().Str("username", username).Msg("login failed: user not found")
		return "", errors.ErrInvalidCredentials
	}
	if err := u.PasswordService.Compare(user.Password, password); err != nil {
		u.logger.Warn().Str("username", username).Uint("user_id", user.ID).Msg("login failed: invalid password")
		return "", errors.ErrInvalidCredentials
	}
	// Extract role names from user roles
	roles := make([]string, len(user.Roles))
	for i, r := range user.Roles {
		roles[i] = r.Name
	}
	token, err := u.JWTService.GenerateToken(user.ID, user.Username, roles)
	if err != nil {
		u.logger.Error().Err(err).Str("username", username).Uint("user_id", user.ID).Msg("failed to generate token")
		return "", errors.ErrInvalidData
	}
	u.logger.Info().Str("username", username).Uint("user_id", user.ID).Msg("login successful")
	return token, nil
}
