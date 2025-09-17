package gorm

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"github.com/edalferes/monogo/internal/modules/auth/repository"
	gormpkg "gorm.io/gorm"
)

type RoleRepositoryGorm struct {
	DB *gormpkg.DB
}

func NewRoleRepositoryGorm(db *gormpkg.DB) *RoleRepositoryGorm {
	return &RoleRepositoryGorm{DB: db}
}

var _ repository.RoleRepository = (*RoleRepositoryGorm)(nil)

func (r *RoleRepositoryGorm) FindByName(name string) (*domain.Role, error) {
	var role domain.Role
	if err := r.DB.Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepositoryGorm) Create(role *domain.Role) error {
	return r.DB.Create(role).Error
}

func (r *RoleRepositoryGorm) ListAll() ([]domain.Role, error) {
	var roles []domain.Role
	if err := r.DB.Preload("Permissions").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *RoleRepositoryGorm) DeleteByName(name string) error {
	return r.DB.Where("name = ?", name).Delete(&domain.Role{}).Error
}

func (r *RoleRepositoryGorm) Update(role *domain.Role) error {
	return r.DB.Save(role).Error
}
