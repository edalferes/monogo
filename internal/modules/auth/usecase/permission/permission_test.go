package permission

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/edalferes/monetics/internal/modules/auth/domain"
	"github.com/edalferes/monetics/internal/modules/auth/errors"
)

// MockPermissionRepository is a mock implementation of Permission interface
type MockPermissionRepository struct {
	mock.Mock
}

func (m *MockPermissionRepository) FindByName(name string) (*domain.Permission, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Permission), args.Error(1)
}

func (m *MockPermissionRepository) FindByID(id uint) (*domain.Permission, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Permission), args.Error(1)
}

func (m *MockPermissionRepository) ListAll() ([]domain.Permission, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Permission), args.Error(1)
}

func (m *MockPermissionRepository) Create(permission *domain.Permission) error {
	args := m.Called(permission)
	return args.Error(0)
}

func (m *MockPermissionRepository) Update(permission *domain.Permission) error {
	args := m.Called(permission)
	return args.Error(0)
}

func (m *MockPermissionRepository) DeleteByName(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func TestCreatePermissionUseCase_Execute(t *testing.T) {
	t.Run("success - create permission", func(t *testing.T) {
		mockRepo := new(MockPermissionRepository)
		uc := CreatePermissionUseCase{
			PermissionRepo: mockRepo,
		}

		permissionName := "users:read"

		expectedPermission := &domain.Permission{
			Name: permissionName,
		}

		mockRepo.On("Create", expectedPermission).Return(nil)

		err := uc.Execute(permissionName)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - repository error", func(t *testing.T) {
		mockRepo := new(MockPermissionRepository)
		uc := CreatePermissionUseCase{
			PermissionRepo: mockRepo,
		}

		permissionName := "users:read"

		expectedPermission := &domain.Permission{
			Name: permissionName,
		}

		mockRepo.On("Create", expectedPermission).Return(errors.ErrInvalidData)

		err := uc.Execute(permissionName)

		assert.Error(t, err)
		assert.Equal(t, errors.ErrInvalidData, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestListPermissionsUseCase_Execute(t *testing.T) {
	t.Run("success - list all permissions", func(t *testing.T) {
		mockRepo := new(MockPermissionRepository)
		uc := ListPermissionsUseCase{
			PermissionRepo: mockRepo,
		}

		expectedPermissions := []domain.Permission{
			{ID: 1, Name: "users:read"},
			{ID: 2, Name: "users:write"},
			{ID: 3, Name: "roles:read"},
		}

		mockRepo.On("ListAll").Return(expectedPermissions, nil)

		result, err := uc.Execute()

		assert.NoError(t, err)
		assert.Equal(t, expectedPermissions, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - repository error", func(t *testing.T) {
		mockRepo := new(MockPermissionRepository)
		uc := ListPermissionsUseCase{
			PermissionRepo: mockRepo,
		}

		mockRepo.On("ListAll").Return(nil, errors.ErrInvalidData)

		result, err := uc.Execute()

		assert.Error(t, err)
		assert.Equal(t, errors.ErrInvalidData, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetPermissionByIDUseCase_Execute(t *testing.T) {
	t.Run("success - find permission by id", func(t *testing.T) {
		mockRepo := new(MockPermissionRepository)
		uc := GetPermissionByIDUseCase{
			PermissionRepo: mockRepo,
		}

		permissionID := uint(1)
		expectedPermission := &domain.Permission{ID: permissionID, Name: "users:read"}

		mockRepo.On("FindByID", permissionID).Return(expectedPermission, nil)

		result, err := uc.Execute(permissionID)

		assert.NoError(t, err)
		assert.Equal(t, expectedPermission, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - permission not found", func(t *testing.T) {
		mockRepo := new(MockPermissionRepository)
		uc := GetPermissionByIDUseCase{
			PermissionRepo: mockRepo,
		}

		permissionID := uint(999)
		mockRepo.On("FindByID", permissionID).Return(nil, errors.ErrUserNotFound)

		result, err := uc.Execute(permissionID)

		assert.Error(t, err)
		assert.Equal(t, errors.ErrUserNotFound, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdatePermissionUseCase_Execute(t *testing.T) {
	t.Run("success - update permission", func(t *testing.T) {
		mockRepo := new(MockPermissionRepository)
		uc := UpdatePermissionUseCase{
			Permission: mockRepo,
		}

		input := UpdatePermissionInput{
			ID:   1,
			Name: "users:write",
		}

		existingPermission := &domain.Permission{
			ID:   1,
			Name: "users:read",
		}

		mockRepo.On("FindByID", input.ID).Return(existingPermission, nil)

		updatedPermission := &domain.Permission{
			ID:   1,
			Name: "users:write",
		}
		mockRepo.On("Update", updatedPermission).Return(nil)

		err := uc.Execute(input)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - permission not found", func(t *testing.T) {
		mockRepo := new(MockPermissionRepository)
		uc := UpdatePermissionUseCase{
			Permission: mockRepo,
		}

		input := UpdatePermissionInput{
			ID:   999,
			Name: "users:write",
		}

		mockRepo.On("FindByID", input.ID).Return(nil, errors.ErrUserNotFound)

		err := uc.Execute(input)

		assert.Error(t, err)
		assert.Equal(t, errors.ErrUserNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeletePermissionUseCase_Execute(t *testing.T) {
	t.Run("success - delete permission", func(t *testing.T) {
		mockRepo := new(MockPermissionRepository)
		uc := DeletePermissionUseCase{
			PermissionRepo: mockRepo,
		}

		permissionName := "old_permission"
		mockRepo.On("DeleteByName", permissionName).Return(nil)

		err := uc.Execute(permissionName)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error - repository error", func(t *testing.T) {
		mockRepo := new(MockPermissionRepository)
		uc := DeletePermissionUseCase{
			PermissionRepo: mockRepo,
		}

		permissionName := "old_permission"
		mockRepo.On("DeleteByName", permissionName).Return(errors.ErrInvalidData)

		err := uc.Execute(permissionName)

		assert.Error(t, err)
		assert.Equal(t, errors.ErrInvalidData, err)
		mockRepo.AssertExpectations(t)
	})
}
