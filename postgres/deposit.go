package postgres

import (
	"context"
	"merryworld/metatradas/postgres/models"
	"time"

	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (pg PgDb) GetWellatByAddress(ctx context.Context, address string) (*models.Wallet, error) {
	return models.Wallets(models.WalletWhere.Address.EQ(address)).One(ctx, pg.Db)
}

func (pg PgDb) CreateDeposit(ctx context.Context, accountID, txHash string, amount int64) error {
	tx, err := pg.Db.Begin()
	if err != nil {
		return err
	}

	date := time.Now().Unix()

	deposit := models.Deposit{
		ID:        uuid.NewString(),
		AccountID: accountID,
		Ref:       txHash,
		Amount:    amount,
		Date:      date,
	}

	if deposit.Insert(ctx, tx, boil.Infer()); err != nil {
		return err
	}

	if err := pg.CreditAccountTx(ctx, tx, accountID, amount, date, "deposit ref: "+txHash); err != nil {
		return err
	}

	return nil
}
