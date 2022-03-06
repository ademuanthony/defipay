package postgres

import (
	"context"
	"errors"
	"fmt"
	"merryworld/metatradas/app"
	"merryworld/metatradas/postgres/models"
	"time"

	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (pg PgDb) Transfer(ctx context.Context, senderID, receiverID string, amount int64) error {
	sender, err := pg.GetAccount(ctx, senderID)
	if err != nil {
		return fmt.Errorf("GetAccount::sender %v", err)
	}

	receiver, err := pg.GetAccount(ctx, receiverID)
	if err != nil {
		return fmt.Errorf("GetAccount::receiver %v", err)
	}

	if receiver.Balance < amount {
		return errors.New("insufficient balance")
	}

	tx, err := pg.Db.Begin()
	if err != nil {
		return fmt.Errorf("Begin %v", err)
	}

	date := time.Now().Unix()

	if err := pg.DebitAccountTx(ctx, tx, senderID, amount, date,
		fmt.Sprintf("direct transfer to %s at %v", receiver.Username, time.Now())); err != nil {
		tx.Rollback()
		return fmt.Errorf("DebitAccountTx %v", err)
	}

	if err := pg.CreditAccountTx(ctx, tx, receiverID, amount, date,
		fmt.Sprintf("direct transfer from %s at %v", sender.Username, time.Now())); err != nil {
		tx.Rollback()
		return fmt.Errorf("CreditAccountTx %v", err)
	}

	transfer := models.Transfer{
		ID:         uuid.NewString(),
		SenderID:   senderID,
		ReceiverID: receiverID,
		Amount:     amount,
		Date:       date,
	}
	if err := transfer.Insert(ctx, pg.Db, boil.Infer()); err != nil {
		return fmt.Errorf("transfer.Insert %v", err)
	}

	return tx.Commit()
}

func (pg PgDb) TransferHistories(ctx context.Context, accountID string, offset, limit int) ([]app.TransferViewModel, int64, error) {

	var queries = []qm.QueryMod{
		qm.Where("sender_id = $1 or receiver_id = $1", accountID),
	}

	count, err := models.Transfers(queries...).Count(ctx, pg.Db)
	if err != nil {
		return nil, 0, err
	}

	queries = append(queries,
		qm.Offset(offset),
		qm.Limit(limit),
		qm.Load(models.TransferRels.Receiver),
		qm.Load(models.TransferRels.Sender),
		qm.OrderBy(models.TransferColumns.Date+" desc"),
	)

	rec, err := models.Transfers(
		queries...,
	).All(ctx, pg.Db)

	if err != nil {
		return nil, 0, err
	}

	var models []app.TransferViewModel
	for _, item := range rec {
		models = append(models, app.TransferViewModel{
			ID: item.ID,
			Sender: item.R.Sender.Username,
			Receiver: item.R.Receiver.Username,
			Amount: item.Amount,
			Date: item.Date,
		})
	}

	return models, count, nil
}
