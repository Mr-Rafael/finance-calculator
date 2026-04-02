package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mr-Rafael/finance-calculator/internal/dto"
	"github.com/Mr-Rafael/finance-calculator/internal/mapper"
	"github.com/Mr-Rafael/finance-calculator/internal/service"
	"github.com/google/uuid"
)

type SavingsHandler struct {
	savingsService *service.SavingsService
}

func NewSavingsHandler(service *service.SavingsService) *SavingsHandler {
	return &SavingsHandler{savingsService: service}
}

func (handler *SavingsHandler) HandleCalculateSavings(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	reqParams := dto.SavingsRequestParams{}
	err := decoder.Decode(&reqParams)
	if err != nil {
		respondWithErrorCode(writer, "received bad savings request", http.StatusBadRequest)
		return
	}

	result, err := handler.savingsService.GetSavingsPlan(context.Background(), mapper.ToSavingsInput(reqParams))
	if err != nil {
		respondWithError(writer, fmt.Sprintf("Error calculating savings plan: %v", err), fmt.Sprintf("Error calculating savings plan: %v", err), http.StatusInternalServerError)
	}
	respondWithJSON(writer, mapper.ToSavingsResponse(result), http.StatusOK)
}

func (handler *SavingsHandler) HandleSaveSavings(writer http.ResponseWriter, request *http.Request) {
	userID := request.Context().Value(userIDKey).(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		respondWithErrorCode(writer, "failed to get user ID from context", http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(request.Body)
	reqParams := dto.SavingsSaveRequestParams{}
	err = decoder.Decode(&reqParams)
	if err != nil {
		respondWithErrorCode(writer, "received bad savings request", http.StatusBadRequest)
		return
	}

	result, err := handler.savingsService.SaveSavingsPlan(context.Background(), mapper.ToSaveSavingsInput(userUUID, reqParams))
	if err != nil {
		respondWithError(writer, fmt.Sprintf("Error saving the plan: %v", err), fmt.Sprintf("Error saving the plan: %v", err), http.StatusInternalServerError)
		return
	}

	respondWithJSON(writer, mapper.ToSavingsSaveResponse(result), http.StatusCreated)
}

func (handler *SavingsHandler) HandleListSavings(writer http.ResponseWriter, request *http.Request) {
	userID := request.Context().Value(userIDKey).(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		respondWithErrorCode(writer, "failed to get user ID from context", http.StatusUnauthorized)
		return
	}

	result, err := handler.savingsService.GetSavingsPlansByUser(context.Background(), userUUID)

	respondWithJSON(writer, mapper.ToSavingsListResponse(result), http.StatusOK)
}

func (handler *SavingsHandler) HandleGetSavings(writer http.ResponseWriter, request *http.Request) {
	userID := request.Context().Value(userIDKey).(string)
	planID := request.PathValue("id")

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		respondWithErrorCode(writer, "failed to get user ID from context", http.StatusUnauthorized)
		return
	}
	planUUID, err := uuid.Parse(planID)
	if err != nil {
		respondWithErrorCode(writer, "failed to get user ID from context", http.StatusUnauthorized)
		return
	}

	result, err := handler.savingsService.GetSavedSavingsPlan(context.Background(), planUUID, userUUID)
	if err != nil {
		respondWithErrorCode(writer, fmt.Sprintf("attempt to fetch plan %v by user %v", planUUID, userUUID), http.StatusUnauthorized)
	}

	respondWithJSON(writer, mapper.ToSavingsResponse(result), http.StatusOK)
}
