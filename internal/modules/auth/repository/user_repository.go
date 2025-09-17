package repository

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByUsername(username string) (*domain.User, error)
	Create(user *domain.User) error
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	if err := r.DB.Preload("Roles").Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *domain.User) error {
	return r.DB.Create(user).Error
}
