package dto

import (
	"github.com/edalferes/monetics/internal/modules/budget/usecase"
)

// MonthlyReportResponse represents a monthly financial report
type MonthlyReportResponse struct {
	Year              int                     `json:"year"`
	Month             int                     `json:"month"`
	MonthName         string                  `json:"month_name"`
	TotalIncome       float64                 `json:"total_income"`
	TotalExpense      float64                 `json:"total_expense"`
	Balance           float64                 `json:"balance"`
	CategoryBreakdown []CategoryBreakdownItem `json:"category_breakdown"`
}

// CategoryBreakdownItem represents spending/income by category
type CategoryBreakdownItem struct {
	CategoryID   uint    `json:"category_id"`
	CategoryName string  `json:"category_name"`
	CategoryType string  `json:"category_type"`
	Amount       float64 `json:"amount"`
	Percentage   float64 `json:"percentage"`
}

// ToMonthlyReportResponse converts usecase.MonthlyReport to MonthlyReportResponse
func ToMonthlyReportResponse(report *usecase.MonthlyReport) MonthlyReportResponse {
	breakdown := make([]CategoryBreakdownItem, len(report.CategoryTotals))

	for i, total := range report.CategoryTotals {
		var percentage float64
		if total.CategoryType == "income" && report.TotalIncome > 0 {
			percentage = (total.Amount / report.TotalIncome) * 100
		} else if total.CategoryType == "expense" && report.TotalExpense > 0 {
			percentage = (total.Amount / report.TotalExpense) * 100
		}

		breakdown[i] = CategoryBreakdownItem{
			CategoryID:   total.CategoryID,
			CategoryName: total.CategoryName,
			CategoryType: string(total.CategoryType),
			Amount:       total.Amount,
			Percentage:   percentage,
		}
	}

	return MonthlyReportResponse{
		Year:              report.Year,
		Month:             int(report.Month),
		MonthName:         report.Month.String(),
		TotalIncome:       report.TotalIncome,
		TotalExpense:      report.TotalExpense,
		Balance:           report.Balance,
		CategoryBreakdown: breakdown,
	}
}
