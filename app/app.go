package app

import (
	"context"
	"merryworld/metatradas/web"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jinzhu/now"
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

	go app.runProcessor(context.Background())
	go app.watchDeposit()

	return nil
}

func (m module) buildRoute() {
	m.server.AddRoute("/", web.GET, welcome)

	// AUTH
	m.server.AddRoute("/api/auth/register", web.POST, m.CreateAccount)
	m.server.AddRoute("/api/auth/login", web.POST, m.Login)

	//ACCOUNT
	m.server.AddRoute("/api/account/update", web.POST, m.UpdateAccountDetail, m.server.RequireLogin)
	m.server.AddRoute("/api/account/me", web.GET, m.GetAccountDetail, m.server.RequireLogin)
	m.server.AddRoute("/api/account/referral-count", web.GET, m.GetReferralCount, m.server.RequireLogin)
	m.server.AddRoute("/api/account/team-info", web.GET, m.TeamInformation, m.server.RequireLogin)
	m.server.AddRoute("/api/account/deposit-address", web.GET, m.GetDepositAddress, m.server.RequireLogin)
	m.server.AddRoute("/api/account/deposits", web.GET, m.DepositHistories, m.server.RequireLogin)
	m.server.AddRoute("/api/account/invest", web.POST, m.Invest, m.server.RequireLogin)
	m.server.AddRoute("/api/account/investments", web.GET, m.MyInvestments, m.server.RequireLogin)
	m.server.AddRoute("/api/account/release-investment", web.POST, m.ReleaseInvestment, m.server.RequireLogin, m.server.NoReentry)
	m.server.AddRoute("/api/account/daily-earnings", web.GET, m.MyDailyEarnings, m.server.RequireLogin)

	// ACCOUNTS
	m.server.AddRoute("/api/accounts/count", web.GET, m.GetAllAccountsCount, m.server.RequireLogin)
	m.server.AddRoute("/api/accounts/list", web.GET, m.GetAllAccounts, m.server.RequireLogin)

	// PACKAGES
	m.server.AddRoute("/api/packages/list", web.GET, m.GetPackages)
	m.server.AddRoute("/api/packages/get", web.GET, m.GetPackage)
	m.server.AddRoute("/api/packages/create", web.POST, m.CreatePackage, m.server.RequireLogin)
	m.server.AddRoute("/api/packages/update", web.POST, m.UpdatePackage, m.server.RequireLogin)
	m.server.AddRoute("/api/packages/buy", web.POST, m.BuyPackage, m.server.RequireLogin)
	m.server.AddRoute("/api/packages/subscription", web.GET, m.GetActiveSubscription, m.server.RequireLogin)

	// TRANSFER
	m.server.AddRoute("/api/transfers/create", web.POST, m.makeTransfer, m.server.RequireLogin, m.server.NoReentry)
	m.server.AddRoute("/api/transfers/history", web.GET, m.transferHistory, m.server.RequireLogin)

	m.server.AddRoute("/api/withdrawals/create", web.POST, m.makeWithdrawal, m.server.RequireLogin, m.server.NoReentry)
	m.server.AddRoute("/api/withdrawals/history", web.GET, m.withdrawalHistory, m.server.RequireLogin)

}

func welcome(w http.ResponseWriter, r *http.Request) {
	web.SendJSON(w, "welcome to metstradas api. download the app from app store to start earning")
}

func (m module) runProcessor(ctx context.Context) {
	log.Info("run processor")
	if err := m.db.PopulateEarnings(ctx); err != nil {
		log.Critical("runProcessor", "PopulateEarnings", err)
	}

	log.Info("run processor")
	if err := m.db.ProcessWeeklyPayout(ctx); err != nil {
		log.Critical("runProcessor", "ProcessWeeklyPayout", err)
	}

	next := now.BeginningOfDay().Add(24 * time.Hour)
	time.Sleep(time.Since(next))

	for {
		log.Info("run processor")
		if err := m.db.PopulateEarnings(ctx); err != nil {
			log.Critical("runProcessor", "PopulateEarnings", err)
		}

		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(24 * time.Hour)
			continue
		}
	}
}
