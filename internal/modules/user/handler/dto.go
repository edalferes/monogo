package handler

type CreateUserDTO struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
