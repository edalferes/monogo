package token

import (
	"time"

	"github.com/edalferes/monetics/internal/modules/auth/usecase/interfaces"
	"github.com/golang-jwt/jwt/v5"
)

type jwtService struct {
	secretKey string
	expiry    time.Duration
}

func NewJWTService(secretKey string, expiry time.Duration) interfaces.JWTService {
	return &jwtService{secretKey, expiry}
}

func (j *jwtService) GenerateToken(userID uint, username string, roles []string) (string, error) {
	claims := jwt.MapClaims{
		"uid":      userID,
		"username": username,
		"roles":    roles,
		"exp":      time.Now().Add(j.expiry).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
}
