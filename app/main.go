//go:generate swagger generate spec
package main

import (
	"billing-api/api"
	"billing-api/payments/payment_gate"
	"billing-api/payments/service"
	"billing-api/payments/storage"
	"billing-api/postgredb"
	"billing-api/pub_sub"
)

func main() {
	db := postgredb.New()
	paymentStorage := storage.NewPaymentStorage(db)
	payPalStorage := payment_gate.New()
	paymentService := service.NewPaymentService(paymentStorage, payPalStorage)

	go api.Init(paymentService)

	ps := pub_sub.Init()
	ps.Listen()

}
