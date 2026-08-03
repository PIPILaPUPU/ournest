package wishlist

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"wishlistapp/internal/platform/auth"
	"wishlistapp/internal/platform/httpx"

	"github.com/go-chi/chi"
)

type WishlistHandler struct {
	repo *WishRepository
}

var fixedGroupColors = map[string]string{
	"Техника": "blue",
	"Еда":     "green",
	"Учеба":   "violet",
	"Работа":  "rose",
	"Общее":   "slate",
}

func NewWishlistHandler(repo *WishRepository) *WishlistHandler {
	return &WishlistHandler{repo: repo}
}

func (h *WishlistHandler) Routes(r chi.Router) {
	r.Get("/", h.GetWishlist)
	r.Post("/", h.CreateWish)
	r.Patch("/{id}", h.UpdateWish)
	r.Delete("/{id}", h.DeleteWish)
}

func (h *WishlistHandler) GetWishlist(w http.ResponseWriter, r *http.Request) {
	wishes, err := h.repo.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, ToWishListResponse(wishes))
}

func (h *WishlistHandler) CreateWish(w http.ResponseWriter, r *http.Request) {
	var req CreateWishRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	user, ok := auth.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "authenticated user is required", http.StatusUnauthorized)
		return
	}
	if strings.TrimSpace(req.Title) == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	status := req.Status
	if status == "" {
		status = "wanted"
	}
	if !isValidStatus(status) {
		http.Error(w, "status must be wanted, reserved or done", http.StatusBadRequest)
		return
	}

	rawGroupName := strings.TrimSpace(req.GroupName)
	groupName := normalizeGroupName(rawGroupName)
	if rawGroupName == "" {
		groupName = "Общее"
	}
	groupColor, ok := fixedGroupColors[groupName]
	if !ok {
		http.Error(w, "group_name must be one of: Техника, Еда, Учеба, Работа, Общее", http.StatusBadRequest)
		return
	}

	input := req.ToInput(user.ID, status, groupName, groupColor)
	wish, err := h.repo.Create(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, ToWishResponse(wish))
}

func (h *WishlistHandler) UpdateWish(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req UpdateWishRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.Description == nil && req.URL == nil && req.Price == nil {
		http.Error(w, "at least one field is required: description, url, price", http.StatusBadRequest)
		return
	}

	input := req.ToInput()

	wish, err := h.repo.Update(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, "wish not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, ToWishResponse(wish))
}

func (h *WishlistHandler) DeleteWish(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, "wish not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func isValidStatus(status string) bool {
	switch status {
	case "wanted", "reserved", "done":
		return true
	default:
		return false
	}
}

func normalizeGroupName(name string) string {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case strings.ToLower("Техника"):
		return "Техника"
	case strings.ToLower("Еда"):
		return "Еда"
	case strings.ToLower("Учеба"):
		return "Учеба"
	case strings.ToLower("Работа"):
		return "Работа"
	case strings.ToLower("Общее"):
		return "Общее"
	default:
		return ""
	}
}
