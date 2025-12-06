package category

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/errors"
	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
	"github.com/edalferes/monetics/pkg/logger"
)

// CreateUseCase handles category creation
type CreateUseCase struct {
	categoryRepo interfaces.CategoryRepository
	logger       logger.Logger
}

// NewCreateUseCase creates a new use case instance
func NewCreateUseCase(categoryRepo interfaces.CategoryRepository, log logger.Logger) *CreateUseCase {
	return &CreateUseCase{
		categoryRepo: categoryRepo,
		logger:       log.With().Str("usecase", "category.create").Logger(),
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
	uc.logger.Debug().
		Uint("user_id", input.UserID).
		Str("name", input.Name).
		Str("type", string(input.Type)).
		Msg("creating category")

	if input.UserID == 0 {
		uc.logger.Error().Msg("validation error: user_id cannot be zero")
		return domain.Category{}, errors.ErrInvalidUserID
	}
	if input.Name == "" {
		uc.logger.Error().Msg("validation error: category name cannot be empty")
		return domain.Category{}, errors.ErrCategoryNameRequired
	}
	if !isValidCategoryType(input.Type) {
		uc.logger.Error().Str("type", string(input.Type)).Msg("validation error: invalid category type")
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

	created, err := uc.categoryRepo.Create(ctx, category)
	if err != nil {
		uc.logger.Error().Err(err).Str("name", input.Name).Msg("failed to create category")
		return domain.Category{}, err
	}

	uc.logger.Info().
		Uint("category_id", created.ID).
		Uint("user_id", input.UserID).
		Str("name", input.Name).
		Msg("category created successfully")
	return created, nil
}
