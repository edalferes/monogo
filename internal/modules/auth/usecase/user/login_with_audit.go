package user

import (
	"github.com/edalferes/monetics/internal/modules/auth/errors"
	"github.com/edalferes/monetics/internal/modules/auth/usecase/interfaces"
)

type LoginWithAuditUseCase struct {
	UserRepo        interfaces.User
	PasswordService interfaces.PasswordService
	JWTService      interfaces.JWTService
	AuditService    interfaces.AuditService
}

func (u *LoginWithAuditUseCase) Execute(username, password, ip string) (string, error) {
	user, err := u.UserRepo.FindByUsername(username)
	if err != nil || user == nil {
		u.AuditService.Log(nil, username, "login_failed", "fail", ip, "user not found or error")
		return "", errors.ErrInvalidCredentials
	}
	if err := u.PasswordService.Compare(user.Password, password); err != nil {
		u.AuditService.Log(&user.ID, username, "login_failed", "fail", ip, "wrong password")
		return "", errors.ErrInvalidCredentials
	}
	roles := make([]string, len(user.Roles))
	for i, r := range user.Roles {
		roles[i] = r.Name
	}
	token, err := u.JWTService.GenerateToken(user.ID, user.Username, roles)
	if err != nil {
		u.AuditService.Log(&user.ID, username, "login_failed", "fail", ip, "token error")
		return "", errors.ErrInvalidData
	}
	u.AuditService.Log(&user.ID, username, "login_success", "ok", ip, "")
	return token, nil
}
