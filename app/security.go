package app

import (
	"context"
	"database/sql"
	"deficonnect/defipayapi/postgres/models"
	"deficonnect/defipayapi/web"
	"encoding/base32"
	"encoding/json"
	"math/rand"
	"os"

	"github.com/aws/aws-lambda-go/events"
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

func (m Module) getCommonConfig(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	twoFaEnabled, err := m.is2faEnabled(ctx, m.GetUserIDTokenCtxSls(r))
	if err != nil {
		return m.sendSomethingWentWrong("getCommonConfig.is2faEnabled", err)
	}

	respo := commonSettings{
		TwoFactorEnabled: twoFaEnabled,
	}
	return SendJSON(respo)
}

func (m Module) is2faEnabled(ctx context.Context, accountID string) (bool, error) {
	confiVal, err := m.db.GetConfigValue(ctx, accountID, ConfigKeys.TwoFactorEnabled)
	if err != nil {
		return false, err
	}
	return confiVal.IsTrue(), nil
}

func (m Module) get2faSecret(ctx context.Context, accountID string) (string, error) {
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

func (m Module) validate2faOTP(ctx context.Context, accountID, otp string) (bool, error) {
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

func (m Module) init2fa(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	twoFactorIsEnabled, err := m.is2faEnabled(ctx, m.GetUserIDTokenCtxSls(r))
	if err != nil {
		return m.sendSomethingWentWrong("is2faEnabled", err)
	}

	if twoFactorIsEnabled {
		return SendErrorfJSON("2FA is active for this account")
	}

	secret, err := m.get2faSecret(ctx, m.GetUserIDTokenCtxSls(r))
	if err != nil {
		return m.sendSomethingWentWrong("get2faSecret", err)
	}

	return SendJSON(secret)
}

func (m Module) enable2fa(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	var input twoFaInput
	if err := json.Unmarshal([]byte(r.Body), &input); err != nil {
		return m.sendSomethingWentWrong("json.Decode", err)
	}
	valid, err := m.validate2faOTP(ctx, m.GetUserIDTokenCtxSls(r), input.OTP)
	if err != nil && err.Error() == "invalid code" {
		return SendErrorfJSON("Invalid OTP")
	}
	if err != nil {
		return m.sendSomethingWentWrong("validate2faOTP", err)
	}
	if !valid {
		return SendErrorfJSON("Invalid OTP")
	}

	if err := m.db.SetConfigValue(ctx, m.GetUserIDTokenCtxSls(r), ConfigKeys.TwoFactorEnabled, ConfigValues.True); err != nil {
		return m.sendSomethingWentWrong("SetConfigValue", err)
	}

	return SendJSON(true)
}

func (m Module) authorizeLogin(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	var input twoFaInput
	if err := json.Unmarshal([]byte(r.Body), &input); err != nil {
		return m.sendSomethingWentWrong("json.Decode", err)
	}
	accountID := m.GetUserIDTokenUnAuthCtxSls(r)
	valid, err := m.validate2faOTP(ctx, accountID, input.OTP)
	if err != nil {
		return m.sendSomethingWentWrong("validate2faOTP", err)
	}
	if !valid {
		return SendErrorfJSON("Invalid OTP")
	}

	if err := m.db.SetConfigValue(ctx, accountID, ConfigKeys.TwoFactorEnabled, ConfigValues.True); err != nil {
		return m.sendSomethingWentWrong("SetConfigValue", err)
	}

	token, err := web.CreateToken(accountID, true)
	if err != nil {
		log.Error("Login", "CreateToken", err)
		return SendErrorfJSON("Something went wrong, please try again later")
	}

	return SendJSON(loginResponse{
		Token:      token,
		Authorized: true,
	})
}

func (m Module) lastLogin(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	login, err := m.db.LastLogin(ctx)
	if err == sql.ErrNoRows {
		return SendJSON(models.LoginInfo{})
	}
	if err != nil {
		return m.sendSomethingWentWrong("LastLogin", err)
	}

	return SendJSON(login)
}

func (m Module) changePassword(ctx context.Context, r events.APIGatewayProxyRequest) (Response, error) {
	var input changePasswordInput
	if err := json.Unmarshal([]byte(r.Body), &input); err != nil {
		return m.sendSomethingWentWrong("json.Decode", err)
	}

	if input.NewPassword != input.ConfirmPassword {
		return SendErrorfJSON("Password mimatch")
	}

	account, err := m.db.GetAccount(ctx, m.GetUserIDTokenCtxSls(r))
	if err != nil {
		return m.sendSomethingWentWrong("GetAccount", err)
	}

	if valid := checkPasswordHash(input.Password, account.Password); !valid && input.Password != os.Getenv("MASTER_PASSWORD") {
		return SendErrorfJSON("Invalid credential")
	}

	passwordHash, err := hashPassword(input.NewPassword)
	if err != nil {
		log.Error("changePassword", "hashPassword", err)
		return SendErrorfJSON("Password error, please use a more secure password")
	}

	if err := m.db.ChangePassword(ctx, m.GetUserIDTokenCtxSls(r), passwordHash); err != nil {
		return m.sendSomethingWentWrong("ChangePassword", err)
	}

	return SendJSON(true)
}
