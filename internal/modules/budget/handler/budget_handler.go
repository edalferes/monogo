package handler

import (
	"net/http"

	"github.com/edalferes/monogo/internal/modules/budget/domain"
	"github.com/edalferes/monogo/internal/modules/budget/handler/dto"
	"github.com/edalferes/monogo/internal/modules/budget/usecase"
	"github.com/labstack/echo/v4"
)

// BudgetHandler handles HTTP requests for budgets
type BudgetHandler struct {
	createBudgetUseCase *usecase.CreateBudgetUseCase
	listBudgetsUseCase  *usecase.ListBudgetsUseCase
}

// NewBudgetHandler creates a new budget handler
func NewBudgetHandler(
	createBudgetUseCase *usecase.CreateBudgetUseCase,
	listBudgetsUseCase *usecase.ListBudgetsUseCase,
) *BudgetHandler {
	return &BudgetHandler{
		createBudgetUseCase: createBudgetUseCase,
		listBudgetsUseCase:  listBudgetsUseCase,
	}
}

// CreateBudget handles budget creation
// @Summary Create a new budget
// @Tags Budget - Budgets
// @Accept json
// @Produce json
// @Param request body dto.CreateBudgetRequest true "Budget creation request"
// @Success 201 {object} dto.BudgetResponse
// @Failure 400 {object} map[string]interface{}
// @Router /budgets [post]
func (h *BudgetHandler) CreateBudget(c echo.Context) error {
	var req dto.CreateBudgetRequest
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

	input := usecase.CreateBudgetInput{
		UserID:      userID,
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Amount:      req.Amount,
		Period:      domain.BudgetPeriod(req.Period),
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		AlertAt:     req.AlertAt,
		Description: req.Description,
	}

	budget, err := h.createBudgetUseCase.Execute(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, dto.ToBudgetResponse(budget))
}

// ListBudgets handles listing user budgets
// @Summary List user budgets
// @Tags Budget - Budgets
// @Produce json
// @Param active_only query bool false "Filter active budgets only"
// @Success 200 {array} dto.BudgetResponse
// @Router /budgets [get]
func (h *BudgetHandler) ListBudgets(c echo.Context) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": err.Error(),
		})
	}

	input := usecase.ListBudgetsInput{
		UserID:     userID,
		ActiveOnly: c.QueryParam("active_only") == "true",
	}

	budgets, err := h.listBudgetsUseCase.Execute(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.ToBudgetResponseList(budgets))
}
