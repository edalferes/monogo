package budget

import (
	"time"

	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"gorm.io/gorm"
)

// Seed populates the database with default budget categories
func Seed(db *gorm.DB, userID uint) error {
	categoryRepo := repository.NewGormCategoryRepository(db)

	// Default income categories based on the spreadsheet
	incomeCategories := []domain.Category{
		{UserID: userID, Name: "Salário", Type: domain.CategoryTypeIncome, Icon: "💰", Color: "#4CAF50"},
		{UserID: userID, Name: "Freelance", Type: domain.CategoryTypeIncome, Icon: "💼", Color: "#2196F3"},
		{UserID: userID, Name: "Aluguel de Imóvel Online", Type: domain.CategoryTypeIncome, Icon: "🏠", Color: "#009688"},
		{UserID: userID, Name: "Investimentos", Type: domain.CategoryTypeIncome, Icon: "📈", Color: "#FF9800"},
		{UserID: userID, Name: "Premiações", Type: domain.CategoryTypeIncome, Icon: "🏆", Color: "#FFC107"},
		{UserID: userID, Name: "Outras Fontes", Type: domain.CategoryTypeIncome, Icon: "💵", Color: "#8BC34A"},
	}

	// Default expense categories based on the spreadsheet
	expenseCategories := []domain.Category{
		// Moradia
		{UserID: userID, Name: "Aluguel", Type: domain.CategoryTypeExpense, Icon: "🏡", Color: "#F44336"},
		{UserID: userID, Name: "Condomínio", Type: domain.CategoryTypeExpense, Icon: "🏢", Color: "#E91E63"},
		{UserID: userID, Name: "Energia", Type: domain.CategoryTypeExpense, Icon: "⚡", Color: "#9C27B0"},
		{UserID: userID, Name: "Água", Type: domain.CategoryTypeExpense, Icon: "💧", Color: "#673AB7"},
		{UserID: userID, Name: "Internet", Type: domain.CategoryTypeExpense, Icon: "🌐", Color: "#3F51B5"},
		{UserID: userID, Name: "Gás", Type: domain.CategoryTypeExpense, Icon: "🔥", Color: "#FF5722"},
		{UserID: userID, Name: "IPTU", Type: domain.CategoryTypeExpense, Icon: "🏘️", Color: "#795548"},
		{UserID: userID, Name: "Manutenção", Type: domain.CategoryTypeExpense, Icon: "🔧", Color: "#607D8B"},

		// Food
		{UserID: userID, Name: "Mercado", Type: domain.CategoryTypeExpense, Icon: "🛒", Color: "#4CAF50"},
		{UserID: userID, Name: "Refeições Fora", Type: domain.CategoryTypeExpense, Icon: "🍽️", Color: "#8BC34A"},
		{UserID: userID, Name: "Lanches/Cafés", Type: domain.CategoryTypeExpense, Icon: "☕", Color: "#CDDC39"},
		{UserID: userID, Name: "Delivery", Type: domain.CategoryTypeExpense, Icon: "🚚", Color: "#FFEB3B"},

		// Transporte
		{UserID: userID, Name: "Combustível", Type: domain.CategoryTypeExpense, Icon: "⛽", Color: "#FF9800"},
		{UserID: userID, Name: "Uber/Táxi", Type: domain.CategoryTypeExpense, Icon: "🚕", Color: "#FF5722"},
		{UserID: userID, Name: "Transporte Público", Type: domain.CategoryTypeExpense, Icon: "🚌", Color: "#F44336"},
		{UserID: userID, Name: "Manutenção Veículo", Type: domain.CategoryTypeExpense, Icon: "🔧", Color: "#E91E63"},
		{UserID: userID, Name: "Seguro Auto", Type: domain.CategoryTypeExpense, Icon: "🚗", Color: "#9C27B0"},
		{UserID: userID, Name: "IPVA", Type: domain.CategoryTypeExpense, Icon: "🚙", Color: "#673AB7"},
		{UserID: userID, Name: "Estacionamento/Pedágios", Type: domain.CategoryTypeExpense, Icon: "🅿️", Color: "#3F51B5"},

		// Health
		{UserID: userID, Name: "Plano de Saúde", Type: domain.CategoryTypeExpense, Icon: "🏥", Color: "#2196F3"},
		{UserID: userID, Name: "Medicamentos", Type: domain.CategoryTypeExpense, Icon: "💊", Color: "#03A9F4"},
		{UserID: userID, Name: "Consultas/Exames", Type: domain.CategoryTypeExpense, Icon: "👨‍⚕️", Color: "#00BCD4"},
		{UserID: userID, Name: "Academia", Type: domain.CategoryTypeExpense, Icon: "💪", Color: "#009688"},
		{UserID: userID, Name: "Terapia/Psicólogo", Type: domain.CategoryTypeExpense, Icon: "🧠", Color: "#4CAF50"},

		// Education
		{UserID: userID, Name: "Cursos", Type: domain.CategoryTypeExpense, Icon: "📚", Color: "#8BC34A"},
		{UserID: userID, Name: "Livros/Material", Type: domain.CategoryTypeExpense, Icon: "📖", Color: "#CDDC39"},
		{UserID: userID, Name: "Assinaturas Educacionais", Type: domain.CategoryTypeExpense, Icon: "🎓", Color: "#FFEB3B"},
		{UserID: userID, Name: "Mensalidades/Escola", Type: domain.CategoryTypeExpense, Icon: "🏫", Color: "#FFC107"},

		// Lazer
		{UserID: userID, Name: "Streaming", Type: domain.CategoryTypeExpense, Icon: "📺", Color: "#FFC107"},
		{UserID: userID, Name: "Viagens/Passeios", Type: domain.CategoryTypeExpense, Icon: "✈️", Color: "#FF9800"},
		{UserID: userID, Name: "Hobbies", Type: domain.CategoryTypeExpense, Icon: "🎮", Color: "#FF5722"},
		{UserID: userID, Name: "Restaurantes", Type: domain.CategoryTypeExpense, Icon: "🍴", Color: "#F44336"},
		{UserID: userID, Name: "Cinema/Teatro", Type: domain.CategoryTypeExpense, Icon: "🎭", Color: "#E91E63"},

		// Pessoal
		{UserID: userID, Name: "Roupas", Type: domain.CategoryTypeExpense, Icon: "👔", Color: "#9C27B0"},
		{UserID: userID, Name: "Beleza/Estética", Type: domain.CategoryTypeExpense, Icon: "💄", Color: "#673AB7"},
		{UserID: userID, Name: "Presentes", Type: domain.CategoryTypeExpense, Icon: "🎁", Color: "#3F51B5"},
		{UserID: userID, Name: "Pets", Type: domain.CategoryTypeExpense, Icon: "🐾", Color: "#2196F3"},
	}

	// Check if categories already exist for this user
	existingCategories, err := categoryRepo.GetByUserID(db.Statement.Context, userID)
	if err != nil {
		return err
	}

	// If user already has categories, skip seeding
	if len(existingCategories) > 0 {
		return nil
	}

	// Create income categories
	for _, category := range incomeCategories {
		category.CreatedAt = time.Now()
		category.UpdatedAt = time.Now()
		if _, err := categoryRepo.Create(db.Statement.Context, category); err != nil {
			return err
		}
	}

	// Create expense categories
	for _, category := range expenseCategories {
		category.CreatedAt = time.Now()
		category.UpdatedAt = time.Now()
		if _, err := categoryRepo.Create(db.Statement.Context, category); err != nil {
			return err
		}
	}

	return nil
}
