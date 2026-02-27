package models

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

type SavingsPlan struct {
	MonthlyInterestRate   string          `json:"monthlyInterestRate"`
	TotalInterestEarnings int             `json:"totalEarnings"`
	RateOfReturn          string          `json:"rateOfReturn"`
	InflationAdjustedROR  string          `json:"inflationAdjustedROR"`
	Plan                  []SavingsStatus `json:"plan"`
}

type SavingsStatus struct {
	Date         time.Time `json:"date"`
	Interest     int       `json:"interest"`
	Tax          int       `json:"tax"`
	Contribution int       `json:"contribution"`
	Increase     int       `json:"increase"`
	Capital      int       `json:"capital"`
}
