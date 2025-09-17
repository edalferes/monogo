package gorm

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"github.com/edalferes/monogo/internal/modules/auth/repository"
	gormpkg "gorm.io/gorm"
)

type UserRepositoryGorm struct {
	DB *gormpkg.DB
}

func NewUserRepositoryGorm(db *gormpkg.DB) *UserRepositoryGorm {
	return &UserRepositoryGorm{DB: db}
}

var _ repository.UserRepository = (*UserRepositoryGorm)(nil)

func (r *UserRepositoryGorm) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	if err := r.DB.Preload("Roles").Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryGorm) Create(user *domain.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepositoryGorm) ListAll() ([]domain.User, error) {
	var users []domain.User
	if err := r.DB.Preload("Roles").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
