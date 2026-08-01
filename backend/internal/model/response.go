package model

import "time"

// WishResponse — JSON-ответ клиенту.
type WishResponse struct {
	ID            int       `json:"id"`
	OwnerID       int       `json:"owner_id"`
	OwnerUsername string    `json:"owner_username"`
	Title         string    `json:"title"`
	Description   *string   `json:"description"`
	URL           *string   `json:"url"`
	Price         *float64  `json:"price"`
	Status        string    `json:"status"`
	GroupName     string    `json:"group_name"`
	GroupColor    string    `json:"group_color"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func ToWishResponse(w Wish) WishResponse {
	return WishResponse{
		ID:            w.ID,
		OwnerID:       w.OwnerID,
		OwnerUsername: w.OwnerUsername,
		Title:         w.Title,
		Description:   w.Description,
		URL:           w.URL,
		Price:         w.Price,
		Status:        w.Status,
		GroupName:     w.GroupName,
		GroupColor:    w.GroupColor,
		CreatedAt:     w.CreatedAt,
		UpdatedAt:     w.UpdatedAt,
	}
}

func ToWishListResponse(wishes []Wish) []WishResponse {
	result := make([]WishResponse, 0, len(wishes))
	for _, wish := range wishes {
		result = append(result, ToWishResponse(wish))
	}
	return result
}
