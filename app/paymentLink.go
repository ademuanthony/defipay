package app

import (
	"context"
	"deficonnect/defipayapi/web"
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
)

type CreatePaymentLinkInput struct {
	Permalink     string `boil:"permalink" json:"permalink" toml:"permalink" yaml:"permalink"`
	Email         string `govalid:"req" json:"email" toml:"email" yaml:"email"`
	AccountID     string `json:"-"`
	AccountName   string `govalid:"req" json:"accountName" toml:"accountname" yaml:"accountname"`
	AccountNumber string `govalid:"req" json:"accountNumber" toml:"accountnumber" yaml:"accountnumber"`
	BankName      string `govalid:"req" json:"bankName" toml:"bankname" yaml:"bankname"`
	FixedAmount   int64  `json:"fixedAmount" toml:"fixamount" yaml:"fixamount"`
	Title         string `boil:"title" json:"title" toml:"title" yaml:"title"`
	Description   string `boil:"description" json:"description" toml:"description" yaml:"description"`
}

type PaymentLinkOutput struct {
	Permalink     string `boil:"permalink" json:"permalink" toml:"permalink" yaml:"permalink"`
	Email         string `boil:"email" json:"email" toml:"email" yaml:"email"`
	AccountName   string `boil:"accountName" json:"accountName"`
	AccountNumber string `boil:"accountNumber" json:"accountNumber" toml:"accountnumber" yaml:"accountnumber"`
	BankName      string `boil:"bankName" json:"bankName" toml:"bankname" yaml:"bankname"`
	FixedAmount   int64  `boil:"fixAmount" json:"fixedAmount" toml:"fixamount" yaml:"fixamount"`
	Title         string `boil:"title" json:"title" toml:"title" yaml:"title"`
	Description   string `boil:"description" json:"description" toml:"description" yaml:"description"`
}

type GetPaymentLinksInput struct {
	Offset    int
	Limit     int
	Email     string
	AccountID string
}

func (m Module) CreatePaymentLink(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	var input CreatePaymentLinkInput
	if err := json.Unmarshal([]byte(r.Body), &input); err != nil {
		log.Error("Login", "json::Decode", err)
		return SendErrorfJSON("cannot decode request")
	}

	if input.Permalink != "" {
		if _, err := m.db.GetPaymentLink(ctx, input.Permalink); err == nil {
			return SendErrorfJSON("permalink not available")
		}
	} else {
		for {
			input.Permalink = strings.ReplaceAll(uuid.NewString(), "-", "")[0:10]
			if p, _ := m.db.GetPaymentLink(ctx, input.Permalink); p == nil {
				break
			}
		}
	}

	input.AccountID = m.server.GetUserIDTokenCtxSls(r)

	if err := m.db.CreatePaymentLink(ctx, input); err != nil {
		return m.handleError(err, "Create Payment Link")
	}

	return SendJSON(input)
}

func (m Module) GetPaymentLink(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	permalink := r.PathParameters["permalink"]

	paymentLink, err := m.db.GetPaymentLink(ctx, permalink)

	if err != nil {
		return m.handleError(err, "Get Payment Link")
	}

	return SendJSON(paymentLink)
}

func (m Module) GetPaymentLinks(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	accountID := m.server.GetUserIDTokenCtxSls(r)
	email := r.QueryStringParameters["email"]

	pagedReq := web.GetPaginationInfoSls(r)

	paymentLinks, count, err := m.db.GetPaymentLinks(ctx, GetPaymentLinksInput{
		Limit: pagedReq.Limit, Offset: pagedReq.Offset, AccountID: accountID, Email: email,
	})

	if err != nil {
		return m.handleError(err, "Get Payment Links")
	}

	return SendPagedJSON(paymentLinks, count)
}
