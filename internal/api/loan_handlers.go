package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mr-Rafael/finance-calculator/internal/dto"
	"github.com/Mr-Rafael/finance-calculator/internal/mapper"
	"github.com/Mr-Rafael/finance-calculator/internal/service"
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

	result, err := handler.loanService.GetLoanPaymentPlan(mapper.ToLoanInput(reqParams))
	if err != nil {
		respondWithError(writer, fmt.Sprintf("Error calculating loan payment plan: %v", err), fmt.Sprintf("Error calculating loan payment plan: %v", err), http.StatusInternalServerError)
	}
	respondWithJSON(writer, mapper.ToLoanResponse(result), http.StatusOK)
}
