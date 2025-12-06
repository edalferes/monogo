package category

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/errors"
	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
	"github.com/edalferes/monetics/pkg/logger"
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
	categoryRepo interfaces.CategoryRepository
	logger       logger.Logger
}

// NewUpdateUseCase creates a new use case instance
func NewUpdateUseCase(categoryRepo interfaces.CategoryRepository, log logger.Logger) *UpdateUseCase {
	return &UpdateUseCase{
		categoryRepo: categoryRepo,
		logger:       log.With().Str("usecase", "category.update").Logger(),
	}
}

// Execute updates an existing category
func (uc *UpdateUseCase) Execute(ctx context.Context, input UpdateInput) (domain.Category, error) {
	uc.logger.Debug().Uint("user_id", input.UserID).Uint("category_id", input.ID).Msg("updating category")

	if input.ID == 0 {
		uc.logger.Error().Msg("invalid category_id: cannot be zero")
		return domain.Category{}, errors.ErrCategoryNotFound
	}

	// Get existing category
	category, err := uc.categoryRepo.GetByID(ctx, input.ID)
	if err != nil {
		uc.logger.Error().Err(err).Uint("category_id", input.ID).Msg("category not found")
		return domain.Category{}, err
	}

	// Ensure the category belongs to the user
	if category.UserID != input.UserID {
		uc.logger.Warn().
			Uint("category_user_id", category.UserID).
			Uint("request_user_id", input.UserID).
			Msg("unauthorized access: category belongs to different user")
		return domain.Category{}, errors.ErrCategoryNotFound
	}

	// Update fields if provided
	if input.Name != nil {
		if *input.Name == "" {
			uc.logger.Error().Msg("validation error: category name cannot be empty")
			return domain.Category{}, errors.ErrCategoryNameRequired
		}
		category.Name = *input.Name
	}

	if input.Type != nil {
		if !isValidCategoryType(*input.Type) {
			uc.logger.Error().Str("type", string(*input.Type)).Msg("validation error: invalid category type")
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
		uc.logger.Error().Err(err).Uint("category_id", input.ID).Msg("failed to update category")
		return domain.Category{}, err
	}

	uc.logger.Info().Uint("category_id", input.ID).Msg("category updated successfully")
	return updatedCategory, nil
}
