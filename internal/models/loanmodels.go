package models

import "time"

type LoanRequestParams struct {
	StartingPrincipal  int    `json:"startingPrincipal" validate:"required"`
	YearlyInterestRate string `json:"yearlyInterestRate" validate:"required"`
	MonthlyPayment     int    `json:"monthlyPayment" validate:"required"`
	EscrowPayment      int    `json:"otherExpenditures"`
	StartDate          string `json:"startDate" validate:"required"`
}

type LoanPaymentPlan struct {
	Plan []LoanStatus `json:"plan"`
}

type LoanStatus struct {
	Date          time.Time `json:"date"`
	Principal     int       `json:"principal"`
	Interest      int       `json:"interest"`
	Payment       int       `json:"payment"`
	EscrowPayment int       `json:"escrowPayment"`
	Paydown       int       `json:"paydown"`
}
