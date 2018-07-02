package storage

import (
	"billing-api/model"
	"errors"
	"fmt"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/satori/go.uuid"
	"log"
	"strconv"
	"time"
	"net/url"
)

const (
	initialBalance = 0
)

func NewPaymentStorage(db *pg.DB) *PaymentStorage {
	return &PaymentStorage{
		db: db,
	}
}

type PaymentStorage struct {
	db *pg.DB
}

func (pm *PaymentStorage) Register(cId, uId string) (model.Wallet, error) {
	log.Println("Register in PaymentStorage")
	var wallet model.Wallet

	err := pm.db.RunInTransaction(func(tx *pg.Tx) error {
		var err error

		var walletExisting model.Wallet
		err = pm.db.Model(&walletExisting).Where("company_id = ?", cId).Select()

		if err != nil {
			if err != pg.ErrNoRows {
				return err
			} else {
				w := model.Wallet{CompanyId: cId, Balance: initialBalance}
				if err = pm.db.Insert(&w); err != nil {
					return err
				}

				tuid, err := getUUID()
				if err != nil {
					return err
				}

				err = pm.db.Insert(&model.Transaction{
					Id:        tuid,
					CompanyId: cId,
					UserId:    uId,
					Type:      model.RegisterWalletTransactionType,
					Date:      timeNowFunc(),
				})

				wallet = w
			}
		}

		if walletExisting.CompanyId != "" {
			return errors.New("Wallet for company already exist")
		}
		return err
	})
	return wallet, err
}

type depositFunc func(model.DepositModel) error

func (pm *PaymentStorage) Deposit(payPalDepositFunc depositFunc, dm model.DepositModel) (model.Wallet, error) {
	log.Println("Deposit in PaymentStorage")

	var wallet model.Wallet
	err := pm.db.RunInTransaction(func(tx *pg.Tx) error {
		var err error

		var w model.Wallet
		err = pm.db.Model(&w).Where("company_id = ?", dm.CompanyId).Select()
		if err != nil {
			return err
		}

		cents, err := cents(dm.Total) //TODO "%.2f"
		if err != nil {
			return err
		}

		w.Balance = w.Balance + cents
		if err = pm.db.Update(&w); err != nil {
			return err
		}

		tuid, err := getUUID()
		if err != nil {
			return err
		}

		err = pm.db.Insert(&model.Transaction{
			Id:        tuid,
			CompanyId: dm.CompanyId,
			UserId:    dm.UserId,
			Type:      model.DepositTransactionType,
			Amount:    cents,
			Date:      timeNowFunc(),
		})

		err = payPalDepositFunc(dm)
		if err != nil {
			return err
		}
		wallet = w
		return err
	})
	return wallet, err
}

func (pm *PaymentStorage) CurrentBalance(cId string,p url.Values) ([]model.Wallet, error) {
	log.Println("CurrentBalance in PaymentStorage")

	var wallets []model.Wallet
	err := pm.db.RunInTransaction(func(tx *pg.Tx) error {
		var err error

		if cId != "" {
			err := pm.db.Model(&wallets).Where("company_id = ?", cId).Select()

			if err != nil {
				return err
			}
		} else {
			var ws []model.Wallet
			query := pm.db.Model(&ws)
			query.Apply(orm.Pagination(p))

			err = query.Select()
			if err != nil {
				return err
			}

			wallets = ws
		}
		return err
	})
	return wallets, err
}

func (pm *PaymentStorage) Transfer(uId, from, to string, amount model.Cents) (*model.Wallet, error) {
	log.Println("Transfer in PaymentStorage")
	var wallet *model.Wallet
	err := pm.db.RunInTransaction(func(tx *pg.Tx) error {
		var err error

		var walletFrom model.Wallet
		if err = pm.db.Model(&walletFrom).Where("company_id = ?", from).Select(); err != nil {
			return err
		}

		walletFrom.Balance = walletFrom.Balance - amount
		if err = pm.db.Update(&walletFrom); err != nil {
			return err
		}

		var walletTo model.Wallet

		if err = pm.db.Model(&walletTo).Where("company_id = ?", to).Select(); err != nil {
			return err
		}

		walletTo.Balance = walletTo.Balance + amount
		if err = pm.db.Update(&walletTo); err != nil {
			return err
		}

		tuid, err := getUUID()
		if err != nil {
			return err
		}

		err = pm.db.Insert(&model.Transaction{
			Id:        tuid,
			CompanyId: walletFrom.CompanyId,
			UserId:    uId,
			Type:      model.TransferTransactionType,
			Amount:    amount,
			Date:      timeNowFunc(),
		})

		wallet = &walletFrom
		return err
	})
	return wallet, err
}

func (pm *PaymentStorage) FilterTransactions(f model.TransactionFilter) ([]model.Transaction, error) {
	log.Println("FilterTransactions in PaymentStorage")

	var ts []model.Transaction

	err := pm.db.RunInTransaction(func(tx *pg.Tx) error {
		var err error

		query := pm.db.Model(&ts)

		if f.CompanyId != "" {
			query.Where("company_id = ?", f.CompanyId)
		}

		if f.DateFrom != 0 {
			query.Where("date >= ?", f.DateFrom)
		}

		if f.DateTo != 0 {
			query.Where("date <= ?", f.DateTo)
		}
		query.Apply(orm.Pagination(f.URL))
		err = query.Select()

		return err
	})
	return ts, err
}

type withdrawalFunc func(model.WithdrawalModel) error

func (pm *PaymentStorage) Withdrawal(f withdrawalFunc, ppw model.WithdrawalModel) error {
	log.Println("Withdrawal in PaymentStorage")

	err := pm.db.RunInTransaction(func(tx *pg.Tx) error {
		var err error
		var walletFrom model.Wallet
		if err = pm.db.Model(&walletFrom).Where("company_id = ?", ppw.CompanyId).Select(); err != nil {
			return err
		}

		c, err := cents(ppw.Amount)
		if err != nil {
			return err
		}

		if walletFrom.Balance < c {
			return errors.New("Not enough money to withdrawal")
		}
		walletFrom.Balance = walletFrom.Balance - c
		if err = pm.db.Update(&walletFrom); err != nil {
			return err
		}

		tuid, err := getUUID()
		if err != nil {
			return err
		}

		err = pm.db.Insert(&model.Transaction{
			Id:        tuid,
			CompanyId: ppw.CompanyId,
			UserId:    ppw.UserId,
			Type:      model.WithdrawalTransactionType,
			Amount:    c,
			Date:      timeNowFunc(),
		})

		err = f(ppw)
		if err != nil {
			return err
		}

		return err
	})
	return err
}

func cents(total string) (model.Cents, error) {
	if total == "" {
		return model.Cents(0), nil
	}

	fl, err := strconv.ParseFloat(total, 64)
	if err != nil {
		return model.Cents(0), err
	}
	intCents := fl * 100
	return model.Cents(intCents), nil
}

func getUUID() (string, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		fmt.Printf("Err during generating uuid: %s", err)
		return "", err
	}

	return uuid.String(), nil
}

var timeNowFunc = func() uint {
	return uint(time.Now().Unix())
}

func (pm *PaymentStorage) TruncateTables() {
	pm.db.Exec(`TRUNCATE TABLE wallets;`)
	pm.db.Exec(`TRUNCATE TABLE transactions;`)
}
