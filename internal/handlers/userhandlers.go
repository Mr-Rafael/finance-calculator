package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/Mr-Rafael/finance-calculator/internal/dto"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *ApiConfig) HandlerUsersCreate(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	reqParams := dto.UserCreateRequestParams{}
	err := decoder.Decode(&reqParams)
	if err != nil {
		respondWithErrorCode(writer, fmt.Sprintf("received bad user creation request: %v", err), http.StatusBadRequest)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(reqParams.Password), bcrypt.DefaultCost)
	user, err := cfg.Queries.CreateUser(context.Background(), db.CreateUserParams{
		Email:        reqParams.Email,
		PasswordHash: string(passwordHash),
		Username:     reqParams.Username,
	})
	if err != nil {
		respondWithError(writer, fmt.Sprintf("failed to save user to database: %v", err), fmt.Sprintf("database error creating the user: %v", err), http.StatusInternalServerError)
	}

	respondWithJSON(writer, dto.UserCreateResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt.Time.Format(time.RFC3339),
	}, http.StatusCreated)
}
