// CalculateFee menghitung biaya berdasarkan persentase (contoh: "1.5" untuk 1.5%)
func (m Money) CalculateFee(percentStr string) (Money, error) {
	rate, err := decimal.NewFromString(percentStr)
	if err != nil {
		return Money{}, fmt.Errorf("invalid fee percentage: %w", err)
	}

	hundred := decimal.NewFromInt(100)
	feeAmount := m.amount.Mul(rate).Div(hundred)

	return NewMoney(feeAmount.String(), m.currency)
}

// ConvertFX mengonversi mata uang berdasarkan rate tertentu
func (m Money) ConvertFX(targetCurrency string, rateStr string) (Money, error) {
	rate, err := decimal.NewFromString(rateStr)
	if err != nil {
		return Money{}, fmt.Errorf("invalid FX rate: %w", err)
	}

	convertedAmount := m.amount.Mul(rate)
	return NewMoney(convertedAmount.String(), targetCurrency)
}

type AmortizationItem struct {
	Period           int
	MonthlyPayment   Money
	PrincipalPayment Money
	InterestPayment  Money
	RemainingBalance Money
}

// CalculateAmortization menghitung jadwal cicilan bulanan (Anuitas)
func CalculateAmortization(principal Money, annualRateStr string, months int) ([]AmortizationItem, error) {
	annualRate, err := decimal.NewFromString(annualRateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid interest rate: %w", err)
	}

	// Bunga bulanan = (Annual Rate / 100) / 12
	monthlyRate := annualRate.Div(decimal.NewFromInt(100)).Div(decimal.NewFromInt(12))
	
	// Rumus Anuitas: PMT = P * [r(1+r)^n] / [(1+r)^n - 1]
	one := decimal.NewFromInt(1)
	onePlusR := one.Add(monthlyRate)
	pow := onePlusR.Pow(decimal.NewFromInt(int64(months)))
	
	numerator := monthlyRate.Mul(pow)
	denominator := pow.Sub(one)
	pmtFactor := numerator.Div(denominator)
	
	monthlyPaymentAmount := principal.amount.Mul(pmtFactor)
	monthlyPayment, _ := NewMoney(monthlyPaymentAmount.String(), principal.currency)

	var schedule []AmortizationItem
	balance := principal.amount

	for i := 1; i <= months; i++ {
		interestAmount := balance.Mul(monthlyRate)
		interestPayment, _ := NewMoney(interestAmount.String(), principal.currency)

		principalAmount := monthlyPayment.amount.Sub(interestPayment.amount)
		principalPayment, _ := NewMoney(principalAmount.String(), principal.currency)

		balance = balance.Sub(principalPayment.amount)
		if i == months {
			balance = decimal.Zero // Memastikan sisa saldo 0 di akhir periode
		}
		remainingBalance, _ := NewMoney(balance.String(), principal.currency)

		schedule = append(schedule, AmortizationItem{
			Period:           i,
			MonthlyPayment:   monthlyPayment,
			PrincipalPayment: principalPayment,
			InterestPayment:  interestPayment,
			RemainingBalance: remainingBalance,
		})
	}

	return schedule, nil
}