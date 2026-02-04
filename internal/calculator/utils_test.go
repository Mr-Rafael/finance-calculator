package calculator

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestGetMonthlyInterestMultiplier(t *testing.T) {
	got, _ := getMonthlyInterestMultiplier("4.75")
	got = got.Round(9)
	want, _ := decimal.NewFromString("0.00387468499")
	want = want.Round(9)
	if !got.Equal(want) {
		t.Errorf("Got interest multipler = %v | Expected %v", got, want)
	}
}
