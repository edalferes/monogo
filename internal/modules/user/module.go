package user

import (
	"database/sql"
	"net/http"

	"github.com/edalferes/monogo/internal/modules/user/handler"
	"github.com/edalferes/monogo/internal/modules/user/repository"
	"github.com/edalferes/monogo/internal/modules/user/service"
	"github.com/edalferes/monogo/internal/modules/user/usecase"
)

func WireUp(mux *http.ServeMux, db *sql.DB) {
	repo := repository.NewUserRepository(db)
	createUserUC := usecase.NewCreateUserUseCase(repo)
	svc := &service.Service{CreateUserUC: createUserUC}
	h := &handler.Handler{Service: svc}
	h.RegisterRoutes(mux)
}
