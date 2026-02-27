package models

import "time"

type LoanRequestParams struct {
	StartingPrincipal  int    `json:"startingPrincipal"`
	YearlyInterestRate string `json:"yearlyInterestRate"`
	MonthlyPayment     int    `json:"monthlyPayment"`
	EscrowPayment      int    `json:"escrowPayment"`
	StartDate          string `json:"startDate"`
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
	OtherPayments int       `json:"otherPayments"`
	Paydown       int       `json:"paydown"`
	Principal     int       `json:"principal"`
}
