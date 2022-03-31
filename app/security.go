package app

import (
	"context"
	"database/sql"
	"encoding/base32"
	"encoding/json"
	"math/rand"
	"merryworld/metatradas/web"
	"net/http"

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

func (m module) is2faEnabled(ctx context.Context, accountID string) (bool, error) {
	confiVal, err := m.db.GetConfigValue(ctx, accountID, ConfigKeys.TwoFactorEnabled)
	if err != nil {
		return false, err
	}
	return confiVal.IsTrue(), nil
}

func (m module) get2faSecret(ctx context.Context, accountID string) (string, error) {
	confiVal, err := m.db.GetConfigValue(ctx, accountID, ConfigKeys.TwoFactorEnabled)
	if err == nil {
		return string(confiVal), nil
	}
	if err != sql.ErrNoRows {
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
	accountID := m.server.GetUserIDTokenCtx(r)
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
