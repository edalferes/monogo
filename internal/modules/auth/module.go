package auth

import (
	"time"

	handler_admin "github.com/edalferes/monogo/internal/modules/auth/adapters/http/handlers/admin"
	"github.com/edalferes/monogo/internal/modules/auth/adapters/http/handlers/login"
	gormrepo "github.com/edalferes/monogo/internal/modules/auth/adapters/repository/gorm"
	"github.com/edalferes/monogo/internal/modules/auth/service"
	permUC "github.com/edalferes/monogo/internal/modules/auth/usecase/permission"
	roleUC "github.com/edalferes/monogo/internal/modules/auth/usecase/role"
	userUC "github.com/edalferes/monogo/internal/modules/auth/usecase/user"
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
	publicHandler := &login.Handler{
		LoginUseCase: &userUC.LoginWithAuditUseCase{
			UserRepo:        userRepo,
			PasswordService: passwordService,
			JWTService:      jwtService,
			AuditService:    auditService,
		},
	}
	group.POST("/auth/login", publicHandler.Login)

	// Use cases para usuário
	listUsersUC := &userUC.ListUsersUseCase{UserRepo: userRepo}
	createUserUC := &userUC.RegisterUseCase{
		UserReader:      userRepo,
		UserWriter:      userRepo,
		RoleReader:      roleRepo,
		PasswordService: passwordService,
	}
	// Use cases para role
	listRolesUC := &roleUC.ListRolesUseCase{RoleRepo: roleRepo}
	createRoleUC := &roleUC.CreateRoleUseCase{RoleRepo: roleRepo}
	deleteRoleUC := &roleUC.DeleteRoleUseCase{RoleRepo: roleRepo}
	// Use cases para permission
	listPermissionsUC := &permUC.ListPermissionsUseCase{PermissionRepo: permRepo}
	createPermissionUC := &permUC.CreatePermissionUseCase{PermissionRepo: permRepo}
	deletePermissionUC := &permUC.DeletePermissionUseCase{PermissionRepo: permRepo}

	adminUserHandler := &handler_admin.AdminUserHandler{
		ListUsersUC:  listUsersUC,
		CreateUserUC: createUserUC,
		// Adicione outros use cases conforme necessário
	}
	adminRolePermHandler := &handler_admin.AdminHandler{
		ListRolesUC:        listRolesUC,
		CreateRoleUC:       createRoleUC,
		DeleteRoleUC:       deleteRoleUC,
		ListPermissionsUC:  listPermissionsUC,
		CreatePermissionUC: createPermissionUC,
		DeletePermissionUC: deletePermissionUC,
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
