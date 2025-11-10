package usecase

import (
	"context"

	"github.com/edalferes/monogo/internal/modules/budget/adapters/repository"
	"github.com/edalferes/monogo/internal/modules/budget/domain"
)

// ListCategoriesUseCase handles listing user categories
type ListCategoriesUseCase struct {
	categoryRepo repository.CategoryRepository
}

// NewListCategoriesUseCase creates a new use case instance
func NewListCategoriesUseCase(categoryRepo repository.CategoryRepository) *ListCategoriesUseCase {
	return &ListCategoriesUseCase{
		categoryRepo: categoryRepo,
	}
}

// Execute lists categories for a user, optionally filtered by type
func (uc *ListCategoriesUseCase) Execute(ctx context.Context, userID uint, categoryType *domain.CategoryType) ([]domain.Category, error) {
	if categoryType != nil {
		return uc.categoryRepo.GetByType(ctx, userID, *categoryType)
	}
	return uc.categoryRepo.GetByUserID(ctx, userID)
}
