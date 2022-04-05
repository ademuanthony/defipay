package app

import (
	"context"
	"merryworld/metatradas/web"
	"net/http"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

type module struct {
	server *web.Server
	db     store
	client *ethclient.Client
	config BlockchainConfig

	MgDomain string
	MgKey    string
}

const (
	TxTypeCredit = "credit"
	TxTypeDebit  = "debit"
)

func Start(server *web.Server, db store, client *ethclient.Client, config BlockchainConfig,
	mgDomain, mgKey string) error {
	log.Info("starting...")

	app := module{
		server:   server,
		db:       db,
		client:   client,
		config:   config,
		MgDomain: mgDomain,
		MgKey:    mgKey,
	}

	app.buildRoute()

	if os.Getenv("RUN_BG_PROCESSES") == "1" {
		go app.runProcessor(context.Background())
		// go app.watchDeposit()
		go app.watchBNBDeposit()
		go app.processReferralPayouts()
	}

	return nil
}

func (m module) buildRoute() {
	m.server.AddRoute("/", web.GET, welcome)

	// AUTH
	m.server.AddRoute("/api/auth/register", web.POST, m.CreateAccount)
	m.server.AddRoute("/api/auth/login", web.POST, m.Login)
	m.server.AddRoute("/api/auth/2fa", web.POST, m.authorizeLogin, m.server.ValidBearerToken)

	//ACCOUNT
	m.server.AddRoute("/api/account/update", web.POST, m.UpdateAccountDetail, m.server.RequireLogin, m.server.NoReentry)
	m.server.AddRoute("/api/account/me", web.GET, m.GetAccountDetail, m.server.RequireLogin)
	m.server.AddRoute("/api/account/referral-count", web.GET, m.GetReferralCount, m.server.RequireLogin)
	m.server.AddRoute("/api/account/downlines", web.GET, m.MyDownlines, m.server.RequireLogin)
	m.server.AddRoute("/api/account/team-info", web.GET, m.TeamInformation, m.server.RequireLogin)
	m.server.AddRoute("/api/account/deposit-address", web.GET, m.GetDepositAddress, m.server.RequireLogin)
	m.server.AddRoute("/api/account/deposits", web.GET, m.DepositHistories, m.server.RequireLogin)
	m.server.AddRoute("/api/account/invest", web.POST, m.Invest, m.server.RequireLogin, m.server.NoReentry)
	m.server.AddRoute("/api/account/investments", web.GET, m.MyInvestments, m.server.RequireLogin)
	m.server.AddRoute("/api/account/release-investment", web.POST, m.ReleaseInvestment, m.server.RequireLogin, m.server.NoReentry)
	m.server.AddRoute("/api/account/daily-earnings", web.GET, m.MyDailyEarnings, m.server.RequireLogin)
	m.server.AddRoute("/api/account/active-trades", web.GET, m.MyActiveTrades, m.server.RequireLogin)

	// C250 sub/upgrade
	m.server.AddRoute("/api/c250/subscribe", web.POST, m.createSubscriptionC250, m.server.ValidAPIKey)
	m.server.AddRoute("/api/c250/upgrade", web.POST, m.upgradeSubscriptionC250, m.server.ValidAPIKey)
	m.server.AddRoute("/api/c250/active-package", web.GET, m.activePackageC250, m.server.ValidAPIKey)

	// ACCOUNTS
	m.server.AddRoute("/api/accounts/count", web.GET, m.GetAllAccountsCount, m.server.RequireLogin)
	m.server.AddRoute("/api/accounts/list", web.GET, m.GetAllAccounts, m.server.RequireLogin)

	// PACKAGES
	m.server.AddRoute("/api/packages/list", web.GET, m.GetPackages)
	m.server.AddRoute("/api/packages/get", web.GET, m.GetPackage)
	m.server.AddRoute("/api/packages/create", web.POST, m.CreatePackage, m.server.RequireLogin, m.server.ValidAPIKey, m.server.NoReentry)
	m.server.AddRoute("/api/packages/update", web.POST, m.UpdatePackage, m.server.RequireLogin, m.server.ValidAPIKey, m.server.NoReentry)
	m.server.AddRoute("/api/packages/buy", web.POST, m.BuyPackage, m.server.RequireLogin, m.server.NoReentry)
	m.server.AddRoute("/api/packages/upgrade", web.POST, m.upgradeSubscription, m.server.RequireLogin, m.server.NoReentry)
	m.server.AddRoute("/api/packages/subscription", web.GET, m.GetActiveSubscription, m.server.RequireLogin)
	m.server.AddRoute("/api/packages/subscription-count", web.GET, m.packageSubscriptions, m.server.RequireLogin)

	// TRANSFER
	m.server.AddRoute("/api/transfers/create", web.POST, m.makeTransfer, m.server.RequireLogin, m.server.NoReentry)
	m.server.AddRoute("/api/transfers/history", web.GET, m.transferHistory, m.server.RequireLogin)

	m.server.AddRoute("/api/withdrawals/create", web.POST, m.makeWithdrawal, m.server.RequireLogin, m.server.NoReentry)
	m.server.AddRoute("/api/withdrawals/history", web.GET, m.withdrawalHistory, m.server.RequireLogin)

	m.server.AddRoute("/api/notifications/send", web.POST, m.sendNotification, m.server.ValidAPIKey)
	m.server.AddRoute("/api/notifications/total-pending", web.GET, m.getUnReadNotificationCount, m.server.RequireLogin)
	m.server.AddRoute("/api/notifications/pending", web.GET, m.getNewNotifications, m.server.RequireLogin)
	m.server.AddRoute("/api/notifications/getall", web.GET, m.getNotifications, m.server.RequireLogin)
	m.server.AddRoute("/api/notifications/get", web.GET, m.getNotification, m.server.RequireLogin)

	m.server.AddRoute("/api/config/common-settings", web.GET, m.getCommonConfig, m.server.RequireLogin)

	m.server.AddRoute("/api/security/init2fa", web.POST, m.init2fa, m.server.RequireLogin, m.server.NoReentry)
	m.server.AddRoute("/api/security/enable2fa", web.POST, m.enable2fa, m.server.RequireLogin, m.server.NoReentry)
	m.server.AddRoute("/api/security/last-login", web.GET, m.lastLogin, m.server.RequireLogin)
}

func welcome(w http.ResponseWriter, r *http.Request) {
	web.SendJSON(w, "welcome to metstradas api. download the app from app store to start earning")
}

func (m module) runProcessor(ctx context.Context) {
	i := 1
	runner := func() {
		log.Info("runners are running", i)
		if err := m.db.BuildTradingSchedule(ctx); err != nil {
			log.Critical("runProcessor", "BuildTradingSchedule", err)
		}

		if err := m.db.PopulateTrades(ctx); err != nil {
			log.Critical("runProcessor", "PopulateTrades", err)
		}

		if err := m.db.PopulateEarnings(ctx); err != nil {
			log.Critical("runProcessor", "PopulateEarnings", err)
		}

		if err := m.db.ProcessWeeklyPayout(ctx); err != nil {
			log.Critical("runProcessor", "ProcessWeeklyPayout", err)
		}

		m.proccessPendingWithdrawal()

		i += 1
	}

	runner()

	next := time.Now().Add(24 * time.Hour)
	time.Sleep(time.Until(next))

	for {
		runner()

		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(24 * time.Hour)
			continue
		}
	}
}
