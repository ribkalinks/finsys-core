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