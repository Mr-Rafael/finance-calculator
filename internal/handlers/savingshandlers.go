package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/calculator"
	"github.com/shopspring/decimal"
)

type SavingsRequestParams struct {
	StartingCapital int    `json:"starting_capital"`
	InterestRate    string `json:"interest_rate"`
	Contribution    int    `json:"contribution"`
	Duration        int    `json:"duration"`
	StartDate       string `json:"start_date"`
}

func (cfg *ApiConfig) HandlerCalculateGet(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	reqParams := SavingsRequestParams{}
	err := decoder.Decode(&reqParams)
	if err != nil {
		respondWithErrorCode(writer, "received bad savings request", http.StatusBadRequest)
	}

	calculatorParams, err := getSavingsInfo(reqParams)
	if err != nil {
		respondWithError(writer, err.Error(), fmt.Sprintf("Parse error: %v", err), http.StatusBadRequest)
	}

	response := calculator.CalculateSavingsPlan(calculatorParams)

	respondWithJSON(writer, response, http.StatusOK)
}

func getSavingsInfo(params SavingsRequestParams) (calculator.SavingsInfo, error) {
	interestRate, err := decimal.NewFromString(params.InterestRate)
	if err != nil {
		return calculator.SavingsInfo{}, fmt.Errorf("failed to parse amount %v to decimal: %v", params.InterestRate, err)
	}

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, params.StartDate)
	if err != nil {
		return calculator.SavingsInfo{}, fmt.Errorf("failed to parse date %v to date: %v", params.StartDate, err)
	}

	return calculator.SavingsInfo{
		Capital:      decimal.NewFromInt(int64(params.StartingCapital)),
		InterestRate: interestRate.Div(decimal.NewFromInt(100)),
		Contribution: decimal.NewFromInt(int64(params.Contribution)),
		Duration:     decimal.NewFromInt(int64(params.Duration)),
		StartDate:    startDate,
	}, nil
}
