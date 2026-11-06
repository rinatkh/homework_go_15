package auth

import "github.com/rinatkh/homework_go_15/internal/user"

// CryptoRefreshTokenGenerator — место реализации HW15-03.
// Generate возвращает криптографически случайный opaque credential, Hash — digest.
type CryptoRefreshTokenGenerator struct{}

func NewCryptoRefreshTokenGenerator() *CryptoRefreshTokenGenerator {
	return &CryptoRefreshTokenGenerator{}
}

func (g *CryptoRefreshTokenGenerator) Generate() (string, error) {
	// TODO HW15-03: контракт и ограничения в tasks/02_tokens.md.
	return "", user.ErrNotImplemented
}

func (g *CryptoRefreshTokenGenerator) Hash(raw string) string {
	// TODO HW15-03: детерминированный отпечаток без хранения исходного credential.
	return ""
}
