package models

// PermissionModel representa a estrutura de dados para GORM
type PermissionModel struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
}

func (PermissionModel) TableName() string {
	return "permissions"
}
