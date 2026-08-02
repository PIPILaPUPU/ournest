package dateideas

import "time"

type DateIdea struct {
	ID             int
	AuthorID       int
	AuthorUsername string
	Title          string
	Description    *string
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type CreateRequest struct {
	AuthorID    int     `json:"author_id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

type UpdateRequest struct {
	Description *string `json:"description"`
	Status      *string `json:"status"`
}

type Response struct {
	ID             int       `json:"id"`
	AuthorID       int       `json:"author_id"`
	AuthorUsername string    `json:"author_username"`
	Title          string    `json:"title"`
	Description    *string   `json:"description"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func toResponse(idea DateIdea) Response {
	return Response{
		ID:             idea.ID,
		AuthorID:       idea.AuthorID,
		AuthorUsername: idea.AuthorUsername,
		Title:          idea.Title,
		Description:    idea.Description,
		Status:         idea.Status,
		CreatedAt:      idea.CreatedAt,
		UpdatedAt:      idea.UpdatedAt,
	}
}

func toListResponse(ideas []DateIdea) []Response {
	result := make([]Response, 0, len(ideas))
	for _, idea := range ideas {
		result = append(result, toResponse(idea))
	}
	return result
}
