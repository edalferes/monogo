package handler

type CreateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
