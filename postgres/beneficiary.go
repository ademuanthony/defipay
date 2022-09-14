package postgres

import (
	"context"
	"deficonnect/defipayapi/app"
	"deficonnect/defipayapi/postgres/models"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (pg PgDb) GetMyBeneficiaryByAccountNumber(ctx context.Context, accountID string, accountNumber string) (*app.BeneficiaryOutput, error) {
	ben, err := models.Beneficiaries(
		models.BeneficiaryWhere.AccountID.EQ(null.StringFrom(accountID)),
		models.BeneficiaryWhere.AccountNumber.EQ(accountNumber),
	).One(ctx, pg.Db)

	if err != nil {
		return nil, err
	}

	return convertBeneficiary(ben), nil
}

func (pg PgDb) CreateBeneficiary(ctx context.Context, input app.CreateBeneficiaryInput) error {
	ben := models.Beneficiary{
		ID:              input.ID,
		AccountID:       null.StringFrom(input.AccountID),
		Bank:            input.Bank,
		AccountNumber:   input.AccountNumber,
		AccountName:     input.AccountName,
		Country:         input.Country,
		BeneficialEmail: input.BeneficialEmail,
	}

	return ben.Insert(ctx, pg.Db, boil.Infer())
}

func (pg PgDb) GetBeneficiaries(ctx context.Context, input app.GetBeneficiariesInput) ([]*app.BeneficiaryOutput, int64, error) {
	query := []qm.QueryMod{
		models.BeneficiaryWhere.AccountID.EQ(null.StringFrom(input.AccountID)),
	}

	count, err := models.Beneficiaries(query...).Count(ctx, pg.Db)
	if err != nil {
		return nil, 0, err
	}

	query = append(query, qm.Limit(input.Limit), qm.Offset(input.Offset))

	bens, err := models.Beneficiaries(query...).All(ctx, pg.Db)
	if err != nil {
		return nil, 0, err
	}

	var result []*app.BeneficiaryOutput
	for _, b := range bens {
		result = append(result, convertBeneficiary(b))
	}

	return result, count, nil
}

func (pg PgDb) GetBeneficiary(ctx context.Context, id string) (*app.BeneficiaryOutput, error) {
	ben, err := models.FindBeneficiary(ctx, pg.Db, id)
	if err != nil {
		return nil, err
	}

	return convertBeneficiary(ben), nil
}

func convertBeneficiary(input *models.Beneficiary) *app.BeneficiaryOutput {
	return &app.BeneficiaryOutput{
		ID:              input.ID,
		AccountID:       input.AccountID.String,
		Bank:            input.Bank,
		AccountNumber:   input.AccountNumber,
		AccountName:     input.AccountName,
		Country:         input.Country,
		BeneficialEmail: input.BeneficialEmail,
	}
}
