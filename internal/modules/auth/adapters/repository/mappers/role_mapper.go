package mappers

import (
	"github.com/edalferes/monetics/internal/modules/auth/adapters/repository/models"
	"github.com/edalferes/monetics/internal/modules/auth/domain"
)

// RoleMapper converte entre domain.Role e models.RoleModel
type RoleMapper struct{}

// ToModel converte domain.Role para models.RoleModel
func (m RoleMapper) ToModel(role domain.Role) models.RoleModel {
	permissionModels := make([]models.PermissionModel, len(role.Permissions))
	for i, permission := range role.Permissions {
		permissionModels[i] = PermissionMapper{}.ToModel(permission)
	}

	return models.RoleModel{
		ID:          role.ID,
		Name:        role.Name,
		Permissions: permissionModels,
	}
}

// ToDomain converte models.RoleModel para domain.Role
func (m RoleMapper) ToDomain(roleModel models.RoleModel) domain.Role {
	permissions := make([]domain.Permission, len(roleModel.Permissions))
	for i, permissionModel := range roleModel.Permissions {
		permissions[i] = PermissionMapper{}.ToDomain(permissionModel)
	}

	return domain.Role{
		ID:          roleModel.ID,
		Name:        roleModel.Name,
		Permissions: permissions,
	}
}

// ToModelSlice converte []domain.Role para []models.RoleModel
func (m RoleMapper) ToModelSlice(roles []domain.Role) []models.RoleModel {
	roleModels := make([]models.RoleModel, len(roles))
	for i, role := range roles {
		roleModels[i] = m.ToModel(role)
	}
	return roleModels
}

// ToDomainSlice converte []models.RoleModel para []domain.Role
func (m RoleMapper) ToDomainSlice(roleModels []models.RoleModel) []domain.Role {
	roles := make([]domain.Role, len(roleModels))
	for i, roleModel := range roleModels {
		roles[i] = m.ToDomain(roleModel)
	}
	return roles
}
