package user

import (
	"context"
	"time"
)

// RefreshSession содержит состояние credential, но не HTTP cookie.
// TODO HW15-02: объясните назначение полей и границы ответственности контракта.
type RefreshSession struct {
	ID                  int64
	UserID              int64
	TokenHash           string
	ExpiresAt           time.Time
	CreatedAt           time.Time
	RevokedAt           *time.Time
	ReplacedBySessionID *int64
}

// Rotate заменяет credential атомарно. Повторное использование старого запрещено.
type RefreshSessionRepository interface {
	Create(context.Context, RefreshSession) (RefreshSession, error)
	GetByTokenHash(context.Context, string) (RefreshSession, error)
	Rotate(context.Context, string, RefreshSession) error
	Revoke(context.Context, string) error
}

type RefreshTokenGenerator interface {
	Generate() (string, error)
	Hash(string) string
}

type SessionConfig struct {
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

// TODO HW15-05: внедрите необходимые зависимости в UseCase и расширьте Login.
func (uc *UseCase) Refresh(ctx context.Context, raw string) (LoginResult, error) {
	return LoginResult{}, ErrNotImplemented
}

func (uc *UseCase) Logout(ctx context.Context, raw string) error {
	return ErrNotImplemented
}
