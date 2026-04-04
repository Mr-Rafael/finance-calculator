package dto

import "time"

type SavingsRequestParams struct {
	StartingCapital     int    `json:"startingCapital"`
	YearlyInterestRate  string `json:"yearlyInterestRate"`
	InterestRateType    string `json:"interestRateType"`
	MonthlyContribution int    `json:"monthlyContribution"`
	DurationYears       int    `json:"durationYears"`
	TaxRate             string `json:"taxRate"`
	YearlyInflationRate string `json:"yearlyInflationRate"`
	StartDate           string `json:"startDate"`
}

type SavingsResponseParams struct {
	MonthlyInterestRate   string          `json:"monthlyInterestRate"`
	TotalInterestEarnings int             `json:"totalEarnings"`
	RateOfReturn          string          `json:"rateOfReturn"`
	InflationAdjustedROR  string          `json:"inflationAdjustedROR"`
	Plan                  []SavingsStatus `json:"plan"`
}

type SavedSavingsResponseParams struct {
	ID             string                `json:"id"`
	Name           string                `json:"name"`
	OriginalData   OriginalSavingsData   `json:"originalData"`
	CalculatedData CalculatedSavingsData `json:"calculatedData"`
	Plan           []SavingsStatus       `json:"plan"`
}

type OriginalSavingsData struct {
	StartingCapital     int    `json:"startingCapital"`
	YearlyInterestRate  string `json:"yearlyInterestRate"`
	InterestRateType    string `json:"interestRateType"`
	MonthlyContribution int    `json:"monthlyContribution"`
	DurationYears       int    `json:"durationYears"`
	TaxRate             string `json:"taxRate"`
	YearlyInflationRate string `json:"yearlyInflationRate"`
	StartDate           string `json:"startDate"`
}

type CalculatedSavingsData struct {
	MonthlyInterestRate   string `json:"monthlyInterestRate"`
	TotalInterestEarnings int    `json:"totalInterestEarnings"`
	RateOfReturn          string `json:"rateOfReturn"`
	InflationAdjustedROR  string `json:"inflationAdjustedROR"`
}

type SavingsStatus struct {
	Date         time.Time `json:"date"`
	Interest     int       `json:"interest"`
	Tax          int       `json:"tax"`
	Contribution int       `json:"contribution"`
	Increase     int       `json:"increase"`
	Capital      int       `json:"capital"`
}

type SavingsSaveResponseParams struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	StartingCapital       int    `json:"startingCapital"`
	YearlyInterestRate    string `json:"yearlyInterestRate"`
	InterestRateType      string `json:"interestRateType"`
	MonthlyContribution   int    `json:"monthlyContribution"`
	DurationYears         int    `json:"durationYears"`
	TaxRate               string `json:"taxRate"`
	YearlyInflationRate   string `json:"yearlyInflationRate"`
	StartDate             string `json:"startDate"`
	MonthlyInterestRate   string `json:"monthlyInterestRate"`
	TotalInterestEarnings int    `json:"totalEarnings"`
	RateOfReturn          string `json:"rateOfReturn"`
	InflationAdjustedROR  string `json:"inflationAdjustedROR"`
}

type SavingsSaveRequestParams struct {
	Name                string `json:"name"`
	StartingCapital     int    `json:"startingCapital"`
	YearlyInterestRate  string `json:"yearlyInterestRate"`
	InterestRateType    string `json:"interestRateType"`
	MonthlyContribution int    `json:"monthlyContribution"`
	DurationYears       int    `json:"durationYears"`
	TaxRate             string `json:"taxRate"`
	YearlyInflationRate string `json:"yearlyInflationRate"`
	StartDate           string `json:"startDate"`
}

type SavingsListResponseParams struct {
	Plans []SavingsInfo `json:"plans"`
}

type SavingsInfo struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	StartingCapital int    `json:"starting_capital"`
}
