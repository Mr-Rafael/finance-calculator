package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mr-Rafael/finance-calculator/internal/dto"
	"github.com/Mr-Rafael/finance-calculator/internal/mapper"
	"github.com/Mr-Rafael/finance-calculator/internal/service"
)

type UsersHandler struct {
	userService service.UserService
}

func NewUsersHandler(service service.UserService) UsersHandler {
	return UsersHandler{userService: service}
}

func (handler *UsersHandler) CreateUser(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	reqParams := dto.UserCreateRequestParams{}
	err := decoder.Decode(&reqParams)
	if err != nil {
		respondWithError(writer, fmt.Sprintf("bad create request: %v", err), fmt.Sprintf("received bad user creation request: %v", err), http.StatusBadRequest)
		return
	}
	result, err := handler.userService.RegisterUser(context.Background(), mapper.ToCreateUserInput(reqParams))
	if err != nil {
		respondWithError(writer, fmt.Sprintf("error creating the user: %v", err), fmt.Sprintf("error creating the user: %v", err), http.StatusInternalServerError)
		return
	}
	respondWithJSON(writer, result, http.StatusCreated)
}
