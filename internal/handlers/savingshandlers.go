package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mr-Rafael/finance-calculator/internal/calculator"
	"github.com/Mr-Rafael/finance-calculator/internal/models"
	"github.com/go-playground/validator"
)

func (cfg *ApiConfig) HandlerCalculateGet(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	reqParams := models.SavingsRequestParams{}
	err := decoder.Decode(&reqParams)
	if err != nil {
		respondWithErrorCode(writer, "received bad savings request", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(reqParams); err != nil {
		respondWithError(writer, err.Error(), fmt.Sprintf("missing required fields: %v", err), http.StatusBadRequest)
		return
	}

	response, err := calculator.CalculateSavingsPlan(reqParams)
	if err != nil {
		respondWithError(writer, err.Error(), fmt.Sprintf("missing required fields: %v", err), http.StatusBadRequest)
		return
	}

	respondWithJSON(writer, response, http.StatusOK)
}
