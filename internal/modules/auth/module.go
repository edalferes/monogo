package auth

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"github.com/edalferes/monogo/internal/modules/auth/handler"
	"github.com/edalferes/monogo/internal/modules/auth/repository"
	"github.com/edalferes/monogo/internal/modules/auth/service"
	"github.com/edalferes/monogo/internal/modules/auth/usecase"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Seed garante que as roles padrão existam no banco
func Seed(db *gorm.DB) error {
	roleRepo := repository.NewRoleRepository(db)
	permRepo := repository.NewPermissionRepository(db)

	defaultRoles := []string{"admin", "user"}
	defaultPerms := []string{"read", "write", "delete"}

	// Seed permissions
	for _, permName := range defaultPerms {
		_, err := permRepo.FindByName(permName)
		if err != nil {
			if err := permRepo.Create(&domain.Permission{Name: permName}); err != nil {
				return err
			}
		}
	}

	// Seed roles
	for _, roleName := range defaultRoles {
		_, err := roleRepo.FindByName(roleName)
		if err != nil {
			if err := roleRepo.Create(&domain.Role{Name: roleName}); err != nil {
				return err
			}
		}
	}
	return nil
}

func WireUp(group *echo.Group, db *gorm.DB, jwtSecret string) {
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	passwordService := service.NewPasswordService()
	jwtService := service.NewJWTService(jwtSecret, 24*60*60) // 24h
	loginUseCase := &usecase.LoginUseCase{
		UserRepo:        userRepo,
		PasswordService: passwordService,
		JWTService:      jwtService,
	}
	registerUseCase := &usecase.RegisterUseCase{
		UserRepo:        userRepo,
		RoleRepo:        roleRepo,
		PasswordService: passwordService,
	}
	h := &handler.Handler{
		LoginUseCase:    loginUseCase,
		RegisterUseCase: registerUseCase,
	}
	group.POST("/auth/login", h.Login)
	group.POST("/auth/register", h.Register)
}
