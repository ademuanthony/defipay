package app

import (
	"deficonnect/defipayapi/web"
	"net/http"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/twharmon/govalid"
)

type Module struct {
	server        *web.Server
	db            store
	bscClient     *ethclient.Client
	polygonClient *ethclient.Client
	config        BlockchainConfig

	currencyProcessors map[string]map[Network]CurrencyProcessor

	MgDomain string
	MgKey    string
}

const (
	TxTypeCredit = "credit"
	TxTypeDebit  = "debit"
)

var v = govalid.New()

func Start(db store, bscClient *ethclient.Client,
	polygonClient *ethclient.Client, currencyProcessors map[string]map[Network]CurrencyProcessor,
	config BlockchainConfig,
	mgDomain, mgKey string) (*Module, error) {
	log.Info("starting...")

	app := Module{
		db:                 db,
		bscClient:          bscClient,
		polygonClient:      polygonClient,
		currencyProcessors: currencyProcessors,
		config:             config,
		MgDomain:           mgDomain,
		MgKey:              mgKey,
	}

	app.buildRoute()

	return &app, nil
}

func (m Module) buildRoute() {
	m.server.AddRoute("/", web.GET, welcome)

	// AUTH
	m.server.AddRoute("/api/auth/register", web.POST, m.CreateAccount)
	m.server.AddRoute("/api/auth/login", web.POST, m.Login)
	m.server.AddRoute("/api/auth/2fa", web.POST, m.authorizeLogin, m.server.ValidBearerToken)
	m.server.AddRoute("/api/auth/init-password-reset", web.POST, m.initPasswordReset)
	m.server.AddRoute("/api/auth/reset-password", web.POST, m.resetPassword)

	//ACCOUNT
	m.server.AddRoute("/api/account/update", web.POST, m.UpdateAccountDetail, m.server.RequireLogin, m.server.NoReentry)
	m.server.AddRoute("/api/account/me", web.GET, m.GetAccountDetail, m.server.RequireLogin)
	m.server.AddRoute("/api/account/referral-link", web.GET, m.referralLink, m.server.RequireLogin)
	m.server.AddRoute("/api/account/referral-count", web.GET, m.GetReferralCount, m.server.RequireLogin)
	//m.server.AddRoute("/api/account/downlines", web.GET, m.MyDownlines, m.server.RequireLogin)

	// ACCOUNTS
	m.server.AddRoute("/api/accounts/count", web.GET, m.GetAllAccountsCount, m.server.RequireLogin)
	m.server.AddRoute("/api/accounts/list", web.GET, m.GetAllAccounts, m.server.RequireLogin)

	m.server.AddRoute("/api/notifications/totalPending", web.GET, m.getUnReadNotificationCount, m.server.RequireLogin)
	m.server.AddRoute("/api/notifications/pending", web.GET, m.getNewNotifications, m.server.RequireLogin)
	m.server.AddRoute("/api/notifications/getAll", web.GET, m.getNotifications, m.server.RequireLogin)
	m.server.AddRoute("/api/notifications/get", web.GET, m.getNotification, m.server.RequireLogin)

	m.server.AddRoute("/api/config/common-settings", web.GET, m.getCommonConfig, m.server.RequireLogin)

	m.server.AddRoute("/api/security/init2fa", web.POST, m.init2fa, m.server.RequireLogin, m.server.NoReentry)
	m.server.AddRoute("/api/security/enable2fa", web.POST, m.enable2fa, m.server.RequireLogin, m.server.NoReentry)
	m.server.AddRoute("/api/security/last-login", web.GET, m.lastLogin, m.server.RequireLogin)
	m.server.AddRoute("/api/security/change-password", web.POST, m.changePassword, m.server.RequireLogin, m.server.NoReentry)

	// TRANSACTION
	m.server.AddRoute("/api/transaction/getByID", web.GET, m.getTransaction)
	m.server.AddRoute("/api/transaction/getAll", web.GET, m.getTransactions)
	m.server.AddRoute("/api/transaction/initFundTransfer", web.POST, m.createFundTransferTransaction, m.server.RequireLogin, m.server.NoReentry)
	m.server.AddRoute("/api/transaction/initTopUp", web.POST, m.createTupUpTransaction, m.server.RequireLogin, m.server.NoReentry)
	m.server.AddRoute("/api/transaction/updateCurrency", web.POST, m.updateTransactionCurrency, m.server.RequireLogin, m.server.NoReentry)
	m.server.AddRoute("/api/transaction/checkStatus", web.POST, m.checkTransactionStatus, m.server.RequireLogin, m.server.NoReentry)

	m.server.AddRoute("/config/supported-currencies", web.GET, m.supportedCurrencies)
}

func welcome(w http.ResponseWriter, r *http.Request) {
	web.SendJSON(w, "welcome to defipay api. download the app from app store to start sending money with ease")
}
