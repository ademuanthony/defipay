package app

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"deficonnect/defipayapi/web"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
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

func (m module) getTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		web.SendErrorfJSON(w, "ID is required")
		return
	}

	transaction, err := m.db.Transaction(r.Context(), id)
	if err != nil {
		m.handleError(w, err, "Get Transaction")
		return
	}

	web.SendJSON(w, transaction)
}

func (m module) getTransactions(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	accountID := m.server.GetUserIDTokenCtx(r)
	pagedReq := web.GetPaginationInfo(r)

	transactions, count, err := m.db.Transactions(r.Context(), GetTransactionsInput{
		Email: email, AccountID: accountID, Offset: pagedReq.Offset, Limit: pagedReq.Limit,
	})

	if err != nil {
		m.handleError(w, err, "Get Transactions")
		return
	}

	web.SendPagedJSON(w, transactions, count)
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

func (m module) updateTransactionCurrency(w http.ResponseWriter, r *http.Request) {
	var input UpdateCurrencyInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Error("Login", "json::Decode", err)
		web.SendErrorfJSON(w, "cannot decode request")
		return
	}

	if m.currencyProcessors[Network(input.Currency)] == nil {
		web.SendErrorfJSON(w, "Unsupported currency")
		return
	}

	transaction, err := m.db.Transaction(r.Context(), input.TransactionID)
	if err != nil {
		m.handleError(w, err, "Get Transaction")
		return
	}

	if transaction.Email != m.server.GetUserIDTokenCtx(r) {
		web.SendErrorfJSON(w, "Invalid operation")
		return
	}

	txOutput, err := m.db.UpdateCurrency(r.Context(), input)
	if err != nil {
		m.handleError(w, err, "Update Currency")
		return
	}

	web.SendJSON(w, txOutput)
}

func (m module) checkTransactionStatus(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	transaction, err := m.db.Transaction(r.Context(), id)
	if err != nil {
		m.handleError(w, err, "Get Transaction")
		return
	}

	if transaction.Status == string(TransactionStatuses.Completed) || transaction.Status == string(TransactionStatuses.Cancelled) ||
		transaction.Status == string(TransactionStatuses.Paid) || transaction.Status == string(TransactionStatuses.Processing) {
		web.SendJSON(w, transaction)
		return
	}

	currencyProcessor := m.currencyProcessors[(transaction.Currency)]

	amountPaid, err := currencyProcessor.BalanceOf(nil, common.HexToAddress(transaction.WalletAddress))
	if err != nil {
		m.handleError(w, err)
		return
	}

	if err := m.db.UpdateTransactionPayment(r.Context(), id, amountPaid.String()); err != nil {
		m.handleError(w, err)
		return
	}

	transaction.AmountPaid = amountPaid.String()
	var tokenAmount *big.Int
	tokenAmount, _ = tokenAmount.SetString(transaction.TokenAmount, 64)

	if c := tokenAmount.Cmp(amountPaid); c == 0 || c == -1 {
		status, err := m.processTransaction(r.Context(), transaction)
		if err != nil {
			m.handleError(w, err, "process transaction")
			return
		}
		transaction.Status = status
	}

	web.SendJSON(w, transaction)
}

func (m module) processTransaction(ctx context.Context, transaction *TransactionOutput) (string, error) {
	currencyProcessor := m.currencyProcessors[(transaction.Currency)]
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

	client := m.bscClient 
	if transaction.Network == string(Networks.Polygon) {
		client = m.polygonClient
	}

	opt, err := getAccountAuth(client, pk, 1)
	if err != nil {
		return "", err
	}

	if _, err := currencyProcessor.Transfer(opt, common.HexToAddress(m.config.MasterAddress), amountPaid); err != nil {
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

func getAccountAuth(client *ethclient.Client, privateKeyString string, gasMultiplyer int64) (*bind.TransactOpts, error) {

	privateKey, err := crypto.HexToECDSA(privateKeyString)
	if err != nil {
		panic(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("invalid key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	//fetch the last use nonce of account
	nounce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	// chainID, err := client.ChainID(context.Background())
	// if err != nil {
	// 	panic(err)
	// }

	chainID := big.NewInt(137)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		panic(err)
	}
	if gasPrice.Cmp(big.NewInt(30001000047)) == -1 {
		gasPrice = big.NewInt(30001000047)
	}
	auth.Nonce = big.NewInt(int64(nounce))
	auth.Value = big.NewInt(0)                     // in wei 10:56
	auth.GasLimit = uint64(384696 * gasMultiplyer) // in units
	auth.GasPrice = gasPrice                       //big.NewInt(30001000047)

	return auth, nil
}

func (m module) assignTransactionToAgent(ctx context.Context, transaction *TransactionOutput) error {
	panic("not implemented")
}
