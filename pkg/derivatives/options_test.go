package derivatives

import (
	"testing"
	"github.com/shopspring/decimal"
)

func TestOptionPayoffAndNetProfit(t *testing.T) {
	// Skenario Call Option dengan Strike 100 dan Premi 5
	callOpt, err := NewOptionContract(Call, decimal.NewFromFloat(100), decimal.NewFromFloat(5))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Kasus 1: Spot ada di bawah strike (S_T = 90) -> Opsi expire worthless, rugi sebesar premi (-5)
	spotLoss := decimal.NewFromFloat(90)
	if !callOpt.CalculatePayoff(spotLoss).Equal(decimal.Zero) {
		t.Errorf("expected zero payoff, got %v", callOpt.CalculatePayoff(spotLoss))
	}
	expectedNetLoss := decimal.NewFromFloat(-5)
	if !callOpt.CalculateNetProfit(spotLoss).Equal(expectedNetLoss) {
		t.Errorf("expected net profit %v, got %v", expectedNetLoss, callOpt.CalculateNetProfit(spotLoss))
	}

	// Kasus 2: Spot di atas strike (S_T = 120) -> Payoff = 20, Net Profit = 20 - 5 = 15
	spotProfit := decimal.NewFromFloat(120)
	expectedPayoff := decimal.NewFromFloat(20)
	if !callOpt.CalculatePayoff(spotProfit).Equal(expectedPayoff) {
		t.Errorf("expected payoff %v, got %v", expectedPayoff, callOpt.CalculatePayoff(spotProfit))
	}
	expectedNetProfit := decimal.NewFromFloat(15)
	if !callOpt.CalculateNetProfit(spotProfit).Equal(expectedNetProfit) {
		t.Errorf("expected net profit %v, got %v", expectedNetProfit, callOpt.CalculateNetProfit(spotProfit))
	}
}