package role

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/edalferes/monetics/internal/modules/auth/domain"
	"github.com/edalferes/monetics/internal/modules/auth/errors"
)

// MockRoleRepository is a mock implementation of Role interface
type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) FindByName(name string) (*domain.Role, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Role), args.Error(1)
}

func (m *MockRoleRepository) FindByID(id uint) (*domain.Role, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Role), args.Error(1)
}

func (m *MockRoleRepository) ListAll() ([]domain.Role, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Role), args.Error(1)
}

func (m *MockRoleRepository) Create(role *domain.Role) error {
	args := m.Called(role)
	return args.Error(0)
}

func (m *MockRoleRepository) Update(role *domain.Role) error {
	args := m.Called(role)
	return args.Error(0)
}

func (m *MockRoleRepository) DeleteByName(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func TestCreateRoleUseCase_Execute(t *testing.T) {
	t.Run("success - create role with permissions", func(t *testing.T) {
		mockRepo := new(MockRoleRepository)
		uc := CreateRoleUseCase{
			RoleRepo: mockRepo,
		}

		roleName := "admin"
		permissionIDs := []uint{1, 2, 3}

		expectedRole := &domain.Role{
			Name: roleName,
			Permissions: []domain.Permission{
				{ID: 1},
				{ID: 2},
				{ID: 3},
			},
		}

		mockRepo.On("Create", expectedRole).Return(nil)

		err := uc.Execute(roleName, permissionIDs)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success - create role without permissions", func(t *testing.T) {
		mockRepo := new(MockRoleRepository)
		uc := CreateRoleUseCase{
			RoleRepo: mockRepo,
		}

		roleName := "viewer"
		permissionIDs := []uint{}

		expectedRole := &domain.Role{
			Name:        roleName,
			Permissions: []domain.Permission{},
		}

		mockRepo.On("Create", expectedRole).Return(nil)

		err := uc.Execute(roleName, permissionIDs)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - repository error", func(t *testing.T) {
		mockRepo := new(MockRoleRepository)
		uc := CreateRoleUseCase{
			RoleRepo: mockRepo,
		}

		roleName := "admin"
		permissionIDs := []uint{1}

		expectedRole := &domain.Role{
			Name: roleName,
			Permissions: []domain.Permission{
				{ID: 1},
			},
		}

		mockRepo.On("Create", expectedRole).Return(errors.ErrInvalidData)

		err := uc.Execute(roleName, permissionIDs)

		assert.Error(t, err)
		assert.Equal(t, errors.ErrInvalidData, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestListRolesUseCase_Execute(t *testing.T) {
	t.Run("success - list all roles", func(t *testing.T) {
		mockRepo := new(MockRoleRepository)
		uc := ListRolesUseCase{
			RoleRepo: mockRepo,
		}

		expectedRoles := []domain.Role{
			{ID: 1, Name: "admin"},
			{ID: 2, Name: "user"},
		}

		mockRepo.On("ListAll").Return(expectedRoles, nil)

		result, err := uc.Execute()

		assert.NoError(t, err)
		assert.Equal(t, expectedRoles, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - repository error", func(t *testing.T) {
		mockRepo := new(MockRoleRepository)
		uc := ListRolesUseCase{
			RoleRepo: mockRepo,
		}

		mockRepo.On("ListAll").Return(nil, errors.ErrInvalidData)

		result, err := uc.Execute()

		assert.Error(t, err)
		assert.Equal(t, errors.ErrInvalidData, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetRoleByIDUseCase_Execute(t *testing.T) {
	t.Run("success - find role by id", func(t *testing.T) {
		mockRepo := new(MockRoleRepository)
		uc := GetRoleByIDUseCase{
			RoleRepo: mockRepo,
		}

		roleID := uint(1)
		expectedRole := &domain.Role{ID: roleID, Name: "admin"}

		mockRepo.On("FindByID", roleID).Return(expectedRole, nil)

		result, err := uc.Execute(roleID)

		assert.NoError(t, err)
		assert.Equal(t, expectedRole, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - role not found", func(t *testing.T) {
		mockRepo := new(MockRoleRepository)
		uc := GetRoleByIDUseCase{
			RoleRepo: mockRepo,
		}

		roleID := uint(999)
		mockRepo.On("FindByID", roleID).Return(nil, errors.ErrUserNotFound)

		result, err := uc.Execute(roleID)

		assert.Error(t, err)
		assert.Equal(t, errors.ErrUserNotFound, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateRoleUseCase_Execute(t *testing.T) {
	t.Run("success - update role name and permissions", func(t *testing.T) {
		mockRepo := new(MockRoleRepository)
		uc := UpdateRoleUseCase{
			Role: mockRepo,
		}

		input := UpdateRoleInput{
			ID:            1,
			Name:          "super_admin",
			PermissionIDs: []uint{1, 2, 3, 4},
		}

		existingRole := &domain.Role{
			ID:   1,
			Name: "admin",
			Permissions: []domain.Permission{
				{ID: 1}, {ID: 2},
			},
		}

		mockRepo.On("FindByID", input.ID).Return(existingRole, nil)

		updatedRole := &domain.Role{
			ID:   1,
			Name: "super_admin",
			Permissions: []domain.Permission{
				{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4},
			},
		}
		mockRepo.On("Update", updatedRole).Return(nil)

		err := uc.Execute(input)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success - update only name", func(t *testing.T) {
		mockRepo := new(MockRoleRepository)
		uc := UpdateRoleUseCase{
			Role: mockRepo,
		}

		input := UpdateRoleInput{
			ID:   1,
			Name: "moderator",
		}

		existingRole := &domain.Role{
			ID:   1,
			Name: "admin",
			Permissions: []domain.Permission{
				{ID: 1}, {ID: 2},
			},
		}

		mockRepo.On("FindByID", input.ID).Return(existingRole, nil)

		updatedRole := &domain.Role{
			ID:   1,
			Name: "moderator",
			Permissions: []domain.Permission{
				{ID: 1}, {ID: 2},
			},
		}
		mockRepo.On("Update", updatedRole).Return(nil)

		err := uc.Execute(input)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - role not found", func(t *testing.T) {
		mockRepo := new(MockRoleRepository)
		uc := UpdateRoleUseCase{
			Role: mockRepo,
		}

		input := UpdateRoleInput{
			ID:   999,
			Name: "newname",
		}

		mockRepo.On("FindByID", input.ID).Return(nil, errors.ErrUserNotFound)

		err := uc.Execute(input)

		assert.Error(t, err)
		assert.Equal(t, errors.ErrUserNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteRoleUseCase_Execute(t *testing.T) {
	t.Run("success - delete role", func(t *testing.T) {
		mockRepo := new(MockRoleRepository)
		uc := DeleteRoleUseCase{
			RoleRepo: mockRepo,
		}

		roleName := "old_role"
		mockRepo.On("DeleteByName", roleName).Return(nil)

		err := uc.Execute(roleName)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - repository error", func(t *testing.T) {
		mockRepo := new(MockRoleRepository)
		uc := DeleteRoleUseCase{
			RoleRepo: mockRepo,
		}

		roleName := "old_role"
		mockRepo.On("DeleteByName", roleName).Return(errors.ErrInvalidData)

		err := uc.Execute(roleName)

		assert.Error(t, err)
		assert.Equal(t, errors.ErrInvalidData, err)
		mockRepo.AssertExpectations(t)
	})
}
