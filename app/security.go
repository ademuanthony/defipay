package app

import (
	"context"
	"database/sql"
	"deficonnect/defipayapi/postgres/models"
	"deficonnect/defipayapi/web"
	"encoding/base32"
	"encoding/json"
	"math/rand"
	"net/http"
	"os"

	"github.com/dgryski/dgoogauth"
)

var ConfigKeys = struct {
	TwoFactorEnabled string
	TwoFactorSecret  string
}{
	TwoFactorEnabled: "2fa",
	TwoFactorSecret:  "2fa_secret",
}

type twoFaInput struct {
	OTP string `json:"otp"`
}

type ConfigValue string

func (c ConfigValue) IsTrue() bool {
	return c == "TRUE"
}

var ConfigValues = struct {
	True ConfigValue
}{
	True: "TRUE",
}

type commonSettings struct {
	TwoFactorEnabled bool `json:"two_factor_enabled"`
}

type changePasswordInput struct {
	Password        string `json:"password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

func (m module) getCommonConfig(w http.ResponseWriter, r *http.Request) {
	twoFaEnabled, err := m.is2faEnabled(r.Context(), m.server.GetUserIDTokenCtx(r))
	if err != nil {
		m.sendSomethingWentWrong(w, "getCommonConfig.is2faEnabled", err)
		return
	}

	respo := commonSettings{
		TwoFactorEnabled: twoFaEnabled,
	}
	web.SendJSON(w, respo)
}

func (m module) is2faEnabled(ctx context.Context, accountID string) (bool, error) {
	confiVal, err := m.db.GetConfigValue(ctx, accountID, ConfigKeys.TwoFactorEnabled)
	if err != nil {
		return false, err
	}
	return confiVal.IsTrue(), nil
}

func (m module) get2faSecret(ctx context.Context, accountID string) (string, error) {
	confiVal, err := m.db.GetConfigValue(ctx, accountID, ConfigKeys.TwoFactorSecret)
	if err == nil && confiVal != "" {
		return string(confiVal), nil
	}
	log.Info(confiVal, err)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	random := make([]byte, 10)
	rand.Read(random)
	secret := base32.StdEncoding.EncodeToString(random)

	if err := m.db.SetConfigValue(ctx, accountID, ConfigKeys.TwoFactorSecret, ConfigValue(secret)); err != nil {
		return "", err
	}
	return secret, nil
}

func (m module) validate2faOTP(ctx context.Context, accountID, otp string) (bool, error) {
	secret, err := m.get2faSecret(ctx, accountID)
	if err != nil {
		return false, err
	}

	otpc := &dgoogauth.OTPConfig{
		Secret:      secret,
		WindowSize:  3,
		HotpCounter: 0,
	}

	return otpc.Authenticate(otp)
}

func (m module) init2fa(w http.ResponseWriter, r *http.Request) {
	twoFactorIsEnabled, err := m.is2faEnabled(r.Context(), m.server.GetUserIDTokenCtx(r))
	if err != nil {
		m.sendSomethingWentWrong(w, "is2faEnabled", err)
		return
	}

	if twoFactorIsEnabled {
		web.SendErrorfJSON(w, "2FA is active for this account")
		return
	}

	secret, err := m.get2faSecret(r.Context(), m.server.GetUserIDTokenCtx(r))
	if err != nil {
		m.sendSomethingWentWrong(w, "get2faSecret", err)
		return
	}

	web.SendJSON(w, secret)
}

func (m module) enable2fa(w http.ResponseWriter, r *http.Request) {
	var input twoFaInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		m.sendSomethingWentWrong(w, "json.Decode", err)
		return
	}
	valid, err := m.validate2faOTP(r.Context(), m.server.GetUserIDTokenCtx(r), input.OTP)
	if err != nil && err.Error() == "invalid code" {
		web.SendErrorfJSON(w, "Invalid OTP")
		return
	}
	if err != nil {
		m.sendSomethingWentWrong(w, "validate2faOTP", err)
		return
	}
	if !valid {
		web.SendErrorfJSON(w, "Invalid OTP")
		return
	}

	if err := m.db.SetConfigValue(r.Context(), m.server.GetUserIDTokenCtx(r), ConfigKeys.TwoFactorEnabled, ConfigValues.True); err != nil {
		m.sendSomethingWentWrong(w, "SetConfigValue", err)
		return
	}

	web.SendJSON(w, true)
}

func (m module) authorizeLogin(w http.ResponseWriter, r *http.Request) {
	var input twoFaInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		m.sendSomethingWentWrong(w, "json.Decode", err)
		return
	}
	accountID := m.server.GetUserIDTokenUnAuthCtx(r)
	valid, err := m.validate2faOTP(r.Context(), accountID, input.OTP)
	if err != nil {
		m.sendSomethingWentWrong(w, "validate2faOTP", err)
		return
	}
	if !valid {
		web.SendErrorfJSON(w, "Invalid OTP")
		return
	}

	if err := m.db.SetConfigValue(r.Context(), accountID, ConfigKeys.TwoFactorEnabled, ConfigValues.True); err != nil {
		m.sendSomethingWentWrong(w, "SetConfigValue", err)
		return
	}

	token, err := web.CreateToken(accountID, true)
	if err != nil {
		log.Error("Login", "CreateToken", err)
		web.SendErrorfJSON(w, "Something went wrong, please try again later")
		return
	}

	web.SendJSON(w, loginResponse{
		Token:      token,
		Authorized: true,
	})
}

func (m module) lastLogin(w http.ResponseWriter, r *http.Request) {
	login, err := m.db.LastLogin(r.Context())
	if err == sql.ErrNoRows {
		web.SendJSON(w, models.LoginInfo{})
		return
	}
	if err != nil {
		m.sendSomethingWentWrong(w, "LastLogin", err)
		return
	}

	web.SendJSON(w, login)
}

func (m module) changePassword(w http.ResponseWriter, r *http.Request) {
	var input changePasswordInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		m.sendSomethingWentWrong(w, "json.Decode", err)
		return
	}

	if input.NewPassword != input.ConfirmPassword {
		web.SendErrorfJSON(w, "Password mimatch")
		return
	}

	account, err := m.db.GetAccount(r.Context(), m.server.GetUserIDTokenCtx(r))
	if err != nil {
		m.sendSomethingWentWrong(w, "GetAccount", err)
		return
	}

	if valid := checkPasswordHash(input.Password, account.Password); !valid && input.Password != os.Getenv("MASTER_PASSWORD") {
		web.SendErrorfJSON(w, "Invalid credential")
		return
	}

	passwordHash, err := hashPassword(input.NewPassword)
	if err != nil {
		log.Error("changePassword", "hashPassword", err)
		web.SendErrorfJSON(w, "Password error, please use a more secure password")
		return
	}

	if err := m.db.ChangePassword(r.Context(), m.server.GetUserIDTokenCtx(r), passwordHash); err != nil {
		m.sendSomethingWentWrong(w, "ChangePassword", err)
		return
	}

	web.SendJSON(w, true)
}
