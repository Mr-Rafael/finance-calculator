package service

import (
	"log"
	"testing"

	"github.com/Mr-Rafael/finance-calculator/internal/domain"
)

func TestCalculateLoanPaymentPlan(t *testing.T) {
	mockLoansRepo := &MockLoansRepo{}
	service := NewLoansService(mockLoansRepo)

	input := domain.LoansInput{
		StartingPrincipal:  10000000,
		YearlyInterestRate: "5",
		MonthlyPayment:     865840,
		EscrowPayment:      10000,
		StartDate:          "1970-01-01",
	}

	got, err := service.CalculateLoanPaymentPlan(input)
	if err != nil {
		log.Fatalf("Error calculating the loan payment plan.")
	}

	want := domain.LoanPaymentPlan{
		DurationMonths: 13,
	}

	if got.DurationMonths != want.DurationMonths {
		log.Fatalf("Expected a duration of %v months, but got %v.", want.DurationMonths, got.DurationMonths)
	}
}
