package service

import (
	"log"
	"strings"
	"testing"

	"github.com/Mr-Rafael/finance-calculator/internal/domain"
)

func TestCalculateLoanPaymentPlan(t *testing.T) {
	mockLoansRepo := &MockLoansRepo{}
	service := NewLoansService(mockLoansRepo)

	input := domain.LoansInput{
		StartingPrincipal:  10000000,
		YearlyInterestRate: "5",
		MonthlyPayment:     900076,
		EscrowPayment:      10000,
		StartDate:          "1970-01-01",
	}

	got, err := service.CalculateLoanPaymentPlan(input)
	if err != nil {
		log.Fatalf("Error calculating the loan payment plan.")
	}

	want := domain.LoanPaymentPlan{
		DurationMonths: 12,
	}

	if got.DurationMonths != want.DurationMonths {
		log.Fatalf("Expected a duration of %v months, but got %v.", want.DurationMonths, got.DurationMonths)
	}
}

func TestCalculateLoanMaxTerm(t *testing.T) {
	mockLoansRepo := &MockLoansRepo{}
	service := NewLoansService(mockLoansRepo)
	maxLoanTerm := 360

	input := domain.LoansInput{
		StartingPrincipal:  10000000,
		YearlyInterestRate: "5",
		MonthlyPayment:     63683,
		EscrowPayment:      10000,
		StartDate:          "1970-01-01",
	}

	got, err := service.CalculateLoanPaymentPlan(input)
	if err != nil {
		log.Fatalf("Error calculating the loan payment plan.")
	}

	want := domain.LoanPaymentPlan{
		DurationMonths: maxLoanTerm,
	}

	if got.DurationMonths != want.DurationMonths {
		log.Fatalf("Expected a duration of %v months, but got %v.", want.DurationMonths, got.DurationMonths)
	}
}

func TestCalculateLoanMinTerm(t *testing.T) {
	mockLoansRepo := &MockLoansRepo{}
	service := NewLoansService(mockLoansRepo)
	minLoanTerm := 1

	input := domain.LoansInput{
		StartingPrincipal:  1000000,
		YearlyInterestRate: "5",
		MonthlyPayment:     1014167,
		EscrowPayment:      10000,
		StartDate:          "1970-01-01",
	}

	got, err := service.CalculateLoanPaymentPlan(input)
	if err != nil {
		log.Fatalf("Error calculating the loan payment plan.")
	}

	want := domain.LoanPaymentPlan{
		DurationMonths: minLoanTerm,
	}

	if got.DurationMonths != want.DurationMonths {
		log.Fatalf("Expected a duration of %v month, but got %v.", want.DurationMonths, got.DurationMonths)
	}
}

func TestCalculateLoanTermTooLong(t *testing.T) {
	mockLoansRepo := &MockLoansRepo{}
	service := NewLoansService(mockLoansRepo)

	input := domain.LoansInput{
		StartingPrincipal:  10000000,
		YearlyInterestRate: "5",
		MonthlyPayment:     63682,
		EscrowPayment:      10000,
		StartDate:          "1970-01-01",
	}

	_, err := service.CalculateLoanPaymentPlan(input)
	if err == nil {
		log.Fatalf("Expected the loan calcuation to fail due to the term being longer than 360 months, but it didn't.")
	}
}

func TestCalculateZeroInterestAndEscrow(t *testing.T) {
	mockLoansRepo := &MockLoansRepo{}
	service := NewLoansService(mockLoansRepo)

	input := domain.LoansInput{
		StartingPrincipal:  10000000,
		YearlyInterestRate: "0",
		MonthlyPayment:     1000000,
		EscrowPayment:      0,
		StartDate:          "1970-01-01",
	}

	got, err := service.CalculateLoanPaymentPlan(input)
	if err != nil {
		log.Fatalf("Error calculating the loan payment plan.")
	}

	want := domain.LoanPaymentPlan{
		DurationMonths: 10,
	}

	if got.DurationMonths != want.DurationMonths {
		log.Fatalf("Expected a duration of %v months, but got %v.", want.DurationMonths, got.DurationMonths)
	}
}

func TestCalculateLoanPrincipalTooHigh(t *testing.T) {
	mockLoansRepo := &MockLoansRepo{}
	service := NewLoansService(mockLoansRepo)

	input := domain.LoansInput{
		StartingPrincipal:  100000000001,
		YearlyInterestRate: "1",
		MonthlyPayment:     1,
		EscrowPayment:      0,
		StartDate:          "1970-01-01",
	}

	_, err := service.CalculateLoanPaymentPlan(input)
	if err == nil {
		log.Fatalf("Expected the loan calculation to fail due to starting principal being larger than the accepted amount (100000000000), but it didn't.")
	}
}

func TestCalculateLoanInterestTooHigh(t *testing.T) {
	mockLoansRepo := &MockLoansRepo{}
	service := NewLoansService(mockLoansRepo)

	input := domain.LoansInput{
		StartingPrincipal:  100,
		YearlyInterestRate: "101",
		MonthlyPayment:     1,
		EscrowPayment:      0,
		StartDate:          "1970-01-01",
	}

	_, err := service.CalculateLoanPaymentPlan(input)
	if err == nil {
		log.Fatalf("Expected the loan calculation to fail due to the interest rate being higher than valid percent (100%% yearly), but it didn't.")
	}
}

func TestCalculateLoanMonthlyPaymentTooLow(t *testing.T) {
	mockLoansRepo := &MockLoansRepo{}
	service := NewLoansService(mockLoansRepo)

	input := domain.LoansInput{
		StartingPrincipal:  10000000,
		YearlyInterestRate: "5",
		MonthlyPayment:     51666,
		EscrowPayment:      10000,
		StartDate:          "1970-01-01",
	}

	_, err := service.CalculateLoanPaymentPlan(input)
	if !strings.Contains(err.Error(), "not enough to cover interest and escrow payment") {
		log.Fatalf("Expected the loan calculation to fail due to the monthly payment (%v cents) not even covering interest and escrow, but it didn't.", input.MonthlyPayment)
	}
}

func TestCalculateLoanEscrowTooHigh(t *testing.T) {
	mockLoansRepo := &MockLoansRepo{}
	service := NewLoansService(mockLoansRepo)

	input := domain.LoansInput{
		StartingPrincipal:  100,
		YearlyInterestRate: "1",
		MonthlyPayment:     1,
		EscrowPayment:      100000000001,
		StartDate:          "1970-01-01",
	}

	_, err := service.CalculateLoanPaymentPlan(input)
	if err == nil {
		log.Fatalf("Expected the loan calculation to fail due to escrow payment being higher than the valid amount (100000000000 cents), but it didn't.")
	}
}

func TestCalculateLoanInvalidDateFormat(t *testing.T) {
	mockLoansRepo := &MockLoansRepo{}
	service := NewLoansService(mockLoansRepo)

	input := domain.LoansInput{
		StartingPrincipal:  100,
		YearlyInterestRate: "1",
		MonthlyPayment:     1,
		EscrowPayment:      100000000001,
		StartDate:          "01/01/1970",
	}

	_, err := service.CalculateLoanPaymentPlan(input)
	if err == nil {
		log.Fatalf("Expected the loan calculation to fail due to invalid start date format, but it didn't.")
	}
}
