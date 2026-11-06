// Package dto описывает внешний JSON contract HTTP API.
package dto

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

// TODO HW15-07: добавьте token_type/expires_in; refresh не входит в JSON.
type TokenResponse struct {
	Token string `json:"access_token"`
}
