package repository

import (
	"context"

	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository/mappers"
	"github.com/edalferes/monetics/internal/modules/budget/adapters/repository/models"
	"github.com/edalferes/monetics/internal/modules/budget/domain"
	"gorm.io/gorm"
)

type gormAccountRepository struct {
	db     *gorm.DB
	mapper mappers.AccountMapper
}

// NewGormAccountRepository creates a new GORM-based account repository
func NewGormAccountRepository(db *gorm.DB) AccountRepository {
	return &gormAccountRepository{
		db:     db,
		mapper: mappers.AccountMapper{},
	}
}

func (r *gormAccountRepository) Create(ctx context.Context, account domain.Account) (domain.Account, error) {
	model := r.mapper.ToModel(account)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return domain.Account{}, err
	}
	return r.mapper.ToDomain(model), nil
}

func (r *gormAccountRepository) GetByID(ctx context.Context, id uint) (domain.Account, error) {
	var model models.AccountModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return domain.Account{}, err
	}
	return r.mapper.ToDomain(model), nil
}

func (r *gormAccountRepository) GetByUserID(ctx context.Context, userID uint) ([]domain.Account, error) {
	var accountModels []models.AccountModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&accountModels).Error; err != nil {
		return nil, err
	}
	return r.mapper.ToDomainSlice(accountModels), nil
}

func (r *gormAccountRepository) Update(ctx context.Context, account domain.Account) (domain.Account, error) {
	model := r.mapper.ToModel(account)
	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return domain.Account{}, err
	}
	return r.mapper.ToDomain(model), nil
}

func (r *gormAccountRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.AccountModel{}, id).Error
}

func (r *gormAccountRepository) ExistsByID(ctx context.Context, id uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.AccountModel{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}
