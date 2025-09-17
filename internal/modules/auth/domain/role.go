package domain

type Role struct {
	ID          uint         `json:"id" gorm:"primaryKey"`
	Name        string       `json:"name" gorm:"unique;not null"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
}
