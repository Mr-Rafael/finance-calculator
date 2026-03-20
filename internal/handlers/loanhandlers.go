package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mr-Rafael/finance-calculator/internal/calculator"
	"github.com/Mr-Rafael/finance-calculator/internal/dto"
)

func (cfg *ApiConfig) HandlerLoansCalculatePost(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	reqParams := dto.LoanRequestParams{}
	err := decoder.Decode(&reqParams)
	if err != nil {
		respondWithErrorCode(writer, "received bad savings request", http.StatusBadRequest)
		return
	}

	response, err := calculator.CalculateLoanPaymentPlan(reqParams)
	if err != nil {
		respondWithError(writer, err.Error(), fmt.Sprintf("error calculating the payment plan: %v", err), http.StatusBadRequest)
		return
	}

	respondWithJSON(writer, response, http.StatusOK)
}
