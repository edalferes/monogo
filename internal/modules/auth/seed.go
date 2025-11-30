package auth

import (
	"gorm.io/gorm"

	gormrepo "github.com/edalferes/monetics/internal/modules/auth/adapters/repository/gorm"
	"github.com/edalferes/monetics/internal/modules/auth/domain"
	"github.com/edalferes/monetics/internal/modules/auth/service"
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
	allPerms, err := permRepo.ListAll()
	if err != nil {
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

		// Check if role already exists
		_, err := roleRepo.FindByName(roleName)
		if err != nil {
			// Role doesn't exist, create it
			role := &domain.Role{Name: roleName, Permissions: permsToAssign}
			if err := roleRepo.Create(role); err != nil {
				return err
			}
		}
		// Note: For simplicity, we're not updating existing roles permissions
		// In a real app, you might want to implement an Update method for this
	}

	// Seed root user
	rootUsername := RootUsername
	rootPassword := RootPassword
	_, err = userRepo.FindByUsername(rootUsername)
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
