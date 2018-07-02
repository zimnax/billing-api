//go:generate swagger generate spec
package api

import (
	"billing-api/model"
	"billing-api/payments/service"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

var api API

type API struct {
	PaymentService *service.PaymentService
}

func Init(ps *service.PaymentService) {
	api = API{
		PaymentService: ps,
	}

	router := mux.NewRouter().PathPrefix("/api/billing/wallet").Subrouter()

	router.HandleFunc("/register", api.RegisterWallet).Methods("POST")
	router.HandleFunc("/deposit", api.DepositMoney).Methods("POST")
	router.HandleFunc("/balance", api.WalletBalance).Methods("GET")
	router.HandleFunc("/transfer", api.TransferMoney).Methods("POST")
	router.HandleFunc("/withdrawal", api.Withdrawal).Methods("POST")

	router.HandleFunc("/transactions", api.GetTransactions).Methods("GET")
	if err := http.ListenAndServe(":8000", router); err != nil {
		panic(err)
	}
}

type RegisterWalletRequest struct {
	CompanyId string
}

func (a *API) RegisterWallet(w http.ResponseWriter, r *http.Request) {
	var d RegisterWalletRequest

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("Error during decode body: %+v", err))
	}

	UserId := "user_id" //TODO FROM TOKEN

	wallet, err := a.PaymentService.Register(d.CompanyId, UserId)
	if err != nil {
		respondWithError(w, model.PaymentError{Message: err.Error(), Code: 500})
	} else {
		respondWithJSON(w, 200, wallet)
	}
}

type DepositRequest struct {
	Number      string
	Type        string
	ExpireMonth string
	ExpireYear  string
	CVV2        string
	FirstName   string
	LastName    string
	Total       string // in "1.00" format
}

func (a *API) DepositMoney(w http.ResponseWriter, r *http.Request) {
	var d DepositRequest

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("Error during decode body: %+v", err))
	}

	dm := model.DepositModel{
		Number:      d.Number,
		Type:        d.Type,
		ExpireMonth: d.ExpireMonth,
		ExpireYear:  d.ExpireYear,
		CVV2:        d.CVV2,
		FirstName:   d.FirstName,
		LastName:    d.LastName,
		Total:       d.Total,

		UserId:    "userId",    //TODO auth request STUB
		CompanyId: "companyId", //TODO auth request STUB
	}
	wallet, err := a.PaymentService.Deposit(dm)

	if err != nil {
		respondWithError(w, model.InternalError)
	} else {
		respondWithJSON(w, 200, wallet)
	}
}

func (a *API) WalletBalance(w http.ResponseWriter, r *http.Request) {
	companyId := r.FormValue("company_id")
	root := false //TODO ROOT STUB

	if !root && companyId == "" {
		respondWithError(w, model.PermissionRestrictions)
	}

	wallets, err := a.PaymentService.Balance(companyId,r.URL.Query())
	if err != nil {
		respondWithError(w, model.InternalError)
	} else {
		respondWithJSON(w, 200, wallets)
	}
}

type TransferRequest struct {
	CompanyIdFrom string

	CompanyIdTo string
	Amount      model.Cents
}

func (a *API) TransferMoney(w http.ResponseWriter, r *http.Request) {
	var tr TransferRequest

	userId := "userId" //TODO auth request STUB

	if err := json.NewDecoder(r.Body).Decode(&tr); err != nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("Error during decode body: %+v", err))
	}

	walletFrom, err := a.PaymentService.Transfer(userId, tr.CompanyIdFrom, tr.CompanyIdTo, tr.Amount)

	if err != nil {
		respondWithError(w, model.InternalError)
	} else {
		respondWithJSON(w, 200, walletFrom)
	}
}

type WithdrawalRequest struct {
	PayPalEmail string
	Amount      string
}

func (a *API) Withdrawal(w http.ResponseWriter, r *http.Request) {
	var wr WithdrawalRequest

	if err := json.NewDecoder(r.Body).Decode(&wr); err != nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("Error during decode body: %+v", err))
	}

	ppw := model.WithdrawalModel{
		PayPalEmail: wr.PayPalEmail,
		Amount:      wr.Amount,
	}

	ppw.UserId = "userId"       //TODO auth request STUB
	ppw.CompanyId = "companyId" //TODO auth request STUB

	err := a.PaymentService.Withdrawal(ppw)
	if err != nil {
		respondWithError(w, model.NewError(err, 500, err.Error()))
	}
}

type TransactionsResponse struct {
	Transactions []model.Transaction
}

func (a *API) GetTransactions(w http.ResponseWriter, r *http.Request) {

	f, err := getRequestFilter(r)
	if err != nil {
		respondWithError(w, model.PermissionRestrictions)
	}

	fmt.Println(fmt.Sprintf("FILTER: %+v", f))

	ts, errs := a.PaymentService.Transactions(*f)

	if errs != nil {
		respondWithError(w, model.InternalError)
	}
	respondWithJSON(w, 200, TransactionsResponse{ts})
}
