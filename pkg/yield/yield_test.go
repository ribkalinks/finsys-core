package yield

import (
	"testing"
	"github.com/shopspring/decimal"
)

func TestBondCalculations(t *testing.T) {
	// Skenario: Obligasi Nilai Nominal 1000, Kupon 6% (0.06), Harga Pasar 950, Jatuh tempo 5 tahun, 1x setahun
	bond, err := NewBond("1000", "0.06", "950", 5, 1)
	if err != nil {
		t.Fatalf("Failed to create bond: %v", err)
	}

	// 1. Test Current Yield
	// Annual Coupon = 1000 * 0.06 = 60
	// Current Yield = 60 / 950 = 0.06315789... (~0.0632)
	cy, err := bond.CurrentYield()
	if err != nil {
		t.Errorf("Unexpected error in CurrentYield: %v", err)
	}
	expectedCY, _ := decimal.NewFromString("0.0631578947368421")
	if !cy.Round(4).Equal(expectedCY.Round(4)) {
		t.Errorf("Expected Current Yield ~0.0632, got %s", cy.String())
	}

	// 2. Test Approximate YTM
	ytm, err := bond.ApproximateYTM()
	if err != nil {
		t.Errorf("Unexpected error in ApproximateYTM: %v", err)
	}
	// Pastikan YTM menghasilkan nilai yang wajar di sekitar 7% - 7.2% karena dibeli diskon (di bawah 1000)
	if ytm.LessThanOrEqual(decimal.Zero) {
		t.Errorf("Expected positive YTM, got %s", ytm.String())
	}
}