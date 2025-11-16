package usecase

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/errors"
)

// CreateCategoryUseCase handles category creation
type CreateCategoryUseCase struct {
	categoryRepo repository.CategoryRepository
}

// NewCreateCategoryUseCase creates a new use case instance
func NewCreateCategoryUseCase(categoryRepo repository.CategoryRepository) *CreateCategoryUseCase {
	return &CreateCategoryUseCase{
		categoryRepo: categoryRepo,
	}
}

// Execute creates a new category
func (uc *CreateCategoryUseCase) Execute(ctx context.Context, input CreateCategoryInput) (domain.Category, error) {
	// Validate input
	if input.UserID == 0 {
		return domain.Category{}, errors.ErrInvalidUserID
	}
	if input.Name == "" {
		return domain.Category{}, errors.ErrCategoryNameRequired
	}
	if !isValidCategoryType(input.Type) {
		return domain.Category{}, errors.ErrInvalidCategoryType
	}

	// Create category domain entity
	category := domain.Category{
		UserID:      input.UserID,
		Name:        input.Name,
		Type:        input.Type,
		Icon:        input.Icon,
		Color:       input.Color,
		Description: input.Description,
	}

	// Save to repository
	return uc.categoryRepo.Create(ctx, category)
}

// CreateCategoryInput represents the input for creating a category
type CreateCategoryInput struct {
	UserID      uint                `json:"user_id"`
	Name        string              `json:"name"`
	Type        domain.CategoryType `json:"type"`
	Icon        string              `json:"icon,omitempty"`
	Color       string              `json:"color,omitempty"`
	Description string              `json:"description,omitempty"`
}

func isValidCategoryType(categoryType domain.CategoryType) bool {
	switch categoryType {
	case domain.CategoryTypeIncome, domain.CategoryTypeExpense:
		return true
	default:
		return false
	}
}
