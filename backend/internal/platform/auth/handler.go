package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"wishlistapp/internal/platform/httpx"

	"github.com/go-chi/chi"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	repo   *UserRepository
	tokens *TokenManager
}

func NewAuthHandler(repo *UserRepository, tokens *TokenManager) *AuthHandler {
	return &AuthHandler{repo: repo, tokens: tokens}
}

func (h *AuthHandler) Routes(r chi.Router) {
	r.Post("/login", h.Login)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || req.Password == "" {
		http.Error(w, "username and password are required", http.StatusBadRequest)
		return
	}

	creds, err := h.repo.FindCredentialsByUsername(r.Context(), req.Username)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			http.Error(w, "invalid username or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(creds.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := h.tokens.Issue(CurrentUser{ID: creds.ID, Username: creds.Username})
	if err != nil {
		http.Error(w, "could not create session", http.StatusInternalServerError)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, LoginResponse{
		UserID:   creds.ID,
		Username: creds.Username,
		Token:    token,
	})
}
