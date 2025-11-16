package handler

import (
	"net/http"

	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"github.com/edalferes/monetics/internal/modules/budget/handler/dto"
	"github.com/edalferes/monetics/internal/modules/budget/usecase"
	"github.com/labstack/echo/v4"
)

// CategoryHandler handles HTTP requests for categories
type CategoryHandler struct {
	createCategoryUseCase *usecase.CreateCategoryUseCase
	listCategoriesUseCase *usecase.ListCategoriesUseCase
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(
	createCategoryUseCase *usecase.CreateCategoryUseCase,
	listCategoriesUseCase *usecase.ListCategoriesUseCase,
) *CategoryHandler {
	return &CategoryHandler{
		createCategoryUseCase: createCategoryUseCase,
		listCategoriesUseCase: listCategoriesUseCase,
	}
}

// CreateCategory handles category creation
// @Summary Create a new category
// @Tags Budget - Categories
// @Accept json
// @Produce json
// @Param request body dto.CreateCategoryRequest true "Category creation request"
// @Success 201 {object} dto.CategoryResponse
// @Failure 400 {object} map[string]interface{}
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	var req dto.CreateCategoryRequest
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

	input := usecase.CreateCategoryInput{
		UserID:      userID,
		Name:        req.Name,
		Type:        domain.CategoryType(req.Type),
		Icon:        req.Icon,
		Color:       req.Color,
		Description: req.Description,
	}

	category, err := h.createCategoryUseCase.Execute(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, dto.ToCategoryResponse(category))
}

// ListCategories handles listing user categories
// @Summary List user categories
// @Tags Budget - Categories
// @Produce json
// @Param type query string false "Filter by type (income or expense)"
// @Success 200 {array} dto.CategoryResponse
// @Router /categories [get]
func (h *CategoryHandler) ListCategories(c echo.Context) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": err.Error(),
		})
	}

	var categoryType *domain.CategoryType
	if typeParam := c.QueryParam("type"); typeParam != "" {
		ct := domain.CategoryType(typeParam)
		categoryType = &ct
	}

	categories, err := h.listCategoriesUseCase.Execute(c.Request().Context(), userID, categoryType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.ToCategoryResponseList(categories))
}
