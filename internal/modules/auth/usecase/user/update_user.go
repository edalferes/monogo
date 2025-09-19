package user

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"github.com/edalferes/monogo/internal/modules/auth/repository"
)

type UpdateUserInput struct {
	ID       uint
	Username string
	Password *string
	RoleIDs  []uint
}

type UpdateUserUseCase struct {
	UserRepo repository.UserRepository
	RoleRepo repository.RoleRepository
}

func (u *UpdateUserUseCase) Execute(input UpdateUserInput) error {
	user, err := u.UserRepo.FindByID(input.ID)
	if err != nil {
		return err
	}
	if input.Username != "" {
		user.Username = input.Username
	}
	if input.Password != nil {
		user.Password = *input.Password
	}
	if input.RoleIDs != nil {
		var roles []domain.Role
		for _, rid := range input.RoleIDs {
			role, err := u.RoleRepo.FindByID(rid)
			if err != nil {
				return err
			}
			roles = append(roles, *role)
		}
		user.Roles = roles
	}
	return u.UserRepo.Update(user)
}
