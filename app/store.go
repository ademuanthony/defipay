package app

import (
	"context"
	"deficonnect/defipayapi/postgres/models"
)

type Trade struct {
	ID        string `boil:"id" json:"id" toml:"id" yaml:"id"`
	AccountID string `boil:"account_id" json:"account_id" toml:"account_id" yaml:"account_id"`
	TradeNo   int    `boil:"trade_no" json:"trade_no" toml:"trade_no" yaml:"trade_no"`
	Date      int64  `boil:"date" json:"date" toml:"date" yaml:"date"`
	StartDate int64  `boil:"start_date" json:"start_date" toml:"start_date" yaml:"start_date"`
	EndDate   int64  `boil:"end_date" json:"end_date" toml:"end_date" yaml:"end_date"`
	Amount    int64  `boil:"amount" json:"amount" toml:"amount" yaml:"amount"`
	Profit    int64  `boil:"profit" json:"profit" toml:"profit" yaml:"profit"`
}

type store interface {
	CreateAccount(ctx context.Context, input CreateAccountInput) error
	GetAccount(ctx context.Context, id string) (*models.Account, error)
	AccountBalance(ctx context.Context, accountId string) (int64, error)
	GetAccounts(ctx context.Context, skip, limit int) ([]*models.Account, error)
	GetPasswordResetCode(ctx context.Context, accountID string) (string, error)
	ValidatePasswordResetCode(ctx context.Context, accountID, code string) (bool, error)
	ChangePassword(ctx context.Context, accountID, password string) error
	GetAccountIDs(ctx context.Context) ([]string, error)
	GetAllAccountsCount(ctx context.Context) (int64, error)
	GetAccountByEmail(ctx context.Context, email string) (*models.Account, error)
	UpdateAccountDetail(ctx context.Context, accountID string, input UpdateDetailInput) error
	GetRefferalCount(ctx context.Context, accountID string) (int64, error)

	CreditAccount(ctx context.Context, accountID string, amount, date int64, ref string) error

	CreateNotification(ctx context.Context, accountID, title, message, actionText, actionLink string, notificationType int) error
	UnReadNotificationCount(ctx context.Context, accountID string, notificationType int) (int64, error)
	GetNotifications(ctx context.Context, accountID string, notificationType int, offset, limit int) (models.NotificationSlice, int64, error)
	GetNewNotifications(ctx context.Context, accountID string, notificationType int, offset, limit int) (models.NotificationSlice, int64, error)
	GetNotification(ctx context.Context, id string) (*models.Notification, error)

	SetConfigValue(ctx context.Context, accountID, key string, value ConfigValue) error
	GetConfigValue(ctx context.Context, accountID, key string) (ConfigValue, error)
	GetConfigs(ctx context.Context, accountID string) (models.UserSettingSlice, error)
	AddLogin(ctx context.Context, accountID, ip, platform string, date int64) error
	LastLogin(ctx context.Context) (*models.LoginInfo, error)

	// Transaction
	CreateTransaction(ctx context.Context, input CreateTransactionInput) (*TransactionOutput, error)
	Transaction(ctx context.Context, ID string) (*TransactionOutput, error)
	Transactions(ctx context.Context, input GetTransactionsInput) ([]TransactionOutput, int64, error)
	UpdateCurrency(txt context.Context, input UpdateCurrencyInput) (*TransactionOutput, error)
	TransactionPK(ctx context.Context, transactionID string) (string, error)
	UpdateTransactionStatus(ctx context.Context, transactionID string, status TransactionStatus) error
	UpdateTransactionPayment(ctx context.Context, transactionID string, amountPaid string) error

	// Payment Link
	CreatePaymentLink(ctx context.Context, input CreatePaymentLinkInput) error
	GetPaymentLink(ctx context.Context, permalink string) (*PaymentLinkOutput, error)
	GetPaymentLinks(ctx context.Context, input GetPaymentLinksInput) ([]*PaymentLinkOutput, int64, error)

	// beneficiary
	GetMyBeneficiaryByAccountNumber(ctx context.Context, accountID string, accountNumber string) (*BeneficiaryOutput, error)
	CreateBeneficiary(ctx context.Context, input CreateBeneficiaryInput) error
	GetBeneficiaries(ctx context.Context, input GetBeneficiariesInput) ([]*BeneficiaryOutput, int64, error)
	GetBeneficiary(ctx context.Context, id string) (*BeneficiaryOutput, error)

	// Agent
	CreateAgent(ctx context.Context, input CreateAgentInput) error
	UpdateAgent(ctx context.Context, slackUsername string, input UpdateAgentInput) (AgentOutput, error)
	AgentExists(ctx context.Context, slackUsername string) (bool, error)
	GetAgent(ctx context.Context, slackUsername string) (*AgentOutput, error)
	GetAgents(ctx context.Context, input GetAgentsInput) ([]AgentOutput, int64, error)
	NextAvailableAgent(ctx context.Context, transactionAmount int64) (*AgentOutput, error)
	AssignAgent(ctx context.Context, agentID string, transactionID string) error
	GetAssignedAgent(ctx context.Context, transactionID string) (AgentOutput, error)
	GetAgentAssignments(ctx context.Context, input GetAgentAssignmentsInput) ([]AssignmentOutput, error)
}
