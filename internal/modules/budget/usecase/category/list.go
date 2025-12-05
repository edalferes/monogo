package category

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
)

// ListUseCase handles listing user categories
type ListUseCase struct {
	categoryRepo interfaces.CategoryRepository
}

// NewListUseCase creates a new use case instance
func NewListUseCase(categoryRepo interfaces.CategoryRepository) *ListUseCase {
	return &ListUseCase{
		categoryRepo: categoryRepo,
	}
}

// Execute lists all categories for a user
func (uc *ListUseCase) Execute(ctx context.Context, userID uint) ([]domain.Category, error) {
	return uc.categoryRepo.GetByUserID(ctx, userID)
}
