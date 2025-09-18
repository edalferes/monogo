package auth

import (
	"github.com/edalferes/monogo/internal/modules/auth/domain"
	gormrepo "github.com/edalferes/monogo/internal/modules/auth/repository/gorm"
	"github.com/edalferes/monogo/internal/modules/auth/service"
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
		userRole, err := roleRepo.FindByName("user")
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
			Roles:    []domain.Role{*adminRole, *userRole},
		}
		if err := userRepo.Create(rootUser); err != nil {
			return err
		}
	}
	return nil
}
