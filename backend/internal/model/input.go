package model

// CreateWishInput — данные для repository (не HTTP).
type CreateWishInput struct {
	OwnerID     int
	Title       string
	Description *string
	URL         *string
	Price       *float64
	Status      string
	GroupName   string
	GroupColor  string
}

// UpdateWishInput — частичное обновление для repository.
type UpdateWishInput struct {
	Description *string
	URL         *string
	Price       *float64
}
