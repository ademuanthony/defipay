package app

import (
	"context"
	"deficonnect/defipayapi/web"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
)

type BeneficiaryOutput struct {
	ID            string `json:"id"`
	AccountID     string `json:"accountID"`
	Bank          string `boil:"bank" json:"bankName" toml:"bank" yaml:"bank"`
	AccountNumber string `boil:"account_number" json:"accountNumber" toml:"account_number" yaml:"account_number"`
	AccountName   string `boil:"account_name" json:"accountName" toml:"account_name" yaml:"account_name"`
	Country       string `boil:"country" json:"country" toml:"country" yaml:"country"`
}

type CreateBeneficiaryInput struct {
	ID            string `json:"id"`
	AccountID     string `json:"-"`
	Bank          string `govalid:"req" json:"bankName" toml:"bank" yaml:"bank"`
	AccountNumber string `govalid:"req" json:"accountNumber" toml:"account_number" yaml:"account_number"`
	AccountName   string `govalid:"req" json:"accountName" toml:"account_name" yaml:"account_name"`
	Country       string `govalid:"req" json:"country" toml:"country" yaml:"country"`
}

type GetBeneficiariesInput struct {
	Limit     int
	Offset    int
	AccountID string
}

func (m Module) CreateBeneficiary(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	var input CreateBeneficiaryInput
	if err := json.Unmarshal([]byte(r.Body), &input); err != nil {
		log.Error("Login", "json::Decode", err)
		return SendErrorfJSON("cannot decode request")
	}

	vio, err := v.Violations(&input)
	if err != nil {
		return m.sendSomethingWentWrong("Violations", err)
	}
	if len(vio) > 0 {
		return m.handleError(NewValidationError(vio))
	}

	accountID := m.GetUserIDTokenCtxSls(r)
	if accountID == "" {
		return SendAuthErrorfJSON("Login required")
	}

	input.AccountID = accountID

	if b, _ := m.db.GetMyBeneficiaryByAccountNumber(ctx, accountID, input.AccountNumber); b != nil {
		return SendErrorfJSON("Beneficiary exists")
	}

	input.ID = uuid.NewString()

	if err := m.db.CreateBeneficiary(ctx, input); err != nil {
		return m.handleError(err)
	}

	return SendJSON(input)
}

func (m Module) GetBeneficiaries(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	accountID := m.GetUserIDTokenCtxSls(r)
	if accountID == "" {
		return SendAuthErrorfJSON("Login required")
	}

	pagedReq := web.GetPaginationInfoSls(r)

	beneficiaries, count, err := m.db.GetBeneficiaries(ctx, GetBeneficiariesInput{
		Limit:     pagedReq.Limit,
		Offset:    pagedReq.Offset,
		AccountID: accountID,
	})

	if err != nil {
		return m.handleError(err)
	}

	return SendPagedJSON(beneficiaries, count)
}

func (m Module) GetBeneficiary(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	id := request.PathParameters["id"]
	transaction, err := m.db.GetBeneficiary(ctx, id)
	if err != nil {
		return m.handleError(err)
	}
	return SendJSON(transaction)
}
