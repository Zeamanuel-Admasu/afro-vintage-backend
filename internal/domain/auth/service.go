package auth

import (
	"context"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/user"
)

type PasswordService interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type JWTService interface {
	GenerateToken(userID, username, role string) (string, error)
	ParseToken(token string) (*TokenClaims, error)
}
type AuthUsecase interface {
	Login(ctx context.Context, creds LoginCredentials) (string, error)
	Register(ctx context.Context, user user.User) (string, error)
}
