package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/rinatkh/homework_go_15/internal/httpapi/dto"
	"github.com/rinatkh/homework_go_15/internal/httpapi/mapper"
	"github.com/rinatkh/homework_go_15/internal/user"
)

type AuthUseCase interface {
	Register(ctx context.Context, input user.RegisterInput) (user.User, error)
	Login(ctx context.Context, input user.LoginInput) (user.LoginResult, error)
	Refresh(context.Context, string) (user.LoginResult, error)
	Logout(context.Context, string) error
}

type AuthHandler struct {
	useCase AuthUseCase
	cookie  CookieConfig
}

func NewAuthHandler(useCase AuthUseCase, cookie CookieConfig) *AuthHandler {
	return &AuthHandler{useCase: useCase, cookie: cookie}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var request dto.RegisterRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	result, err := h.useCase.Register(r.Context(), user.RegisterInput{
		Email: request.Email, Password: request.Password,
	})
	switch {
	case errors.Is(err, user.ErrInvalidEmail), errors.Is(err, user.ErrWeakPassword):
		writeError(w, http.StatusBadRequest, err.Error())
		return
	case errors.Is(err, user.ErrAlreadyExists):
		writeError(w, http.StatusConflict, err.Error())
		return
	case err != nil:
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusCreated, mapper.UserToResponse(result))
}

// TODO HW15-06: дополните ответ Login refresh-cookie; raw credential не входит в JSON.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request dto.LoginRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	result, err := h.useCase.Login(r.Context(), user.LoginInput{
		Email: request.Email, Password: request.Password,
	})
	if errors.Is(err, user.ErrInvalidCredentials) {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	writeJSON(w, http.StatusOK, dto.TokenResponse{Token: result.AccessToken})
}

// TODO HW15-06: реализуйте HTTP-границу для существующих методов UseCase.
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "TODO HW15-06: refresh handler")
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "TODO HW15-06: logout handler")
}
