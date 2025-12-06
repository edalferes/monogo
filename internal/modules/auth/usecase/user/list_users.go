package user

import (
	"github.com/edalferes/monetics/internal/modules/auth/domain"
	"github.com/edalferes/monetics/internal/modules/auth/usecase/interfaces"
	"github.com/edalferes/monetics/pkg/logger"
)

type ListUsersUseCase struct {
	UserRepo interfaces.User
	logger   logger.Logger
}

func NewListUsersUseCase(userRepo interfaces.User, log logger.Logger) *ListUsersUseCase {
	return &ListUsersUseCase{
		UserRepo: userRepo,
		logger:   log.With().Str("usecase", "auth.list_users").Logger(),
	}
}

func (u *ListUsersUseCase) Execute() ([]domain.User, error) {
	u.logger.Debug().Msg("listing all users")

	users, err := u.UserRepo.ListAll()
	if err != nil {
		u.logger.Error().Err(err).Msg("failed to list users")
		return nil, err
	}

	u.logger.Info().Int("count", len(users)).Msg("users listed successfully")
	return users, nil
}
