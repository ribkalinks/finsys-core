package money

import (
	"testing"
)

func TestCalculateFeeAndFX(t *testing.T) {
    // Fee 1.5% dari USD 100.00 -> 1.5
    usd, _ := NewMoney("100.00", "USD")
    fee, err := usd.CalculateFee("1.5")
    if err != nil || fee.Amount() != "1.5" {
        t.Errorf("Expected fee 1.5, got %s", fee.Amount())
    }

    // Konversi USD 100 ke IDR dengan rate 15500.50 -> 1550050 IDR
    idr, err := usd.ConvertFX("IDR", "15500.50")
    if err != nil || idr.Amount() != "1550050" {
        t.Errorf("Expected 1550050 IDR, got %s", idr.Amount())
    }
}