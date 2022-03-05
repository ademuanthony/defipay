package app

import (
	"encoding/json"
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
		web.RenderErrorfJSON(w, "cannot decode request")
		return
	}

	if input.Email == "" {
		web.RenderErrorfJSON(w, "Email is required")
		return
	}

	if input.Username == "" {
		web.RenderErrorfJSON(w, "Username is required")
		return
	}

	if _, err := m.db.GetAccountByUsername(r.Context(), input.Username); err == nil {
		web.RenderErrorfJSON(w, "Username is not available")
		return
	}

	if input.Password == "" {
		web.RenderErrorfJSON(w, "Password is required")
		return
	}

	passwordHash, err := hashPassword(input.Password)
	if err != nil {
		log.Error("CreateAccount", "hashPassword", err)
		web.RenderErrorfJSON(w, "Password error, please use a more secure password")
		return
	}
	input.Password = passwordHash

	wallet, privateKey, err := GenerateWallet()
	if err != nil {
		log.Critical("CreateAccount", "GenerateWallet", err)
		web.RenderErrorfJSON(w, "Something went wrong. Please try again later")
		return
	}

	input.WalletAddress = wallet
	input.PrivateKey = privateKey

	if err := m.db.CreateAccount(r.Context(), input); err != nil {
		log.Error("CreeateAccount", "db.CreateAccount", err)
		web.RenderErrorfJSON(w, "Error in creating account. Please try again later")
		return
	}

	web.RenderJSON(w, true)
}

func (m module) Login(w http.ResponseWriter, r *http.Request) {
	var input LoginRequst
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Error("Login", "json::Decode", err)
		web.RenderErrorfJSON(w, "cannot decode request")
		return
	}

	account, err := m.db.GetAccountByUsername(r.Context(), input.Username)
	if err != nil {
		log.Error("Login", "GetAccountByUsername", err)
		web.RenderErrorfJSON(w, "Invalid credential")
		return
	}

	if valid := checkPasswordHash(input.Password, account.Password); !valid {
		web.RenderErrorfJSON(w, "Invalid credential")
		return
	}

	token, err := web.CreateToken(account.ID)
	if err != nil {
		log.Error("Login", "CreateToken", err)
		web.RenderErrorfJSON(w, "Something went wrong, please try again later")
		return
	}

	web.RenderJSON(w, token)
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
		web.RenderErrorfJSON(w, "cannot decode request")
		return
	}

	accountID := m.server.GetUserIDTokenCtx(r)

	if err := m.db.UpdateAccountDetail(r.Context(), accountID, input); err != nil {
		log.Error("UpdateAccountDetail", "UpdateAccountDetail", err)
		web.RenderErrorfJSON(w, "Something went wrong. Please try again later")
		return
	}

	web.RenderJSON(w, true)
}

func (m module) GetAccountDetail(w http.ResponseWriter, r *http.Request) {
	account, err := m.db.GetAccount(r.Context(), m.server.GetUserIDTokenCtx(r))
	if err != nil {
		log.Critical("GetAccountDetail", "m.db.GetAccount", err)
		web.RenderErrorfJSON(w, "Error in getting account detail. Please try again later")
		return
	}

	account.Password = ""
	web.RenderJSON(w, account)
}
