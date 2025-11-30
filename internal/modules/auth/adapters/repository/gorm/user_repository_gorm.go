package gorm

import (
	gormpkg "gorm.io/gorm"

	"github.com/edalferes/monetics/internal/modules/auth/adapters/repository/mappers"
	"github.com/edalferes/monetics/internal/modules/auth/adapters/repository/models"
	"github.com/edalferes/monetics/internal/modules/auth/domain"
	"github.com/edalferes/monetics/internal/modules/auth/usecase/interfaces"
)

type UserRepositoryGorm struct {
	DB     *gormpkg.DB
	mapper mappers.UserMapper
}

func (r *UserRepositoryGorm) FindByID(id uint) (*domain.User, error) {
	var userModel models.UserModel
	if err := r.DB.Preload("Roles.Permissions").Preload("Roles").First(&userModel, id).Error; err != nil {
		return nil, err
	}
	user := r.mapper.ToDomain(userModel)
	return &user, nil
}

func (r *UserRepositoryGorm) Update(user *domain.User) error {
	userModel := r.mapper.ToModel(*user)
	return r.DB.Session(&gormpkg.Session{FullSaveAssociations: true}).Updates(&userModel).Error
}

func (r *UserRepositoryGorm) Delete(id uint) error {
	return r.DB.Delete(&models.UserModel{}, id).Error
}

func NewUserRepositoryGorm(db *gormpkg.DB) *UserRepositoryGorm {
	return &UserRepositoryGorm{
		DB:     db,
		mapper: mappers.UserMapper{},
	}
}

var _ interfaces.User = (*UserRepositoryGorm)(nil)

func (r *UserRepositoryGorm) FindByUsername(username string) (*domain.User, error) {
	var userModel models.UserModel
	if err := r.DB.Preload("Roles.Permissions").Preload("Roles").Where("username = ?", username).First(&userModel).Error; err != nil {
		return nil, err
	}
	user := r.mapper.ToDomain(userModel)
	return &user, nil
}

func (r *UserRepositoryGorm) Create(user *domain.User) error {
	userModel := r.mapper.ToModel(*user)
	return r.DB.Create(&userModel).Error
}

func (r *UserRepositoryGorm) ListAll() ([]domain.User, error) {
	var userModels []models.UserModel
	if err := r.DB.Preload("Roles.Permissions").Preload("Roles").Find(&userModels).Error; err != nil {
		return nil, err
	}
	return r.mapper.ToDomainSlice(userModels), nil
}
