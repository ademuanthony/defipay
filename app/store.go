package app

import (
	"context"
	"merryworld/metatradas/postgres/models"
)

type store interface {
	CreateAccount(ctx context.Context, input CreateAccountInput) error
	GetAccount(ctx context.Context, id string) (*models.Account, error)
	GetAccountByUsername(ctx context.Context, username string) (*models.Account, error)
	UpdateAccountDetail(ctx context.Context, accountID string, input UpdateDetailInput) error

	GetDepositAddress(ctx context.Context, accountID string) (*models.Wallet, error)
}
