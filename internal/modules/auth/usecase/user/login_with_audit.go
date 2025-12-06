package user

import (
	"github.com/edalferes/monetics/internal/modules/auth/errors"
	"github.com/edalferes/monetics/internal/modules/auth/usecase/interfaces"
	"github.com/edalferes/monetics/pkg/logger"
)

type LoginWithAuditUseCase struct {
	UserRepo        interfaces.User
	PasswordService interfaces.PasswordService
	JWTService      interfaces.JWTService
	AuditService    interfaces.AuditService
	logger          logger.Logger
}

func NewLoginWithAuditUseCase(userRepo interfaces.User, passwordService interfaces.PasswordService, jwtService interfaces.JWTService, auditService interfaces.AuditService, log logger.Logger) *LoginWithAuditUseCase {
	return &LoginWithAuditUseCase{
		UserRepo:        userRepo,
		PasswordService: passwordService,
		JWTService:      jwtService,
		AuditService:    auditService,
		logger:          log.With().Str("usecase", "auth.login_with_audit").Logger(),
	}
}

func (u *LoginWithAuditUseCase) Execute(username, password, ip string) (string, error) {
	u.logger.Debug().Str("username", username).Str("ip", ip).Msg("login attempt with audit")

	user, err := u.UserRepo.FindByUsername(username)
	if err != nil || user == nil {
		u.logger.Warn().Str("username", username).Str("ip", ip).Msg("login failed: user not found")
		u.AuditService.Log(nil, username, "login_failed", "fail", ip, "user not found or error")
		return "", errors.ErrInvalidCredentials
	}
	if err := u.PasswordService.Compare(user.Password, password); err != nil {
		u.logger.Warn().Str("username", username).Uint("user_id", user.ID).Str("ip", ip).Msg("login failed: invalid password")
		u.AuditService.Log(&user.ID, username, "login_failed", "fail", ip, "wrong password")
		return "", errors.ErrInvalidCredentials
	}
	roles := make([]string, len(user.Roles))
	for i, r := range user.Roles {
		roles[i] = r.Name
	}
	token, err := u.JWTService.GenerateToken(user.ID, user.Username, roles)
	if err != nil {
		u.logger.Error().Err(err).Str("username", username).Uint("user_id", user.ID).Msg("failed to generate token")
		u.AuditService.Log(&user.ID, username, "login_failed", "fail", ip, "token error")
		return "", errors.ErrInvalidData
	}
	u.logger.Info().Str("username", username).Uint("user_id", user.ID).Str("ip", ip).Msg("login successful with audit")
	u.AuditService.Log(&user.ID, username, "login_success", "ok", ip, "")
	return token, nil
}
