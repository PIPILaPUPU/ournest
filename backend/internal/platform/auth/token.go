package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const tokenTTL = 30 * 24 * time.Hour

type CurrentUser struct {
	ID       int
	Username string
}

type currentUserContextKey struct{}

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type TokenManager struct {
	secret []byte
}

func NewTokenManager(secret string) *TokenManager {
	return &TokenManager{secret: []byte(secret)}
}

func (m *TokenManager) Issue(user CurrentUser) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Username,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenTTL)),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
}

func (m *TokenManager) Parse(raw string) (CurrentUser, error) {
	token, err := jwt.ParseWithClaims(
		raw,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("unexpected signing method")
			}
			return m.secret, nil
		},
	)
	if err != nil || !token.Valid {
		return CurrentUser{}, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || claims.UserID <= 0 || claims.Username == "" {
		return CurrentUser{}, errors.New("invalid token claims")
	}

	return CurrentUser{ID: claims.UserID, Username: claims.Username}, nil
}

func (m *TokenManager) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := strings.TrimSpace(r.Header.Get("Authorization"))
		const prefix = "Bearer "
		if !strings.HasPrefix(header, prefix) {
			http.Error(w, "authorization token is required", http.StatusUnauthorized)
			return
		}

		user, err := m.Parse(strings.TrimSpace(strings.TrimPrefix(header, prefix)))
		if err != nil {
			http.Error(w, "invalid or expired authorization token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), currentUserContextKey{}, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserFromContext(ctx context.Context) (CurrentUser, bool) {
	user, ok := ctx.Value(currentUserContextKey{}).(CurrentUser)
	return user, ok
}
