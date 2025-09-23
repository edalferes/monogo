package models

// UserModel representa a estrutura de dados para GORM
type UserModel struct {
	ID       uint        `gorm:"primaryKey"`
	Username string      `gorm:"unique;not null"`
	Password string      `gorm:"not null"`
	Roles    []RoleModel `gorm:"many2many:user_roles;"`
}

func (UserModel) TableName() string {
	return "users"
}
