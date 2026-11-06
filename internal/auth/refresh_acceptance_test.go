//go:build acceptance

package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"testing"
)

func TestAcceptanceRefreshGenerator(t *testing.T) {
	generator := NewCryptoRefreshTokenGenerator()
	seen := make(map[string]bool)
	for i := 0; i < 32; i++ {
		raw, err := generator.Generate()
		if err != nil {
			t.Fatalf("HW15-03 Generate: %v", err)
		}
		decoded, err := base64.RawURLEncoding.DecodeString(raw)
		if err != nil || len(decoded) < 32 || seen[raw] {
			t.Fatal("expected distinct base64url credentials with at least 32 random bytes")
		}
		seen[raw] = true
	}
	// Уникальные примеры не являются статистическим доказательством энтропии.
	// Использование crypto/rand дополнительно проверяется при ревью кода.
}

func TestAcceptanceRefreshHash(t *testing.T) {
	generator := NewCryptoRefreshTokenGenerator()
	for _, raw := range []string{"", "opaque-refresh-example", "a different credential"} {
		digest := sha256.Sum256([]byte(raw))
		want := hex.EncodeToString(digest[:])
		if got := generator.Hash(raw); got != want || generator.Hash(raw) != got {
			t.Fatal("HW15-03: expected deterministic SHA-256 hex digest")
		}
	}
}
