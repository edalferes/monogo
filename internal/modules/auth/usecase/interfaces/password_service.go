package interfaces

type PasswordService interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}
