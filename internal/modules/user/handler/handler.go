package handler

import (
	"encoding/json"
	"net/http"

	"github.com/edalferes/monogo/internal/modules/user/service"
)

type Handler struct {
	Service *service.Service
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/users", h.CreateUser)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a user with name and email
// @Tags users
// @Accept json
// @Produce json
// @Param user body handler.CreateUserInput true "User data"
// @Success 201 {string} string "created"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post]
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, err := h.Service.Register(input.Name, input.Email); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
