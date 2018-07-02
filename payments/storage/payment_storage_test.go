package storage

import (
	"billing-api/model"
	"fmt"
	"github.com/go-pg/pg"
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/satori/go.uuid"
)

var dbc *pg.DB
var ps *PaymentStorage

func init() {
	db := pg.Connect(&pg.Options{
		Database: "b_api",
		User:     "b_api",
		Password: "wEfNbjayEbxuv7oaUbM7kosU",
	})

	//db := pg.Connect(&pg.Options{
	//	User:     "postgres",
	//	Password: "postgres",
	//	Database: "postgres",
	//})

	dropTables(db)

	dbc = db

	err := createSchema(db)
	if err != nil {
		panic(err)
	}

	paymentStorage := NewPaymentStorage(db)
	paymentStorage.db = db

	ps = paymentStorage
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{&model.Wallet{}, &model.Transaction{}} {
		err := db.CreateTable(model, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func dropTables(db *pg.DB) {
	for _, model := range []interface{}{&model.Wallet{}, &model.Transaction{}} {
		db.DropTable(model, nil)
	}
}

func truncateTables(db *pg.DB) {
	db.Exec(`TRUNCATE TABLE wallets;`)
	db.Exec(`TRUNCATE TABLE transactions;`)
}

func TestRegisterWallet(t *testing.T) {

	companyId := "test-company_Id"
	userId := "test-user_Id"

	w, err := ps.Register(companyId, userId)

	assert.NoError(t, err)
	assert.NotEmpty(t, w.CompanyId)

	fmt.Println(w)


	pagin:= map[string][]string{
		"page":{"1"},
		"limit":{"10"},
	}

	wallet, err := ps.CurrentBalance(w.CompanyId,pagin)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(wallet))

	assert.Equal(t, companyId, wallet[0].CompanyId)
	assert.Equal(t, model.Cents(0), wallet[0].Balance)

	truncateTables(dbc)
}

func TestDepositMoney(t *testing.T) {
	companyId := "test-company_Id"
	userId := "test-user_Id"

	w, err := ps.Register(companyId, userId)

	assert.NoError(t, err)
	assert.NotEmpty(t, w)

	f := func(model.DepositModel) error {
		return nil
	}

	w1, err := ps.Deposit(f, model.DepositModel{
		CompanyId: companyId,
		Total:     "8.88",
	})

	assert.NoError(t, err)
	assert.Equal(t, w1.CompanyId, companyId)

	pagin:= map[string][]string{
		"page":{"1"},
		"limit":{"10"},
	}

	wallet, err := ps.CurrentBalance(w1.CompanyId,pagin)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(wallet))

	assert.Equal(t, companyId, wallet[0].CompanyId)
	assert.Equal(t, model.Cents(888), wallet[0].Balance)

	truncateTables(dbc)
}

func TestTransferMoney(t *testing.T) {
	companyId1 := "test-company_Id_1"
	userId1 := "test-user_Id_1"

	w, err := ps.Register(companyId1, userId1)

	assert.NoError(t, err)
	assert.NotEmpty(t, w)

	f := func(model.DepositModel) error {
		return nil
	}
	w1, err := ps.Deposit(f, model.DepositModel{
		CompanyId: companyId1,
		Total:     "8.88",
	})
	assert.NoError(t, err)
	assert.Equal(t, companyId1, w1.CompanyId)

	companyId2 := "test-company_Id_2"
	userId2 := "test-user_Id_2"

	w2, err := ps.Register(companyId2, userId2)

	assert.NoError(t, err)
	assert.NotEmpty(t, w2)

	wFrom, err := ps.Transfer(userId1, companyId1, companyId2, 444)
	assert.NoError(t, err)
	assert.Equal(t, model.Cents(444), wFrom.Balance)

	pagin:= map[string][]string{
		"page":{"1"},
		"limit":{"10"},
	}

	wallet1, err := ps.CurrentBalance(companyId1,pagin)
	assert.NoError(t, err)

	assert.Equal(t, model.Cents(444), wallet1[0].Balance)

	wallet2, err := ps.CurrentBalance(companyId2,pagin)
	assert.NoError(t, err)

	assert.Equal(t, model.Cents(444), wallet2[0].Balance)

	truncateTables(dbc)
}

func TestFilterTransactions_RegisterWallet(t *testing.T) {
	companyId1 := "test-company_Id_1"
	userId1 := "test-user_Id_1"

	wId1, err := ps.Register(companyId1, userId1)

	assert.NoError(t, err)
	assert.NotEmpty(t, wId1)

	ts, err := ps.FilterTransactions(model.TransactionFilter{
		CompanyId: companyId1,

	})

	assert.NoError(t, err)

	assert.Equal(t, 1, len(ts))
	assert.Equal(t, model.RegisterWalletTransactionType, ts[0].Type)

	truncateTables(dbc)
}

func TestFilterTransactions_RegisterWallet_WithPagination(t *testing.T) {
	companyId1 := "test-company_Id_1"
	userId1 := "test-user_Id_1"

	wId1, err := ps.Register(companyId1, userId1)

	assert.NoError(t, err)
	assert.NotEmpty(t, wId1)

	for i:=0; i<15; i++{
		f := func(model.DepositModel) error {
			return nil
		}

		w1, err := ps.Deposit(f, model.DepositModel{
			CompanyId: companyId1,
			Total:     "8.88",
		})
		assert.NoError(t, err)
		assert.Equal(t, w1.CompanyId, companyId1)
	}

	ts, err := ps.FilterTransactions(model.TransactionFilter{
		CompanyId: companyId1,
		URL: map[string][]string{
			"page":{"1"},
			"limit":{"10"},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, 10, len(ts))


	ts2, err := ps.FilterTransactions(model.TransactionFilter{
		CompanyId: companyId1,
		URL: map[string][]string{
			"page":{"2"},
			"limit":{"10"},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, 6, len(ts2))

	truncateTables(dbc)
}

func TestWithdrawalMoney(t *testing.T) {
	companyId1 := "test-company_Id_1"
	userId1 := "test-user_Id_1"

	wId1, err := ps.Register(companyId1, userId1)

	assert.NoError(t, err)
	assert.NotEmpty(t, wId1)

	fd := func(model.DepositModel) error {
		return nil
	}
	w1, err := ps.Deposit(fd, model.DepositModel{
		CompanyId: wId1.CompanyId,
		Total:     "8.88",
	})
	assert.NoError(t, err)

	f := func(model.WithdrawalModel) error {
		return nil
	}

	ps.Withdrawal(f, model.WithdrawalModel{
		UserId:      userId1,
		CompanyId:   companyId1,
		PayPalEmail: "user@domen.com",
		Amount:      "3.33", // in "1.00" format
	})

	pagin:= map[string][]string{
		"page":{"1"},
		"limit":{"10"},
	}

	wallet1, err := ps.CurrentBalance(w1.CompanyId,pagin)

	assert.NoError(t, err)

	fmt.Println(wallet1)

	assert.Equal(t, model.Cents(555), wallet1[0].Balance)
}

func TestConvertCents(t *testing.T) {
	stringCents := "10.88"

	c, err := cents(stringCents)
	assert.NoError(t, err)

	assert.Equal(t, model.Cents(1088), c)
}