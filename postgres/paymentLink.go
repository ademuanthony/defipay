package postgres

import (
	"context"
	"deficonnect/defipayapi/app"
	"deficonnect/defipayapi/postgres/models"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (pg PgDb) CreatePaymentLink(ctx context.Context, input app.CreatePaymentLinkInput) error {
	paymentLink := models.PaymentLink{
		Permalink:     input.Permalink,
		AccountID:     null.StringFrom(input.AccountID),
		Email:         input.Email,
		Accountname:   input.AccountName,
		Accountnumber: input.AccountNumber,
		Bankname:      input.BankName,
		Fixamount:     input.FixedAmount,
		Title:         input.Title,
		Description:   input.Description,
	}

	return paymentLink.Insert(ctx, pg.Db, boil.Infer())
}

func (pg PgDb) GetPaymentLink(ctx context.Context, permalink string) (*app.PaymentLinkOutput, error) {
	paymentLink, err := models.FindPaymentLink(ctx, pg.Db, permalink)
	if err != nil {
		return nil, err
	}

	return convertPaymentLink(paymentLink), nil
}

func (pg PgDb) GetPaymentLinks(ctx context.Context, input app.GetPaymentLinksInput) ([]*app.PaymentLinkOutput, int64, error) {
	query := []qm.QueryMod{
		qm.Where("account_id = $1 or email = $2", input.AccountID, input.Email),
	}

	count, err := models.PaymentLinks(query...).Count(ctx, pg.Db)
	if err != nil {
		return nil, 0, err
	}

	query = append(query, qm.Offset(input.Offset), qm.Limit(input.Limit))

	links, err := models.PaymentLinks(query...).All(ctx, pg.Db)
	if err != nil {
		return nil, 0, err
	}

	var outputs []*app.PaymentLinkOutput
	for _, pl := range links {
		outputs = append(outputs, convertPaymentLink(pl))
	}

	return outputs, count, nil
}

func convertPaymentLink(input *models.PaymentLink) *app.PaymentLinkOutput {
	return &app.PaymentLinkOutput{
		Permalink:     input.Permalink,
		Email:         input.Email,
		AccountName:   input.Accountname,
		AccountNumber: input.Accountnumber,
		BankName:      input.Bankname,
		FixedAmount:   input.Fixamount,
		Title:         input.Title,
		Description:   input.Description,
	}
}
