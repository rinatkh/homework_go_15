// Package postgres реализует user.Repository через SQL.
// TODO HW15-04: добавить RefreshSessionRepository без SQL во внутренних слоях.
package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// DBTX реализуют *pgxpool.Pool и pgx.Tx. Это инфраструктурный контракт,
// он не выходит в user и не меняет application interfaces.
type DBTX interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}
