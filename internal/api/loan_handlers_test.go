package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Mr-Rafael/finance-calculator/internal/service"
	"github.com/stretchr/testify/require"
)

func TestCalculate(t *testing.T) {
	mockLoansRepo := &service.MockLoansRepo{}
	service := service.NewLoansService(mockLoansRepo)
	handler := NewLoanHandler(service)

	req := httptest.NewRequest(
		http.MethodPost,
		"/app/loans/calculate",
		strings.NewReader(`{
			"startingPrincipal": 10000000,
			"yearlyInterestRate": "5",
			"monthlyPayment": 1500000,
			"escrowPayment": 10000,
			"startDate": "1970-01-01"
		}`),
	)
	rr := httptest.NewRecorder()

	handler.HandleCalculateLoan(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}
