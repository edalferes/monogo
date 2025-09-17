package usecase

type CheckPermissionUseCase struct {
	// TODO: added dependencies like UserRepository, PermissionRepository
}

func (u *CheckPermissionUseCase) Execute(userID uint, permissionName string) (bool, error) {
	// TODO: implement permission check
	return false, nil
}
