//go:build acceptance

package httpapi_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/rinatkh/homework_go_15/internal/auth"
	"github.com/rinatkh/homework_go_15/internal/httpapi"
	"github.com/rinatkh/homework_go_15/internal/httpapi/handler"
	"github.com/rinatkh/homework_go_15/internal/httpapi/middleware"
	"github.com/rinatkh/homework_go_15/internal/repository/memory"
	"github.com/rinatkh/homework_go_15/internal/user"
	"golang.org/x/crypto/bcrypt"
)

func acceptanceRouter() http.Handler {
	tokens := auth.NewJWTManager("0123456789abcdef0123456789abcdef", 15*time.Minute)
	// HW15-09: замените nil своим test double sessions после реализации wiring.
	uc := user.NewUseCase(memory.NewUserRepository(), auth.NewBcryptHasher(bcrypt.MinCost), tokens, nil, auth.NewCryptoRefreshTokenGenerator(), user.SessionConfig{AccessTTL: 15 * time.Minute, RefreshTTL: 24 * time.Hour})
	return httpapi.NewRouter(handler.NewAuthHandler(uc, handler.CookieConfig{Name: "refresh_token", Secure: true}), middleware.NewAuth(tokens))
}

func TestAcceptanceRefreshWithoutCookie(t *testing.T) {
	router := acceptanceRouter()
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest("POST", "/auth/refresh", nil))
	if w.Code != 401 || w.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("HW15-06: got %d, expected private 401", w.Code)
	}
}

func TestAcceptanceLogoutWithoutCookie(t *testing.T) {
	router := acceptanceRouter()
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest("POST", "/auth/logout", nil))
	cookies := w.Result().Cookies()
	if w.Code != 204 || len(cookies) != 1 || cookies[0].Path != "/auth" || !cookies[0].HttpOnly || cookies[0].MaxAge != -1 {
		t.Fatal("HW15-06: expected idempotent logout and matching cookie deletion")
	}
}

func TestAcceptanceCookieDoesNotAuthorizeProfile(t *testing.T) {
	req := httptest.NewRequest("GET", "/me", nil)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "opaque"})
	w := httptest.NewRecorder()
	acceptanceRouter().ServeHTTP(w, req)
	if w.Code != 401 || strings.Contains(w.Body.String(), "opaque") {
		t.Fatal("refresh cookie must not authenticate /me")
	}
}
