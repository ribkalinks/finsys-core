package derivatives

import (
	"errors"
	"github.com/shopspring/decimal"
)

// OptionType mendefinisikan jenis opsi (Call atau Put)
type OptionType string

const (
	Call OptionType = "CALL"
	Put  OptionType = "PUT"
)

// OptionContract merepresentasikan parameter kontrak opsi
type OptionContract struct {
	Type        OptionType      // CALL atau PUT
	StrikePrice decimal.Decimal // K: Strike price
	Premium     decimal.Decimal // Harga premi opsi yang dibayar di awal
}

// NewOptionContract membuat instance opsi baru dengan validasi
func NewOptionContract(optType OptionType, strike, premium decimal.Decimal) (*OptionContract, error) {
	if optType != Call && optType != Put {
		return nil, errors.New("invalid option type: must be CALL or PUT")
	}
	if strike.LessThanOrEqual(decimal.Zero) {
		return nil, errors.New("strike price must be greater than zero")
	}
	if premium.LessThan(decimal.Zero) {
		return nil, errors.New("option premium cannot be negative")
	}

	return &OptionContract{
		Type:        optType,
		StrikePrice: strike,
		Premium:     premium,
	}, nil
}

// CalculatePayoff menghitung gross payoff saat jatuh tempo berdasarkan harga spot akhir (S_T)
func (o *OptionContract) CalculatePayoff(spotAtMaturity decimal.Decimal) decimal.Decimal {
	if o.Type == Call {
		// Call Payoff = max(S_T - K, 0)
		diff := spotAtMaturity.Sub(o.StrikePrice)
		if diff.IsPositive() {
			return diff
		}
		return decimal.Zero
	} else {
		// Put Payoff = max(K - S_T, 0)
		diff := o.StrikePrice.Sub(spotAtMaturity)
		if diff.IsPositive() {
			return diff
		}
		return decimal.Zero
	}
}

// CalculateNetProfit menghitung keuntungan/kerugian bersih setelah dikurangi premi
// Net Profit = Gross Payoff - Premium
func (o *OptionContract) CalculateNetProfit(spotAtMaturity decimal.Decimal) decimal.Decimal {
	grossPayoff := o.CalculatePayoff(spotAtMaturity)
	return grossPayoff.Sub(o.Premium).Round(4)
}