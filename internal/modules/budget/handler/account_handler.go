package handler

import (
	"net/http"
	"strconv"

	"github.com/edalferes/monogo/internal/modules/budget/domain"
	"github.com/edalferes/monogo/internal/modules/budget/handler/dto"
	"github.com/edalferes/monogo/internal/modules/budget/usecase"
	"github.com/labstack/echo/v4"
)

// AccountHandler handles HTTP requests for accounts
type AccountHandler struct {
	createAccountUseCase     *usecase.CreateAccountUseCase
	listAccountsUseCase      *usecase.ListAccountsUseCase
	getAccountBalanceUseCase *usecase.GetAccountBalanceUseCase
}

// NewAccountHandler creates a new account handler
func NewAccountHandler(
	createAccountUseCase *usecase.CreateAccountUseCase,
	listAccountsUseCase *usecase.ListAccountsUseCase,
	getAccountBalanceUseCase *usecase.GetAccountBalanceUseCase,
) *AccountHandler {
	return &AccountHandler{
		createAccountUseCase:     createAccountUseCase,
		listAccountsUseCase:      listAccountsUseCase,
		getAccountBalanceUseCase: getAccountBalanceUseCase,
	}
}

// CreateAccount handles account creation
// @Summary Create a new account
// @Tags Budget - Accounts
// @Accept json
// @Produce json
// @Param request body dto.CreateAccountRequest true "Account creation request"
// @Success 201 {object} dto.AccountResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /accounts [post]
func (h *AccountHandler) CreateAccount(c echo.Context) error {
	var req dto.CreateAccountRequest
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

	// Get user ID from context (set by auth middleware)
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": err.Error(),
		})
	}

	input := usecase.CreateAccountInput{
		UserID:         userID,
		Name:           req.Name,
		Type:           domain.AccountType(req.Type),
		InitialBalance: req.InitialBalance,
		Currency:       req.Currency,
		Description:    req.Description,
	}

	account, err := h.createAccountUseCase.Execute(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, dto.ToAccountResponse(account))
}

// ListAccounts handles listing user accounts
// @Summary List user accounts
// @Tags Budget - Accounts
// @Produce json
// @Success 200 {array} dto.AccountResponse
// @Failure 500 {object} map[string]interface{}
// @Router /accounts [get]
func (h *AccountHandler) ListAccounts(c echo.Context) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": err.Error(),
		})
	}

	accounts, err := h.listAccountsUseCase.Execute(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.ToAccountResponseList(accounts))
}

// GetAccount handles getting a single account
// @Summary Get account by ID with calculated balance
// @Tags Budget - Accounts
// @Produce json
// @Param id path int true "Account ID"
// @Success 200 {object} dto.AccountBalanceResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /accounts/{id} [get]
func (h *AccountHandler) GetAccount(c echo.Context) error {
	accountID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid account ID",
		})
	}

	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": err.Error(),
		})
	}

	result, err := h.getAccountBalanceUseCase.Execute(c.Request().Context(), userID, uint(accountID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
	}

	response := dto.AccountBalanceResponse{
		Account:        dto.ToAccountResponse(result.Account),
		CurrentBalance: result.CurrentBalance,
		TotalIncome:    result.TotalIncome,
		TotalExpense:   result.TotalExpense,
		TotalTransfers: result.TotalTransfers,
	}

	return c.JSON(http.StatusOK, response)
}
