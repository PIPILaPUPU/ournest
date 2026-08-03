package notifications

import "testing"

func TestShouldNotifyDateIdea(t *testing.T) {
	if ShouldNotifyDateIdea(true) {
		t.Fatal("secret date idea must not create a notification")
	}
	if !ShouldNotifyDateIdea(false) {
		t.Fatal("public date idea must create a notification")
	}
}

func TestCreatedPayloadsDoNotExposePrivateFields(t *testing.T) {
	wish := WishCreatedPayload("alice", "Книга")
	if wish.Type != EventWishCreated || wish.URL != "/?section=wishlist" {
		t.Fatalf("unexpected wish payload: %#v", wish)
	}

	idea := DateIdeaCreatedPayload("bob", "Прогулка")
	if idea.Type != EventDateIdeaCreated || idea.URL != "/?section=date-ideas" {
		t.Fatalf("unexpected date idea payload: %#v", idea)
	}
}
