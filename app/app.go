package app

import (
	"context"
	"merryworld/metatradas/web"
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

	// AUTH
	app.server.AddRoute("/api/auth/register", web.POST, app.CreateAccount)
	app.server.AddRoute("/api/auth/login", web.POST, app.Login)

	//ACCOUNT
	app.server.AddRoute("/api/account/update", web.POST, app.UpdateAccountDetail, app.server.RequireLogin)
	app.server.AddRoute("/api/account/me", web.GET, app.GetAccountDetail, app.server.RequireLogin)
	app.server.AddRoute("/api/account/deposit-address", web.GET, app.GetDepositAddress, app.server.RequireLogin)
	app.server.AddRoute("/api/account/deposits", web.GET, app.DepositHistories, app.server.RequireLogin)
	app.server.AddRoute("/api/account/invest", web.POST, app.Invest, app.server.RequireLogin)
	app.server.AddRoute("/api/account/investments", web.GET, app.MyInvestments, app.server.RequireLogin)
	app.server.AddRoute("/api/account/daily-earnings", web.GET, app.MyDailyEarnings, app.server.RequireLogin)

	// ACCOUNTS
	app.server.AddRoute("/api/accounts/count", web.GET, app.GetAllAccountsCount, app.server.RequireLogin)
	app.server.AddRoute("/api/accounts/list", web.GET, app.GetAllAccounts, app.server.RequireLogin)

	// PACKAGES
	app.server.AddRoute("/api/packages/list", web.GET, app.GetPackages)
	app.server.AddRoute("/api/packages/get", web.GET, app.GetPackage)
	app.server.AddRoute("/api/packages/create", web.POST, app.CreatePackage, server.RequireLogin)
	app.server.AddRoute("/api/packages/update", web.POST, app.UpdatePackage, server.RequireLogin)
	app.server.AddRoute("/api/packages/buy", web.POST, app.BuyPackage, server.RequireLogin)
	app.server.AddRoute("/api/packages/subscription", web.GET, app.GetActiveSubscription, server.RequireLogin)

	// TRANSFER
	app.server.AddRoute("/api/transfers/create", web.POST, app.makeTransfer, server.RequireLogin, server.NoReentry)
	app.server.AddRoute("/api/transfers/history", web.GET, app.transferHistory, server.RequireLogin)

	app.server.AddRoute("/api/withdrawals/create", web.POST, app.makeWithdrawal, server.RequireLogin, server.NoReentry)
	app.server.AddRoute("/api/withdrawals/history", web.GET, app.withdrawalHistory, server.RequireLogin)

	go app.runProcessor(context.Background())

	return nil
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
