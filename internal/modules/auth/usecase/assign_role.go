package usecase

type AssignRoleUseCase struct {
	// Adicione dependências como UserRepository, RoleRepository
}

func (u *AssignRoleUseCase) Execute(userID uint, roleName string) error {
	// TODO: implementar atribuição de role
	return nil
}
