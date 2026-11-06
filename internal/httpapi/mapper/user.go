// Package mapper преобразует внутренние модели в HTTP response DTO.
package mapper

import (
	"github.com/rinatkh/homework_go_15/internal/httpapi/dto"
	"github.com/rinatkh/homework_go_15/internal/user"
)

func UserToResponse(value user.User) dto.UserResponse {
	return dto.UserResponse{ID: value.ID, Email: value.Email}
}
