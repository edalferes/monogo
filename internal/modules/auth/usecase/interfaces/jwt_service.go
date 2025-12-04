package interfaces

import "github.com/golang-jwt/jwt/v5"

type JWTService interface {
	GenerateToken(userID uint, username string, roles []string) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}
