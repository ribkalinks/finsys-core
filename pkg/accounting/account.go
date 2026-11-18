package accounting

import "fmt"

type AccountType string

const (
	Asset     AccountType = "ASSET"
	Liability AccountType = "LIABILITY"
	Equity    AccountType = "EQUITY"
	Revenue   AccountType = "REVENUE"
	Expense   AccountType = "EXPENSE"
)

type NormalBalance string

const (
	DebitBalance  NormalBalance = "DEBIT"
	CreditBalance NormalBalance = "CREDIT"
)

type Account struct {
	ID           string
	Name         string
	Type         AccountType
	NormalBalance NormalBalance
}

// NewAccount membuat akun baru dengan otomatis menentukan saldo normalnya
func NewAccount(id, name string, accType AccountType) (*Account, error) {
	var normal NormalBalance
	switch accType {
	case Asset, Expense:
		normal = DebitBalance
	case Liability, Equity, Revenue:
		normal = CreditBalance
	default:
		return nil, fmt.Errorf("tipe akun tidak valid: %s", accType)
	}

	return &Account{
		ID:            id,
		Name:          name,
		Type:          accType,
		NormalBalance: normal,
	}, nil
}