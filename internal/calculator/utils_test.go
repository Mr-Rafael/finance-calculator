package calculator

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestMonthlyInterestMultiplier(t *testing.T) {
	got, _ := getMonthlyInterestMultiplier("4.75")
	got = got.Round(9)
	want, _ := decimal.NewFromString("0.00387468499")
	want = want.Round(9)
	if !got.Equal(want) {
		t.Errorf("Got interest multipler = %v | Expected %v", got, want)
	}
}

func TestMonthlyInterestMultiplierError(t *testing.T) {
	_, got := getMonthlyInterestMultiplier("0.0.3")
	if got == nil {
		t.Errorf("Got nil error | Expected error")
	}
}

func TestTaxMultiplier(t *testing.T) {
	got, _ := getTaxMultiplier("12")
	want, _ := decimal.NewFromString("0.12")
	if !got.Equal(want) {
		t.Errorf("Got tax multiplier = %v | Expected %v", got, want)
	}
}

func TestTaxMultiplierError(t *testing.T) {
	_, got := getTaxMultiplier("12%")
	if got == nil {
		t.Errorf("Got nil error | Expected error")
	}
}

func TestInflationMultiplier(t *testing.T) {
	got, _ := getYearlyInflationMultiplier("6")
	want, _ := decimal.NewFromString("1.06")
	if !got.Equal(want) {
		t.Errorf("Got inflation multiplier = %v | Expected %v", got, want)
	}
}

func TestInflationMultiplierError(t *testing.T) {
	_, got := getTaxMultiplier("6%")
	if got == nil {
		t.Errorf("Got nil error | Expected error")
	}
}
