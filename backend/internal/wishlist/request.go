package wishlist

// CreateWishRequest — тело POST /wishlist.
type CreateWishRequest struct {
	Title       string   `json:"title"`
	Description *string  `json:"description"`
	URL         *string  `json:"url"`
	Price       *float64 `json:"price"`
	Status      string   `json:"status"`
	GroupName   string   `json:"group_name"`
	GroupColor  string   `json:"group_color"`
}

func (r CreateWishRequest) ToInput(ownerID int, status, groupName, groupColor string) CreateWishInput {
	return CreateWishInput{
		OwnerID:     ownerID,
		Title:       r.Title,
		Description: r.Description,
		URL:         r.URL,
		Price:       r.Price,
		Status:      status,
		GroupName:   groupName,
		GroupColor:  groupColor,
	}
}

// UpdateWishRequest — тело PATCH /wishlist/{id}.
type UpdateWishRequest struct {
	Description *string  `json:"description"`
	URL         *string  `json:"url"`
	Price       *float64 `json:"price"`
}

func (r UpdateWishRequest) ToInput() UpdateWishInput {
	return UpdateWishInput{
		Description: r.Description,
		URL:         r.URL,
		Price:       r.Price,
	}
}
