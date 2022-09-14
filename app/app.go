package app

import (
	"deficonnect/defipayapi/web"

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

	return &app, nil
}
