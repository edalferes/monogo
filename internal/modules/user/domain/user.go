package domain

import (
	"github.com/google/uuid"
)

type User struct {
	ID    uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name  string    `json:"name" gorm:"size:255;not null"`
	Email string    `json:"email" gorm:"size:255;not null;unique"`
}
