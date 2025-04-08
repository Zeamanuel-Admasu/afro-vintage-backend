package auth

type LoginCredentials struct {
	Email    string
	Password string
}

type TokenClaims struct {
	UserID string
	Role   string
	Expiry int64
}
