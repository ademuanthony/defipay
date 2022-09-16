package app

import (
	"deficonnect/defipayapi/app/processors"
	"deficonnect/defipayapi/web"

	"github.com/ethereum/go-ethereum/common"
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

func Start(db store,
	cfg BlockchainConfig,
	connectBlockchain bool,
	mgDomain, mgKey string) (*Module, error) {
	log.Info("starting...")

	app := Module{
		db:       db,
		config:   cfg,
		MgDomain: mgDomain,
		MgKey:    mgKey,
		server:   &web.Server{},
	}

	if connectBlockchain {
		bscClient, err := ethclient.Dial(cfg.BSCNode)
		if err != nil {
			return nil, err
		}

		polygonClient, err := ethclient.Dial(cfg.PolygonNode)
		if err != nil {
			return nil, err
		} else {
			defer bscClient.Close()
		}

		currencyProcessors := map[string]map[Network]CurrencyProcessor{}

		dfcProcessor, err := processors.NewDfcProcessor(bscClient, common.HexToAddress(cfg.DFCBscContractAddress))
		if err != nil {
			return nil, err
		}

		currencyProcessors[DFC.Name] = map[Network]CurrencyProcessor{}
		currencyProcessors[DFC.Name][Networks.BSC] = dfcProcessor

		app.bscClient = bscClient
		app.polygonClient = polygonClient
		app.currencyProcessors = currencyProcessors

	}

	return &app, nil
}
