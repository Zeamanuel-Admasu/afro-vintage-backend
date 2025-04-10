// infrastructure/auth/jwt_service.go
package authinfra

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtService struct {
	secretKey string
}

func NewJWTService(secretKey string) *jwtService {
	return &jwtService{secretKey: secretKey}
}

func (s *jwtService) GenerateToken(userID, username, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(72 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *jwtService) ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.secretKey), nil
	})
}
