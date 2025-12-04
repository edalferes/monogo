package category

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
	"github.com/edalferes/monetics/internal/modules/budget/errors"
)

// DeleteUseCase handles category deletion (soft delete)
type DeleteUseCase struct {
	categoryRepo interfaces.CategoryRepository
}

// NewDeleteUseCase creates a new use case instance
func NewDeleteUseCase(categoryRepo interfaces.CategoryRepository) *DeleteUseCase {
	return &DeleteUseCase{
		categoryRepo: categoryRepo,
	}
}

// Execute soft-deletes a category
func (uc *DeleteUseCase) Execute(ctx context.Context, categoryID, userID uint) error {
	if categoryID == 0 {
		return errors.ErrCategoryNotFound
	}

	// Get existing category
	category, err := uc.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return err
	}

	// Ensure the category belongs to the user
	if category.UserID != userID {
		return errors.ErrCategoryNotFound
	}

	// Soft delete by marking as inactive
	return uc.categoryRepo.Delete(ctx, categoryID)
}
