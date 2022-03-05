package app

import (
	"encoding/json"
	"merryworld/metatradas/postgres/models"
	"merryworld/metatradas/web"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type CreateAccountInput struct {
	ReferralID string `json:"referralId"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Password   string `json:"password"`

	WalletAddress string `json:"_"`
	PrivateKey    string `json:"_"`
}

type LoginRequst struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateDetailInput struct {
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	PhoneNumber       string `json:"phone_number"`
	WithdrawalAddress string `json:"withdrawal_addresss"`
}

func (m module) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var input CreateAccountInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Error("CreateAccount", "json::Decode", err)
		web.SendErrorfJSON(w, "cannot decode request")
		return
	}

	if input.Email == "" {
		web.SendErrorfJSON(w, "Email is required")
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

	privateKey, wallet, err := GenerateWallet()
	if err != nil {
		log.Critical("CreateAccount", "GenerateWallet", err)
		web.SendErrorfJSON(w, "Something went wrong. Please try again later")
		return
	}

	if input.ReferralID != "" {
		_, err := m.db.GetAccountByUsername(r.Context(), input.ReferralID)
		if err != nil {
			web.SendErrorfJSON(w, "Invalid referral ID, please try again")
			return
		}
	}

	input.WalletAddress = wallet
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

	account, err := m.db.GetAccountByUsername(r.Context(), input.Username)
	if err != nil {
		log.Error("Login", "GetAccountByUsername", err)
		web.SendErrorfJSON(w, "Invalid credential")
		return
	}

	if valid := checkPasswordHash(input.Password, account.Password); !valid {
		web.SendErrorfJSON(w, "Invalid credential")
		return
	}

	token, err := web.CreateToken(account.ID)
	if err != nil {
		log.Error("Login", "CreateToken", err)
		web.SendErrorfJSON(w, "Something went wrong, please try again later")
		return
	}

	web.SendJSON(w, token)
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

func (m module) UpdateAccountDetail(w http.ResponseWriter, r *http.Request) {
	var input UpdateDetailInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Error("UpdateAccountDetail", "json::Decode", err)
		web.SendErrorfJSON(w, "cannot decode request")
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

func (m module) currentAccount(r *http.Request) (*models.Account, error) {
	acc, err := m.db.GetAccount(r.Context(), m.server.GetUserIDTokenCtx(r))
	acc.Password = ""
	return acc, err
}