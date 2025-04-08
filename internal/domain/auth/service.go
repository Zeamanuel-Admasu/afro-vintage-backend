package auth

import "context"

type AuthService interface {
	Login(ctx context.Context, creds LoginCredentials) (string, error) // returns JWT
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
	GenerateToken(userID, role string) (string, error)
	ParseToken(token string) (*TokenClaims, error)
}
