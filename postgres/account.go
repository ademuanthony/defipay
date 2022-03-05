package postgres

import (
	"context"
	"merryworld/metatradas/app"
	"merryworld/metatradas/postgres/models"

	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (pg PgDb) CreateAccount(ctx context.Context, input app.CreateAccountInput) error {
	account := models.Account{
		ID:       uuid.NewString(),
		Username: input.Username,
		Password: input.Password,
		Email:    input.Email,
	}

	tx, err := pg.Db.Begin()
	if err != nil {
		return err
	}

	err = account.Insert(ctx, tx, boil.Infer())
	if err != nil {
		tx.Rollback()
		return err
	}

	wallet := models.Wallet{
		AccountID:  account.ID,
		Address:    input.WalletAddress,
		PrivateKey: input.PrivateKey,
		CoinSymbol: "BEP20-USDT",
	}

	if err = wallet.Insert(ctx, tx, boil.Infer()); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pg PgDb) GetAccount(ctx context.Context, id string) (*models.Account, error) {
	return models.FindAccount(ctx, pg.Db, id)
}

func (pg PgDb) GetAccountByUsername(ctx context.Context, username string) (*models.Account, error) {
	return models.Accounts(
		models.AccountWhere.Username.EQ(username),
	).One(ctx, pg.Db)
}

func (pg PgDb) UpdateAccountDetail(ctx context.Context, accountID string, input app.UpdateDetailInput) error {
	var upCol = models.M{}
	if input.FirstName != "" {
		upCol[models.AccountColumns.FirstName] = input.FirstName
	}
	if input.PhoneNumber != "" {
		upCol[models.AccountColumns.PhoneNumber] = input.PhoneNumber
	}
	if input.LastName != "" {
		upCol[models.AccountColumns.LastName] = input.LastName
	}
	if input.WithdrawalAddress != "" {
		upCol[models.AccountColumns.WithdrawalAddresss] = input.WithdrawalAddress
	}

	_, err := models.Accounts(models.AccountWhere.ID.EQ(accountID)).UpdateAll(ctx, pg.Db, upCol)
	return err
}

func (pg PgDb) GetDepositAddress(ctx context.Context, accountID string) (*models.Wallet, error) {
	return models.Wallets(models.WalletWhere.AccountID.EQ(accountID)).One(ctx, pg.Db)
}
