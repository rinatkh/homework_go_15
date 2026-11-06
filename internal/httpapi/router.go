// Package httpapi собирает routes небольшого Sessions API.
package httpapi

import (
	"github.com/rinatkh/homework_go_15/internal/httpapi/handler"
	"github.com/rinatkh/homework_go_15/internal/httpapi/middleware"
	"net/http"
)

func NewRouter(auth *handler.AuthHandler, guard *middleware.Auth) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("POST /register", auth.Register)
	mux.HandleFunc("POST /login", auth.Login)
	mux.Handle("GET /me", guard.RequireUser(http.HandlerFunc(handler.Me)))
	// Доступ к refresh не требует ещё действующего access JWT.
	mux.HandleFunc("POST /auth/refresh", auth.Refresh)
	mux.HandleFunc("POST /auth/logout", auth.Logout)
	return mux
}
