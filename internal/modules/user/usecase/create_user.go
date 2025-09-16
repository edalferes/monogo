package usecase

import (
	"github.com/edalferes/monogo/internal/modules/user/domain"
	"github.com/edalferes/monogo/internal/modules/user/repository"
)

type CreateUserUseCase struct {
	repo repository.UserRepository
}

func NewCreateUserUseCase(repo repository.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{repo: repo}
}

func (uc *CreateUserUseCase) Execute(name, email string) (*domain.User, error) {
	id, err := uc.repo.CreateUser(name, email)
	if err != nil {
		return nil, err
	}
	user, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
