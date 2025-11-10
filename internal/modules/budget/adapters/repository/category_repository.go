package repository

import (
	"context"

	"github.com/edalferes/monogo/internal/modules/budget/domain"
)

// CategoryRepository defines the contract for category persistence operations
type CategoryRepository interface {
	Create(ctx context.Context, category domain.Category) (domain.Category, error)
	GetByID(ctx context.Context, id uint) (domain.Category, error)
	GetByUserID(ctx context.Context, userID uint) ([]domain.Category, error)
	GetByType(ctx context.Context, userID uint, categoryType domain.CategoryType) ([]domain.Category, error)
	Update(ctx context.Context, category domain.Category) (domain.Category, error)
	Delete(ctx context.Context, id uint) error
	ExistsByID(ctx context.Context, id uint) (bool, error)
}
