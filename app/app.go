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
	config        AppConfig

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
	cfg AppConfig,
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

		currencyProcessors[DFC.Symbol] = map[Network]CurrencyProcessor{}
		currencyProcessors[DFC.Symbol][Networks.BSC] = dfcProcessor

		// ADD USDT processor

		usdtProcessor, err := processors.NewUsdtProcessor(bscClient, common.HexToAddress(cfg.USDTBscContractAddress))
		if err != nil {
			return nil, err
		}

		currencyProcessors[USDT.Symbol] = map[Network]CurrencyProcessor{}
		currencyProcessors[USDT.Symbol][Networks.BSC] = usdtProcessor

		// ADD BUSD processor

		busdProcessor, err := processors.NewBusdProcessor(bscClient, common.HexToAddress(cfg.BUSDContractAddress))
		if err != nil {
			return nil, err
		}

		currencyProcessors[BUSD.Symbol] = map[Network]CurrencyProcessor{}
		currencyProcessors[BUSD.Symbol][Networks.BSC] = busdProcessor

		app.bscClient = bscClient
		app.polygonClient = polygonClient
		app.currencyProcessors = currencyProcessors

	}

	return &app, nil
}
