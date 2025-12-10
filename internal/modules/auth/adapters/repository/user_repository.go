package repository

import (
	"gorm.io/gorm"

	"github.com/edalferes/monetics/internal/modules/auth/domain"
	"github.com/edalferes/monetics/internal/modules/auth/usecase/interfaces"
)

type UserRepository struct {
	DB *gorm.DB
}

func (r *UserRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	if err := r.DB.Preload("Roles.Permissions").Preload("Roles").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *domain.User) error {
	return r.DB.Session(&gorm.Session{FullSaveAssociations: true}).Updates(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.DB.Delete(&domain.User{}, id).Error
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

var _ interfaces.User = (*UserRepository)(nil)

func (r *UserRepository) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	if err := r.DB.Preload("Roles.Permissions").Preload("Roles").Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *domain.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) ListAll() ([]domain.User, error) {
	var users []domain.User
	if err := r.DB.Preload("Roles.Permissions").Preload("Roles").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
