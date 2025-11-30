package gorm

import (
	gormpkg "gorm.io/gorm"

	"github.com/edalferes/monetics/internal/modules/auth/adapters/repository/mappers"
	"github.com/edalferes/monetics/internal/modules/auth/adapters/repository/models"
	"github.com/edalferes/monetics/internal/modules/auth/domain"
	"github.com/edalferes/monetics/internal/modules/auth/usecase/interfaces"
)

type PermissionRepositoryGorm struct {
	DB     *gormpkg.DB
	mapper mappers.PermissionMapper
}

func NewPermissionRepositoryGorm(db *gormpkg.DB) *PermissionRepositoryGorm {
	return &PermissionRepositoryGorm{
		DB:     db,
		mapper: mappers.PermissionMapper{},
	}
}

var _ interfaces.Permission = (*PermissionRepositoryGorm)(nil)

func (r *PermissionRepositoryGorm) FindByID(id uint) (*domain.Permission, error) {
	var permModel models.PermissionModel
	if err := r.DB.First(&permModel, id).Error; err != nil {
		return nil, err
	}
	perm := r.mapper.ToDomain(permModel)
	return &perm, nil
}

func (r *PermissionRepositoryGorm) FindByName(name string) (*domain.Permission, error) {
	var permissionModel models.PermissionModel
	if err := r.DB.Where("name = ?", name).First(&permissionModel).Error; err != nil {
		return nil, err
	}
	permission := r.mapper.ToDomain(permissionModel)
	return &permission, nil
}

func (r *PermissionRepositoryGorm) Create(permission *domain.Permission) error {
	permissionModel := r.mapper.ToModel(*permission)
	return r.DB.Create(&permissionModel).Error
}

func (r *PermissionRepositoryGorm) ListAll() ([]domain.Permission, error) {
	var permModels []models.PermissionModel
	if err := r.DB.Find(&permModels).Error; err != nil {
		return nil, err
	}
	return r.mapper.ToDomainSlice(permModels), nil
}

func (r *PermissionRepositoryGorm) DeleteByName(name string) error {
	return r.DB.Where("name = ?", name).Delete(&models.PermissionModel{}).Error
}

func (r *PermissionRepositoryGorm) Update(permission *domain.Permission) error {
	permissionModel := r.mapper.ToModel(*permission)
	return r.DB.Save(&permissionModel).Error
}
