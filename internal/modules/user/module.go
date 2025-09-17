package user

import (
	"github.com/edalferes/monogo/internal/modules/user/handler"
	"github.com/edalferes/monogo/internal/modules/user/repository"
	"github.com/edalferes/monogo/internal/modules/user/service"
	"github.com/edalferes/monogo/internal/modules/user/usecase"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// WireUpEcho registra as rotas do módulo user usando Echo
func WireUpEcho(g *echo.Group, db *gorm.DB) {
	repo := repository.NewUserRepository(db)
	createUserUC := usecase.NewCreateUserUseCase(repo)
	svc := &service.Service{CreateUserUC: createUserUC}
	h := &handler.Handler{Service: svc}
	h.RegisterRoutes(g)
}
