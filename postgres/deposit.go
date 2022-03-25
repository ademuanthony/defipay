package postgres

import (
	"context"
	"fmt"
	"merryworld/metatradas/postgres/models"
	"time"

	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (pg PgDb) GetWellatByAddress(ctx context.Context, address string) (*models.Wallet, error) {
	where := fmt.Sprintf("lower(%s) = lower($1)", models.WalletColumns.Address)
	return models.Wallets(qm.Where(where, address)).One(ctx, pg.Db)
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
		tx.Rollback()
		return err
	}

	if err := pg.CreditAccountTx(ctx, tx, accountID, amount, date, "deposit ref: "+txHash); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pg PgDb) GetWalletByAddresses(ctx context.Context) ([]string, error) {
	wallets, err := models.Wallets(
		qm.Select(models.WalletColumns.Address),
	).All(ctx, pg.Db)
	if err != nil {
		return nil, err
	}

	var addresses []string
	for _, wallet := range wallets {
		addresses = append(addresses, wallet.Address)
	}

	return addresses, nil
}
