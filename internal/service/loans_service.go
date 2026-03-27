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
	DurationMonths      int
	TotalExpenditure    decimal.Decimal
	TotalPaid           decimal.Decimal
	CostOfCreditPercent decimal.Decimal
	Plan                []LoanStatus
}

type LoanStatus struct {
	Date          time.Time
	Payment       decimal.Decimal
	Interest      decimal.Decimal
	OtherPayments decimal.Decimal
	Paydown       decimal.Decimal
	Principal     decimal.Decimal
}

const minLoanCents = "1"
const maxLoanCents = "100000000000"
const minInterestRate = "0"
const maxInterestRate = "100"
const minMonthlyPaymentCents = "1"
const maxMonthlyPaymentCents = "100000000000"
const minEscrowCents = "0"
const maxEscrowCents = "100000000000"
const maxPaymentYears = 30

func (s *LoansService) GetLoanPaymentPlan(input LoansInput) (LoanPaymentPlan, error) {
	plan, err := initializePaymentPlan(input)
	if err != nil {
		return LoanPaymentPlan{}, fmt.Errorf("failed to initialize the payment plan struct: %v", err)
	}

	i := 0
	for plan.CurrentPrincipal.Compare(decimal.Zero) == 1 {
		i++
		if i > maxPaymentYears*12 {
			remainder := plan.CurrentPrincipal.Div(decimal.NewFromInt(100)).Round(2).String()
			return LoanPaymentPlan{}, fmt.Errorf("The payment plan exceeds the year limit (%v years), with a remaining %v to pay", maxPaymentYears, remainder)
		}
		payment := plan.passMonth()
		payment = plan.generateInterest(payment)
		payment = plan.chargeEscrow(payment)
		payment = plan.makePayment(payment)
		plan.Plan = append(plan.Plan, payment)
	}
	plan.finalCalculations()

	return plan, nil
}

func initializePaymentPlan(input LoansInput) (LoanPaymentPlan, error) {
	plan := LoanPaymentPlan{}
	oneHundred := decimal.NewFromInt(100)

	startingPrincipal := decimal.NewFromInt(int64(input.StartingPrincipal))
	if !decimalIsBetween(startingPrincipal, minLoanCents, maxLoanCents) {
		return LoanPaymentPlan{}, fmt.Errorf("invalid starting principal: '%v'. the accepted range is 0.01 - 1,000,000,000", startingPrincipal.Div(oneHundred).Round(2))
	}
	plan.StartingPrincipal = startingPrincipal
	plan.CurrentPrincipal = startingPrincipal

	monthlyInterestRate, err := getMonthlyAPRMultiplier(input.YearlyInterestRate)
	if !stringNumberBetween(input.YearlyInterestRate, minInterestRate, maxInterestRate) {
		return LoanPaymentPlan{}, fmt.Errorf("invalid interest rate: '%v'. the accepted range is 0%% - 100%%", input.YearlyInterestRate)
	}
	if err != nil {
		return LoanPaymentPlan{}, fmt.Errorf("invalid interest rate: '%v'", input.YearlyInterestRate)
	}
	plan.InterestMultiplierM = monthlyInterestRate

	monthlyPayment := decimal.NewFromInt(int64(input.MonthlyPayment))
	if !decimalIsBetween(monthlyPayment, minMonthlyPaymentCents, maxMonthlyPaymentCents) {
		return LoanPaymentPlan{}, fmt.Errorf("invalid monthly payments: '%v'. the accepted range is 0.01 - 1,000,000,000", monthlyPayment.Div(oneHundred).Round(2))
	}
	plan.PaymentM = monthlyPayment

	escrow := decimal.NewFromInt(int64(input.EscrowPayment))
	if !decimalIsBetween(escrow, minEscrowCents, maxEscrowCents) {
		return LoanPaymentPlan{}, fmt.Errorf("invalid escrow payment: '%v'. the accepted range is 0.01 - 1,000,000,000", escrow.Div(oneHundred).Round(2))
	}
	plan.EscrowM = escrow

	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		return LoanPaymentPlan{}, fmt.Errorf("invalid start date: %v", input.StartDate)
	}
	plan.Date = startDate

	return plan, nil
}

func (p *LoanPaymentPlan) passMonth() LoanStatus {
	p.Date = p.Date.AddDate(0, 1, 0)
	p.DurationMonths += 1
	return LoanStatus{
		Date: p.Date,
	}
}

func (p *LoanPaymentPlan) generateInterest(s LoanStatus) LoanStatus {
	interest := p.CurrentPrincipal.Mul(p.InterestMultiplierM)
	p.CurrentPrincipal = p.CurrentPrincipal.Add(interest)
	p.TotalExpenditure = p.TotalExpenditure.Add(interest)

	s.Interest = interest
	return s
}

func (p *LoanPaymentPlan) chargeEscrow(s LoanStatus) LoanStatus {
	p.CurrentPrincipal = p.CurrentPrincipal.Add(p.EscrowM)
	p.TotalExpenditure = p.TotalExpenditure.Add(p.EscrowM)

	s.OtherPayments = p.EscrowM
	return s
}

func (p *LoanPaymentPlan) makePayment(s LoanStatus) LoanStatus {
	paydown := p.PaymentM.Sub(s.Interest).Sub(s.OtherPayments)
	if p.CurrentPrincipal.Cmp(paydown) == -1 {
		payment := p.CurrentPrincipal.Add(s.Interest).Add(s.OtherPayments)
		p.TotalPaid = p.TotalPaid.Add(payment)
		s.Payment = payment
		s.Paydown = p.CurrentPrincipal
		p.CurrentPrincipal = decimal.Zero
		s.Principal = p.CurrentPrincipal
	} else {
		p.TotalPaid = p.TotalPaid.Add(p.PaymentM)
		s.Payment = p.PaymentM
		s.Paydown = paydown
		p.CurrentPrincipal = p.CurrentPrincipal.Sub(paydown)
		s.Principal = p.CurrentPrincipal
	}

	return s
}

func (p *LoanPaymentPlan) finalCalculations() {
	p.CostOfCreditPercent = p.TotalPaid.Div(p.StartingPrincipal)
}
