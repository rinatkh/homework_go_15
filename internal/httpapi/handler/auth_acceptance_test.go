//go:build acceptance

package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/rinatkh/homework_go_15/internal/httpapi/handler"
	"github.com/rinatkh/homework_go_15/internal/user"
)

type fakeAuthApplication struct{ refreshError error }

func (f *fakeAuthApplication) Register(context.Context, user.RegisterInput) (user.User, error) {
	return user.User{ID: 1}, nil
}
func (f *fakeAuthApplication) Login(context.Context, user.LoginInput) (user.LoginResult, error) {
	return user.LoginResult{AccessToken: "test-access", RefreshToken: "test-refresh", ExpiresIn: 900, RefreshExpiresAt: time.Now().Add(time.Hour)}, nil
}
func (f *fakeAuthApplication) Refresh(context.Context, string) (user.LoginResult, error) {
	return user.LoginResult{}, f.refreshError
}
func (f *fakeAuthApplication) Logout(context.Context, string) error { return nil }

func TestAcceptanceLoginCookieContract(t *testing.T) {
	h := handler.NewAuthHandler(&fakeAuthApplication{}, handler.CookieConfig{Name: "custom_refresh", Secure: true})
	w := httptest.NewRecorder()
	h.Login(w, httptest.NewRequest("POST", "/login", strings.NewReader(`{"email":"alice@example.test","password":"secret123"}`)))
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	cookies := w.Result().Cookies()
	if w.Code != 200 || len(body) != 3 || body["access_token"] != "test-access" || body["expires_in"] != float64(900) || body["token_type"] != "Bearer" || strings.Contains(w.Body.String(), "test-refresh") {
		t.Fatal("HW15-06/07: token DTO contract")
	}
	if len(cookies) != 1 || cookies[0].Name != "custom_refresh" || !cookies[0].Secure || !cookies[0].HttpOnly || cookies[0].SameSite != http.SameSiteLaxMode || cookies[0].Path != "/auth" || cookies[0].MaxAge <= 0 || w.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("HW15-06: cookie/cache contract")
	}
}

func TestAcceptanceRefreshErrorsAreNeutral(t *testing.T) {
	for _, cause := range []error{user.ErrInvalidRefreshToken, user.ErrRefreshExpired, user.ErrRefreshRevoked} {
		h := handler.NewAuthHandler(&fakeAuthApplication{refreshError: cause}, handler.CookieConfig{Name: "refresh_token"})
		req := httptest.NewRequest("POST", "/auth/refresh", nil)
		req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "test-refresh"})
		w := httptest.NewRecorder()
		h.Refresh(w, req)
		if w.Code != 401 || strings.TrimSpace(w.Body.String()) != `{"error":"invalid refresh token"}` || len(w.Result().Cookies()) != 0 {
			t.Fatal("HW15-06: expected one neutral refresh 401 without cookie replacement")
		}
	}
}
