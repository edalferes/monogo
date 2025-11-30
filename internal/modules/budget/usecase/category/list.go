package category

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
)

// ListUseCase handles listing user categories
type ListUseCase struct {
	categoryRepo repository.CategoryRepository
}

// NewListUseCase creates a new use case instance
func NewListUseCase(categoryRepo repository.CategoryRepository) *ListUseCase {
	return &ListUseCase{
		categoryRepo: categoryRepo,
	}
}

// Execute lists all categories for a user
func (uc *ListUseCase) Execute(ctx context.Context, userID uint) ([]domain.Category, error) {
	return uc.categoryRepo.GetByUserID(ctx, userID)
}
