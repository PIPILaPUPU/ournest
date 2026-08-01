package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"wishlistapp/internal/model"
	"wishlistapp/internal/repository"

	"github.com/go-chi/chi"
)

type WishlistHandler struct {
	repo *repository.WishRepository
}

var fixedGroupColors = map[string]string{
	"Техника": "blue",
	"Еда":     "green",
	"Учеба":   "violet",
	"Работа":  "rose",
	"Общее":   "slate",
}

func NewWishlistHandler(repo *repository.WishRepository) *WishlistHandler {
	return &WishlistHandler{repo: repo}
}

func (h *WishlistHandler) GetWishlist(w http.ResponseWriter, r *http.Request) {
	wishes, err := h.repo.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, model.ToWishListResponse(wishes))
}

func (h *WishlistHandler) CreateWish(w http.ResponseWriter, r *http.Request) {
	var req model.CreateWishRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.OwnerID <= 0 {
		http.Error(w, "owner_id is required", http.StatusBadRequest)
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

	wish, err := h.repo.Create(r.Context(), req.ToInput(status, groupName, groupColor))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, model.ToWishResponse(wish))
}

func (h *WishlistHandler) UpdateWish(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req model.UpdateWishRequest
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
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, "wish not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, model.ToWishResponse(wish))
}

func (h *WishlistHandler) DeleteWish(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, "wish not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func parseID(raw string) (int, error) {
	id, err := strconv.Atoi(raw)
	if err != nil || id <= 0 {
		return 0, err
	}
	return id, nil
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
