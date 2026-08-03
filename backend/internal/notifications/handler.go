package notifications

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"wishlistapp/internal/platform/auth"
	"wishlistapp/internal/platform/httpx"

	"github.com/go-chi/chi"
)

type Handler struct {
	repo      *Repository
	publicKey string
	enabled   bool
}

func NewHandler(repo *Repository, publicKey string, enabled bool) *Handler {
	return &Handler{repo: repo, publicKey: publicKey, enabled: enabled}
}

func (h *Handler) Routes(r chi.Router) {
	r.Post("/subscriptions", h.Subscribe)
	r.Delete("/subscriptions", h.Unsubscribe)
}

func (h *Handler) PublicKey(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"public_key": h.publicKey,
		"enabled":    h.enabled,
	})
}

func (h *Handler) Subscribe(w http.ResponseWriter, r *http.Request) {
	if !h.enabled {
		http.Error(w, "push notifications are not configured", http.StatusServiceUnavailable)
		return
	}

	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "authenticated user is required", http.StatusUnauthorized)
		return
	}

	var request SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	request.Endpoint = strings.TrimSpace(request.Endpoint)
	request.Keys.P256DH = strings.TrimSpace(request.Keys.P256DH)
	request.Keys.Auth = strings.TrimSpace(request.Keys.Auth)
	if !validSubscription(request) {
		http.Error(w, "valid endpoint, p256dh and auth keys are required", http.StatusBadRequest)
		return
	}

	if err := h.repo.SaveSubscription(
		r.Context(),
		user.ID,
		request,
		r.UserAgent(),
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "authenticated user is required", http.StatusUnauthorized)
		return
	}

	var request UnsubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	request.Endpoint = strings.TrimSpace(request.Endpoint)
	if request.Endpoint == "" {
		http.Error(w, "endpoint is required", http.StatusBadRequest)
		return
	}

	if err := h.repo.RemoveSubscription(r.Context(), user.ID, request.Endpoint); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func validSubscription(request SubscribeRequest) bool {
	if request.Endpoint == "" || request.Keys.P256DH == "" || request.Keys.Auth == "" {
		return false
	}
	endpoint, err := url.ParseRequestURI(request.Endpoint)
	return err == nil && endpoint.Scheme == "https" && endpoint.Host != ""
}
