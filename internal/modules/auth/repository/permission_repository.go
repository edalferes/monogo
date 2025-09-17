package repository

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"gorm.io/gorm"
)

type PermissionRepository interface {
	FindByName(name string) (*domain.Permission, error)
	Create(permission *domain.Permission) error
	ListAll() ([]domain.Permission, error)
	DeleteByName(name string) error
	Update(permission *domain.Permission) error
}

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

func (r *permissionRepository) ListAll() ([]domain.Permission, error) {
	var perms []domain.Permission
	if err := r.DB.Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}

func (r *permissionRepository) DeleteByName(name string) error {
	return r.DB.Where("name = ?", name).Delete(&domain.Permission{}).Error
}

func (r *permissionRepository) Update(permission *domain.Permission) error {
	return r.DB.Save(permission).Error
}
