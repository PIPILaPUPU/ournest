package auth

import "testing"

func TestTokenRoundTrip(t *testing.T) {
	manager := NewTokenManager("test-secret-with-at-least-thirty-two-characters")
	expected := CurrentUser{ID: 7, Username: "alice"}

	raw, err := manager.Issue(expected)
	if err != nil {
		t.Fatalf("issue token: %v", err)
	}
	actual, err := manager.Parse(raw)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}
	if actual != expected {
		t.Fatalf("unexpected user: got %#v want %#v", actual, expected)
	}
}

func TestTokenRejectsTampering(t *testing.T) {
	manager := NewTokenManager("test-secret-with-at-least-thirty-two-characters")
	raw, err := manager.Issue(CurrentUser{ID: 7, Username: "alice"})
	if err != nil {
		t.Fatalf("issue token: %v", err)
	}

	if _, err := manager.Parse(raw + "tampered"); err == nil {
		t.Fatal("tampered token was accepted")
	}
}
