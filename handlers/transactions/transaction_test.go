package transactions

import (
	"context"
	"deficonnect/defipayapi/app"
	"deficonnect/defipayapi/handlers"
	"deficonnect/defipayapi/postgres/models"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestCreateTransaction(t *testing.T) {
	return
	m, err := handlers.InitSlsApp(true)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	tran, err := m.CreateTransaction(context.Background(), app.CreateTransactionInput{
		BankName:      "FBN",
		AccountNumber: "0078912239",
		AccountName:   "Ademu Anthony",
		Amount:        1250000,
		Email:         "a@b.c",
		Network:       string(app.Networks.BSC),
		Currency:      app.DFC.Symbol,
		Type:          "fund transfer",
	}, &models.Account{
		ID: m.MasterAccountID(),
	})

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if tran.TokenAmount != "125000" {
		t.Fail()
	}

}

func TestDollarToToken(t *testing.T) {
	amount, valid := common.Big1.SetString("1000000000000000000", 10)
	if !valid {
		t.Fail()
	}

	m := app.Module{}

}
