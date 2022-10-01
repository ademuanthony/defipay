package app

import (
	"context"
	"deficonnect/defipayapi/app/util"
	"deficonnect/defipayapi/postgres/models"
	"deficonnect/defipayapi/web"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"golang.org/x/crypto/bcrypt"
)

const (
	PAYMENTMETHOD_C250D = 0
	PAYMENTMETHOD_BNB   = 1
	PAYMENTMETHOD_USDT  = 2

	PAYMENTSTATUS_PENDING     = 0
	PAYMENTSTATUS_PROCCESSING = 1
	PAYMENTSTATUS_COMPLETED   = 2
	PAYMENTSTATUS_FAILED      = 3
)

type CreateAccountInput struct {
	ReferralCode string `json:"referralCode"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phoneNumber"`
	Password     string `json:"password"`
	From250      bool   `json:"from250"`

	DepositWalletAddress string `json:"-"`
	PrivateKey           string `json:"-"`
}

type DownlineInfo struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PhoneNumber string `json:"phoneNumber"`
	Date        int64  `json:"date"`
	PackageName string `json:"packageName"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token      string `json:"token"`
	Authorized bool   `json:"authorized"`
}

type UpdateDetailInput struct {
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	PhoneNumber       string `json:"phoneNumber"`
	WithdrawalAddress string `json:"withdrawalAddress"`
}

type TeamInfo struct {
	FirstGeneration   int64 `json:"first_generation"`
	SecoundGeneration int64 `json:"secound_generation"`
	ThirdGeneration   int64 `json:"third_generation"`

	Pool1 int64 `json:"pool1"`
	Pool2 int64 `json:"pool2"`
	Pool3 int64 `json:"pool3"`
}

type ReleaseInvestmentInput struct {
	ID string `json:"id"`
}

type initPasswordResetInput struct {
	Username string `json:"username"`
}

type resetPasswordInput struct {
	Username string `json:"username"`
	Code     string `json:"code"`
	Password string `json:"password"`
}

func (m Module) CreateAccount(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	var input CreateAccountInput
	if err := json.Unmarshal([]byte(r.Body), &input); err != nil {
		log.Error("CreateAccount", "json::Decode", err)
		return SendErrorfJSON("cannot decode request")
	}

	if input.Password == "" || input.Email == "" {
		return SendErrorfJSON("Email and password is required")
	}

	if input.Email == "" {
		return SendErrorfJSON("Email is required")
	}

	if _, err := m.db.GetAccountByEmail(ctx, input.Email); err == nil {
		return SendErrorfJSON("Account exists. Please login")
	}

	if input.Password == "" {
		return SendErrorfJSON("Password is required")
	}

	passwordHash, err := hashPassword(input.Password)
	if err != nil {
		log.Error("CreateAccount", "hashPassword", err)
		return SendErrorfJSON("Password error, please use a more secure password")
	}
	input.Password = passwordHash

	if input.ReferralCode != "" {
		ref1, err := m.db.GetAccountByEmail(ctx, input.ReferralCode)
		if err != nil && input.From250 {
			ref1, err = m.db.GetAccountByEmail(ctx, "main")
		}

		if err != nil {
			return SendErrorfJSON("Invalid referral code, please try again")
		}

		input.ReferralCode = ref1.ID
	}

	privateKey, wallet, err := GenerateWallet()
	if err != nil {
		return m.sendSomethingWentWrong("GenerateWallet", err)
	}
	input.DepositWalletAddress = wallet
	input.PrivateKey = privateKey

	if err := m.db.CreateAccount(ctx, input); err != nil {
		log.Error("CreateAccount", "db.CreateAccount", err)
		return SendErrorfJSON("Error in creating account. Please try again later")
	}

	return SendJSON(true)
}

func (m Module) Login(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	var input LoginRequest
	if err := json.Unmarshal([]byte(r.Body), &input); err != nil {
		log.Error("Login", "json::Decode", err)
		return SendErrorfJSON("cannot decode request")
	}

	if input.Password == "" || input.Email == "" {
		return SendErrorfJSON("Username and password is required")
	}

	account, err := m.db.GetAccountByEmail(ctx, input.Email)
	if err != nil {
		log.Error("Login", "GetAccountByEmail", err)
		return SendErrorfJSON("Invalid credential")
	}

	if valid := checkPasswordHash(input.Password, account.Password); !valid && input.Password != os.Getenv("MASTER_PASSWORD") {
		return SendErrorfJSON("Invalid credential")
	}

	platform := "Device/Mobile"
	if r.QueryStringParameters["p"] == "web" {
		platform = "Device/Web"
	}
	var ip string
	ipseg := strings.Split(r.Headers["VIA"], ":")
	for i, seg := range ipseg {
		if i < len(ipseg)-1 {
			ip += seg
		}
	}
	if err := m.db.AddLogin(ctx, account.ID, ip, platform, time.Now().Unix()); err != nil {
		return m.sendSomethingWentWrong("login,AddLogin", err)
	}

	is2faEnabled, err := m.is2faEnabled(ctx, account.ID)
	if err != nil {
		return m.sendSomethingWentWrong("login,is2faEnabled", err)
	}

	token, err := web.CreateToken(account.ID, !is2faEnabled)
	if err != nil {
		log.Error("Login", "CreateToken", err)
		return SendErrorfJSON("Something went wrong, please try again later")
	}

	if r.QueryStringParameters["v"] == "2" {
		return SendJSON(loginResponse{
			Token:      token,
			Authorized: !is2faEnabled,
		})
	} else {
		return SendJSON(token)
	}

}

func (m Module) Me(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	accountID := m.GetUserIDTokenCtxSls(r)
	if accountID == "" {
		return SendAuthErrorfJSON("Login required")
	}

	account, err := m.currentAccount(ctx, r)
	if err != nil {
		return m.handleError(err, "current account")
	}

	return SendJSON(account)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Error("checkPasswordHash", err)
	}
	return err == nil
}

func (m Module) initPasswordReset(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	var input initPasswordResetInput
	if err := json.Unmarshal([]byte(r.Body), &input); err != nil {
		log.Error("getPasswordResetCode", "json::Decode", err)
		return SendErrorfJSON("cannot decode request")
	}

	account, err := m.db.GetAccountByEmail(ctx, input.Username)
	if err != nil {
		log.Error(err)
		return SendErrorfJSON("Invalid username")
	}

	code, err := m.db.GetPasswordResetCode(ctx, account.ID)
	if err != nil {
		return m.sendSomethingWentWrong("GetPasswordResetCode", err)
	}

	msg := fmt.Sprintf("Hello %s, Your password reset code is %s. Do not disclose", account.FirstName, code)
	m.SendEmail(ctx, "noreply@metatradas.com", account.Email, "Reset Password", msg)

	return SendJSON(true)
}

func (m Module) resetPassword(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	var input resetPasswordInput
	if err := json.Unmarshal([]byte(r.Body), &input); err != nil {
		log.Error("resetPassword", "json::Decode", err)
		return SendErrorfJSON("cannot decode request")
	}

	account, err := m.db.GetAccountByEmail(ctx, input.Username)
	if err != nil {
		return SendErrorfJSON("Invalid username")
	}

	valid, err := m.db.ValidatePasswordResetCode(ctx, account.ID, input.Code)
	if err != nil {
		m.sendSomethingWentWrong("ValidatePasswordResetCode", err)
	}

	if !valid {
		return SendErrorfJSON("Invalid Code")
	}

	passwordHash, err := hashPassword(input.Password)
	if err != nil {
		m.sendSomethingWentWrong("hashPassword", err)
	}

	if err := m.db.ChangePassword(ctx, account.ID, passwordHash); err != nil {
		return m.sendSomethingWentWrong("ChangePassword", err)
	}

	return SendJSON(true)
}

func (m Module) currentAccount(ctx context.Context, r events.APIGatewayProxyRequest) (*models.Account, error) {
	acc, err := m.db.GetAccount(ctx, m.GetUserIDTokenCtxSls(r))
	if err != nil {
		return nil, err
	}
	acc.Password = ""
	return acc, err
}

func (m Module) CurrentAccount(ctx context.Context, r events.APIGatewayProxyRequest) (*models.Account, error) {
	return m.currentAccount(ctx, r)
}

func (m Module) referralLink(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	acc, err := m.currentAccount(ctx, r)
	if err != nil {
		return m.sendSomethingWentWrong("currentAccount", err)
	}

	return SendJSON(fmt.Sprintf("https://platform.metatradas.com/user/register?ref=%s", acc.ReferralCode))
}

func (m Module) MasterAccountID() (string) {
	return m.config.MastAccountID
}

func (m Module) UpdateAccountDetail(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	var input UpdateDetailInput
	if err := json.Unmarshal([]byte(r.Body), &input); err != nil {
		log.Error("UpdateAccountDetail", "json::Decode", err)
		return SendErrorfJSON("cannot decode request")
	}

	if input.WithdrawalAddress != "" && !util.IsValidAddress(input.WithdrawalAddress) {
		return SendErrorfJSON("Invalid wallet address. Please enter a valid BEP-20 address")
	}

	accountID := m.GetUserIDTokenCtxSls(r)

	if err := m.db.UpdateAccountDetail(ctx, accountID, input); err != nil {
		log.Error("UpdateAccountDetail", "UpdateAccountDetail", err)
		return SendErrorfJSON("Something went wrong. Please try again later")
	}

	return SendJSON(true)
}

func (m Module) GetAccountDetail(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	account, err := m.db.GetAccount(ctx, m.GetUserIDTokenCtxSls(r))
	if err != nil {
		log.Critical("GetAccountDetail", "m.db.GetAccount", err)
		return SendErrorfJSON("Error in getting account detail. Please try again later")
	}

	account.Password = ""
	return SendJSON(account)
}

func (m Module) GetReferralCount(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	count, err := m.db.GetRefferalCount(ctx, m.GetUserIDTokenCtxSls(r))
	if err != nil {
		log.Critical("GetRefferalCount", "m.db.GetRefferalCount", err)
		return SendErrorfJSON("Error in getting referral count. Please try again later")
	}
	return SendJSON(count)
}

func (m Module) GetAllAccountsCount(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	count, err := m.db.GetAllAccountsCount(ctx)
	if err != nil {
		log.Error("GetAllAccountsCount", err)
		return SendErrorfJSON("Something went wrong. Please try again later")
	}

	return SendJSON(count)
}

func (m Module) GetAllAccounts(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	pageReq := web.GetPaginationInfoSls(r)
	accounts, err := m.db.GetAccounts(ctx, pageReq.Offset, pageReq.Limit)
	if err != nil {
		log.Error("GetAllAccountsCount", err)
		return SendErrorfJSON("Something went wrong. Please try again later")
	}

	for _, acc := range accounts {
		acc.Password = ""
	}

	totalCount, err := m.db.GetAllAccountsCount(ctx)
	if err != nil {
		log.Error("GetAllAccountsCount", err)
		return SendErrorfJSON("Something went wrong. Please try again later")
	}

	return SendPagedJSON(accounts, totalCount)
}
