package auth

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"github.com/edalferes/monogo/internal/modules/auth/handler"
	handler_admin "github.com/edalferes/monogo/internal/modules/auth/handler/admin"
	"github.com/edalferes/monogo/internal/modules/auth/repository"
	"github.com/edalferes/monogo/internal/modules/auth/service"
	"github.com/edalferes/monogo/internal/modules/auth/usecase"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// const user and strong passwd
const (
	RootUsername string = "root"
	RootPassword string = "ZDcxMDUxZmM4M2Jl"
)

// Seed garante que as roles padrão existam no banco
func Seed(db *gorm.DB) error {
	roleRepo := repository.NewRoleRepository(db)
	permRepo := repository.NewPermissionRepository(db)
	userRepo := repository.NewUserRepository(db)
	passwordService := service.NewPasswordService()

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

	// Seed root user
	rootUsername := RootUsername
	rootPassword := RootPassword
	_, err := userRepo.FindByUsername(rootUsername)
	if err != nil {
		adminRole, err := roleRepo.FindByName("admin")
		if err != nil {
			return err
		}
		hash, err := passwordService.Hash(rootPassword)
		if err != nil {
			return err
		}
		rootUser := &domain.User{
			Username: rootUsername,
			Password: hash,
			Roles:    []domain.Role{*adminRole},
		}
		if err := userRepo.Create(rootUser); err != nil {
			return err
		}
	}
	return nil
}

func WireUp(group *echo.Group, db *gorm.DB, jwtSecret string) {
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	passwordService := service.NewPasswordService()
	jwtService := service.NewJWTService(jwtSecret, 24*60*60) // 24h

	// Handler para rotas públicas (apenas login)
	publicHandler := &handler.Handler{
		LoginUseCase: &usecase.LoginUseCase{
			UserRepo:        userRepo,
			PasswordService: passwordService,
			JWTService:      jwtService,
		},
	}
	group.POST("/auth/login", publicHandler.Login)

	// Handler para rotas administrativas
	adminUserHandler := &handler_admin.AdminUserHandler{
		UserRepo:        userRepo,
		RoleRepo:        roleRepo,
		PasswordService: passwordService,
	}
	adminGroup := group.Group("/admin")
	adminGroup.POST("/users", adminUserHandler.CreateUser)
}
