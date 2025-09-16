package service

import (
	"github.com/edalferes/monogo/internal/modules/user/domain"
	"github.com/edalferes/monogo/internal/modules/user/usecase"
)

type Service struct {
	CreateUserUC *usecase.CreateUserUseCase
}

func (s *Service) Register(name, email string) (*domain.User, error) {
	return s.CreateUserUC.Execute(name, email)
}
