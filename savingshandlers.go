package main

import (
	"encoding/json"
	"net/http"
)

type SavingsRequestParams struct {
	StartingCapital     int    `json:"starting_capital"`
	InterestRate        string `json:"interest_rate"`
	MonthlyContribution int    `json:"monthly_contribution"`
	Duration            int    `json:"duration"`
	StartDate           string `json:"start_date"`
}

type SavingsPlan struct {
	Plan []SavingsStatus `json:"plan"`
}

type SavingsStatus struct {
	Capital      int `json:"capital"`
	Interest     int `json:"interest"`
	Contribution int `json:"contribution"`
	Increase     int `json:"increase"`
}

func (c *apiConfig) handlerCalculateGet(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	reqParams := SavingsRequestParams{}
	err := decoder.Decode(&reqParams)
	if err != nil {
		respondWithErrorCode(writer, "received bad savings request", http.StatusBadRequest)
	}
	respondWithJSON(writer, reqParams, http.StatusOK)
}
