package models

import "time"

type SavingsRequestParams struct {
	StartingCapital     int    `json:"startingCapital" validate:"required"`
	YearlyInterestRate  string `json:"yearlyInterestRate" validate:"required"`
	MonthlyContribution int    `json:"monthlyContribution" validate:"required"`
	DurationYears       int    `json:"durationYears" validate:"required"`
	TaxRate             string `json:"taxRate"`
	YearlyInflationRate string `json:"yearlyInflationRate"`
	StartDate           string `json:"startDate" validate:"required"`
}

type SavingsPlan struct {
	TotalInterestEarnings int             `json:"totalEarnings"`
	RateOfReturn          string          `json:"rateOfReturn"`
	InflationAdjustedROR  string          `json:"inflationAdjustedROR"`
	Plan                  []SavingsStatus `json:"plan"`
}

type SavingsStatus struct {
	Date         time.Time
	Interest     int `json:"interest"`
	Tax          int `json:"tax"`
	Contribution int `json:"contribution"`
	Increase     int `json:"increase"`
	Capital      int `json:"capital"`
}
