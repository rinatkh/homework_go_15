package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/rinatkh/homework_go_15/internal/user"
)

type SessionDB interface {
	DBTX
	BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error)
}

type RefreshSessionRepository struct{ db SessionDB }

var _ user.RefreshSessionRepository = (*RefreshSessionRepository)(nil)

func NewRefreshSessionRepository(db SessionDB) *RefreshSessionRepository {
	return &RefreshSessionRepository{db: db}
}

func (r *RefreshSessionRepository) Create(ctx context.Context, value user.RefreshSession) (user.RefreshSession, error) {
	// TODO HW15-04: сохраняйте данные через переданный ctx.
	return user.RefreshSession{}, user.ErrNotImplemented
}
func (r *RefreshSessionRepository) GetByTokenHash(ctx context.Context, hash string) (user.RefreshSession, error) {
	return user.RefreshSession{}, user.ErrNotImplemented
}
func (r *RefreshSessionRepository) Rotate(ctx context.Context, oldHash string, next user.RefreshSession) error {
	// TODO HW15-04: контракт атомарной замены описан в tasks/03_repository.md.
	return user.ErrNotImplemented
}
func (r *RefreshSessionRepository) Revoke(ctx context.Context, hash string) error {
	return user.ErrNotImplemented
}
