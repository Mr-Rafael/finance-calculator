package calculator

import (
	"fmt"
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/models"
	"github.com/shopspring/decimal"
)

type LoanInfo struct {
	startingPrincipal   decimal.Decimal
	monthlyInterestRate decimal.Decimal
	monthlyPayment      decimal.Decimal
	escrowPayment       decimal.Decimal
	startDate           time.Time
}

func getLoanInfoFromRequest(request models.LoanRequestParams) (LoanInfo, error) {
	info := LoanInfo{}
	info.startingPrincipal = decimal.NewFromInt(int64(request.StartingPrincipal))
	monthlyInterestRate, err := getMonthlyAPRMultiplier(request.YearlyInterestRate)
	if err != nil {
		return LoanInfo{}, fmt.Errorf("invalid interest rate: %v", request.YearlyInterestRate)
	}
	info.monthlyInterestRate = monthlyInterestRate
	info.monthlyPayment = decimal.NewFromInt(int64(request.MonthlyPayment))
	info.escrowPayment = decimal.NewFromInt(int64(request.EscrowPayment))
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
	totalExpenditure := decimal.NewFromInt(0)
	i := 0

	for currentPrincipal.Compare(decimal.NewFromInt(0)) != -1 && i < 360 {
		currentInterest := currentPrincipal.Mul(loanInfo.monthlyInterestRate)
		currentExpenditure := currentInterest.Add(loanInfo.escrowPayment)
		totalExpenditure = totalExpenditure.Add(currentExpenditure)
		currentPaydown := loanInfo.monthlyPayment.Sub(currentExpenditure)
		currentPrincipal = currentPrincipal.Sub(currentPaydown)
		i++
		currentStatus := models.LoanStatus{
			Date:          loanInfo.startDate.AddDate(0, i, 0),
			Principal:     int(currentPrincipal.Round(0).IntPart()),
			Interest:      int(currentInterest.Round(0).IntPart()),
			Payment:       int(loanInfo.monthlyPayment.Round(0).IntPart()),
			EscrowPayment: int(currentInterest.Add(loanInfo.escrowPayment).Round(0).IntPart()),
			Paydown:       int(currentPaydown.Round(0).IntPart()),
		}
		plan.Plan = append(plan.Plan, currentStatus)
	}

	return plan, nil
}
