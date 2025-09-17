package auth

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	"github.com/edalferes/monogo/internal/modules/auth/handler"
	handler_admin "github.com/edalferes/monogo/internal/modules/auth/handler/admin"
	gormrepo "github.com/edalferes/monogo/internal/modules/auth/repository/gorm"
	"github.com/edalferes/monogo/internal/modules/auth/service"
	"github.com/edalferes/monogo/internal/modules/auth/usecase"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

const (
	RootUsername string = "root"
	RootPassword string = "ZDcxMDUxZmM4M2Jl"
)

func Seed(db *gorm.DB) error {
	roleRepo := gormrepo.NewRoleRepositoryGorm(db)
	permRepo := gormrepo.NewPermissionRepositoryGorm(db)
	userRepo := gormrepo.NewUserRepositoryGorm(db)
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

	// Seed roles and assign permissions
	var allPerms []domain.Permission
	if err := db.Find(&allPerms).Error; err != nil {
		return err
	}
	for _, roleName := range defaultRoles {
		var permsToAssign []domain.Permission
		if roleName == "admin" {
			permsToAssign = allPerms // admin allows all permissions
		} else {
			// user only read
			for _, p := range allPerms {
				if p.Name == "read" {
					permsToAssign = append(permsToAssign, p)
				}
			}
		}
		var role domain.Role
		if err := db.Where("name = ?", roleName).First(&role).Error; err != nil {
			role = domain.Role{Name: roleName, Permissions: permsToAssign}
			if err := db.Create(&role).Error; err != nil {
				return err
			}
		} else {
			// Update permissions if role already exists
			db.Model(&role).Association("Permissions").Replace(permsToAssign)
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

	userRepo := gormrepo.NewUserRepositoryGorm(db)
	roleRepo := gormrepo.NewRoleRepositoryGorm(db)
	permRepo := gormrepo.NewPermissionRepositoryGorm(db)
	passwordService := service.NewPasswordService()
	jwtService := service.NewJWTService(jwtSecret, 24*60*60) // 24h

	// Handler public (only login)
	publicHandler := &handler.Handler{
		LoginUseCase: &usecase.LoginUseCase{
			UserRepo:        userRepo,
			PasswordService: passwordService,
			JWTService:      jwtService,
		},
	}
	group.POST("/auth/login", publicHandler.Login)

	// Handler admin (user management, roles, permissions)
	adminUserHandler := &handler_admin.AdminUserHandler{
		UserRepo:        userRepo,
		RoleRepo:        roleRepo,
		PasswordService: passwordService,
	}
	adminRolePermHandler := &handler_admin.AdminHandler{
		RoleRepo:       roleRepo,
		PermissionRepo: permRepo,
	}
	adminGroup := group.Group("/admin")
	adminGroup.Use(JWTMiddleware(jwtSecret))
	adminGroup.POST("/users", adminUserHandler.CreateUser)
	adminGroup.GET("/users", adminUserHandler.ListUsers)
	// Roles
	adminGroup.GET("/roles", adminRolePermHandler.ListRoles)
	adminGroup.POST("/roles", adminRolePermHandler.CreateRole)
	adminGroup.DELETE("/roles/:name", adminRolePermHandler.DeleteRole)
	// Permissions
	adminGroup.GET("/permissions", adminRolePermHandler.ListPermissions)
	adminGroup.POST("/permissions", adminRolePermHandler.CreatePermission)
	adminGroup.DELETE("/permissions/:name", adminRolePermHandler.DeletePermission)
}
