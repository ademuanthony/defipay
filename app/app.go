package app

import (
	"merryworld/metatradas/web"

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
	//app.server.AddRoute("/api/account/deposits", web.GET, app.GetDeposits, app.server.RequireLogin)

	// ACCOUNTS
	app.server.AddRoute("/api/accounts/count", web.GET, app.GetAllAccountsCount, app.server.RequireLogin)
	app.server.AddRoute("/api/accounts/list", web.GET, app.GetAllAccounts, app.server.RequireLogin)

	return nil
}
