package api

import (
	"billing-api/model"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

const (
	basicURL = "http://localhost:8000/api/billing/wallet"
)

func TestAPI_RegisterWallet(t *testing.T) {
	cId := getRandomString()
	b, err1 := json.Marshal(RegisterWalletRequest{CompanyId: cId})
	if err1 != nil {
		panic(err1)
	}

	resp, err := http.Post(basicURL+"/register", "application/json", bytes.NewBuffer(b))
	if err != nil {
		panic(err)
	}

	assert.Equal(t, 200, resp.StatusCode)

	w := &model.Wallet{}
	bb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(bb, w)

	assert.NotEmpty(t, w)
	assert.NotEmpty(t, w.Id)
	assert.Equal(t, cId, w.CompanyId)
}

func TestAPI_WalletBalance(t *testing.T) {
	cId := getRandomString()
	b, err1 := json.Marshal(RegisterWalletRequest{CompanyId: cId})
	if err1 != nil {
		panic(err1)
	}

	resp, err := http.Post(basicURL+"/register", "application/json", bytes.NewBuffer(b))
	if err != nil {
		panic(err)
	}

	w := &model.Wallet{}
	bb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(bb, w)

	assert.NotEmpty(t, w)
	assert.NotEmpty(t, w.Id)
	assert.NotEmpty(t, w.CompanyId)

	resp1, err := http.Get(basicURL + "/balance?company_id="+cId)
	if err != nil {
		panic(err)
	}

	ws := []model.Wallet{}

	wsResp, err := ioutil.ReadAll(resp1.Body)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(wsResp, &ws)

	assert.NotEmpty(t, ws)

	fmt.Println(fmt.Sprintf("%+v", ws))
}

func TestAPI_TransferMoney(t *testing.T) {

	cId := getRandomString()
	w1request, err1 := json.Marshal(RegisterWalletRequest{CompanyId: cId})
	if err1 != nil {
		panic(err1)
	}

	respW1, err := http.Post(basicURL+"/register", "application/json", bytes.NewBuffer(w1request))
	if err != nil {
		panic(err)
	}

	w1 := &model.Wallet{}
	w1body, err := ioutil.ReadAll(respW1.Body)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(w1body, w1)

	assert.NotEmpty(t, w1)
	assert.NotEmpty(t, w1.Id)
	assert.NotEmpty(t, w1.CompanyId)

	cId_2 := getRandomString()

	w2request, err2 := json.Marshal(RegisterWalletRequest{CompanyId: cId_2})
	if err2 != nil {
		panic(err2)
	}

	respW2, err := http.Post(basicURL+"/register", "application/json", bytes.NewBuffer(w2request))
	if err != nil {
		panic(err)
	}

	w2 := &model.Wallet{}
	w2body, err := ioutil.ReadAll(respW2.Body)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(w2body, w2)

	assert.NotEmpty(t, w2)
	assert.NotEmpty(t, w2.Id)
	assert.NotEmpty(t, w2.CompanyId)

	tr := TransferRequest{
		CompanyIdFrom: cId,

		CompanyIdTo: cId_2,
		Amount:      0,
	}

	trRequest, err1 := json.Marshal(tr)

	respW3, err := http.Post(basicURL+"/transfer", "application/json", bytes.NewBuffer(trRequest))
	if err != nil {
		panic(err)
	}

	w3 := &model.Wallet{}
	w3body, err := ioutil.ReadAll(respW3.Body)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(w3body, w3)

	fmt.Println("-->", w3)
	assert.NotEmpty(t, w3)
	assert.Equal(t, w1.Id, w3.Id)
	fmt.Println(w3)
}

func TestAPI_GetTransactions(t *testing.T) {
	cId := getRandomString()
	w1request, err1 := json.Marshal(RegisterWalletRequest{CompanyId: cId})
	if err1 != nil {
		panic(err1)
	}

	respW1, err := http.Post(basicURL+"/register", "application/json", bytes.NewBuffer(w1request))
	if err != nil {
		panic(err)
	}

	w1 := &model.Wallet{}
	w1body, err := ioutil.ReadAll(respW1.Body)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(w1body, w1)

	assert.NotEmpty(t, w1)
	assert.NotEmpty(t, w1.Id)
	assert.NotEmpty(t, w1.CompanyId)

	respTr, err := http.Get(basicURL + "/transactions?company_id=" + cId+"&page=1&size=10")
	if err != nil {
		panic(err)
	}

	tResponse := TransactionsResponse{}
	t1body, err := ioutil.ReadAll(respTr.Body)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(t1body, &tResponse)
	assert.True(t, len(tResponse.Transactions) >= 1)
}

//func TestAPI_DepositMoney(t *testing.T) {
//
//	b, err1 := json.Marshal(RegisterWalletRequest{UserId: "userId", CompanyId: "company_id"})
//	if err1 != nil {
//		panic(err1)
//	}
//
//	resp, err := http.Post(basicURL+"/wallet/register", "application/json", bytes.NewBuffer(b))
//	if err != nil {
//		panic(err)
//	}
//
//	w := &model.Wallet{}
//	bb, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		panic(err)
//	}
//	json.Unmarshal(bb, w)
//
//	assert.NotEmpty(t, w)
//	assert.NotEmpty(t, w.Id)
//	assert.NotEmpty(t, w.CompanyId)
//
//
//	dm := DepositRequest{
//		Number: "4032031547160837",
//		Type: "visa",
//		ExpireMonth: "03",
//		ExpireYear: "2023",
//		CVV2: "777",
//		FirstName: "John",
//		LastName:"Doe",
//		Total:"7.00",
//
//		//UserId:"userId", //TODO auth request
//		//CompanyId:"companyId",
//		//WalletId:"WalletId",
//	}
//
//
//
//	depositRequest, err := json.Marshal(dm)
//	if err!=nil{
//		panic(err)
//	}
//
//
//	depositResp, err := http.Post(basicURL+"/wallet/deposit", "application/json", bytes.NewBuffer(depositRequest))
//	if err != nil {
//		panic(err)
//	}
//
//	depositWallet := &model.Wallet{}
//	dwb, err := ioutil.ReadAll(depositResp.Body)
//	if err != nil {
//		panic(err)
//	}
//	json.Unmarshal(dwb, depositWallet)
//
//	assert.NotEmpty(t, depositWallet)
//	assert.Equal(t,700,depositWallet.Balance)
//
//}

func getRandomString() string {
	uuid, err := uuid.NewV1()
	if err != nil {

		panic(err)
	}
	return uuid.String()
}
