package money

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestDurationAndConvexity(t *testing.T) {
	// Contoh arus kas obligasi sederhana (misal: kupon 5 per tahun selama 2 tahun + principal 100)
	cashFlows := []CashFlow{
		{Time: decimal.NewFromFloat(1.0), Amount: decimal.NewFromInt(5)},
		{Time: decimal.NewFromFloat(2.0), Amount: decimal.NewFromInt(105)},
	}

	yield := decimal.NewFromFloat(0.05) // 5% yield
	price := decimal.NewFromFloat(100.0) // Asumsi harga pasar 100

	// 1. Test Macaulay Duration
	macDuration, err := CalculateMacaulayDuration(cashFlows, yield, price)
	if err != nil {
		t.Fatalf("unexpected error in Macaulay Duration: %v", err)
	}
	if macDuration.LessThanOrEqual(decimal.Zero) {
		t.Errorf("expected positive Macaulay Duration, got %s", macDuration.String())
	}

	// 2. Test Modified Duration
	modDuration := CalculateModifiedDuration(macDuration, yield)
	if modDuration.LessThanOrEqual(decimal.Zero) {
		t.Errorf("expected positive Modified Duration, got %s", modDuration.String())
	}

	// 3. Test Convexity
	convexity, err := CalculateConvexity(cashFlows, yield, price)
	if err != nil {
		t.Fatalf("unexpected error in Convexity: %v", err)
	}
	if convexity.LessThanOrEqual(decimal.Zero) {
		t.Errorf("expected positive Convexity, got %s", convexity.String())
	}
}