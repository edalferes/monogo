package mappers

import (
	"github.com/edalferes/monogo/internal/modules/auth/adapters/repository/models"
	"github.com/edalferes/monogo/internal/modules/auth/domain"
)

// PermissionMapper converte entre domain.Permission e models.PermissionModel
type PermissionMapper struct{}

// ToModel converte domain.Permission para models.PermissionModel
func (m PermissionMapper) ToModel(permission domain.Permission) models.PermissionModel {
	return models.PermissionModel{
		ID:   permission.ID,
		Name: permission.Name,
	}
}

// ToDomain converte models.PermissionModel para domain.Permission
func (m PermissionMapper) ToDomain(permissionModel models.PermissionModel) domain.Permission {
	return domain.Permission{
		ID:   permissionModel.ID,
		Name: permissionModel.Name,
	}
}

// ToModelSlice converte []domain.Permission para []models.PermissionModel
func (m PermissionMapper) ToModelSlice(permissions []domain.Permission) []models.PermissionModel {
	permissionModels := make([]models.PermissionModel, len(permissions))
	for i, permission := range permissions {
		permissionModels[i] = m.ToModel(permission)
	}
	return permissionModels
}

// ToDomainSlice converte []models.PermissionModel para []domain.Permission
func (m PermissionMapper) ToDomainSlice(permissionModels []models.PermissionModel) []domain.Permission {
	permissions := make([]domain.Permission, len(permissionModels))
	for i, permissionModel := range permissionModels {
		permissions[i] = m.ToDomain(permissionModel)
	}
	return permissions
}
