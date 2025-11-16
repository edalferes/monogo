package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/handler/dto"
	"github.com/edalferes/monetics/internal/modules/budget/usecase"
	"github.com/labstack/echo/v4"
)

// TransactionHandler handles HTTP requests for transactions
type TransactionHandler struct {
	createTransactionUseCase *usecase.CreateTransactionUseCase
	listTransactionsUseCase  *usecase.ListTransactionsUseCase
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(
	createTransactionUseCase *usecase.CreateTransactionUseCase,
	listTransactionsUseCase *usecase.ListTransactionsUseCase,
) *TransactionHandler {
	return &TransactionHandler{
		createTransactionUseCase: createTransactionUseCase,
		listTransactionsUseCase:  listTransactionsUseCase,
	}
}

// CreateTransaction handles transaction creation
// @Summary Create a new transaction
// @Tags Budget - Transactions
// @Accept json
// @Produce json
// @Param request body dto.CreateTransactionRequest true "Transaction creation request"
// @Success 201 {object} dto.TransactionResponse
// @Failure 400 {object} map[string]interface{}
// @Router /transactions [post]
func (h *TransactionHandler) CreateTransaction(c echo.Context) error {
	var req dto.CreateTransactionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": err.Error(),
		})
	}

	input := usecase.CreateTransactionInput{
		UserID:               userID,
		AccountID:            req.AccountID,
		CategoryID:           req.CategoryID,
		Type:                 domain.TransactionType(req.Type),
		Amount:               req.Amount,
		Description:          req.Description,
		Date:                 req.Date,
		Tags:                 req.Tags,
		Attachments:          req.Attachments,
		IsRecurring:          req.IsRecurring,
		RecurrenceRule:       req.RecurrenceRule,
		RecurrenceEnd:        req.RecurrenceEnd,
		DestinationAccountID: req.DestinationAccountID,
		TransferFee:          req.TransferFee,
	}

	transaction, err := h.createTransactionUseCase.Execute(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, dto.ToTransactionResponse(transaction))
}

// ListTransactions handles listing user transactions
// @Summary List user transactions
// @Tags Budget - Transactions
// @Produce json
// @Param account_id query int false "Filter by account ID"
// @Param category_id query int false "Filter by category ID"
// @Param type query string false "Filter by type (income, expense, transfer)"
// @Param start_date query string false "Filter by start date (RFC3339)"
// @Param end_date query string false "Filter by end date (RFC3339)"
// @Success 200 {array} dto.TransactionResponse
// @Router /transactions [get]
func (h *TransactionHandler) ListTransactions(c echo.Context) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": err.Error(),
		})
	}

	input := usecase.ListTransactionsInput{
		UserID: userID,
	}

	// Parse optional filters
	if accountIDStr := c.QueryParam("account_id"); accountIDStr != "" {
		if accountID, err := strconv.ParseUint(accountIDStr, 10, 32); err == nil {
			id := uint(accountID)
			input.AccountID = &id
		}
	}

	if categoryIDStr := c.QueryParam("category_id"); categoryIDStr != "" {
		if categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32); err == nil {
			id := uint(categoryID)
			input.CategoryID = &id
		}
	}

	if typeParam := c.QueryParam("type"); typeParam != "" {
		t := domain.TransactionType(typeParam)
		input.Type = &t
	}

	if startDateStr := c.QueryParam("start_date"); startDateStr != "" {
		if startDate, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			input.StartDate = &startDate
		}
	}

	if endDateStr := c.QueryParam("end_date"); endDateStr != "" {
		if endDate, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			input.EndDate = &endDate
		}
	}

	transactions, err := h.listTransactionsUseCase.Execute(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.ToTransactionResponseList(transactions))
}
