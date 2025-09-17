package gorm

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"github.com/edalferes/monogo/internal/modules/auth/repository"
	gormpkg "gorm.io/gorm"
)

type PermissionRepositoryGorm struct {
	DB *gormpkg.DB
}

func NewPermissionRepositoryGorm(db *gormpkg.DB) *PermissionRepositoryGorm {
	return &PermissionRepositoryGorm{DB: db}
}

var _ repository.PermissionRepository = (*PermissionRepositoryGorm)(nil)

func (r *PermissionRepositoryGorm) FindByName(name string) (*domain.Permission, error) {
	var permission domain.Permission
	if err := r.DB.Where("name = ?", name).First(&permission).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *PermissionRepositoryGorm) Create(permission *domain.Permission) error {
	return r.DB.Create(permission).Error
}

func (r *PermissionRepositoryGorm) ListAll() ([]domain.Permission, error) {
	var perms []domain.Permission
	if err := r.DB.Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}

func (r *PermissionRepositoryGorm) DeleteByName(name string) error {
	return r.DB.Where("name = ?", name).Delete(&domain.Permission{}).Error
}

func (r *PermissionRepositoryGorm) Update(permission *domain.Permission) error {
	return r.DB.Save(permission).Error
}
