package auth

import (
	"time"

	handler_admin "github.com/edalferes/monogo/internal/modules/auth/adapters/http/handlers/admin"
	"github.com/edalferes/monogo/internal/modules/auth/adapters/http/handlers/login"
	handler_user "github.com/edalferes/monogo/internal/modules/auth/adapters/http/handlers/user"
	gormrepo "github.com/edalferes/monogo/internal/modules/auth/adapters/repository/gorm"
	"github.com/edalferes/monogo/internal/modules/auth/service"
	permUC "github.com/edalferes/monogo/internal/modules/auth/usecase/permission"
	roleUC "github.com/edalferes/monogo/internal/modules/auth/usecase/role"
	userUC "github.com/edalferes/monogo/internal/modules/auth/usecase/user"
	"github.com/edalferes/monogo/pkg/logger"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func WireUp(group *echo.Group, db *gorm.DB, jwtSecret string, log logger.Logger) {

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
		Logger: log,
	}
	group.POST("/auth/login", publicHandler.Login)

	// User use cases
	listUsersUC := &userUC.ListUsersUseCase{UserRepo: userRepo}
	createUserUC := &userUC.RegisterUseCase{
		User:            userRepo,
		Role:            roleRepo,
		PasswordService: passwordService,
	}
	getUserByIDUC := &userUC.GetUserByIDUseCase{UserRepo: userRepo}
	updateUserUC := &userUC.UpdateUserByAdminUseCase{
		UserRepo:        userRepo,
		PasswordService: passwordService,
	}
	deleteUserUC := &userUC.DeleteUserUseCase{UserRepo: userRepo}
	assignRoleUC := &userUC.AssignRoleUseCase{
		UserRepo: userRepo,
		RoleRepo: roleRepo,
	}
	removeRoleUC := &userUC.RemoveRoleUseCase{
		UserRepo: userRepo,
		RoleRepo: roleRepo,
	}
	changePasswordUC := &userUC.ChangePasswordUseCase{
		UserRepo:        userRepo,
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
		ListUsersUC:   listUsersUC,
		CreateUserUC:  createUserUC,
		GetUserByIDUC: getUserByIDUC,
		UpdateUserUC:  updateUserUC,
		DeleteUserUC:  deleteUserUC,
		AssignRoleUC:  assignRoleUC,
		RemoveRoleUC:  removeRoleUC,
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
	adminGroup.GET("/users/:id", adminUserHandler.GetUser)
	adminGroup.PUT("/users/:id", adminUserHandler.UpdateUser)
	adminGroup.DELETE("/users/:id", adminUserHandler.DeleteUser)
	adminGroup.POST("/users/:id/roles", adminUserHandler.AssignRoleToUser)
	adminGroup.DELETE("/users/:id/roles/:roleName", adminUserHandler.RemoveRoleFromUser)
	// Roles
	adminGroup.GET("/roles", adminRolePermHandler.ListRoles)
	adminGroup.POST("/roles", adminRolePermHandler.CreateRole)
	adminGroup.DELETE("/roles/:name", adminRolePermHandler.DeleteRole)
	// Permissions
	adminGroup.GET("/permissions", adminRolePermHandler.ListPermissions)
	adminGroup.POST("/permissions", adminRolePermHandler.CreatePermission)
	adminGroup.DELETE("/permissions/:name", adminRolePermHandler.DeletePermission)

	// User endpoints (protected, user can access their own data)
	userHandler := &handler_user.UserHandler{
		ChangePasswordUC: changePasswordUC,
	}
	userGroup := group.Group("/user")
	userGroup.Use(JWTMiddleware(jwtSecret))
	userGroup.PUT("/password", userHandler.ChangePassword)
}
