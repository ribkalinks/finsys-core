package yield

import (
	"errors"
	"github.com/shopspring/decimal"
)

// Bond represents a fixed-income bond instrument.
type Bond struct {
	FaceValue    decimal.Decimal // Nilai nominal / par value (misal: 1000)
	CouponRate   decimal.Decimal // Tingkat kupon tahunan dalam bentuk desimal (misal: 0.05 untuk 5%)
	CurrentPrice decimal.Decimal // Harga pasar obligasi saat ini
	YearsToMaturity int          // Sisa tahun menuju jatuh tempo
	PaymentsPerYear int          // Frekuensi pembayaran kupon per tahun (misal: 2 untuk semiannual, 1 untuk annual)
}

// NewBond initializes a new Bond instance with decimal validations.
func NewBond(faceValue, couponRate, currentPrice string, yearsToMaturity, paymentsPerYear int) (*Bond, error) {
	fv, err := decimal.NewFromString(faceValue)
	if err != nil {
		return nil, errors.New("invalid face value format")
	}
	cr, err := decimal.NewFromString(couponRate)
	if err != nil {
		return nil, errors.New("invalid coupon rate format")
	}
	cp, err := decimal.NewFromString(currentPrice)
	if err != nil {
		return nil, errors.New("invalid current price format")
	}

	if paymentsPerYear <= 0 {
		paymentsPerYear = 1 // Default ke annual jika tidak valid
	}

	return &Bond{
		FaceValue:       fv,
		CouponRate:      cr,
		CurrentPrice:    cp,
		YearsToMaturity: yearsToMaturity,
		PaymentsPerYear: paymentsPerYear,
	}, nil
}

// CurrentYield calculates the annual coupon payment divided by the current market price.
// Formula: Current Yield = Annual Coupon Payment / Current Market Price
func (b *Bond) CurrentYield() (decimal.Decimal, error) {
	if b.CurrentPrice.IsZero() {
		return decimal.Zero, errors.New("current price cannot be zero")
	}
	// Annual Coupon Payment = FaceValue * CouponRate
	annualCoupon := b.FaceValue.Mul(b.CouponRate)
	return annualCoupon.Div(b.CurrentPrice), nil
}

// ApproximateYTM calculates the approximate Yield to Maturity for a bond.
// Formula: [C + (F - P) / n] / [(F + P) / 2]
// C = Annual Coupon, F = Face Value, P = Current Price, n = Years to Maturity
func (b *Bond) ApproximateYTM() (decimal.Decimal, error) {
	if b.CurrentPrice.IsZero() || b.YearsToMaturity <= 0 {
		return decimal.Zero, errors.New("invalid price or years to maturity for YTM calculation")
	}

	annualCoupon := b.FaceValue.Mul(b.CouponRate)
	n := decimal.NewFromInt(int64(b.YearsToMaturity))

	// (FaceValue - CurrentPrice) / n
	priceDiff := b.FaceValue.Sub(b.CurrentPrice).Div(n)
	numerator := annualCoupon.Add(priceDiff)

	// (FaceValue + CurrentPrice) / 2
	denominator := b.FaceValue.Add(b.CurrentPrice).Div(decimal.NewFromInt(2))

	if denominator.IsZero() {
		return decimal.Zero, errors.New("denominator cannot be zero in YTM calculation")
	}

	return numerator.Div(denominator), nil
}