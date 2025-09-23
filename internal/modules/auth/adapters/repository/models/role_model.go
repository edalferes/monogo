package models

// RoleModel representa a estrutura de dados para GORM
type RoleModel struct {
	ID          uint              `gorm:"primaryKey"`
	Name        string            `gorm:"unique;not null"`
	Permissions []PermissionModel `gorm:"many2many:role_permissions;"`
}

func (RoleModel) TableName() string {
	return "roles"
}
