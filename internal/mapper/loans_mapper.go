package mapper

import (
	"github.com/Mr-Rafael/finance-calculator/internal/dto"
	"github.com/Mr-Rafael/finance-calculator/internal/service"
)

func ToLoanResponse(plan service.LoanPaymentPlan) dto.LoanResponseParams {
	response := dto.LoanResponseParams{}

	response.DurationMonths = plan.DurationMonths
	response.TotalExpenditure = int(plan.TotalExpenditure.Round(0).IntPart())
	response.TotalPaid = int(plan.TotalPaid.Round(0).IntPart())
	response.CostOfCreditPercent = plan.CostOfCreditPercent.Round(2).String()
	for _, status := range plan.Plan {
		response.Plan = append(response.Plan, dto.LoanStatus{
			Date:          status.Date,
			Payment:       int(status.Payment.Round(0).IntPart()),
			Interest:      int(status.Interest.Round(0).IntPart()),
			OtherPayments: int(status.OtherPayments.Round(0).IntPart()),
			Paydown:       int(status.Paydown.Round(0).IntPart()),
			Principal:     int(status.Principal.Round(0).IntPart()),
		})
	}
	return response
}

func ToLoanInput(input dto.LoanRequestParams) service.LoansInput {
	loan := service.LoansInput{
		StartingPrincipal:  input.StartingPrincipal,
		YearlyInterestRate: input.YearlyInterestRate,
		MonthlyPayment:     input.MonthlyPayment,
		EscrowPayment:      input.EscrowPayment,
		StartDate:          input.StartDate,
	}

	return loan
}
