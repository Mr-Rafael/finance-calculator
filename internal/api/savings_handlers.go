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
	respondWithJSON(writer, result, http.StatusOK)
}
