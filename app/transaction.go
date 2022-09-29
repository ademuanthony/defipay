package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"deficonnect/defipayapi/postgres/models"
	"deficonnect/defipayapi/web"

	"github.com/aws/aws-lambda-go/events"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type CreateTransactionInput struct {
	BankName        string `json:"bankName" toml:"bank_name" yaml:"bank_name"`
	AccountNumber   string `json:"accountNumber" toml:"account_number" yaml:"account_number"`
	AccountName     string `json:"accountName" toml:"account_name" yaml:"account_name"`
	Amount          int64  `govalid:"req|min:10|max:100000000" json:"amount" toml:"amount" yaml:"amount"`
	TokenAmount     string `json:"-"`
	Email           string `govalid:"req" json:"email" toml:"email" yaml:"email"`
	Network         string `json:"network" toml:"network" yaml:"network"`
	Currency        string `json:"currency" toml:"currency" yaml:"currency"`
	PaymentLink     string `boil:"payment_link" json:"paymentLink" toml:"payment_link" yaml:"payment_link"`
	Type            string `govalid:"req" json:"type" toml:"type" yaml:"type"`
	PaymentMethod   string `json:"paymentMethod"`
	SaveBeneficiary bool   `json:"saveBeneficiary"`

	WalletAddress string `json:"-"`
	PrivateKey    string `json:"-"`
}

