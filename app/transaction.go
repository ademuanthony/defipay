package app

import (
	"context"
	"deficonnect/defipayapi/web"
	"encoding/json"
	"net/http"
)

type CreateTransactionInput struct {
	BankName      string          `govalid:"req" json:"bank_name" toml:"bank_name" yaml:"bank_name"`
	AccountNumber string          `govalid:"req" json:"account_number" toml:"account_number" yaml:"account_number"`
	AccountName   string          `govalid:"req" json:"account_name" toml:"account_name" yaml:"account_name"`
	Amount        int64           `govalid:"req|min:10|max:10000" json:"amount" toml:"amount" yaml:"amount"`
	Email         string          `govalid:"req" json:"email" toml:"email" yaml:"email"`
	Network       string          `govalid:"req" json:"network" toml:"network" yaml:"network"`
	Currency      string          `govalid:"req" json:"currency" toml:"currency" yaml:"currency"`
	PaymentLink   string          `boil:"payment_link" json:"payment_link" toml:"payment_link" yaml:"payment_link"`
	Type          transactionType `boil:"type" json:"type" toml:"type" yaml:"type"`

	WalletAddress string `json:"-"`
	PrivateKey    string `json:"-"`
}

type TransactionOutput struct {
	ID            string `boil:"id" json:"id" toml:"id" yaml:"id"`
	BankName      string `boil:"bank_name" json:"bank_name" toml:"bank_name" yaml:"bank_name"`
	AccountNumber string `boil:"account_number" json:"account_number" toml:"account_number" yaml:"account_number"`
	AccountName   string `boil:"account_name" json:"account_name" toml:"account_name" yaml:"account_name"`
	Amount        int64  `boil:"amount" json:"amount" toml:"amount" yaml:"amount"`
	Email         string `boil:"email" json:"email" toml:"email" yaml:"email"`
	Network       string `boil:"network" json:"network" toml:"network" yaml:"network"`
	Currency      string `boil:"currency" json:"currency" toml:"currency" yaml:"currency"`
	WalletAddress string `boil:"wallet_address" json:"wallet_address" toml:"wallet_address" yaml:"wallet_address"`
	PaymentLink   string `boil:"payment_link" json:"payment_link" toml:"payment_link" yaml:"payment_link"`
	Type          string `boil:"type" json:"type" toml:"type" yaml:"type"`
}

type transactionType string

var transactionTypes = struct {
	TopUp        transactionType
	FundTransfer transactionType
}{
	TopUp:        "top up",
	FundTransfer: "fund transfer",
}

func (m module) createFundTransferTransaction(w http.ResponseWriter, r *http.Request) {
	var input CreateTransactionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Error("Login", "json::Decode", err)
		web.SendErrorfJSON(w, "cannot decode request")
		return
	}

	input.Type = transactionTypes.FundTransfer

	tran, err := m.createTransaction(r.Context(), input)
	if err != nil {
		log.Error("Create Transaction", err)
		msg := "Cannot create transaction. Please try again"
		if messenger, ok := err.(ErrorMessenger); ok {
			msg = messenger.ErrorMessage()
		}
		web.SendErrorfJSON(w, msg)
		return
	}

	web.SendJSON(w, tran)
}

func (m module) createTupUpTransaction(w http.ResponseWriter, r *http.Request) {
	var input CreateTransactionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Error("Login", "json::Decode", err)
		web.SendErrorfJSON(w, "cannot decode request")
		return
	}

	input.Type = transactionTypes.TopUp

	tran, err := m.createTransaction(r.Context(), input)
	if err != nil {
		log.Error("Create Transaction", err)
		msg := "Cannot create transaction. Please try again"
		if messenger, ok := err.(ErrorMessenger); ok {
			msg = messenger.ErrorMessage()
		}
		web.SendErrorfJSON(w, msg)
		return
	}

	web.SendJSON(w, tran)
}

func (m module) createTransaction(ctx context.Context, input CreateTransactionInput) (*TransactionOutput, error) {
	vio, err := v.Violations(&input)
	if err != nil {
		return nil, err
	}
	if len(vio) > 0 {
		return nil, newValidationError(vio)
	}

	privateKey, wallet, err := GenerateWallet()
	if err != nil {
		return nil, err
	}

	input.PrivateKey = privateKey
	input.WalletAddress = wallet

	return m.db.CreateTransaction(ctx, input)
}
