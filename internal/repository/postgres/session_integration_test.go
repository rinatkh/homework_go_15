//go:build integration

package postgres

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rinatkh/homework_go_15/internal/user"
)

func sessionTestPool(t *testing.T) (*pgxpool.Pool, context.Context) {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Fatal("задайте отдельную TEST_DATABASE_URL для integration tests")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	t.Cleanup(cancel)
	admin, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(admin.Close)
	var random [8]byte
	if _, err := rand.Read(random[:]); err != nil {
		t.Fatal(err)
	}
	schema := "homework15_test_" + hex.EncodeToString(random[:])
	quoted := pgx.Identifier{schema}.Sanitize()
	if _, err := admin.Exec(ctx, "CREATE SCHEMA "+quoted); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		cleanup, stop := context.WithTimeout(context.Background(), 5*time.Second)
		defer stop()
		if _, err := admin.Exec(cleanup, "DROP SCHEMA "+quoted+" CASCADE"); err != nil {
			t.Error(err)
		}
	})
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatal(err)
	}
	cfg.ConnConfig.RuntimeParams["search_path"] = schema
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	db := stdlib.OpenDB(*cfg.ConnConfig)
	defer db.Close()
	provider, err := goose.NewProvider(goose.DialectPostgres, db, os.DirFS("../../../migrations"), goose.WithDisableGlobalRegistry(true))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := provider.Up(ctx); err != nil {
		t.Fatalf("HW15-01: migrations are incomplete: %v", err)
	}
	return pool, ctx
}

func TestIntegrationSessionLifecycle(t *testing.T) {
	pool, ctx := sessionTestPool(t)
	owner, err := NewUserRepository(pool).Create(ctx, user.User{Email: "owner@example.test", PasswordHash: "test-hash"})
	if err != nil {
		t.Fatal(err)
	}
	r := NewRefreshSessionRepository(pool)
	oldHash, nextHash := strings.Repeat("a", 64), strings.Repeat("b", 64)
	old, err := r.Create(ctx, user.RefreshSession{UserID: owner.ID, TokenHash: oldHash, ExpiresAt: time.Now().Add(time.Hour)})
	if err != nil || old.ID <= 0 {
		t.Fatalf("HW15-04 Create: %v", err)
	}
	found, err := r.GetByTokenHash(ctx, oldHash)
	if err != nil || found.UserID != owner.ID {
		t.Fatal("session lookup failed")
	}
	if err := r.Rotate(ctx, oldHash, user.RefreshSession{UserID: owner.ID, TokenHash: nextHash, ExpiresAt: time.Now().Add(time.Hour)}); err != nil {
		t.Fatal(err)
	}
	if err := r.Rotate(ctx, oldHash, user.RefreshSession{UserID: owner.ID, TokenHash: strings.Repeat("c", 64), ExpiresAt: time.Now().Add(time.Hour)}); !errors.Is(err, user.ErrRefreshRevoked) {
		t.Fatal("old refresh must be rejected")
	}
	for i := 0; i < 2; i++ {
		if err := r.Revoke(ctx, nextHash); err != nil {
			t.Fatal(err)
		}
	}
	revoked, err := r.GetByTokenHash(ctx, nextHash)
	if err != nil || revoked.RevokedAt == nil {
		t.Fatal("revoke not persisted")
	}
	if _, err := pool.Exec(ctx, "DELETE FROM users WHERE id=$1", owner.ID); err != nil {
		t.Fatal(err)
	}
	var count int
	if err := pool.QueryRow(ctx, "SELECT count(*) FROM refresh_sessions").Scan(&count); err != nil || count != 0 {
		t.Fatal("FK cascade must remove sessions")
	}
}

func TestIntegrationRotationRollbackAndRace(t *testing.T) {
	pool, ctx := sessionTestPool(t)
	owner, err := NewUserRepository(pool).Create(ctx, user.User{Email: "race@example.test", PasswordHash: "test-hash"})
	if err != nil {
		t.Fatal(err)
	}
	r := NewRefreshSessionRepository(pool)
	oldHash, existingHash := strings.Repeat("a", 64), strings.Repeat("b", 64)
	for _, hash := range []string{oldHash, existingHash} {
		if _, err := r.Create(ctx, user.RefreshSession{UserID: owner.ID, TokenHash: hash, ExpiresAt: time.Now().Add(time.Hour)}); err != nil {
			t.Fatal(err)
		}
	}
	if err := r.Rotate(ctx, oldHash, user.RefreshSession{UserID: owner.ID, TokenHash: existingHash, ExpiresAt: time.Now().Add(time.Hour)}); !errors.Is(err, user.ErrSessionAlreadyExists) {
		t.Fatalf("duplicate hash error=%v", err)
	}
	old, err := r.GetByTokenHash(ctx, oldHash)
	if err != nil || old.RevokedAt != nil {
		t.Fatal("failed insertion must preserve old session")
	}
	start := make(chan struct{})
	results := make(chan error, 2)
	var wg sync.WaitGroup
	for _, hash := range []string{strings.Repeat("c", 64), strings.Repeat("d", 64)} {
		wg.Add(1)
		go func(hash string) {
			defer wg.Done()
			<-start
			results <- r.Rotate(ctx, oldHash, user.RefreshSession{UserID: owner.ID, TokenHash: hash, ExpiresAt: time.Now().Add(time.Hour)})
		}(hash)
	}
	close(start)
	wg.Wait()
	close(results)
	successes := 0
	for err := range results {
		if err == nil {
			successes++
		} else if !errors.Is(err, user.ErrRefreshRevoked) {
			t.Fatal(err)
		}
	}
	if successes != 1 {
		t.Fatalf("expected one successful rotation, got %d", successes)
	}
}
