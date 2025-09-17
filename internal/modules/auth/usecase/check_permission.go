package usecase

type CheckPermissionUseCase struct {
	// Adicione dependências como UserRepository, PermissionRepository
}

func (u *CheckPermissionUseCase) Execute(userID uint, permissionName string) (bool, error) {
	// TODO: implementar checagem de permissão
	return false, nil
}
