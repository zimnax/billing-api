package payment_gate

import (
	"billing-api/cfg"
	"billing-api/model"
	"fmt"
	"github.com/logpacker/PayPal-Go-SDK"
)

type PayGate interface {
	Withdrawal(model.WithdrawalModel) error
	Deposit(model.DepositModel) error
}

type PayPalClient struct {
	Client *paypalsdk.Client
}

func New() *PayPalClient {
	c, err := paypalsdk.NewClient(cfg.PayPal.ClientId, cfg.PayPal.ApiSecret, paypalsdk.APIBaseSandBox)
	if err != nil {
		panic(fmt.Sprintf("Error while creating PayPal client. Err: %+v", err))

	}
	accessToken, err := c.GetAccessToken()
	fmt.Println("accessToken: ", accessToken.Token)

	return &PayPalClient{
		Client: c,
	}
}

func (ppc PayPalClient) Withdrawal(ppw model.WithdrawalModel) error {
	//developer.paypal.com/developer
	payout := paypalsdk.Payout{
		SenderBatchHeader: &paypalsdk.SenderBatchHeader{
			EmailSubject: "Withdrawal money to personal account",
		},
		Items: []paypalsdk.PayoutItem{
			{
				RecipientType: "EMAIL",
				Receiver:      ppw.PayPalEmail,
				Amount: &paypalsdk.AmountPayout{
					Value:    ppw.Amount,
					Currency: "USD",
				},
				//Note:         "Optional note",
				//SenderItemID: "Optional Item ID",
			},
		},
	}

	payoutResp, err := ppc.Client.CreateSinglePayout(payout)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("payoutResp: %+v", payoutResp))
	return nil
}

func (ppc PayPalClient) Deposit(dm model.DepositModel) error {
	p := paypalsdk.Payment{
		Intent: "sale",
		Payer: &paypalsdk.Payer{
			PaymentMethod: "credit_card",
			FundingInstruments: []paypalsdk.FundingInstrument{
				{
					CreditCard: &paypalsdk.CreditCard{
						Number:      dm.Number,
						Type:        dm.Type,
						ExpireMonth: dm.ExpireMonth,
						ExpireYear:  dm.ExpireYear,
						CVV2:        dm.CVV2,
						FirstName:   dm.FirstName,
						LastName:    dm.LastName,
					},
				}},
		},
		Transactions: []paypalsdk.Transaction{
			{
				Amount: &paypalsdk.Amount{
					Currency: "USD",
					Total:    dm.Total,
				},
				Description: "Deposit wallet",
			}},
	}
	_, err := ppc.Client.CreatePayment(p)
	if err != nil {
		return err
	}

	return nil
}
