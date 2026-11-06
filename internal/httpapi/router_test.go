package httpapi_test

import (
	"encoding/json"
	"github.com/rinatkh/homework_go_15/internal/auth"
	"github.com/rinatkh/homework_go_15/internal/httpapi"
	"github.com/rinatkh/homework_go_15/internal/httpapi/handler"
	"github.com/rinatkh/homework_go_15/internal/httpapi/middleware"
	"github.com/rinatkh/homework_go_15/internal/repository/memory"
	"github.com/rinatkh/homework_go_15/internal/user"
	"golang.org/x/crypto/bcrypt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestStarterRegisterLoginProtectedProfile(t *testing.T) {
	tokens := auth.NewJWTManager("0123456789abcdef0123456789abcdef", 15*time.Minute)
	router := httpapi.NewRouter(handler.NewAuthHandler(user.NewUseCase(memory.NewUserRepository(), auth.NewBcryptHasher(bcrypt.MinCost), tokens, nil, nil, user.SessionConfig{}), handler.CookieConfig{Name: "refresh_token"}), middleware.NewAuth(tokens))
	body := `{"email":"alice@example.test","password":"secret123"}`
	call := func(method, path, token, body string, want int) *httptest.ResponseRecorder {
		t.Helper()
		req := httptest.NewRequest(method, path, strings.NewReader(body))
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		out := httptest.NewRecorder()
		router.ServeHTTP(out, req)
		if out.Code != want {
			t.Fatalf("%s %s: %d want %d", method, path, out.Code, want)
		}
		return out
	}
	call("GET", "/me", "", "", 401)
	call("POST", "/register", "", body, 201)
	login := call("POST", "/login", "", body, 200)
	var response struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(login.Body.Bytes(), &response); err != nil || response.AccessToken == "" {
		t.Fatal("missing access")
	}
	call("GET", "/me", response.AccessToken, "", 200)
	// Новые контракты проверяются отдельно: make test-acceptance и make test-db.
}
