package postgres

import (
	"context"
	"deficonnect/defipayapi/app"
	"deficonnect/defipayapi/postgres/models"

	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
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

func (pg PgDb) Transaction(ctx context.Context, ID string) (*app.TransactionOutput, error) {
	transaction, err := models.FindTransaction(ctx, pg.Db, ID)
	if err != nil {
		return nil, err
	}

	txOutput := convertTransaction(transaction)

	return &txOutput, nil

}

func convertTransaction(transaction *models.Transaction) app.TransactionOutput {
	return app.TransactionOutput{
		ID:            transaction.ID,
		Amount:        transaction.Amount,
		BankName:      transaction.BankName,
		AccountNumber: transaction.AccountNumber,
		AccountName:   transaction.AccountName,
		Email:         transaction.Email,
		Network:       transaction.Network,
		Currency:      transaction.Currency,
		WalletAddress: transaction.WalletAddress,
		PaymentLink:   transaction.PaymentLink,
		Type:          string(transaction.Type),
	}
}

func (pg PgDb) Transactions(ctx context.Context, input app.GetTransactionsInput) ([]app.TransactionOutput, int64, error) {
	query := []qm.QueryMod{
		qm.Where("email = $1 or account_id = $2", input.Email, input.AccountID),
	}

	count, err := models.Transactions(query...).Count(ctx, pg.Db)
	if err != nil {
		return nil, 0, err
	}

	query = append(query, qm.Offset(input.Offset), qm.Limit(input.Limit))

	transactions, err := models.Transactions(query...).All(ctx, pg.Db)
	if err != nil {
		return nil, 0, err
	}

	var result []app.TransactionOutput
	for _, transaction := range transactions {
		result = append(result, convertTransaction(transaction))
	}

	return result, count, nil
}

func (pg PgDb) UpdateCurrency(ctx context.Context, input app.UpdateCurrencyInput) (*app.TransactionOutput, error) {
	updateCol := models.M{
		models.TransactionColumns.Network:  input.Network,
		models.TransactionColumns.Currency: input.Currency,
	}
	_, err := models.Transactions(models.TransactionWhere.ID.EQ(input.TransactionID)).UpdateAll(ctx, pg.Db, updateCol)

	if err != nil {
		return nil, err
	}

	return pg.Transaction(ctx, input.TransactionID)
}

func (pg PgDb) TransactionPK(ctx context.Context, transactionID string) (string, error) {
	tx, err := models.Transactions(models.TransactionWhere.ID.EQ(transactionID),
		qm.Select(models.TransactionColumns.PrivateKey)).One(ctx, pg.Db)

	if err != nil {
		return "", err
	}

	return tx.PrivateKey, nil
}

func (pg PgDb) UpdateTransactionStatus(ctx context.Context, transactionID string, status app.TransactionStatus) error {
	updateCol := models.M{
		models.TransactionColumns.Status: string(status),
	}
	_, err := models.Transactions(models.TransactionWhere.ID.EQ(transactionID)).UpdateAll(ctx, pg.Db, updateCol)
	return err
}
