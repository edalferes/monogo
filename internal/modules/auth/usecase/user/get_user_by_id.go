package user

import (
	"github.com/edalferes/monetics/internal/modules/auth/domain"
	"github.com/edalferes/monetics/internal/modules/auth/usecase/interfaces"
	"github.com/edalferes/monetics/pkg/logger"
)

type GetUserByIDUseCase struct {
	UserRepo interfaces.User
	logger   logger.Logger
}

func NewGetUserByIDUseCase(userRepo interfaces.User, log logger.Logger) *GetUserByIDUseCase {
	return &GetUserByIDUseCase{
		UserRepo: userRepo,
		logger:   log.With().Str("usecase", "auth.get_user_by_id").Logger(),
	}
}

func (u *GetUserByIDUseCase) Execute(id uint) (*domain.User, error) {
	u.logger.Debug().Uint("user_id", id).Msg("getting user by id")

	user, err := u.UserRepo.FindByID(id)
	if err != nil {
		u.logger.Error().Err(err).Uint("user_id", id).Msg("user not found")
		return nil, err
	}

	u.logger.Info().Uint("user_id", id).Str("username", user.Username).Msg("user retrieved successfully")
	return user, nil
}
