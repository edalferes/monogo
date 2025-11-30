package category

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/errors"
)

// CreateUseCase handles category creation
type CreateUseCase struct {
	categoryRepo repository.CategoryRepository
}

// NewCreateUseCase creates a new use case instance
func NewCreateUseCase(categoryRepo repository.CategoryRepository) *CreateUseCase {
	return &CreateUseCase{
		categoryRepo: categoryRepo,
	}
}

// CreateInput represents the input for creating a category
type CreateInput struct {
	UserID      uint                `json:"user_id" validate:"required"`
	Name        string              `json:"name" validate:"required"`
	Type        domain.CategoryType `json:"type" validate:"required"`
	Icon        string              `json:"icon,omitempty"`
	Color       string              `json:"color,omitempty"`
	Description string              `json:"description,omitempty"`
}

// Execute creates a new category
func (uc *CreateUseCase) Execute(ctx context.Context, input CreateInput) (domain.Category, error) {
	if input.UserID == 0 {
		return domain.Category{}, errors.ErrInvalidUserID
	}
	if input.Name == "" {
		return domain.Category{}, errors.ErrCategoryNameRequired
	}
	if !isValidCategoryType(input.Type) {
		return domain.Category{}, errors.ErrInvalidCategoryType
	}

	category := domain.Category{
		UserID:      input.UserID,
		Name:        input.Name,
		Type:        input.Type,
		Icon:        input.Icon,
		Color:       input.Color,
		Description: input.Description,
		IsActive:    true,
	}

	return uc.categoryRepo.Create(ctx, category)
}
