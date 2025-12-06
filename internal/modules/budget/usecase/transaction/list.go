package transaction

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/usecase/interfaces"
	"github.com/edalferes/monetics/pkg/logger"
)

type ListUseCase struct {
	transactionRepo interfaces.TransactionRepository
	logger          logger.Logger
}

func NewListUseCase(transactionRepo interfaces.TransactionRepository, log logger.Logger) *ListUseCase {
	return &ListUseCase{
		transactionRepo: transactionRepo,
		logger:          log.With().Str("usecase", "transaction.list").Logger(),
	}
}

type ListInput struct {
	UserID   uint
	Page     int
	PageSize int
}

type ListOutput struct {
	Transactions []domain.Transaction
	Total        int64
	Page         int
	PageSize     int
	TotalPages   int
}

func (uc *ListUseCase) Execute(ctx context.Context, input ListInput) (ListOutput, error) {
	uc.logger.Debug().
		Uint("user_id", input.UserID).
		Int("page", input.Page).
		Int("page_size", input.PageSize).
		Msg("listing transactions")

	// Default pagination values
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 20 // Default 20 items per page
	}
	if input.PageSize > 100 {
		input.PageSize = 100 // Max 100 items per page
	}

	// Calculate offset
	offset := (input.Page - 1) * input.PageSize

	// Get paginated transactions
	transactions, err := uc.transactionRepo.GetByUserIDPaginated(ctx, input.UserID, input.PageSize, offset)
	if err != nil {
		uc.logger.Error().Err(err).Uint("user_id", input.UserID).Msg("failed to get paginated transactions")
		return ListOutput{}, err
	}

	// Get total count
	total, err := uc.transactionRepo.CountByUserID(ctx, input.UserID)
	if err != nil {
		uc.logger.Error().Err(err).Uint("user_id", input.UserID).Msg("failed to count transactions")
		return ListOutput{}, err
	}

	// Calculate total pages
	totalPages := int(total) / input.PageSize
	if int(total)%input.PageSize > 0 {
		totalPages++
	}

	uc.logger.Info().
		Uint("user_id", input.UserID).
		Int("count", len(transactions)).
		Int("total", int(total)).
		Msg("transactions listed successfully")

	return ListOutput{
		Transactions: transactions,
		Total:        total,
		Page:         input.Page,
		PageSize:     input.PageSize,
		TotalPages:   totalPages,
	}, nil
}
