package mapper

import (
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/dto"
	"github.com/Mr-Rafael/finance-calculator/internal/service"
)

func ToCreateUserParams(reqParams dto.UserCreateRequestParams) service.RegisterUserInput {
	return service.RegisterUserInput{
		Email:    reqParams.Email,
		Password: reqParams.Password,
		Username: reqParams.Username,
	}
}

func ToCreateUserResponse(user service.User) dto.UserCreateResponse {
	return dto.UserCreateResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt.Time.Format(time.RFC3339),
	}
}
