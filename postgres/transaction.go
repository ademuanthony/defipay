package postgres

import (
	"context"
	"deficonnect/defipayapi/app"
	"deficonnect/defipayapi/postgres/models"

	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (pg PgDb) CreateTransaction(ctx context.Context, input app.CreateTransactionInput) (*app.TransactionOutput, error) {
	transaction := models.Transaction{
		ID:            uuid.NewString(),
		Amount:        input.Amount,
		BankName:      input.BankName,
		AccountNumber: input.AccountNumber,
		AccountName:   input.AccountName,
		Email:         input.Email,
		Network:       input.Network,
		Currency:      input.Currency,
		WalletAddress: input.WalletAddress,
		PrivateKey:    input.PrivateKey,
		PaymentLink:   input.PaymentLink,
		Type:          string(input.Type),
	}

	if err := transaction.Insert(ctx, pg.Db, boil.Infer()); err != nil {
		return nil, err
	}

	return &app.TransactionOutput{
		ID:            transaction.ID,
		Amount:        input.Amount,
		BankName:      input.BankName,
		AccountNumber: input.AccountNumber,
		AccountName:   input.AccountName,
		Email:         input.Email,
		Network:       input.Network,
		Currency:      input.Currency,
		WalletAddress: input.WalletAddress,
		PaymentLink:   input.PaymentLink,
		Type:          string(input.Type),
	}, nil
}
