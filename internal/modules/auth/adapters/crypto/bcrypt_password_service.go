package crypto

import (
	"github.com/edalferes/monetics/internal/modules/auth/usecase/interfaces"
	"golang.org/x/crypto/bcrypt"
)

type bcryptPasswordService struct{}

func NewBcryptPasswordService() interfaces.PasswordService {
	return &bcryptPasswordService{}
}

func (p *bcryptPasswordService) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (p *bcryptPasswordService) Compare(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
