package repository

import (
	"gorm.io/gorm"

	"github.com/edalferes/monetics/internal/modules/auth/domain"
	"github.com/edalferes/monetics/internal/modules/auth/usecase/interfaces"
)

type RoleRepository struct {
	DB *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{
		DB: db,
	}
}

var _ interfaces.Role = (*RoleRepository)(nil)

func (r *RoleRepository) FindByID(id uint) (*domain.Role, error) {
	var role domain.Role
	if err := r.DB.Preload("Permissions").First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) FindByName(name string) (*domain.Role, error) {
	var role domain.Role
	if err := r.DB.Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) Create(role *domain.Role) error {
	return r.DB.Create(role).Error
}

func (r *RoleRepository) ListAll() ([]domain.Role, error) {
	var roles []domain.Role
	if err := r.DB.Preload("Permissions").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *RoleRepository) DeleteByName(name string) error {
	return r.DB.Where("name = ?", name).Delete(&domain.Role{}).Error
}

func (r *RoleRepository) Update(role *domain.Role) error {
	return r.DB.Save(role).Error
}
