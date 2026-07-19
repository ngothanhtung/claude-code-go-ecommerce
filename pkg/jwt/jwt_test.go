package jwt

import "testing"

func TestAccessTokenRoundTrip(t *testing.T) {
	m := New("test-secret", 15, 168)
	tok, _, err := m.GenerateAccessToken("user-1", "a@b.com", []string{"user"})
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	claims, err := m.Parse(tok)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if claims.UserID != "user-1" {
		t.Fatalf("expected user-1 got %s", claims.UserID)
	}
	if claims.Email != "a@b.com" {
		t.Fatalf("expected a@b.com got %s", claims.Email)
	}
}

func TestParseRejectsWrongSecret(t *testing.T) {
	m := New("secret-A", 15, 168)
	tok, _, _ := m.GenerateAccessToken("u", "e", nil)
	other := New("secret-B", 15, 168)
	if _, err := other.Parse(tok); err == nil {
		t.Fatal("expected error parsing token signed with different secret")
	}
}

func TestGenerateRefreshTokenWithIDEmbedsID(t *testing.T) {
	m := New("test-secret", 15, 168)
	claims := Claims{UserID: "u-1", Email: "a@b.com", Roles: []string{"user"}}
	tok, _, err := m.GenerateRefreshTokenWithID("rt-id-123", claims)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	got, err := m.Parse(tok)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got.ID != "rt-id-123" {
		t.Fatalf("expected id rt-id-123 got %q", got.ID)
	}
}
