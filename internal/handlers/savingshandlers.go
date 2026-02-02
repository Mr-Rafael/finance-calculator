package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/calculator"
	"github.com/go-playground/validator"
	"github.com/shopspring/decimal"
)

type SavingsRequestParams struct {
	StartingCapital     int    `json:"startingCapital" validate:"required"`
	YearlyInterestRate  string `json:"yearlyInterestRate" validate:"required"`
	MonthlyContribution int    `json:"monthlyContribution" validate:"required"`
	DurationYears       int    `json:"durationYears" validate:"required"`
	TaxRate             string `json:"taxRate"`
	YearlyInflationRate string `json:"yearlyInflationRate"`
	StartDate           string `json:"startDate" validate:"required"`
}

func (cfg *ApiConfig) HandlerCalculateGet(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	reqParams := SavingsRequestParams{}
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

	calculatorParams, err := getSavingsInfo(reqParams)
	if err != nil {
		respondWithError(writer, err.Error(), fmt.Sprintf("Parse error: %v", err), http.StatusBadRequest)
		return
	}

	response, err := calculator.CalculateSavingsPlan(calculatorParams)

	respondWithJSON(writer, response, http.StatusOK)
}

func getSavingsInfo(params SavingsRequestParams) (calculator.SavingsInfo, error) {
	var err error

	taxRate := decimal.NewFromInt(0)
	if len(params.TaxRate) > 0 {
		taxRate, err = decimal.NewFromString(params.TaxRate)
		if err != nil {
			return calculator.SavingsInfo{}, fmt.Errorf("failed to parse amount %v to decimal: %v", params.TaxRate, err)
		}
	}

	inflationRate := decimal.NewFromInt(0)
	if len(params.YearlyInflationRate) > 0 {
		inflationRate, err = decimal.NewFromString(params.YearlyInflationRate)
		if err != nil {
			return calculator.SavingsInfo{}, fmt.Errorf("failed to parse amount %v to decimal: %v", params.YearlyInflationRate, err)
		}
	}

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, params.StartDate)
	if err != nil {
		return calculator.SavingsInfo{}, fmt.Errorf("failed to parse date %v to date: %v", params.StartDate, err)
	}

	return calculator.SavingsInfo{
		Capital:             params.StartingCapital,
		YearlyInterestRate:  params.YearlyInterestRate,
		MonthlyContribution: decimal.NewFromInt(int64(params.MonthlyContribution)),
		DurationYears:       decimal.NewFromInt(int64(params.DurationYears)),
		TaxRate:             taxRate.Div(decimal.NewFromInt(100)),
		InflationRate:       inflationRate.Div(decimal.NewFromInt(100)),
		StartDate:           startDate,
	}, nil
}
