package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/handler/dto"
	"github.com/edalferes/monetics/internal/modules/budget/usecase/transaction"
	"github.com/labstack/echo/v4"
)

// convertTransactionType converts string pointer to TransactionType pointer
func convertTransactionType(t *string) *domain.TransactionType {
	if t == nil {
		return nil
	}
	txType := domain.TransactionType(*t)
	return &txType
}

// TransactionHandler handles HTTP requests for transactions
type TransactionHandler struct {
	createTransactionUseCase  *transaction.CreateUseCase
	listTransactionsUseCase   *transaction.ListUseCase
	getTransactionByIDUseCase *transaction.GetByIDUseCase
	updateTransactionUseCase  *transaction.UpdateUseCase
	deleteTransactionUseCase  *transaction.DeleteUseCase
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(
	createTransactionUseCase *transaction.CreateUseCase,
	listTransactionsUseCase *transaction.ListUseCase,
	getTransactionByIDUseCase *transaction.GetByIDUseCase,
	updateTransactionUseCase *transaction.UpdateUseCase,
	deleteTransactionUseCase *transaction.DeleteUseCase,
) *TransactionHandler {
	return &TransactionHandler{
		createTransactionUseCase:  createTransactionUseCase,
		listTransactionsUseCase:   listTransactionsUseCase,
		getTransactionByIDUseCase: getTransactionByIDUseCase,
		updateTransactionUseCase:  updateTransactionUseCase,
		deleteTransactionUseCase:  deleteTransactionUseCase,
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

	input := transaction.CreateInput{
		UserID:               userID,
		AccountID:            req.AccountID,
		CategoryID:           req.CategoryID,
		Type:                 domain.TransactionType(req.Type),
		Amount:               req.Amount,
		Description:          req.Description,
		Date:                 req.Date.Format(time.RFC3339),
		DestinationAccountID: req.DestinationAccountID,
	}

	tx, err := h.createTransactionUseCase.Execute(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, dto.ToTransactionResponse(tx))
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

	// TODO: Implement filters (account_id, category_id, type, date range) if needed
	transactions, err := h.listTransactionsUseCase.Execute(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.ToTransactionResponseList(transactions))
}

// GetTransactionByID handles getting transaction by ID
// @Summary Get transaction by ID
// @Tags Budget - Transactions
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} dto.TransactionResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /transactions/{id} [get]
func (h *TransactionHandler) GetTransactionByID(c echo.Context) error {
	transactionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid transaction ID",
		})
	}

	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": err.Error(),
		})
	}

	tx, err := h.getTransactionByIDUseCase.Execute(c.Request().Context(), userID, uint(transactionID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.ToTransactionResponse(tx))
}

// UpdateTransaction handles transaction update
// @Summary Update a transaction
// @Tags Budget - Transactions
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Param request body dto.UpdateTransactionRequest true "Transaction update request"
// @Success 200 {object} dto.TransactionResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /transactions/{id} [put]
func (h *TransactionHandler) UpdateTransaction(c echo.Context) error {
	transactionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid transaction ID",
		})
	}

	var req dto.UpdateTransactionRequest
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

	var dateStr *string
	if req.Date != nil {
		ds := req.Date.Format(time.RFC3339)
		dateStr = &ds
	}

	input := transaction.UpdateInput{
		ID:          uint(transactionID),
		UserID:      userID,
		AccountID:   req.AccountID,
		CategoryID:  req.CategoryID,
		Type:        convertTransactionType(req.Type),
		Amount:      req.Amount,
		Description: req.Description,
		Date:        dateStr,
	}

	tx, err := h.updateTransactionUseCase.Execute(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.ToTransactionResponse(tx))
}

// DeleteTransaction handles transaction deletion (soft delete)
// @Summary Delete a transaction
// @Tags Budget - Transactions
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /transactions/{id} [delete]
func (h *TransactionHandler) DeleteTransaction(c echo.Context) error {
	transactionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid transaction ID",
		})
	}

	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": err.Error(),
		})
	}

	err = h.deleteTransactionUseCase.Execute(c.Request().Context(), userID, uint(transactionID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}
