package app

import (
	"context"
	"merryworld/metatradas/postgres/models"
)

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
	PopulateEarnings(ctx context.Context) error
	DailyEarnings(ctx context.Context, accountId string, offset, limit int) ([]*models.DailyEarning, int64, error)
	ProcessWeeklyPayout(ctx context.Context) error

	Transfer(ctx context.Context, senderID, receiverID string, amount int64) error
	TransferHistories(ctx context.Context, accountID string, offset, limit int) ([]TransferViewModel, int64, error)

	Withdraw(ctx context.Context, accountID string, amount int64) error
	Withdrawals(ctx context.Context, accountID string, offset, limit int) ([]*models.Withdrawal, int64, error)

	GetWalletByAddresses(ctx context.Context) ([]string, error)
	GetWellatByAddress(ctx context.Context, address string) (*models.Wallet, error)
	CreateDeposit(ctx context.Context, accountID, txHash string, amount int64) error

	CreateNotification(ctx context.Context, accountID, title, message string) error
	UnReadNotificationCount(ctx context.Context, accountID string) (int64, error)
	GetNotifications(ctx context.Context, accountID string, offset, limit int) (models.NotificationSlice, int64, error)
	GetNewNotifications(ctx context.Context, accountID string, offset, limit int) (models.NotificationSlice, int64, error)
	GetNotification(ctx context.Context, id string) (*models.Notification, error)
}
