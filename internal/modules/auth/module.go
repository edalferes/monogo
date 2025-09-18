package auth

import (
	"time"

	"github.com/edalferes/monogo/internal/modules/auth/handler"
	handler_admin "github.com/edalferes/monogo/internal/modules/auth/handler/admin"
	gormrepo "github.com/edalferes/monogo/internal/modules/auth/repository/gorm"
	"github.com/edalferes/monogo/internal/modules/auth/service"
	"github.com/edalferes/monogo/internal/modules/auth/usecase"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func WireUp(group *echo.Group, db *gorm.DB, jwtSecret string) {

	userRepo := gormrepo.NewUserRepositoryGorm(db)
	roleRepo := gormrepo.NewRoleRepositoryGorm(db)
	permRepo := gormrepo.NewPermissionRepositoryGorm(db)
	passwordService := service.NewPasswordService()
	jwtService := service.NewJWTService(jwtSecret, time.Hour) // 1 hour expiration token
	auditLogRepo := gormrepo.NewAuditLogRepositoryGorm(db)
	auditService := service.NewAuditService(auditLogRepo)

	// Handler public (only login)
	publicHandler := &handler.Handler{
		LoginUseCase: &usecase.LoginWithAuditUseCase{
			UserRepo:        userRepo,
			PasswordService: passwordService,
			JWTService:      jwtService,
			AuditService:    auditService,
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
