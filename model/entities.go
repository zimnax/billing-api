package model

import "net/url"

const (
	RegisterWalletTransactionType = TransactionType("WALLET_REGISTER")
	DepositTransactionType        = TransactionType("DEPOSIT")
	TransferTransactionType       = TransactionType("TRANSFER_MONEY")
	WithdrawalTransactionType     = TransactionType("WITHDRAWAL_MONEY")
)

type Cents uint64

type User struct {
	Id        string
	CompanyId string
	Email     string
	FirstName string
	LastName  string
	Activated bool
}

type Company struct {
	Id          string
	Description string
	Enabled     bool
}

type CreditCard struct {
	Number      string
	Type        string
	ExpireMonth string
	ExpireYear  string
	CVV2        string
}

type Wallet struct {
	Id        int64
	CompanyId string
	Balance   Cents
}

type Transaction struct {
	Id        string
	UserId    string
	CompanyId string
	Amount    Cents
	Type      TransactionType
	Date      uint
}

type TransactionType string

type TransactionFilter struct {
	CompanyId string
	DateFrom  uint64
	DateTo    uint64

	Page int
	Limit int

	URL url.Values
}

type WithdrawalModel struct {
	UserId      string
	CompanyId   string
	PayPalEmail string
	Amount      string // in "1.00" format
}

type DepositModel struct {
	Number      string
	Type        string
	ExpireMonth string
	ExpireYear  string
	CVV2        string
	FirstName   string
	LastName    string
	Currency    string
	Total       string // in "1.00" format

	UserId    string
	CompanyId string
}
