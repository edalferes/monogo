package interfaces

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/domain"
)

// AccountRepository defines the contract for account persistence operations
type AccountRepository interface {
	Create(ctx context.Context, account domain.Account) (domain.Account, error)
	GetByID(ctx context.Context, id uint) (domain.Account, error)
	GetByUserID(ctx context.Context, userID uint) ([]domain.Account, error)
	Update(ctx context.Context, account domain.Account) (domain.Account, error)
	Delete(ctx context.Context, id uint) error
	ExistsByID(ctx context.Context, id uint) (bool, error)
}
