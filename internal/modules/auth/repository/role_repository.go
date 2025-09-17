package repository

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"gorm.io/gorm"
)

type RoleRepository interface {
	FindByName(name string) (*domain.Role, error)
	Create(role *domain.Role) error
	ListAll() ([]domain.Role, error)
	DeleteByName(name string) error
	Update(role *domain.Role) error
}

type roleRepository struct {
	DB *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{DB: db}
}

func (r *roleRepository) FindByName(name string) (*domain.Role, error) {
	var role domain.Role
	if err := r.DB.Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) Create(role *domain.Role) error {
	return r.DB.Create(role).Error
}

func (r *roleRepository) ListAll() ([]domain.Role, error) {
	var roles []domain.Role
	if err := r.DB.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *roleRepository) DeleteByName(name string) error {
	return r.DB.Where("name = ?", name).Delete(&domain.Role{}).Error
}

func (r *roleRepository) Update(role *domain.Role) error {
	return r.DB.Save(role).Error
}
