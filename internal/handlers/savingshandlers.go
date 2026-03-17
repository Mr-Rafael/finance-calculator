package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/calculator"
	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/Mr-Rafael/finance-calculator/internal/models"
	"github.com/Mr-Rafael/finance-calculator/internal/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

func (cfg *ApiConfig) HandlerSavingsCalculatePost(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	reqParams := models.SavingsRequestParams{}
	err := decoder.Decode(&reqParams)
	if err != nil {
		respondWithErrorCode(writer, "received bad savings request", http.StatusBadRequest)
		return
	}

	response, err := calculator.CalculateSavingsPlan(reqParams)
	if err != nil {
		respondWithError(writer, err.Error(), fmt.Sprintf("error calculating the savings plan: %v", err), http.StatusBadRequest)
		return
	}

	respondWithJSON(writer, response, http.StatusOK)
}

func (cfg *ApiConfig) HandlerSavingsSavePost(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	reqParams := models.SavingsSaveRequestParams{}
	userIDString, ok := GetUserID(request.Context())
	if !ok {
		respondWithErrorCode(writer, "failed to extract used ID from access token.", http.StatusUnauthorized)
		return
	}
	userID, err := utils.StringToUUID(userIDString)
	if err != nil {
		respondWithErrorCode(writer, fmt.Sprintf("failed to extract used ID from access token: %v", err), http.StatusUnauthorized)
		return
	}

	err = decoder.Decode(&reqParams)
	if err != nil {
		respondWithErrorCode(writer, "received bad savings request", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("2006-01-02", reqParams.StartDate)
	if err != nil {
		respondWithError(writer, fmt.Sprintf("Invalid start date: %v", err), fmt.Sprintf("Invalid start date '%v'", reqParams.StartDate), http.StatusBadRequest)
		return
	}
	savingsPlan, err := calculator.CalculateSavingsPlan(saveReqToCalcReq(reqParams))
	if err != nil {
		respondWithError(writer, fmt.Sprintf("Error calculating savings plan: %v", err), fmt.Sprintf("Error calculating savings plan: %v", err), http.StatusInternalServerError)
		return
	}

	queryParams := db.CreateSavingsParams{UserID: userID,
		Name:                reqParams.Name,
		StartingCapital:     int32(reqParams.StartingCapital),
		YearlyInterestRate:  reqParams.YearlyInterestRate,
		InterestRateType:    reqParams.InterestRateType,
		MonthlyContribution: int32(reqParams.MonthlyContribution),
		DurationYears:       int32(reqParams.DurationYears),
		TaxRate:             reqParams.TaxRate,
		YearlyInflationRate: pgtype.Text{
			String: reqParams.YearlyInflationRate,
			Valid:  true,
		},
		StartDate: pgtype.Timestamptz{
			Time:  startDate,
			Valid: true,
		},
		MonthlyInterestRate:   savingsPlan.MonthlyInterestRate,
		TotalInterestEarnings: int32(savingsPlan.TotalInterestEarnings),
		RateOfReturn:          savingsPlan.RateOfReturn,
		InflationAdjustedRor:  savingsPlan.InflationAdjustedROR,
	}
	queryResult, err := cfg.Queries.CreateSavings(context.Background(), queryParams)
	if err != nil {
		respondWithError(writer, fmt.Sprintf("error saving savings data to database: %v", err), "Error saving the savings data on database.", http.StatusInternalServerError)
		return
	}

	respondWithJSON(writer, queryResult, http.StatusOK)
}

func saveReqToCalcReq(originalRequest models.SavingsSaveRequestParams) models.SavingsRequestParams {
	savingsRequest := models.SavingsRequestParams{
		StartingCapital:     originalRequest.StartingCapital,
		YearlyInterestRate:  originalRequest.YearlyInterestRate,
		InterestRateType:    originalRequest.InterestRateType,
		MonthlyContribution: originalRequest.MonthlyContribution,
		DurationYears:       originalRequest.DurationYears,
		TaxRate:             originalRequest.TaxRate,
		YearlyInflationRate: originalRequest.YearlyInflationRate,
		StartDate:           originalRequest.StartDate,
	}
	return savingsRequest
}
