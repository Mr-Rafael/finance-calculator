package mapper

import (
	"time"

	"github.com/Mr-Rafael/finance-calculator/internal/db"
	"github.com/Mr-Rafael/finance-calculator/internal/domain"
	"github.com/Mr-Rafael/finance-calculator/internal/dto"
	"github.com/google/uuid"
)

func ToLoanResponse(plan domain.LoanPaymentPlan) dto.LoanResponseParams {
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

func ToSaveLoanResponse(loan db.Loan) dto.LoanSaveResponseParams {
	return dto.LoanSaveResponseParams{
		ID:                  loan.ID.String(),
		Name:                loan.Name,
		StartingPrincipal:   int(loan.StartingPrincipal),
		YearlyInterestRate:  loan.YearlyInterestRate,
		MonthlyPayment:      int(loan.MonthlyPayment),
		EscrowPayment:       int(loan.EscrowPayment),
		StartDate:           loan.StartDate.Time.Format(time.RFC3339),
		DurationMonths:      int(loan.DurationMonths),
		TotalExpenditure:    int(loan.TotalExpenditure),
		TotalPaid:           int(loan.TotalPaid),
		CostOfCreditPercent: loan.CostOfCredit,
	}
}

func ToLoanInput(input dto.LoanRequestParams) domain.LoansInput {
	loan := domain.LoansInput{
		StartingPrincipal:  input.StartingPrincipal,
		YearlyInterestRate: input.YearlyInterestRate,
		MonthlyPayment:     input.MonthlyPayment,
		EscrowPayment:      input.EscrowPayment,
		StartDate:          input.StartDate,
	}

	return loan
}

func ToSaveLoanInput(userId uuid.UUID, input dto.LoanSaveRequestParams) domain.SaveLoanInput {
	loan := domain.SaveLoanInput{
		UserID:             userId,
		LoanName:           input.Name,
		StartingPrincipal:  input.StartingPrincipal,
		YearlyInterestRate: input.YearlyInterestRate,
		MonthlyPayment:     input.MonthlyPayment,
		EscrowPayment:      input.EscrowPayment,
		StartDate:          input.StartDate,
	}
	return loan
}
