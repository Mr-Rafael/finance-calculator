package service

import (
	"log"
	"testing"

	"github.com/Mr-Rafael/finance-calculator/internal/domain"
	"github.com/shopspring/decimal"
)

func TestCalculateSimpleSavingsPlan(t *testing.T) {
	mockSavingsRepo := &MockSavingsRepo{}
	service := NewSavingsService(mockSavingsRepo)

	input := domain.SavingsInput{
		StartingCapital:     10000000,
		YearlyInterestRate:  "5",
		InterestRateType:    "APY",
		MonthlyContribution: 10000,
		DurationYears:       1,
		TaxRate:             "0",
		YearlyInflationRate: "0",
		StartDate:           "1970-01-01",
	}

	got, err := service.CalculateSavingsPlan(input)
	if err != nil {
		log.Fatalf("Error calculating the Savings Plan: %v", err)
	}

	want := domain.SavingsPlan{
		CurrentCapital: decimal.NewFromInt32(10622726),
	}

	if got.CurrentCapital.Round(0).Compare(want.CurrentCapital) != 0 {
		log.Fatalf("Expected calculated principal (%v) to match expected principal (%v) at the end of loan, but it didn't.", got.CurrentCapital.Round(0), want.CurrentCapital)
	}
}
