package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rinatkh/homework_go_15/internal/auth"
	"github.com/rinatkh/homework_go_15/internal/httpapi"
	"github.com/rinatkh/homework_go_15/internal/httpapi/handler"
	"github.com/rinatkh/homework_go_15/internal/httpapi/middleware"
	"github.com/rinatkh/homework_go_15/internal/repository/postgres"
	"github.com/rinatkh/homework_go_15/internal/user"
)

func main() {
	secret := os.Getenv("JWT_SECRET")
	if len(secret) < 32 {
		log.Fatal("JWT_SECRET должен содержать не менее 32 байт; см. README")
	}
	ttl, err := time.ParseDuration(envOr("ACCESS_TOKEN_TTL", "15m"))
	if err != nil || ttl <= 0 {
		log.Fatal("ACCESS_TOKEN_TTL должен быть положительной duration")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL обязателен; см. README.md")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatal("не удалось создать pool PostgreSQL; проверьте DATABASE_URL")
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		log.Fatal("PostgreSQL недоступен; проверьте подключение и docs/setup.md")
	}

	// TODO HW15-09: подключите session repository, generator и cookie config.
	userRepository := postgres.NewUserRepository(pool)
	passwordHasher := auth.NewBcryptHasher(0)
	tokenManager := auth.NewJWTManager(secret, ttl)
	authUseCase := user.NewUseCase(userRepository, passwordHasher, tokenManager, nil, nil, user.SessionConfig{})
	authHandler := handler.NewAuthHandler(authUseCase, handler.CookieConfig{Name: "refresh_token", Secure: true})
	authMiddleware := middleware.NewAuth(tokenManager)
	router := httpapi.NewRouter(authHandler, authMiddleware)
	server := &http.Server{Addr: envOr("HTTP_ADDR", ":8081"), Handler: router, ReadHeaderTimeout: 5 * time.Second}
	log.Printf("Sessions API starter запущен на %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func envOr(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
