package budget

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/edalferes/monetics/internal/modules/auth"
	"github.com/edalferes/monetics/internal/modules/auth/service"
	"github.com/edalferes/monetics/internal/modules/budget/adapters/http/handlers"
	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository"
	"github.com/edalferes/monetics/internal/modules/budget/usecase/account"
	budgetUseCase "github.com/edalferes/monetics/internal/modules/budget/usecase/budget"
	"github.com/edalferes/monetics/internal/modules/budget/usecase/category"
	"github.com/edalferes/monetics/internal/modules/budget/usecase/report"
	"github.com/edalferes/monetics/internal/modules/budget/usecase/transaction"
	"github.com/edalferes/monetics/pkg/logger"
)

// Module represents the budget module
type Module struct {
	db *gorm.DB

	// Repositories
	accountRepo     repository.AccountRepository
	categoryRepo    repository.CategoryRepository
	transactionRepo repository.TransactionRepository
	budgetRepo      repository.BudgetRepository

	// Use cases - Account
	createAccountUseCase  *account.CreateUseCase
	listAccountsUseCase   *account.ListUseCase
	getAccountByIDUseCase *account.GetByIDUseCase
	updateAccountUseCase  *account.UpdateUseCase
	deleteAccountUseCase  *account.DeleteUseCase

	// Use cases - Category
	createCategoryUseCase  *category.CreateUseCase
	listCategoriesUseCase  *category.ListUseCase
	getCategoryByIDUseCase *category.GetByIDUseCase
	updateCategoryUseCase  *category.UpdateUseCase
	deleteCategoryUseCase  *category.DeleteUseCase

	// Use cases - Transaction
	createTransactionUseCase  *transaction.CreateUseCase
	listTransactionsUseCase   *transaction.ListUseCase
	getTransactionByIDUseCase *transaction.GetByIDUseCase
	updateTransactionUseCase  *transaction.UpdateUseCase
	deleteTransactionUseCase  *transaction.DeleteUseCase

	// Use cases - Budget
	createBudgetUseCase  *budgetUseCase.CreateUseCase
	listBudgetsUseCase   *budgetUseCase.ListUseCase
	getBudgetByIDUseCase *budgetUseCase.GetByIDUseCase
	updateBudgetUseCase  *budgetUseCase.UpdateUseCase
	deleteBudgetUseCase  *budgetUseCase.DeleteUseCase

	// Use cases - Report
	getAccountBalanceUseCase *report.GetAccountBalanceUseCase
	getMonthlyReportUseCase  *report.GetMonthlyReportUseCase

	// Handlers
	accountHandler     *handlers.AccountHandler
	categoryHandler    *handlers.CategoryHandler
	transactionHandler *handlers.TransactionHandler
	budgetHandler      *handlers.BudgetHandler
	reportHandler      *handlers.ReportHandler
}

// NewModule creates a new budget module instance
func NewModule(db *gorm.DB) *Module {
	module := &Module{
		db: db,
	}

	// Initialize repositories
	module.accountRepo = repository.NewGormAccountRepository(db)
	module.categoryRepo = repository.NewGormCategoryRepository(db)
	module.transactionRepo = repository.NewGormTransactionRepository(db)
	module.budgetRepo = repository.NewGormBudgetRepository(db)

	// Initialize use cases - Account
	module.createAccountUseCase = account.NewCreateUseCase(module.accountRepo)
	module.listAccountsUseCase = account.NewListUseCase(module.accountRepo)
	module.getAccountByIDUseCase = account.NewGetByIDUseCase(module.accountRepo)
	module.updateAccountUseCase = account.NewUpdateUseCase(module.accountRepo)
	module.deleteAccountUseCase = account.NewDeleteUseCase(module.accountRepo)

	// Initialize use cases - Category
	module.createCategoryUseCase = category.NewCreateUseCase(module.categoryRepo)
	module.listCategoriesUseCase = category.NewListUseCase(module.categoryRepo)
	module.getCategoryByIDUseCase = category.NewGetByIDUseCase(module.categoryRepo)
	module.updateCategoryUseCase = category.NewUpdateUseCase(module.categoryRepo)
	module.deleteCategoryUseCase = category.NewDeleteUseCase(module.categoryRepo)

	// Initialize use cases - Transaction
	module.createTransactionUseCase = transaction.NewCreateUseCase(
		module.transactionRepo,
		module.accountRepo,
		module.categoryRepo,
	)
	module.listTransactionsUseCase = transaction.NewListUseCase(module.transactionRepo)
	module.getTransactionByIDUseCase = transaction.NewGetByIDUseCase(module.transactionRepo)
	module.updateTransactionUseCase = transaction.NewUpdateUseCase(
		module.transactionRepo,
		module.accountRepo,
		module.categoryRepo,
	)
	module.deleteTransactionUseCase = transaction.NewDeleteUseCase(module.transactionRepo)

	// Initialize use cases - Budget
	module.createBudgetUseCase = budgetUseCase.NewCreateUseCase(
		module.budgetRepo,
		module.categoryRepo,
	)
	module.listBudgetsUseCase = budgetUseCase.NewListUseCase(module.budgetRepo)
	module.getBudgetByIDUseCase = budgetUseCase.NewGetByIDUseCase(module.budgetRepo)
	module.updateBudgetUseCase = budgetUseCase.NewUpdateUseCase(module.budgetRepo)
	module.deleteBudgetUseCase = budgetUseCase.NewDeleteUseCase(module.budgetRepo)

	// Initialize use cases - Report
	module.getAccountBalanceUseCase = report.NewGetAccountBalanceUseCase(
		module.accountRepo,
		module.transactionRepo,
	)
	module.getMonthlyReportUseCase = report.NewGetMonthlyReportUseCase(
		module.transactionRepo,
		module.budgetRepo,
		module.categoryRepo,
	)

	// Initialize handlers
	module.accountHandler = handlers.NewAccountHandler(
		module.createAccountUseCase,
		module.listAccountsUseCase,
		module.getAccountByIDUseCase,
		module.updateAccountUseCase,
		module.deleteAccountUseCase,
		module.getAccountBalanceUseCase,
	)
	module.categoryHandler = handlers.NewCategoryHandler(
		module.createCategoryUseCase,
		module.listCategoriesUseCase,
		module.getCategoryByIDUseCase,
		module.updateCategoryUseCase,
		module.deleteCategoryUseCase,
	)
	module.transactionHandler = handlers.NewTransactionHandler(
		module.createTransactionUseCase,
		module.listTransactionsUseCase,
		module.getTransactionByIDUseCase,
		module.updateTransactionUseCase,
		module.deleteTransactionUseCase,
	)
	module.budgetHandler = handlers.NewBudgetHandler(
		module.createBudgetUseCase,
		module.listBudgetsUseCase,
		module.getBudgetByIDUseCase,
		module.updateBudgetUseCase,
		module.deleteBudgetUseCase,
	)
	module.reportHandler = handlers.NewReportHandler(
		module.getMonthlyReportUseCase,
	)

	return module
}

