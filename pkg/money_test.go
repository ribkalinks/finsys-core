func TestCalculateFeeAndFX(t *testing.T) {
	// Fee 1.5% dari USD 100.00 -> 1.50
	usd, _ := NewMoney("100.00", "USD")
	fee, err := usd.CalculateFee("1.5")
	if err != nil || fee.Amount() != "1.50" {
		t.Errorf("Expected fee 1.50, got %s", fee.Amount())
	}

	// Konversi USD 100 ke IDR dengan rate 15500.50 -> 1550050 IDR
	idr, err := usd.ConvertFX("IDR", "15500.50")
	if err != nil || idr.Amount() != "1550050" {
		t.Errorf("Expected 1550050 IDR, got %s", idr.Amount())
	}
}

func TestAmortizationSchedule(t *testing.T) {
	principal, _ := NewMoney("12000000", "IDR") // Pinjaman 12 Juta
	schedule, err := CalculateAmortization(principal, "12", 12) // Bunga 12% per tahun, tenor 12 bulan

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(schedule) != 12 {
		t.Errorf("Expected 12 periods, got %d", len(schedule))
	}

	// Cek apakah saldo akhir di bulan ke-12 bernilai 0
	lastPeriod := schedule[len(schedule)-1]
	if lastPeriod.RemainingBalance.Amount() != "0" {
		t.Errorf("Expected final balance 0, got %s", lastPeriod.RemainingBalance.Amount())
	}
}