package mappers

import (
	"github.com/edalferes/monetics/internal/modules/auth/adapters/repository/models"
	"github.com/edalferes/monetics/internal/modules/auth/domain"
)

// UserMapper converte entre domain.User e models.UserModel
type UserMapper struct{}

// ToModel converte domain.User para models.UserModel
func (m UserMapper) ToModel(user domain.User) models.UserModel {
	roleModels := make([]models.RoleModel, len(user.Roles))
	for i, role := range user.Roles {
		roleModels[i] = RoleMapper{}.ToModel(role)
	}

	return models.UserModel{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Roles:    roleModels,
	}
}

// ToDomain converte models.UserModel para domain.User
func (m UserMapper) ToDomain(userModel models.UserModel) domain.User {
	roles := make([]domain.Role, len(userModel.Roles))
	for i, roleModel := range userModel.Roles {
		roles[i] = RoleMapper{}.ToDomain(roleModel)
	}

	return domain.User{
		ID:       userModel.ID,
		Username: userModel.Username,
		Password: userModel.Password,
		Roles:    roles,
	}
}

// ToModelSlice converte []domain.User para []models.UserModel
func (m UserMapper) ToModelSlice(users []domain.User) []models.UserModel {
	userModels := make([]models.UserModel, len(users))
	for i, user := range users {
		userModels[i] = m.ToModel(user)
	}
	return userModels
}

// ToDomainSlice converte []models.UserModel para []domain.User
func (m UserMapper) ToDomainSlice(userModels []models.UserModel) []domain.User {
	users := make([]domain.User, len(userModels))
	for i, userModel := range userModels {
		users[i] = m.ToDomain(userModel)
	}
	return users
}
