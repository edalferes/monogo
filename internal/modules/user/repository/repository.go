package repository

import (
	"github.com/edalferes/monogo/internal/modules/user/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(name string, email string) (uuid.UUID, error)
	FindByID(id uuid.UUID) (*domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(name string, email string) (uuid.UUID, error) {
	user := &domain.User{
		ID:    uuid.New(),
		Name:  name,
		Email: email,
	}
	if err := r.db.Create(user).Error; err != nil {
		return uuid.Nil, err
	}
	return user.ID, nil
}

func (r *userRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

var _ UserRepository = (*userRepository)(nil)
