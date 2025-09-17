package repository

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"gorm.io/gorm"
)

type permissionRepository struct {
	DB *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{DB: db}
}

func (r *permissionRepository) FindByName(name string) (*domain.Permission, error) {
	var permission domain.Permission
	if err := r.DB.Where("name = ?", name).First(&permission).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) Create(permission *domain.Permission) error {
	return r.DB.Create(permission).Error
}

type PermissionRepository interface {
	FindByName(name string) (*domain.Permission, error)
	Create(permission *domain.Permission) error
}
