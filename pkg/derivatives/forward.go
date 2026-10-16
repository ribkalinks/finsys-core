package derivatives

import (
	"errors"
	"github.com/shopspring/decimal"
)

// ForwardContract merepresentasikan parameter kontrak forward
type ForwardContract struct {
	SpotPrice    decimal.Decimal // S_0: Harga spot aset dasar
	RiskFreeRate decimal.Decimal // r: Risk-free interest rate (annualized, diskrit)
	TimeYears    decimal.Decimal // T: Waktu menuju jatuh tempo (dalam tahun)
	StorageCost  decimal.Decimal // Biaya penyimpanan / convenience yield bersih
}

// NewForwardContract membuat instance baru dengan validasi dasar
func NewForwardContract(spot, rate, time decimal.Decimal) (*ForwardContract, error) {
	if spot.LessThanOrEqual(decimal.Zero) {
		return nil, errors.New("spot price must be greater than zero")
	}
	if time.LessThan(decimal.Zero) {
		return nil, errors.New("time to maturity cannot be negative")
	}

	return &ForwardContract{
		SpotPrice:    spot,
		RiskFreeRate: rate,
		TimeYears:    time,
		StorageCost:  decimal.Zero,
	}, nil
}

// CalculateForwardPrice menghitung harga forward menggunakan model cost-of-carry (diskrit)
// Rumus: F_0 = S_0 * (1 + r - q)^T (di mana q adalah storage/convenience cost net)
func (f *ForwardContract) CalculateForwardPrice() decimal.Decimal {
	netRate := f.RiskFreeRate.Sub(f.StorageCost)
	onePlusRate := decimal.NewFromInt(1).Add(netRate)
	compoundingFactor := onePlusRate.Pow(f.TimeYears)
	
	return f.SpotPrice.Mul(compoundingFactor).Round(4)
}

// CalculateForwardValue menghitung nilai kontrak forward (V_t) di tengah masa berlaku
// Rumus: V_t = [S_t - F_0 * (1 + r)^(- (T - t))] atau setara dengan diskonto selisih harga spot saat ini dengan forward price awal
func CalculateForwardValue(currentSpot, initialForwardPrice, riskFreeRate, timeRemaining decimal.Decimal) (decimal.Decimal, error) {
	if timeRemaining.LessThan(decimal.Zero) {
		return decimal.Zero, errors.New("remaining time cannot be negative")
	}

	// Faktor diskonto: (1 + r)^(-t_remaining)
	onePlusRate := decimal.NewFromInt(1).Add(riskFreeRate)
	discountFactor := onePlusRate.Pow(timeRemaining.Neg())

	// Present value dari initial forward price
	pvForwardPrice := initialForwardPrice.Mul(discountFactor)

	// Nilai kontrak forward bagi pemegang posisi long
	forwardValue := currentSpot.Sub(pvForwardPrice)

	return forwardValue.Round(4), nil
}