package budget

import (
	"github.com/edalferes/monetics/internal/modules/auth"
	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository"
	"github.com/edalferes/monetics/internal/modules/budget/handler"
	"github.com/edalferes/monetics/internal/modules/budget/usecase"
	"github.com/edalferes/monetics/pkg/logger"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Module represents the budget module
type Module struct {
	db *gorm.DB

	// Repositories
	accountRepo     repository.AccountRepository
	categoryRepo    repository.CategoryRepository
	transactionRepo repository.TransactionRepository
	budgetRepo      repository.BudgetRepository

	// Use cases
	createAccountUseCase     *usecase.CreateAccountUseCase
	listAccountsUseCase      *usecase.ListAccountsUseCase
	getAccountBalanceUseCase *usecase.GetAccountBalanceUseCase
	createCategoryUseCase    *usecase.CreateCategoryUseCase
	listCategoriesUseCase    *usecase.ListCategoriesUseCase
	createTransactionUseCase *usecase.CreateTransactionUseCase
	listTransactionsUseCase  *usecase.ListTransactionsUseCase
	createBudgetUseCase      *usecase.CreateBudgetUseCase
	listBudgetsUseCase       *usecase.ListBudgetsUseCase
	getMonthlyReportUseCase  *usecase.GetMonthlyReportUseCase

	// Handlers
	accountHandler     *handler.AccountHandler
	categoryHandler    *handler.CategoryHandler
	transactionHandler *handler.TransactionHandler
	budgetHandler      *handler.BudgetHandler
	reportHandler      *handler.ReportHandler
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

	// Initialize use cases
	module.createAccountUseCase = usecase.NewCreateAccountUseCase(module.accountRepo)
	module.listAccountsUseCase = usecase.NewListAccountsUseCase(module.accountRepo)
	module.getAccountBalanceUseCase = usecase.NewGetAccountBalanceUseCase(
		module.accountRepo,
		module.transactionRepo,
	)
	module.createCategoryUseCase = usecase.NewCreateCategoryUseCase(module.categoryRepo)
	module.listCategoriesUseCase = usecase.NewListCategoriesUseCase(module.categoryRepo)
	module.createTransactionUseCase = usecase.NewCreateTransactionUseCase(
		module.transactionRepo,
		module.accountRepo,
		module.categoryRepo,
	)
	module.listTransactionsUseCase = usecase.NewListTransactionsUseCase(module.transactionRepo)
	module.createBudgetUseCase = usecase.NewCreateBudgetUseCase(
		module.budgetRepo,
		module.categoryRepo,
	)
	module.listBudgetsUseCase = usecase.NewListBudgetsUseCase(module.budgetRepo)
	module.getMonthlyReportUseCase = usecase.NewGetMonthlyReportUseCase(
		module.transactionRepo,
		module.budgetRepo,
		module.categoryRepo,
	)

	// Initialize handlers
	module.accountHandler = handler.NewAccountHandler(
		module.createAccountUseCase,
		module.listAccountsUseCase,
		module.getAccountBalanceUseCase,
	)
	module.categoryHandler = handler.NewCategoryHandler(
		module.createCategoryUseCase,
		module.listCategoriesUseCase,
	)
	module.transactionHandler = handler.NewTransactionHandler(
		module.createTransactionUseCase,
		module.listTransactionsUseCase,
	)
	module.budgetHandler = handler.NewBudgetHandler(
		module.createBudgetUseCase,
		module.listBudgetsUseCase,
	)
	module.reportHandler = handler.NewReportHandler(
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

	// Category routes
	categories := budget.Group("/categories")
	categories.POST("", m.categoryHandler.CreateCategory)
	categories.GET("", m.categoryHandler.ListCategories)

	// Transaction routes
	transactions := budget.Group("/transactions")
	transactions.POST("", m.transactionHandler.CreateTransaction)
	transactions.GET("", m.transactionHandler.ListTransactions)

	// Budget routes
	budgets := budget.Group("/budgets")
	budgets.POST("", m.budgetHandler.CreateBudget)
	budgets.GET("", m.budgetHandler.ListBudgets)

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

	log.Info().Msg("Budget module initialized successfully")
}
