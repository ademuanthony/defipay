package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"deficonnect/defipayapi/web"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ethereum/go-ethereum/common"
)

type CreateTransactionInput struct {
	BankName      string          `govalid:"req" json:"bankName" toml:"bank_name" yaml:"bank_name"`
	AccountNumber string          `govalid:"req" json:"accountNumber" toml:"account_number" yaml:"account_number"`
	AccountName   string          `govalid:"req" json:"accountName" toml:"account_name" yaml:"account_name"`
	Amount        int64           `govalid:"req|min:10|max:10000" json:"amount" toml:"amount" yaml:"amount"`
	Email         string          `govalid:"req" json:"email" toml:"email" yaml:"email"`
	Network       string          `govalid:"req" json:"network" toml:"network" yaml:"network"`
	Currency      string          `govalid:"req" json:"currency" toml:"currency" yaml:"currency"`
	PaymentLink   string          `boil:"payment_link" json:"paymentLink" toml:"payment_link" yaml:"payment_link"`
	Type          transactionType `govalid:"req" json:"type" toml:"type" yaml:"type"`

	WalletAddress string `json:"-"`
	PrivateKey    string `json:"-"`
}

type TransactionOutput struct {
	ID            string `boil:"id" json:"id" toml:"id" yaml:"id"`
	BankName      string `boil:"bank_name" json:"bank_name" toml:"bank_name" yaml:"bank_name"`
	AccountNumber string `boil:"account_number" json:"account_number" toml:"account_number" yaml:"account_number"`
	AccountName   string `boil:"account_name" json:"account_name" toml:"account_name" yaml:"account_name"`
	Amount        int64  `boil:"amount" json:"amount" toml:"amount" yaml:"amount"`
	AmountPaid    string `json:"amount_paid"`
	TokenAmount   string `json:"token_amount"`
	Email         string `boil:"email" json:"email" toml:"email" yaml:"email"`
	Network       string `boil:"network" json:"network" toml:"network" yaml:"network"`
	Currency      string `boil:"currency" json:"currency" toml:"currency" yaml:"currency"`
	WalletAddress string `boil:"wallet_address" json:"wallet_address" toml:"wallet_address" yaml:"wallet_address"`
	PaymentLink   string `boil:"payment_link" json:"payment_link" toml:"payment_link" yaml:"payment_link"`
	Type          string `boil:"type" json:"type" toml:"type" yaml:"type"`
	Status        string `json:"status"`
}

type UpdateCurrencyInput struct {
	TransactionID string `json:"transaction_id"`
	Network       string `json:"network"`
	Currency      string `json:"currency"`
}

type transactionType string

var transactionTypes = struct {
	TopUp        transactionType
	FundTransfer transactionType
}{
	TopUp:        "top up",
	FundTransfer: "fund transfer",
}

type TransactionStatus string

var TransactionStatuses = struct {
	Pending       TransactionStatus
	PartiallyPaid TransactionStatus
	Paid          TransactionStatus
	Processing    TransactionStatus
	Completed     TransactionStatus
	Cancelled     TransactionStatus
}{
	Pending:       "pending",
	PartiallyPaid: "partial payment",
	Paid:          "paid",
	Processing:    "processing",
	Completed:     "completed",
	Cancelled:     "cancelled",
}

type GetTransactionsInput struct {
	Email     string
	AccountID string
	Offset    int
	Limit     int
}

type Response events.APIGatewayProxyResponse

func (m Module) GetTransaction(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	id := request.PathParameters["id"]
	transaction, err := m.db.Transaction(ctx, id)
	if err != nil {
		return Response{StatusCode: 400}, err
	}
	return SendJSON(transaction)
}

func (m Module) GetTransactions(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	email := r.QueryStringParameters["email"]
	accountID := m.server.GetUserIDTokenCtxSls(r)
	pagedReq := web.GetPaginationInfoSls(r)

	transactions, count, err := m.db.Transactions(ctx, GetTransactionsInput{
		Email: email, AccountID: accountID, Offset: pagedReq.Offset, Limit: pagedReq.Limit,
	})

	if err != nil {
		return Response{StatusCode: 400}, err
	}

	return SendPagedJSON(transactions, count)
}

func (m Module) CreateTransaction(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	var input CreateTransactionInput
	if err := json.Unmarshal([]byte(r.Body), &input); err != nil {
		log.Error("Login", "json::Decode", err)
		return SendErrorfJSON("cannot decode request")
	}

	switch(input.Type) {
	case "0":
		input.Type = transactionTypes.FundTransfer
	case "1":
		input.Type = transactionTypes.TopUp
	default:
		return SendErrorfJSON("unsupported transaction type")
	}

	tran, err := m.createTransaction(ctx, input)
	if err != nil {
		log.Error("Create Transaction", err)
		msg := "Cannot create transaction. Please try again"
		if messenger, ok := err.(ErrorMessenger); ok {
			msg = messenger.ErrorMessage()
		}
		return SendErrorfJSON(msg)
	}

	return SendJSON(tran)
}

