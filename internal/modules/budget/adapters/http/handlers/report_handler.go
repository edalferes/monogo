package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/edalferes/monetics/internal/modules/budget/adapters/http/dto"
	"github.com/edalferes/monetics/internal/modules/budget/usecase/report"
)

// ReportHandler handles HTTP requests for reports
type ReportHandler struct {
	getMonthlyReportUseCase *report.GetMonthlyReportUseCase
}

// NewReportHandler creates a new report handler
func NewReportHandler(
	getMonthlyReportUseCase *report.GetMonthlyReportUseCase,
) *ReportHandler {
	return &ReportHandler{
		getMonthlyReportUseCase: getMonthlyReportUseCase,
	}
}

// GetMonthlyReport handles getting a monthly financial report
// @Summary Get monthly financial report
// @Tags Budget - Reports
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month (1-12)"
// @Success 200 {object} dto.MonthlyReportResponse
// @Failure 400 {object} map[string]interface{}
// @Router /reports/monthly [get]
func (h *ReportHandler) GetMonthlyReport(c echo.Context) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": err.Error(),
		})
	}

	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")

	if yearStr == "" || monthStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "year and month are required",
		})
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid year",
		})
	}

	monthInt, err := strconv.Atoi(monthStr)
	if err != nil || monthInt < 1 || monthInt > 12 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid month",
		})
	}

	month := time.Month(monthInt)

	report, err := h.getMonthlyReportUseCase.Execute(c.Request().Context(), userID, year, month)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.ToMonthlyReportResponse(report))
}
