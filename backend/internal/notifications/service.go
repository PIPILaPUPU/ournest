package notifications

import "fmt"

func WishCreatedPayload(username, title string) EventPayload {
	return EventPayload{
		Title: "Новое желание",
		Body:  fmt.Sprintf("%s добавил(а): %s", username, title),
		Type:  EventWishCreated,
		URL:   "/?section=wishlist",
	}
}

func DateIdeaCreatedPayload(username, title string) EventPayload {
	return EventPayload{
		Title: "Новая идея для свидания",
		Body:  fmt.Sprintf("%s добавил(а): %s", username, title),
		Type:  EventDateIdeaCreated,
		URL:   "/?section=date-ideas",
	}
}

func ShouldNotifyDateIdea(isSecret bool) bool {
	return !isSecret
}
