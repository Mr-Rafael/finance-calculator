package calculator

import (
	"fmt"
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/models"
	"github.com/shopspring/decimal"
)

const minLoanCents = "1"
const maxLoanCents = "100000000000"
const minInterestRate = "0"
const maxInterestRate = "100"
const minMonthlyPaymentCents = "1"
const maxMonthlyPaymentCents = "100000000000"
const minEscrowCents = "0"
const maxEscrowCents = "100000000000"
const maxPaymentYears = 30

type LoanInfo struct {
	startingPrincipal   decimal.Decimal
	monthlyInterestRate decimal.Decimal
	monthlyPayment      decimal.Decimal
	escrowPayment       decimal.Decimal
	startDate           time.Time
}

func getLoanInfoFromRequest(request models.LoanRequestParams) (LoanInfo, error) {
	info := LoanInfo{}
	aHundred := decimal.NewFromInt(100)

	info.startingPrincipal = decimal.NewFromInt(int64(request.StartingPrincipal))
	if !decimalIsBetween(info.startingPrincipal, minLoanCents, maxLoanCents) {
		return LoanInfo{}, fmt.Errorf("invalid starting principal: '%v'. the accepted range is 0.01 - 1,000,000,000", info.startingPrincipal.Div(aHundred).Round(2))
	}

	if !stringNumberBetween(request.YearlyInterestRate, minInterestRate, maxInterestRate) {
		return LoanInfo{}, fmt.Errorf("invalid interest rate: '%v'. the accepted range is 0%% - 100%%", request.YearlyInterestRate)
	}
	monthlyInterestRate, err := getMonthlyAPRMultiplier(request.YearlyInterestRate)
	if err != nil {
		return LoanInfo{}, fmt.Errorf("invalid interest rate: '%v'", request.YearlyInterestRate)
	}
	info.monthlyInterestRate = monthlyInterestRate

	info.monthlyPayment = decimal.NewFromInt(int64(request.MonthlyPayment))
	if !decimalIsBetween(info.monthlyPayment, minMonthlyPaymentCents, maxMonthlyPaymentCents) {
		return LoanInfo{}, fmt.Errorf("invalid monthly payments: '%v'. the accepted range is 0.01 - 1,000,000,000", info.monthlyPayment.Div(aHundred).Round(2))
	}

	info.escrowPayment = decimal.NewFromInt(int64(request.EscrowPayment))
	if !decimalIsBetween(info.escrowPayment, minEscrowCents, maxEscrowCents) {
		return LoanInfo{}, fmt.Errorf("invalid escrow payment: '%v'. the accepted range is 0.01 - 1,000,000,000", info.escrowPayment.Div(aHundred).Round(2))
	}

	startDate, err := time.Parse("2006-01-02", request.StartDate)
	if err != nil {
		return LoanInfo{}, fmt.Errorf("invalid start date: %v", request.StartDate)
	}
	info.startDate = startDate

	return info, nil
}

func CalculateLoanPaymentPlan(info models.LoanRequestParams) (models.LoanPaymentPlan, error) {
	loanInfo, err := getLoanInfoFromRequest(info)
	if err != nil {
		return models.LoanPaymentPlan{}, err
	}

	plan := models.LoanPaymentPlan{}
	currentPrincipal := loanInfo.startingPrincipal
	totalPaid := decimal.NewFromInt(0)
	totalExpenditure := decimal.NewFromInt(0)
	i := 0

	isPaymentEnough, minPayment := loanInfo.isPaymentEnough()
	if !isPaymentEnough {
		return models.LoanPaymentPlan{}, fmt.Errorf("entered payment amount does not cover the first month's interest and expenditures (%v). please enter a higher payment amount",
			minPayment.Div(decimal.NewFromInt(100)).Round(2).IntPart())
	}

	for currentPrincipal.Compare(decimal.NewFromInt(0)) == 1 && i < (maxPaymentYears*12) {
		currentPayment := loanInfo.monthlyPayment
		currentInterest := currentPrincipal.Mul(loanInfo.monthlyInterestRate)
		currentExpenditure := currentInterest.Add(loanInfo.escrowPayment)
		totalExpenditure = totalExpenditure.Add(currentExpenditure)
		currentPaydown := loanInfo.monthlyPayment.Sub(currentExpenditure)
		currentPrincipal = currentPrincipal.Sub(currentPaydown)
		if currentPrincipal.Compare(decimal.NewFromInt(0)) == -1 {
			currentPayment = currentPayment.Add(currentPrincipal)
			currentPaydown = currentPaydown.Add(currentPrincipal)
			currentPrincipal = decimal.NewFromInt(0)
		}
		totalPaid = totalPaid.Add(currentPayment)
		i++
		currentStatus := models.LoanStatus{
			Date:          loanInfo.startDate.AddDate(0, i, 0),
			Principal:     int(currentPrincipal.Round(0).IntPart()),
			Interest:      int(currentInterest.Round(0).IntPart()),
			Payment:       int(currentPayment.Round(0).IntPart()),
			EscrowPayment: int(loanInfo.escrowPayment.Round(0).IntPart()),
			Paydown:       int(currentPaydown.Round(0).IntPart()),
		}
		plan.Plan = append(plan.Plan, currentStatus)
	}
	if currentPrincipal.GreaterThan(decimal.Zero) {
		return models.LoanPaymentPlan{}, fmt.Errorf("loan term surpasses the accepted limit (%v years), with a remaining %v principal. please enter a higher monthly payment.",
			maxPaymentYears,
			currentPrincipal.Round(0).IntPart())
	}

	plan.DurationMonths = i
	plan.TotalExpenditure = int(totalExpenditure.Round(0).IntPart())
	plan.TotalPaid = int(totalPaid.Round(0).IntPart())
	plan.CostOfCreditPercent = getReturnPercent(totalPaid.Div(loanInfo.startingPrincipal))
	return plan, nil
}

func (info LoanInfo) isPaymentEnough() (bool, decimal.Decimal) {
	minFirstPayment := info.startingPrincipal.Mul(info.monthlyInterestRate).Add(info.escrowPayment)
	isEnough := !minFirstPayment.GreaterThanOrEqual(info.monthlyPayment)
	return isEnough, minFirstPayment
}
