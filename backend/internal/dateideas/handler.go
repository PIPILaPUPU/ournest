package dateideas

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"wishlistapp/internal/platform/httpx"

	"github.com/go-chi/chi"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Routes(r chi.Router) {
	r.Get("/", h.List)
	r.Get("/secret/random", h.RandomSecret)
	r.Post("/", h.Create)
	r.Patch("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ideas, err := h.repo.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, toListResponse(ideas))
}

func (h *Handler) RandomSecret(w http.ResponseWriter, r *http.Request) {
	idea, err := h.repo.RandomSecret(r.Context())
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, "secret date ideas not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, toResponse(idea))
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	req.Title = strings.TrimSpace(req.Title)
	if req.AuthorID <= 0 {
		http.Error(w, "author_id is required", http.StatusBadRequest)
		return
	}
	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	idea, err := h.repo.Create(
		r.Context(),
		req.AuthorID,
		req.Title,
		req.Description,
		req.IsSecret,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, toResponse(idea))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	if req.Description == nil && req.Status == nil {
		http.Error(w, "at least one field is required: description, status", http.StatusBadRequest)
		return
	}
	if req.Status != nil && *req.Status != "planned" && *req.Status != "done" {
		http.Error(w, "status must be planned or done", http.StatusBadRequest)
		return
	}

	idea, err := h.repo.Update(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, "date idea not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, toResponse(idea))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.repo.Delete(r.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, "date idea not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
