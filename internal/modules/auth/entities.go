package auth

import "github.com/edalferes/monogo/internal/modules/auth/adapters/repository/models"

func Entities() []interface{} {
	return []interface{}{
		&models.UserModel{},
		&models.RoleModel{},
		&models.PermissionModel{},
		&models.AuditLogModel{},
	}
}