// RegisterRoutes registers all budget module routes
func (m *Module) RegisterRoutes(api *echo.Group, authMiddleware echo.MiddlewareFunc) {
	budget := api.Group("/budget")
	budget.Use(authMiddleware)

	// Account routes
	accounts := budget.Group("/accounts")
	accounts.POST("", m.accountHandler.CreateAccount)
	accounts.GET("", m.accountHandler.ListAccounts)
	accounts.GET("/:id", m.accountHandler.GetAccount)
	accounts.GET("/:id/detail", m.accountHandler.GetAccountByID)
	accounts.PUT("/:id", m.accountHandler.UpdateAccount)
	accounts.DELETE("/:id", m.accountHandler.DeleteAccount)

	// Category routes
	categories := budget.Group("/categories")
	categories.POST("", m.categoryHandler.CreateCategory)
	categories.GET("", m.categoryHandler.ListCategories)
	categories.GET("/:id", m.categoryHandler.GetCategoryByID)
	categories.PUT("/:id", m.categoryHandler.UpdateCategory)
	categories.DELETE("/:id", m.categoryHandler.DeleteCategory)

	// Transaction routes
	transactions := budget.Group("/transactions")
	transactions.POST("", m.transactionHandler.CreateTransaction)
	transactions.GET("", m.transactionHandler.ListTransactions)
	transactions.GET("/:id", m.transactionHandler.GetTransactionByID)
	transactions.PUT("/:id", m.transactionHandler.UpdateTransaction)
	transactions.DELETE("/:id", m.transactionHandler.DeleteTransaction)

	// Budget routes
	budgets := budget.Group("/budgets")
	budgets.POST("", m.budgetHandler.CreateBudget)
	budgets.GET("", m.budgetHandler.ListBudgets)
	budgets.GET("/:id", m.budgetHandler.GetBudgetByID)
	budgets.PUT("/:id", m.budgetHandler.UpdateBudget)
	budgets.DELETE("/:id", m.budgetHandler.DeleteBudget)

	// Report routes
	reports := budget.Group("/reports")
	reports.GET("/monthly", m.reportHandler.GetMonthlyReport)
}

// WireUp initializes and registers the budget module
func WireUp(group *echo.Group, db *gorm.DB, jwtSecret string, log logger.Logger) {
	log.Info().Msg("Initializing budget module")

	module := NewModule(db)
	authMiddleware := auth.JWTMiddleware(jwtSecret)
	module.RegisterRoutes(group, authMiddleware)

	log.Info().Msg("Budget module initialized")
}

// WireUpWithHTTP inicializa budget usando HTTP para comunicar com auth remoto
func WireUpWithHTTP(group *echo.Group, db *gorm.DB, jwtSecret string, log logger.Logger, authURL string) {
	log.Info().Str("auth_url", authURL).Msg("Initializing budget module with HTTP auth service")

	// Cria UserService HTTP para se comunicar com auth remoto
	userService := service.NewUserServiceHTTP(authURL)

	// Currently not using userService directly, but it's available
	_ = userService

	module := NewModule(db)
	authMiddleware := auth.JWTMiddleware(jwtSecret)
	module.RegisterRoutes(group, authMiddleware)

	log.Info().Msg("Budget module initialized with HTTP communication")
}
