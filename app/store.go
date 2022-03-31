package app

import (
	"context"
	"merryworld/metatradas/postgres/models"
)

type Trade struct {
	ID             string `boil:"id" json:"id" toml:"id" yaml:"id"`
	AccountID      string `boil:"account_id" json:"account_id" toml:"account_id" yaml:"account_id"`
	TradeNo        int    `boil:"trade_no" json:"trade_no" toml:"trade_no" yaml:"trade_no"`
	Date           int64  `boil:"date" json:"date" toml:"date" yaml:"date"`
	StartDate      int64  `boil:"start_date" json:"start_date" toml:"start_date" yaml:"start_date"`
	EndDate        int64  `boil:"end_date" json:"end_date" toml:"end_date" yaml:"end_date"`
	Amount         int64  `boil:"amount" json:"amount" toml:"amount" yaml:"amount"`
	Profit         int64  `boil:"profit" json:"profit" toml:"profit" yaml:"profit"`
}

type store interface {
	CreateAccount(ctx context.Context, input CreateAccountInput) error
	GetAccount(ctx context.Context, id string) (*models.Account, error)
	GetAccounts(ctx context.Context, skip, limit int) ([]*models.Account, error)
	GetAccountIDs(ctx context.Context) ([]string, error)
	GetAllAccountsCount(ctx context.Context) (int64, error)
	GetAccountByUsername(ctx context.Context, username string) (*models.Account, error)
	UpdateAccountDetail(ctx context.Context, accountID string, input UpdateDetailInput) error
	MyDownlines(ctx context.Context, accountID string, generation int64, offset, limit int) ([]DownlineInfo, int64, error) 
	GetRefferalCount(ctx context.Context, accountID string) (int64, error)
	GetTeamInformation(ctx context.Context, accountID string) (*TeamInfo, error)

	CreditAccount(ctx context.Context, accountID string, amount, date int64, ref string) error
	CreateDepositWallet(ctx context.Context, accountID, address, privateKey string) (*models.Wallet, error)
	GetDepositAddress(ctx context.Context, accountID string) (*models.Wallet, error)
	GetDeposits(ctx context.Context, accountID string, offset, limit int) ([]*models.Deposit, int64, error)

	CreatePackage(ctx context.Context, pkg models.Package) error
	PatchPackage(ctx context.Context, id string, input UpdatePackageInput) error
	GetPackages(ctx context.Context) ([]*models.Package, error)
	GetPackage(ctx context.Context, id string) (*models.Package, error)
	GetPackageByName(ctx context.Context, name string) (*models.Package, error)
	CreateSubscription(ctx context.Context, accountID, packageID string, c250 bool) error
	ActiveSubscription(ctx context.Context, accountID string) (*models.Subscription, error)
	PendingReferralPayouts(ctx context.Context) (models.ReferralPayoutSlice, error) 
	UpdateReferralPayout(ctx context.Context, payout *models.ReferralPayout) error
	Invest(ctx context.Context, accountID string, amount int64) error
	Investments(ctx context.Context, accountId string, offset, limit int) ([]*models.Investment, int64, error)
	Investment(ctx context.Context, id string) (*models.Investment, error)
	ReleaseInvestment(ctx context.Context, id string) error
	BuildTradingSchedule(ctx context.Context) error 
	PopulateTrades(ctx context.Context) error
	PopulateEarnings(ctx context.Context) error
	ActiveTrades(ctx context.Context, accountID string) ([]Trade, error)
	DailyEarnings(ctx context.Context, accountId string, offset, limit int) ([]*models.DailyEarning, int64, error)
	ProcessWeeklyPayout(ctx context.Context) error

	Transfer(ctx context.Context, senderID, receiverID string, amount int64) error
	TransferHistories(ctx context.Context, accountID string, offset, limit int) ([]TransferViewModel, int64, error)

	Withdraw(ctx context.Context, accountID string, amount int64) error
	Withdrawals(ctx context.Context, accountID string, offset, limit int) ([]*models.Withdrawal, int64, error)

	GetWalletByAddresses(ctx context.Context) ([]string, error)
	GetWellatByAddress(ctx context.Context, address string) (*models.Wallet, error)
	CreateDeposit(ctx context.Context, accountID, txHash string, amount int64) error

	CreateNotification(ctx context.Context, accountID, title, message, actionText, actionLink string, notificationType int) error
	NotifyAll(ctx context.Context, titile string, content, actionText, actionLink string, notificationType int) error
	UnReadNotificationCount(ctx context.Context, accountID string, notificationType int) (int64, error)
	GetNotifications(ctx context.Context, accountID string, notificationType int, offset, limit int) (models.NotificationSlice, int64, error)
	GetNewNotifications(ctx context.Context, accountID string, notificationType int, offset, limit int) (models.NotificationSlice, int64, error)
	GetNotification(ctx context.Context, id string) (*models.Notification, error)

	SetConfigValue(ctx context.Context, accountID, key string, value ConfigValue) error
	GetConfigValue(ctx context.Context, accountID, key string) (ConfigValue, error)
	GetConfigs(ctx context.Context, accountID string) (models.UserSettingSlice, error)
}
