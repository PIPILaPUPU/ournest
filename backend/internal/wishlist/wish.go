package wishlist

import "time"

// Wish — сущность из БД. Без json-тегов: наружу отдаём WishResponse.
type Wish struct {
	ID            int
	OwnerID       int
	OwnerUsername string
	Title         string
	Description   *string
	URL           *string
	Price         *float64
	Status        string
	GroupName     string
	GroupColor    string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
