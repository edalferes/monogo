package domain

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"-"`
	Roles    []Role `json:"roles" gorm:"many2many:user_roles;"`
}
