package repository

import (
	"gorm.io/gorm"

	"github.com/edalferes/monetics/internal/modules/auth/domain"
	"github.com/edalferes/monetics/internal/modules/auth/usecase/interfaces"
)

type PermissionRepository struct {
	DB *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{
		DB: db,
	}
}

var _ interfaces.Permission = (*PermissionRepository)(nil)

func (r *PermissionRepository) FindByID(id uint) (*domain.Permission, error) {
	var perm domain.Permission
	if err := r.DB.First(&perm, id).Error; err != nil {
		return nil, err
	}
	return &perm, nil
}

func (r *PermissionRepository) FindByName(name string) (*domain.Permission, error) {
	var permission domain.Permission
	if err := r.DB.Where("name = ?", name).First(&permission).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *PermissionRepository) Create(permission *domain.Permission) error {
	return r.DB.Create(permission).Error
}

func (r *PermissionRepository) ListAll() ([]domain.Permission, error) {
	var permissions []domain.Permission
	if err := r.DB.Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *PermissionRepository) DeleteByName(name string) error {
	return r.DB.Where("name = ?", name).Delete(&domain.Permission{}).Error
}

func (r *PermissionRepository) Update(permission *domain.Permission) error {
	return r.DB.Save(permission).Error
}
