package auth

type LoginCredentials struct {
	Username string
	Password string
}

type TokenClaims struct {
	UserID string
	Role   string
	Expiry int64
}
