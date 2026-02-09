package models

import "time"

type LoanRequestParams struct {
	StartingPrincipal  int    `json:"startingPrincipal" validate:"required"`
	YearlyInterestRate string `json:"yearlyInterestRate" validate:"required"`
	MonthlyPayment     int    `json:"monthlyPayment" validate:"required"`
	EscrowPayment      int    `json:"escrowPayment"`
	StartDate          string `json:"startDate" validate:"required"`
}

type LoanPaymentPlan struct {
	DurationMonths      int          `json:"durationMonths"`
	TotalExpenditure    int          `json:"totalExpenditure"`
	TotalPaid           int          `json:"totalPaid"`
	CostOfCreditPercent string       `json:"costOfCreditPercent"`
	Plan                []LoanStatus `json:"plan"`
}

type LoanStatus struct {
	Date          time.Time `json:"date"`
	Payment       int       `json:"payment"`
	Interest      int       `json:"interest"`
	EscrowPayment int       `json:"escrowPayment"`
	Paydown       int       `json:"paydown"`
	Principal     int       `json:"principal"`
}
