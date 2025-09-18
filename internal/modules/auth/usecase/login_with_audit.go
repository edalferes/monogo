package usecase

import (
	"github.com/edalferes/monogo/internal/modules/auth/errors"
	"github.com/edalferes/monogo/internal/modules/auth/repository"
	"github.com/edalferes/monogo/internal/modules/auth/service"
)

type LoginWithAuditUseCase struct {
	UserRepo        repository.UserRepository
	PasswordService service.PasswordService
	JWTService      service.JWTService
	AuditService    service.AuditService
}

// Execute performs login and logs audit
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
