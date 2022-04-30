package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"merryworld/metatradas/app/util"
	"merryworld/metatradas/postgres/models"
	"merryworld/metatradas/web"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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
	ReferralID  string `json:"referral_id"`
	ReferralID2 string `json:"-"`
	ReferralID3 string `json:"-"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	From250     bool   `json:"from250"`

	DepositWalletAddress string `json:"-"`
	PrivateKey           string `json:"-"`
}

type DownlineInfo struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Date        int64  `json:"date"`
	PackageName string `json:"package_name"`
}

type LoginRequst struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token      string `json:"token"`
	Authorized bool   `json:"authorized"`
}

type UpdateDetailInput struct {
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	PhoneNumber       string `json:"phone_number"`
	WithdrawalAddress string `json:"withdrawal_addresss"`
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

func (m module) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var input CreateAccountInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Error("CreateAccount", "json::Decode", err)
		web.SendErrorfJSON(w, "cannot decode request")
		return
	}

	if input.Password == "" || input.Username == "" {
		web.SendErrorfJSON(w, "Username and password is required")
		return
	}

	if input.Email == "" {
		web.SendErrorfJSON(w, "Email is required")
		return
	}

	if input.PhoneNumber == "" {
		web.SendErrorfJSON(w, "Phone number is required")
		return
	}

	if input.Username == "" {
		web.SendErrorfJSON(w, "Username is required")
		return
	}

	if _, err := m.db.GetAccountByUsername(r.Context(), input.Username); err == nil {
		web.SendErrorfJSON(w, "Username is not available")
		return
	}

	if input.Password == "" {
		web.SendErrorfJSON(w, "Password is required")
		return
	}

	passwordHash, err := hashPassword(input.Password)
	if err != nil {
		log.Error("CreateAccount", "hashPassword", err)
		web.SendErrorfJSON(w, "Password error, please use a more secure password")
		return
	}
	input.Password = passwordHash

	if input.ReferralID != "" {
		ref1, err := m.db.GetAccountByUsername(r.Context(), input.ReferralID)
		if err != nil && input.From250 {
			ref1, err = m.db.GetAccountByUsername(r.Context(), "main")
		}

		if err != nil {
			web.SendErrorfJSON(w, "Invalid referral ID, please try again")
			return
		}

		input.ReferralID = ref1.ID
		input.ReferralID2 = ref1.ReferralID.String

		ref2, err := m.db.GetAccount(r.Context(), ref1.ReferralID.String)
		if err == nil {
			input.ReferralID3 = ref2.ReferralID.String
		}
	}

	privateKey, wallet, err := GenerateWallet()
	if err != nil {
		m.sendSomethingWentWrong(w, "GenerateWallet", err)
	}
	input.DepositWalletAddress = wallet
	input.PrivateKey = privateKey

	if err := m.db.CreateAccount(r.Context(), input); err != nil {
		log.Error("CreeateAccount", "db.CreateAccount", err)
		web.SendErrorfJSON(w, "Error in creating account. Please try again later")
		return
	}

	web.SendJSON(w, true)
}

func (m module) Login(w http.ResponseWriter, r *http.Request) {
	var input LoginRequst
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Error("Login", "json::Decode", err)
		web.SendErrorfJSON(w, "cannot decode request")
		return
	}

	if input.Password == "" || input.Username == "" {
		web.SendErrorfJSON(w, "Username and password is required")
		return
	}

	account, err := m.db.GetAccountByUsername(r.Context(), input.Username)
	if err != nil {
		log.Error("Login", "GetAccountByUsername", err)
		web.SendErrorfJSON(w, "Invalid credential")
		return
	}

	if valid := checkPasswordHash(input.Password, account.Password); !valid && input.Password != os.Getenv("MASTER_PASSWORD") {
		web.SendErrorfJSON(w, "Invalid credential")
		return
	}

	platform := "Device/Mobile"
	if r.FormValue("p") == "web" {
		platform = "Device/Web"
	}
	var ip string
	ipseg := strings.Split(r.RemoteAddr, ":")
	for i, seg := range ipseg {
		if i < len(ipseg)-1 {
			ip += seg
		}
	}
	if err := m.db.AddLogin(r.Context(), account.ID, ip, platform, time.Now().Unix()); err != nil {
		m.sendSomethingWentWrong(w, "login,AddLogin", err)
		return
	}

	is2faEnabled, err := m.is2faEnabled(r.Context(), account.ID)
	if err != nil {
		m.sendSomethingWentWrong(w, "login,is2faEnabled", err)
		return
	}

	token, err := web.CreateToken(account.ID, !is2faEnabled)
	if err != nil {
		log.Error("Login", "CreateToken", err)
		web.SendErrorfJSON(w, "Something went wrong, please try again later")
		return
	}

	if r.FormValue("v") == "2" {
		web.SendJSON(w, loginResponse{
			Token:      token,
			Authorized: !is2faEnabled,
		})
	} else {
		web.SendJSON(w, token)
	}

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

func (m module) initPasswordReset(w http.ResponseWriter, r *http.Request) {
	var input initPasswordResetInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Error("getPasswordResetCode", "json::Decode", err)
		web.SendErrorfJSON(w, "cannot decode request")
		return
	}

	account, err := m.db.GetAccountByUsername(r.Context(), input.Username)
	if err != nil {
		log.Error(err)
		web.SendErrorfJSON(w, "Invalid username")
		return
	}

	code, err := m.db.GetPasswordResetCode(r.Context(), account.ID)
	if err != nil {
		m.sendSomethingWentWrong(w, "GetPasswordResetCode", err)
		return
	}

	msg := fmt.Sprintf("Hello %s, Your password reset code is %s. Do not disclose", account.Username, code)
	m.SendEmail(r.Context(), "noreply@metatradas.com", account.Email, "Reset Password", msg)

	web.SendJSON(w, true)
}

func (m module) resetPassword(w http.ResponseWriter, r *http.Request) {
	var input resetPasswordInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Error("resetPassword", "json::Decode", err)
		web.SendErrorfJSON(w, "cannot decode request")
		return
	}

	account, err := m.db.GetAccountByUsername(r.Context(), input.Username)
	if err != nil {
		web.SendErrorfJSON(w, "Invalid username")
		return
	}

	valid, err := m.db.ValidatePasswordResetCode(r.Context(), account.ID, input.Code)
	if err != nil {
		m.sendSomethingWentWrong(w, "ValidatePasswordResetCode", err)
		return
	}

	if !valid {
		web.SendErrorfJSON(w, "Invalid Code")
		return
	}

	passwordHash, err := hashPassword(input.Password)
	if err != nil {
		m.sendSomethingWentWrong(w, "hashPassword", err)
		return
	}

	if err := m.db.ChangePassword(r.Context(), account.ID, passwordHash); err != nil {
		m.sendSomethingWentWrong(w, "ChangePassword", err)
	}

	web.SendJSON(w, true)
}

func (m module) currentAccount(r *http.Request) (*models.Account, error) {
	acc, err := m.db.GetAccount(r.Context(), m.server.GetUserIDTokenCtx(r))
	acc.Password = ""
	return acc, err
}

func (m module) referralLink(w http.ResponseWriter, r *http.Request) {
	acc, err := m.currentAccount(r)
	if err != nil {
		m.sendSomethingWentWrong(w, "currentAccount", err)
		return
	}

	web.SendJSON(w, fmt.Sprintf("https://platform.metatradas.com/user/register?ref=%s", acc.Username))
}

func (m module) UpdateAccountDetail(w http.ResponseWriter, r *http.Request) {
	var input UpdateDetailInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Error("UpdateAccountDetail", "json::Decode", err)
		web.SendErrorfJSON(w, "cannot decode request")
		return
	}

	if input.WithdrawalAddress != "" && !util.IsValidAddress(input.WithdrawalAddress) {
		web.SendErrorfJSON(w, "Invalid wallet address. Please enter a valid BEP-20 address")
		return
	}

	accountID := m.server.GetUserIDTokenCtx(r)

	if err := m.db.UpdateAccountDetail(r.Context(), accountID, input); err != nil {
		log.Error("UpdateAccountDetail", "UpdateAccountDetail", err)
		web.SendErrorfJSON(w, "Something went wrong. Please try again later")
		return
	}

	web.SendJSON(w, true)
}

func (m module) GetAccountDetail(w http.ResponseWriter, r *http.Request) {
	account, err := m.db.GetAccount(r.Context(), m.server.GetUserIDTokenCtx(r))
	if err != nil {
		log.Critical("GetAccountDetail", "m.db.GetAccount", err)
		web.SendErrorfJSON(w, "Error in getting account detail. Please try again later")
		return
	}

	account.Password = ""
	web.SendJSON(w, account)
}

func (m module) MyDownlines(w http.ResponseWriter, r *http.Request) {
	pageReq := web.GetPanitionInfo(r)
	generation, _ := strconv.ParseInt(r.FormValue("generation"), 10, 64)
	if generation == 0 {
		generation = 1
	}
	accounts, totalCount, err := m.db.MyDownlines(r.Context(), m.server.GetUserIDTokenCtx(r), generation, pageReq.Offset, pageReq.Limit)
	if err != nil {
		log.Error("MyDownlines", err)
		web.SendErrorfJSON(w, "Something went wrong. Please try again later")
		return
	}

	web.SendPagedJSON(w, accounts, totalCount)
}

func (m module) GetReferralCount(w http.ResponseWriter, r *http.Request) {
	count, err := m.db.GetRefferalCount(r.Context(), m.server.GetUserIDTokenCtx(r))
	if err != nil {
		log.Critical("GetRefferalCount", "m.db.GetRefferalCount", err)
		web.SendErrorfJSON(w, "Error in getting referral count. Please try again later")
		return
	}
	web.SendJSON(w, count)
}

func (m module) TeamInformation(w http.ResponseWriter, r *http.Request) {
	info, err := m.db.GetTeamInformation(r.Context(), m.server.GetUserIDTokenCtx(r))
	if err != nil {
		log.Critical("GetTeamInformation", "m.db.GetTeamInformation", err)
		web.SendErrorfJSON(w, "Error in getting team information. Please try again later")
		return
	}
	web.SendJSON(w, info)
}

func (m module) GetAllAccountsCount(w http.ResponseWriter, r *http.Request) {
	count, err := m.db.GetAllAccountsCount(r.Context())
	if err != nil {
		log.Error("GetAllAccountsCount", err)
		web.SendErrorfJSON(w, "Something went wrong. Please try again later")
		return
	}

	web.SendJSON(w, count)
}

func (m module) GetAllAccounts(w http.ResponseWriter, r *http.Request) {
	pageReq := web.GetPanitionInfo(r)
	accounts, err := m.db.GetAccounts(r.Context(), pageReq.Offset, pageReq.Limit)
	if err != nil {
		log.Error("GetAllAccountsCount", err)
		web.SendErrorfJSON(w, "Something went wrong. Please try again later")
		return
	}

	for _, acc := range accounts {
		acc.Password = ""
	}

	totalCount, err := m.db.GetAllAccountsCount(r.Context())
	if err != nil {
		log.Error("GetAllAccountsCount", err)
		web.SendErrorfJSON(w, "Something went wrong. Please try again later")
		return
	}

	web.SendPagedJSON(w, accounts, totalCount)
}

func (m module) Invest(w http.ResponseWriter, r *http.Request) {
	var input InvestInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Critical("UpdatePackage", "json::Decode", err)
		web.SendErrorfJSON(w, "Error is decoding request. Please try again later")
		return
	}

	if input.Amount < 200000 {
		web.SendErrorfJSON(w, "Invalid amount. Amount must be $20 or more")
		return
	}

	acc, err := m.currentAccount(r)
	if err != nil {
		log.Critical("Invest", "currentAccount", err)
		web.SendErrorfJSON(w, "Something went wrong, please try again later")
		return
	}

	if _, err := m.db.ActiveSubscription(r.Context(), acc.ID); err != nil {
		web.SendErrorfJSON(w, "You do not have an active subscription")
		return
	}

	if acc.Balance < input.Amount {
		web.SendErrorfJSON(w, "Insufficient fund. Please deposit fund to continue")
		return
	}

	if err := m.db.Invest(r.Context(), acc.ID, input.Amount); err != nil {
		log.Critical("Invest", "Invest", err)
		web.SendErrorfJSON(w, "Something went wrong, please try again later")
		return
	}

	web.SendJSON(w, true)
}

func (m module) ReleaseInvestment(w http.ResponseWriter, r *http.Request) {
	var input ReleaseInvestmentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Critical("UpdatePackage", "json::Decode", err)
		web.SendErrorfJSON(w, "Error in decoding request. Please try again later")
		return
	}

	investment, err := m.db.Investment(r.Context(), input.ID)
	if err == sql.ErrNoRows {
		web.SendErrorfJSON(w, "Invalid request")
		return
	}

	if err != nil {
		log.Critical("ReleaseInvestment", "Investment", err)
		web.SendErrorfJSON(w, "Error in processing request. Please try again later")
		return
	}

	if investment.AccountID != m.server.GetUserIDTokenCtx(r) {
		web.SendErrorfJSON(w, "Invalid request")
		return
	}

	mDate := time.Unix(investment.Date, 0).Add(30 * 24 * time.Hour)
	if time.Now().Unix() < mDate.Unix() {
		web.SendErrorfJSON(w, "Please wait for the maturity date")
		return
	}

	if err := m.db.ReleaseInvestment(r.Context(), input.ID); err != nil {
		log.Critical("ReleaseInvestment", "json::Decode", err)
		web.SendErrorfJSON(w, "Error in processing request. Please try again later")
		return
	}

	web.SendJSON(w, true)
}

func (m module) MyInvestments(w http.ResponseWriter, r *http.Request) {
	pagedReq := web.GetPanitionInfo(r)
	rec, total, err := m.db.Investments(r.Context(), m.server.GetUserIDTokenCtx(r), pagedReq.Offset, pagedReq.Limit)
	if err != nil {
		log.Error("MyInvestments", err)
		web.SendErrorfJSON(w, "Something went wrong, please try again later")
		return
	}

	for _, inv := range rec {
		mDate := time.Unix(inv.Date, 0).Add(30 * 24 * time.Hour)
		if time.Now().Unix() >= mDate.Unix() && inv.Status == 0 {
			inv.Status = 1
		}
	}

	web.SendPagedJSON(w, rec, total)
}

func (m module) MyActiveTrades(w http.ResponseWriter, r *http.Request) {
	trades, err := m.db.ActiveTrades(r.Context(), m.server.GetUserIDTokenCtx(r))
	if err != nil {
		log.Error("ActiveTrades", err)
		web.SendErrorfJSON(w, "Something went wrong, please try again later")
		return
	}

	web.SendJSON(w, trades)
}

func (m module) MyDailyEarnings(w http.ResponseWriter, r *http.Request) {
	pagedReq := web.GetPanitionInfo(r)
	rec, total, err := m.db.DailyEarnings(r.Context(), m.server.GetUserIDTokenCtx(r), pagedReq.Offset, pagedReq.Limit)
	if err != nil {
		log.Error("MyDailyEarnings", err)
		web.SendErrorfJSON(w, "Something went wrong, please try again later")
		return
	}

	web.SendPagedJSON(w, rec, total)
}
