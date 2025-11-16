package gorm

import (
	"github.com/edalferes/monetics/internal/modules/auth/adapters/repository/mappers"
	"github.com/edalferes/monetics/internal/modules/auth/adapters/repository/models"
	"github.com/edalferes/monetics/internal/modules/auth/domain"
	"github.com/edalferes/monetics/internal/modules/auth/usecase/interfaces"
	gormpkg "gorm.io/gorm"
)

type RoleRepositoryGorm struct {
	DB     *gormpkg.DB
	mapper mappers.RoleMapper
}

func NewRoleRepositoryGorm(db *gormpkg.DB) *RoleRepositoryGorm {
	return &RoleRepositoryGorm{
		DB:     db,
		mapper: mappers.RoleMapper{},
	}
}

var _ interfaces.Role = (*RoleRepositoryGorm)(nil)

func (r *RoleRepositoryGorm) FindByID(id uint) (*domain.Role, error) {
	var roleModel models.RoleModel
	if err := r.DB.Preload("Permissions").First(&roleModel, id).Error; err != nil {
		return nil, err
	}
	role := r.mapper.ToDomain(roleModel)
	return &role, nil
}

func (r *RoleRepositoryGorm) FindByName(name string) (*domain.Role, error) {
	var roleModel models.RoleModel
	if err := r.DB.Where("name = ?", name).First(&roleModel).Error; err != nil {
		return nil, err
	}
	role := r.mapper.ToDomain(roleModel)
	return &role, nil
}

func (r *RoleRepositoryGorm) Create(role *domain.Role) error {
	roleModel := r.mapper.ToModel(*role)
	return r.DB.Create(&roleModel).Error
}

func (r *RoleRepositoryGorm) ListAll() ([]domain.Role, error) {
	var roleModels []models.RoleModel
	if err := r.DB.Preload("Permissions").Find(&roleModels).Error; err != nil {
		return nil, err
	}
	return r.mapper.ToDomainSlice(roleModels), nil
}

func (r *RoleRepositoryGorm) DeleteByName(name string) error {
	return r.DB.Where("name = ?", name).Delete(&models.RoleModel{}).Error
}

func (r *RoleRepositoryGorm) Update(role *domain.Role) error {
	roleModel := r.mapper.ToModel(*role)
	return r.DB.Save(&roleModel).Error
}
