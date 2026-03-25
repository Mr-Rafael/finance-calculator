package service

import (
	"fmt"
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/repository"
	"github.com/shopspring/decimal"
)

type LoansService struct {
	loansRepo *repository.LoansRepo
}

func NewLoansService(repo *repository.LoansRepo) *LoansService {
	return &LoansService{loansRepo: repo}
}

type LoansInput struct {
	StartingPrincipal  int
	YearlyInterestRate string
	MonthlyPayment     int
	EscrowPayment      int
	StartDate          string
}

type LoanPaymentPlan struct {
	StartingPrincipal   decimal.Decimal
	CurrentPrincipal    decimal.Decimal
	InterestMultiplierM decimal.Decimal
	PaymentM            decimal.Decimal
	EscrowM             decimal.Decimal
	Date                time.Time
	DurationMonths      decimal.Decimal
	TotalExpenditure    decimal.Decimal
	TotalPaid           decimal.Decimal
	CostOfCreditPercent decimal.Decimal
	Plan                []LoanStatus
}

type LoanStatus struct {
	Date          time.Time
	Payment       int
	Interest      int
	OtherPayments int
	Paydown       int
	Principal     int
}

func (s *LoansService) GetLoanPaymentPlan(input LoansInput) (LoanPaymentPlan, error) {
	plan, err := initializePaymentPlan(input)
	if err != nil {
		return LoanPaymentPlan{}, fmt.Errorf("failed to initialize the payment plan struct: %v", err)
	}
	return plan, nil
}

func initializePaymentPlan(LoansInput) (LoanPaymentPlan, error) {
	return LoanPaymentPlan{}, nil
}
