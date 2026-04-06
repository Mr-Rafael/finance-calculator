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

type LoanHandler struct {
	loanService *service.LoansService
}

func NewLoanHandler(service *service.LoansService) *LoanHandler {
	return &LoanHandler{loanService: service}
}

func (handler *LoanHandler) HandleCalculateLoan(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	reqParams := dto.LoanRequestParams{}
	err := decoder.Decode(&reqParams)
	if err != nil {
		respondWithErrorCode(writer, "received bad savings request", http.StatusBadRequest)
		return
	}

	result, err := handler.loanService.CalculateLoanPaymentPlan(mapper.ToLoanInput(reqParams))
	if err != nil {
		respondWithError(writer, fmt.Sprintf("Error calculating loan payment plan: %v", err), fmt.Sprintf("Error calculating loan payment plan: %v", err), http.StatusInternalServerError)
	}
	respondWithJSON(writer, mapper.ToLoanResponse(result), http.StatusOK)
}

func (handler *LoanHandler) HandleSaveLoan(writer http.ResponseWriter, request *http.Request) {
	userID := request.Context().Value(userIDKey).(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		respondWithErrorCode(writer, "failed to get user ID from context", http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(request.Body)
	reqParams := dto.LoanSaveRequestParams{}
	err = decoder.Decode(&reqParams)
	if err != nil {
		respondWithErrorCode(writer, "received bad savings request", http.StatusBadRequest)
		return
	}

	result, err := handler.loanService.SaveLoanPaymentPlan(context.Background(), mapper.ToSaveLoanInput(userUUID, reqParams))
	if err != nil {
		respondWithError(writer, fmt.Sprintf("Error saving the plan: %v", err), fmt.Sprintf("Error saving the plan: %v", err), http.StatusInternalServerError)
		return
	}

	respondWithJSON(writer, mapper.ToSaveLoanResponse(result), http.StatusCreated)
}
