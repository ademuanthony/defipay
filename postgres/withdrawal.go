package postgres

import (
	"context"
	"errors"
	"fmt"
	"merryworld/metatradas/postgres/models"
	"time"

	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (pg PgDb) Withdraw(ctx context.Context, accountID string, amount int64) error {
	account, err := pg.GetAccount(ctx, accountID)
	if err != nil {
		return fmt.Errorf("GetAccount::sender %v", err)
	}

	if account.Balance < amount {
		return errors.New("insufficient balance")
	}

	tx, err := pg.Db.Begin()
	if err != nil {
		return fmt.Errorf("Begin %v", err)
	}

	date := time.Now().Unix()

	if err := pg.DebitAccountTx(ctx, tx, accountID, amount, date,
		fmt.Sprintf("fund withdrawal to %s at %v", account.WithdrawalAddresss, time.Now())); err != nil {
		tx.Rollback()
		return fmt.Errorf("DebitAccountTx %v", err)
	}

	model := models.Withdrawal{
		ID:          uuid.NewString(),
		Amount:      amount,
		Date:        date,
		AccountID:   accountID,
		Destination: account.WithdrawalAddresss,
	}
	if err := model.Insert(ctx, pg.Db, boil.Infer()); err != nil {
		return fmt.Errorf("model.Insert %v", err)
	}

	return tx.Commit()
}

func (pg PgDb) Withdrawals(ctx context.Context, accountID string, offset, limit int) ([]*models.Withdrawal, int64, error) {

	var queries = []qm.QueryMod{
		models.WithdrawalWhere.AccountID.EQ(accountID),
	}

	count, err := models.Withdrawals(queries...).Count(ctx, pg.Db)
	if err != nil {
		return nil, 0, err
	}

	queries = append(queries,
		qm.Offset(offset),
		qm.Limit(limit),
		qm.OrderBy(models.WithdrawalColumns.Date+" desc"),
	)

	rec, err := models.Withdrawals(
		queries...,
	).All(ctx, pg.Db)

	if err != nil {
		return nil, 0, err
	}

	return rec, count, nil
}

func (pg PgDb) GetWithdrawalsForProcessing(ctx context.Context) (models.WithdrawalSlice, error) {
	statement := `select withdrawal.id, destination, amount, account_id from withdrawal 
		inner join account on account.id = withdrawal.account_id
		where withdrawal.ref = '' and account.balance >= 0 and destination <> ''`
	return models.Withdrawals(
		qm.SQL(statement),
	).All(ctx, pg.Db)
}

func (pg PgDb) SetWithdrawalTxHash(ctx context.Context, withdarwalID, txHash string) error {
	col := models.M{
		models.WithdrawalColumns.Ref: txHash,
	}
	_, err := models.Withdrawals(models.WithdrawalWhere.ID.EQ(withdarwalID)).UpdateAll(ctx, pg.Db, col)
	return err
}