func (m Module) createTransaction(ctx context.Context, input CreateTransactionInput) (*TransactionOutput, error) {
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

func (m Module) UpdateTransactionCurrency(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	var input UpdateCurrencyInput
	if err := json.Unmarshal([]byte(r.Body), &input); err != nil {
		log.Error("Login", "json::Decode", err)
		return SendErrorfJSON("cannot decode request")
	}

	if m.currencyProcessors[input.Currency][Network(input.Network)] == nil {
		return SendErrorfJSON("Unsupported currency")
	}

	transaction, err := m.db.Transaction(ctx, input.TransactionID)
	if err != nil {
		return m.handleError(err, "Get Transaction")
	}

	if transaction.Email != m.server.GetUserIDTokenCtxSls(r) {
		return SendErrorfJSON("Invalid operation")
	}

	txOutput, err := m.db.UpdateCurrency(ctx, input)
	if err != nil {
		return m.handleError(err, "Update Currency")
	}

	return SendJSON(txOutput)
}

func (m Module) CheckTransactionStatus(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	id := r.QueryStringParameters["id"]
	transaction, err := m.db.Transaction(ctx, id)
	if err != nil {
		return m.handleError(err, "Get Transaction")
	}

	if transaction.Status == string(TransactionStatuses.Completed) || transaction.Status == string(TransactionStatuses.Cancelled) ||
		transaction.Status == string(TransactionStatuses.Paid) || transaction.Status == string(TransactionStatuses.Processing) {
		return SendJSON(transaction)
	}

	currencyProcessor := m.currencyProcessors[(transaction.Currency)][Network(transaction.Network)]

	amountPaid, err := currencyProcessor.BalanceOf(nil, common.HexToAddress(transaction.WalletAddress))
	if err != nil {
		return m.handleError(err)
	}

	if err := m.db.UpdateTransactionPayment(ctx, id, amountPaid.String()); err != nil {
		return m.handleError(err)
	}

	transaction.AmountPaid = amountPaid.String()
	var tokenAmount *big.Int
	tokenAmount, _ = tokenAmount.SetString(transaction.TokenAmount, 64)

	if c := tokenAmount.Cmp(amountPaid); c == 0 || c == -1 {
		status, err := m.processTransaction(ctx, transaction)
		if err != nil {
			return m.handleError(err, "process transaction")
		}
		transaction.Status = status
	}

	return SendJSON(transaction)
}

func (m Module) processTransaction(ctx context.Context, transaction *TransactionOutput) (string, error) {
	currencyProcessor := m.currencyProcessors[(transaction.Currency)][Network(transaction.Network)]
	if transaction.Status == string(TransactionStatuses.Completed) {
		return "", errors.New("already completed")
	}
	amountPaid, err := currencyProcessor.BalanceOf(nil, common.HexToAddress(transaction.WalletAddress))
	if err != nil {
		return "", err
	}

	var tokenAmount *big.Int
	tokenAmount, _ = tokenAmount.SetString(transaction.TokenAmount, 64)

	if c := amountPaid.Cmp(tokenAmount); c == -1 {
		return "", errors.New("incomplete payment")
	}

	if transaction.Status == string(TransactionStatuses.Processing) { // 111000+102500+10000
		return "", errors.New("processing")
	}

	if err := m.db.UpdateTransactionStatus(ctx, transaction.ID, TransactionStatuses.Processing); err != nil {
		return "", err
	}

	pk, err := m.db.TransactionPK(ctx, transaction.ID)
	if err != nil {
		return "", err
	}

	if _, err := currencyProcessor.Transfer(ctx, pk, common.HexToAddress(m.config.MasterAddress), amountPaid); err != nil {
		return "", err
	}

	if transaction.Type == string(transactionTypes.FundTransfer) {
		if err := m.assignTransactionToAgent(ctx, transaction); err != nil {
			return "", err
		}
	} else {
		if err := m.db.UpdateTransactionStatus(ctx, transaction.ID, TransactionStatuses.Completed); err != nil {
			return "", err
		}

		if err := m.db.CreditAccount(ctx, transaction.Email, transaction.Amount, time.Now().Unix(),
			fmt.Sprintf("direct %s deposit", transaction.Currency)); err != nil {
			return "", err
		}
		return string(TransactionStatuses.Completed), nil
	}

	return string(TransactionStatuses.Processing), nil
}

func (m Module) assignTransactionToAgent(ctx context.Context, transaction *TransactionOutput) error {
	panic("not implemented")
}