type TransactionOutput struct {
	ID            string `boil:"id" json:"id" toml:"id" yaml:"id"`
	BankName      string `boil:"bank_name" json:"bankName" toml:"bank_name" yaml:"bank_name"`
	AccountNumber string `boil:"account_number" json:"accountNumber" toml:"account_number" yaml:"account_number"`
	AccountName   string `boil:"account_name" json:"accountName" toml:"account_name" yaml:"account_name"`
	Amount        int64  `boil:"amount" json:"amount" toml:"amount" yaml:"amount"`
	AmountPaid    string `json:"amountPaid"`
	TokenAmount   string `json:"tokenAmount"`
	Email         string `boil:"email" json:"email" toml:"email" yaml:"email"`
	Network       string `boil:"network" json:"network" toml:"network" yaml:"network"`
	Currency      string `boil:"currency" json:"currency" toml:"currency" yaml:"currency"`
	WalletAddress string `boil:"wallet_address" json:"walletAddress" toml:"wallet_address" yaml:"wallet_address"`
	PaymentLink   string `boil:"payment_link" json:"paymentLink" toml:"payment_link" yaml:"payment_link"`
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

type PaymentMethod string

var PaymentMethods = struct {
	Wallet PaymentMethod
	Crypto PaymentMethod
}{
	Wallet: "wallet",
	Crypto: "crypto",
}

type GetTransactionsInput struct {
	Email     string
	AccountID string
	Offset    int
	Limit     int
}

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

func (m Module) GetTransaction(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	id := request.PathParameters["id"]
	transaction, err := m.db.Transaction(ctx, id)
	if err != nil {
		return m.handleError(err)
	}
	return SendJSON(transaction)
}

func (m Module) GetTransactions(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	email := r.QueryStringParameters["email"]
	account, err := m.currentAccount(ctx, r)
	if err == nil {
		email = account.Email
	}
	pagedReq := web.GetPaginationInfoSls(r)

	transactions, count, err := m.db.Transactions(ctx, GetTransactionsInput{
		Email: email, Offset: pagedReq.Offset, Limit: pagedReq.Limit,
	})

	if err != nil {
		return m.handleError(err)
	}

	return SendPagedJSON(transactions, count)
}

func (m Module) CreateTransaction(ctx context.Context, input CreateTransactionInput, account *models.Account) (Response, error) {
	currencyProcessor := m.currencyProcessors[(input.Currency)][Network(input.Network)]
	if currencyProcessor == nil {
		return SendErrorfJSON("Unsupported currency or network in transaction")
	}

	switch input.Type {
	case "0":
		input.Type = string(transactionTypes.FundTransfer)
	case "1":
		input.Type = string(transactionTypes.TopUp)
	default:
		return SendErrorfJSON("unsupported transaction type")
	}

	if input.Amount > 10000*1e4 {
		return SendErrorfJSON("Please enter an amount below $10,000")
	}

	inputAmount, valid := common.Big0.SetString(fmt.Sprintf("%d", input.Amount), 10)
	if !valid {
		return SendErrorfJSON("Cannot set big string")
	}
	// compute token amount
	tokenAmount, err := currencyProcessor.DollarToToken(ctx, inputAmount)
	if err != nil {
		return SendErrorfJSON("Unable to get token conversion. Please contact the customer service")
	}

	input.TokenAmount = tokenAmount.String()

	tran, err := m.createTransaction(ctx, input)
	if err != nil {
		log.Error("Create Transaction", err)
		msg := "Cannot create transaction. Please try again"
		if messenger, ok := err.(ErrorMessenger); ok {
			msg = messenger.ErrorMessage()
		}
		return SendErrorfJSON(msg)
	}

	if account != nil && input.SaveBeneficiary {
		m.db.CreateBeneficiary(ctx, CreateBeneficiaryInput{
			ID:            uuid.NewString(),
			AccountID:     account.ID,
			Bank:          input.BankName,
			AccountNumber: input.AccountNumber,
			AccountName:   input.AccountName,
			Country:       input.Currency,
		})
	}

	return SendJSON(tran)
}

func (m Module) createTransaction(ctx context.Context, input CreateTransactionInput) (*TransactionOutput, error) {
	vio, err := v.Violations(&input)
	if err != nil {
		return nil, err
	}
	if len(vio) > 0 {
		return nil, NewValidationError(vio)
	}

	if input.Type == string(transactionTypes.FundTransfer) {
		if input.AccountName == "" || input.AccountNumber == "" || input.BankName == "" {
			return nil, NewValidationError([]string{"Account details required"})
		}
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

	if transaction.Email != m.GetUserIDTokenCtxSls(r) {
		return SendErrorfJSON("Invalid operation")
	}

	txOutput, err := m.db.UpdateCurrency(ctx, input)
	if err != nil {
		return m.handleError(err, "Update Currency")
	}

	return SendJSON(txOutput)
}

func (m Module) CheckTransactionStatus(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	id := r.PathParameters["id"]
	transaction, err := m.db.Transaction(ctx, id)
	if err != nil {
		return m.handleError(err, "Get Transaction")
	}

	if transaction.Status == string(TransactionStatuses.Completed) || transaction.Status == string(TransactionStatuses.Cancelled) ||
		transaction.Status == string(TransactionStatuses.Paid) || transaction.Status == string(TransactionStatuses.Processing) {
		return SendJSON(transaction)
	}

	currencyProcessor := m.currencyProcessors[(transaction.Currency)][Network(transaction.Network)]
	if currencyProcessor == nil {
		return SendErrorfJSON("Unsupported currency or network in transaction")
	}

	amountPaid, err := currencyProcessor.BalanceOf(context.TODO(), common.HexToAddress(transaction.WalletAddress))
	if err != nil {
		return m.handleError(err)
	}

	if err := m.db.UpdateTransactionPayment(ctx, id, amountPaid.String()); err != nil {
		return m.handleError(err)
	}

	transaction.AmountPaid = amountPaid.String()
	var tokenAmount *big.Int = common.Big0
	tokenAmount, valid := tokenAmount.SetString(transaction.TokenAmount, 10)
	if !valid {
		return SendErrorfJSON("Invalid token amount. Please contact the admin for resolution")
	}

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
	amountPaid, err := currencyProcessor.BalanceOf(context.TODO(), common.HexToAddress(transaction.WalletAddress))
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}

	var tokenAmount *big.Int = common.Big0
	tokenAmount, _ = tokenAmount.SetString(transaction.TokenAmount, 10)

	if c := amountPaid.Cmp(tokenAmount); c == -1 {
		return "", errors.New("incomplete payment")
	}

	if transaction.Status == string(TransactionStatuses.Processing) { // 111000+102500+10000
		return "", errors.New("processing")
	}

	if err := m.db.UpdateTransactionStatus(ctx, transaction.ID, TransactionStatuses.Processing); err != nil {
		return "", err
	}

	// pk, err := m.db.TransactionPK(ctx, transaction.ID)
	// if err != nil {
	// 	return "", err
	// }

	// if _, err := currencyProcessor.Transfer(ctx, pk, common.HexToAddress(m.config.MasterAddress), amountPaid); err != nil {
	// 	return "", err
	// }

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
	agent, err := m.db.NextAvailableAgent(ctx, transaction.Amount)
	if err != nil {
		return err
	}
	err = m.db.AssignAgent(ctx, agent.ID, transaction.ID, transaction.Amount)
	if err != nil {
		return err
	}
	webHookUrl, err := m.db.GetConfigValue(ctx, m.config.MastAccountID, "SLACK_WEB_HOOK_URL")
	if err != nil {
		return err
	}

	conversionRateStr, err := m.db.GetConfigValue(ctx, m.config.MastAccountID, "CONVERSION_RATE")
	if err != nil {
		return err
	}

	conversionRate, err := strconv.Atoi(string(conversionRateStr))
	if err != nil {
		return err
	}

	message := fmt.Sprintf(
		`New Transaction
		Dollar Amount: %.2f
		Naira Amount: %.2f
		Account Name: %s
		Account Number %s
		Bank Name: %s

		Agent: %s (@%s)
		`,
		(float64(transaction.Amount))/float64(1e4),
		(float64(transaction.Amount*int64(conversionRate)))/float64(1e4),
		transaction.AccountName,
		transaction.AccountNumber,
		transaction.BankName,
		agent.Name,
		agent.SlackUsername,
	)

	var input = struct {
		Text string `json:"text"`
	}{
		Text: message,
	}

	postBody, err := json.Marshal(input)
	if err != nil {
		return err
	}
	requestBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(string(webHookUrl), "application/json", requestBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	// b, err := ioutil.ReadAll(resp.Body)  Go.1.15 and earlier
	if err != nil {
		return err
	}

	bodyStr := string(b)
	if strings.ToLower(bodyStr) != "ok" {
		return errors.New(bodyStr)
	}
	return nil
}
