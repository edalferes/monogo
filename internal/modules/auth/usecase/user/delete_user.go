package user

import (
	"github.com/edalferes/monetics/internal/modules/auth/usecase/interfaces"
	"github.com/edalferes/monetics/pkg/logger"
)

type DeleteUserUseCase struct {
	UserRepo interfaces.User
	logger   logger.Logger
}

func NewDeleteUserUseCase(userRepo interfaces.User, log logger.Logger) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		UserRepo: userRepo,
		logger:   log.With().Str("usecase", "auth.delete_user").Logger(),
	}
}

func (u *DeleteUserUseCase) Execute(id uint) error {
	u.logger.Debug().Uint("user_id", id).Msg("deleting user")

	err := u.UserRepo.Delete(id)
	if err != nil {
		u.logger.Error().Err(err).Uint("user_id", id).Msg("failed to delete user")
		return err
	}

	u.logger.Info().Uint("user_id", id).Msg("user deleted successfully")
	return nil
}
