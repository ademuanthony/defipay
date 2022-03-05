package app

import (
	"context"
	"merryworld/metatradas/postgres/models"
)

type store interface {
	CreateAccount(ctx context.Context, input CreateAccountInput) error
	GetAccount(ctx context.Context, id string) (*models.Account, error)
	GetAccounts(ctx context.Context, skip, limit int) ([]*models.Account, error)
	GetAllAccountsCount(ctx context.Context) (int64, error)
	GetAccountByUsername(ctx context.Context, username string) (*models.Account, error)
	UpdateAccountDetail(ctx context.Context, accountID string, input UpdateDetailInput) error

	GetDepositAddress(ctx context.Context, accountID string) (*models.Wallet, error)
	GetDeposits(ctx context.Context, accountID string, offset, limit int) ([]*models.Deposit, int64, error)

	CreatePackage(ctx context.Context, pkg models.Package) error
	PatchPackage(ctx context.Context, id string, input UpdatePackageInput) error
	GetPackages(ctx context.Context, offset, limit int) ([]*models.Package, int64, error)
	GetPackage(ctx context.Context, id string) (*models.Package, error)
	GetPackageByName(ctx context.Context, name string) (*models.Package, error)
	CreateSubscription(ctx context.Context, accountID, packageID string) error
	ActiveSubscription(ctx context.Context, accountID string) (*models.Subscription, error)
}
