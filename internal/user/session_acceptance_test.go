//go:build acceptance

package user

import (
	"context"
	"errors"
	"testing"
	"time"
)

type acceptanceSessions struct {
	created int
	value   RefreshSession
}

func (r *acceptanceSessions) Create(_ context.Context, s RefreshSession) (RefreshSession, error) {
	r.created++
	r.value = s
	return s, nil
}
func (r *acceptanceSessions) GetByTokenHash(context.Context, string) (RefreshSession, error) {
	return RefreshSession{}, ErrInvalidRefreshToken
}
func (r *acceptanceSessions) Rotate(context.Context, string, RefreshSession) error {
	return ErrInvalidRefreshToken
}
func (r *acceptanceSessions) Revoke(context.Context, string) error { return nil }

type acceptanceGenerator struct{}

func (acceptanceGenerator) Generate() (string, error) { return "raw-test-credential", nil }
func (acceptanceGenerator) Hash(string) string        { return "digest-only" }

func TestAcceptanceLoginSavesSession(t *testing.T) {
	sessions := &acceptanceSessions{}
	repo := &fakeRepository{byEmail: map[string]User{"alice@example.test": {ID: 7, Email: "alice@example.test", PasswordHash: "hash:secret123"}}}
	uc := NewUseCase(repo, fakeHasher{}, fakeTokens{}, sessions, acceptanceGenerator{}, SessionConfig{AccessTTL: 15 * time.Minute, RefreshTTL: 24 * time.Hour})
	result, err := uc.Login(context.Background(), LoginInput{Email: "alice@example.test", Password: "secret123"})
	if err != nil {
		t.Fatal(err)
	}
	if result.AccessToken == "" || result.RefreshToken == "" || result.ExpiresIn != 900 || sessions.created != 1 || sessions.value.UserID != 7 || sessions.value.TokenHash != "digest-only" || !sessions.value.ExpiresAt.After(time.Now()) {
		t.Fatal("HW15-05: login must save a hash-only session before returning both credentials")
	}
}

func TestAcceptanceMissingRefresh(t *testing.T) {
	uc := NewUseCase(&fakeRepository{}, fakeHasher{}, fakeTokens{}, &acceptanceSessions{}, acceptanceGenerator{}, SessionConfig{})
	if _, err := uc.Refresh(context.Background(), ""); !errors.Is(err, ErrInvalidRefreshToken) {
		t.Fatalf("HW15-05: empty refresh error=%v", err)
	}
}

func TestAcceptanceLogoutIsIdempotent(t *testing.T) {
	uc := NewUseCase(&fakeRepository{}, fakeHasher{}, fakeTokens{}, &acceptanceSessions{}, acceptanceGenerator{}, SessionConfig{})
	for _, raw := range []string{"", "unknown-credential", "unknown-credential"} {
		if err := uc.Logout(context.Background(), raw); err != nil {
			t.Fatalf("HW15-05 logout: %v", err)
		}
	}
}
