package category

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/errors"
	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
	"github.com/edalferes/monetics/pkg/logger"
)

// GetByIDUseCase handles retrieving a specific category
type GetByIDUseCase struct {
	categoryRepo interfaces.CategoryRepository
	logger       logger.Logger
}

// NewGetByIDUseCase creates a new use case instance
func NewGetByIDUseCase(categoryRepo interfaces.CategoryRepository, log logger.Logger) *GetByIDUseCase {
	return &GetByIDUseCase{
		categoryRepo: categoryRepo,
		logger:       log.With().Str("usecase", "category.getbyid").Logger(),
	}
}

// Execute retrieves a category by ID
func (uc *GetByIDUseCase) Execute(ctx context.Context, userID, categoryID uint) (domain.Category, error) {
	uc.logger.Debug().Uint("user_id", userID).Uint("category_id", categoryID).Msg("getting category by id")

	if categoryID == 0 {
		uc.logger.Error().Msg("invalid category_id: cannot be zero")
		return domain.Category{}, errors.ErrCategoryNotFound
	}

	category, err := uc.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		uc.logger.Error().Err(err).Uint("category_id", categoryID).Msg("category not found")
		return domain.Category{}, err
	}

	// Ensure the category belongs to the user
	if category.UserID != userID {
		uc.logger.Warn().
			Uint("category_user_id", category.UserID).
			Uint("request_user_id", userID).
			Msg("unauthorized access: category belongs to different user")
		return domain.Category{}, errors.ErrCategoryNotFound
	}

	uc.logger.Info().Uint("category_id", categoryID).Msg("category retrieved successfully")
	return category, nil
}
