package category

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/errors"
)

// GetByIDUseCase handles retrieving a specific category
type GetByIDUseCase struct {
	categoryRepo interfaces.CategoryRepository
}

// NewGetByIDUseCase creates a new use case instance
func NewGetByIDUseCase(categoryRepo interfaces.CategoryRepository) *GetByIDUseCase {
	return &GetByIDUseCase{
		categoryRepo: categoryRepo,
	}
}

// Execute retrieves a category by ID
func (uc *GetByIDUseCase) Execute(ctx context.Context, userID, categoryID uint) (domain.Category, error) {
	if categoryID == 0 {
		return domain.Category{}, errors.ErrCategoryNotFound
	}

	category, err := uc.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return domain.Category{}, err
	}

	// Ensure the category belongs to the user
	if category.UserID != userID {
		return domain.Category{}, errors.ErrCategoryNotFound
	}

	return category, nil
}
