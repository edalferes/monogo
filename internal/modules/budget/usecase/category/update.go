package category

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/errors"
)

// UpdateInput represents the input for updating a category
type UpdateInput struct {
	ID          uint
	UserID      uint
	Name        *string
	Type        *domain.CategoryType
	Icon        *string
	Color       *string
	Description *string
}

// UpdateUseCase handles category updates
type UpdateUseCase struct {
	categoryRepo repository.CategoryRepository
}

// NewUpdateUseCase creates a new use case instance
func NewUpdateUseCase(categoryRepo repository.CategoryRepository) *UpdateUseCase {
	return &UpdateUseCase{
		categoryRepo: categoryRepo,
	}
}

// Execute updates an existing category
func (uc *UpdateUseCase) Execute(ctx context.Context, input UpdateInput) (domain.Category, error) {
	if input.ID == 0 {
		return domain.Category{}, errors.ErrCategoryNotFound
	}

	// Get existing category
	category, err := uc.categoryRepo.GetByID(ctx, input.ID)
	if err != nil {
		return domain.Category{}, err
	}

	// Ensure the category belongs to the user
	if category.UserID != input.UserID {
		return domain.Category{}, errors.ErrCategoryNotFound
	}

	// Update fields if provided
	if input.Name != nil {
		if *input.Name == "" {
			return domain.Category{}, errors.ErrCategoryNameRequired
		}
		category.Name = *input.Name
	}

	if input.Type != nil {
		if !isValidCategoryType(*input.Type) {
			return domain.Category{}, errors.ErrInvalidCategoryType
		}
		category.Type = *input.Type
	}

	if input.Icon != nil {
		category.Icon = *input.Icon
	}

	if input.Color != nil {
		category.Color = *input.Color
	}

	if input.Description != nil {
		category.Description = *input.Description
	}

	// Save updates
	updatedCategory, err := uc.categoryRepo.Update(ctx, category)
	if err != nil {
		return domain.Category{}, err
	}

	return updatedCategory, nil
}
