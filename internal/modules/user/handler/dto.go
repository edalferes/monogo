package handler

import "github.com/go-playground/validator/v10"

var validate = validator.New()

type CreateUserDTO struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}
