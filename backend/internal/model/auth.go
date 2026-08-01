package model

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
}

type UserCredentials struct {
	ID           int
	Username     string
	PasswordHash string
}
