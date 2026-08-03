package notifications

import "time"

const (
	EventWishCreated     = "wish.created"
	EventDateIdeaCreated = "date_idea.created"
)

type SubscriptionKeys struct {
	P256DH string `json:"p256dh"`
	Auth   string `json:"auth"`
}

type SubscribeRequest struct {
	Endpoint string           `json:"endpoint"`
	Keys     SubscriptionKeys `json:"keys"`
}

type UnsubscribeRequest struct {
	Endpoint string `json:"endpoint"`
}

type EventPayload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Type  string `json:"type"`
	URL   string `json:"url"`
}

type Delivery struct {
	ID             int64
	SubscriptionID int64
	Endpoint       string
	P256DH         string
	Auth           string
	Payload        []byte
	Attempts       int
}

type SendResult struct {
	StatusCode int
}

type Retry struct {
	DeliveryID int64
	NextTry    time.Time
	Message    string
}
