package service

import (
	"billing-api/model"
	"billing-api/payments/payment_gate"
	"billing-api/payments/storage"
	"net/url"
)

type PaymentService struct {
	PaymentStorage *storage.PaymentStorage
	PaymentGate    payment_gate.PayGate
}

func NewPaymentService(ps *storage.PaymentStorage, pg payment_gate.PayGate) *PaymentService {
	return &PaymentService{
		PaymentStorage: ps,
		PaymentGate:    pg,
	}
}

func (ps *PaymentService) Register(cId, uId string) (model.Wallet, error) {
	return ps.PaymentStorage.Register(cId, uId)
}

func (ps *PaymentService) Deposit(dm model.DepositModel) (model.Wallet, error) {
	return ps.PaymentStorage.Deposit(ps.PaymentGate.Deposit, dm)
}

func (ps *PaymentService) Balance(wid string,pagination url.Values) ([]model.Wallet, error) {
	return ps.PaymentStorage.CurrentBalance(wid,pagination)
}

func (ps *PaymentService) Transfer(uId, from, to string, amount model.Cents) (*model.Wallet, error) {
	return ps.PaymentStorage.Transfer(uId, from, to, amount)
}

func (ps *PaymentService) Transactions(f model.TransactionFilter) ([]model.Transaction, error) {
	return ps.PaymentStorage.FilterTransactions(f)
}

func (ps *PaymentService) Withdrawal(ppw model.WithdrawalModel) error {
	return ps.PaymentStorage.Withdrawal(ps.PaymentGate.Withdrawal, ppw)
}
