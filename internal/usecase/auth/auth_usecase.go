package auth

import (
	"context"
	"errors"
	"time"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/auth"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/user"

	// "github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type authUsecase struct {
	userRepo        user.Repository
	passwordService auth.PasswordService
	jwtService      auth.JWTService
}

func NewAuthUsecase(
	userRepo user.Repository,
	passwordService auth.PasswordService,
	jwtService auth.JWTService,
) auth.AuthUsecase {
	return &authUsecase{
		userRepo:        userRepo,
		passwordService: passwordService,
		jwtService:      jwtService,
	}
}

func (uc *authUsecase) Login(ctx context.Context, creds auth.LoginCredentials) (string, error) {
	u, err := uc.userRepo.FindUserByUsername(ctx, creds.Username)
	if err != nil || !uc.passwordService.CheckPasswordHash(creds.Password, u.Password) {
		return "", errors.New("invalid username or password")
	}

	if creds.Role != "" && creds.Role != string(u.Role) {
		return "", errors.New("access denied: user is not a " + creds.Role)
	}

	token, err := uc.jwtService.GenerateToken(u.ID, u.Username, string(u.Role))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (uc *authUsecase) Register(ctx context.Context, newUser user.User) (string, error) {
	// Check if user already exists
	existing, _ := uc.userRepo.FindUserByUsername(ctx, newUser.Username)
	if existing != nil {
		return "", errors.New("user already exists")
	}

	// Hash password
	hashed, err := uc.passwordService.HashPassword(newUser.Password)
	if err != nil {
		return "", err
	}

	// Generate a valid ObjectID for the user
	objectID := primitive.NewObjectID()
	newUser.ID = objectID.Hex()
	newUser.Password = hashed
	newUser.CreatedAt = time.Now()

	if newUser.Role == "supplier" || newUser.Role == "reseller" {
		newUser.TrustScore = 100
	}

	// Default role to "consumer" if not set
	if newUser.Role == "" {
		newUser.Role = "consumer"
	}

	// Save user
	if err := uc.userRepo.CreateUser(ctx, &newUser); err != nil {
		return "", err
	}

	// Generate token
	token, err := uc.jwtService.GenerateToken(newUser.ID, newUser.Username, string(newUser.Role))
	if err != nil {
		return "", err
	}

	return token, nil
}
