package handlers

import (
	"deficonnect/defipayapi/app"
	"deficonnect/defipayapi/app/processors"
	"deficonnect/defipayapi/postgres"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func InitSlsApp() (*app.Module, error) {
	// Parse the configuration file, and setup logger.
	cfg, err := loadConfig()
	if err != nil {
		fmt.Printf("Failed to load pdanalytics config: %s\n", err.Error())
		return nil, err
	}

	db, err := postgres.NewPgDb(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName, os.Getenv("DEBUG_SQL") == "1")

	if err != nil {
		return nil, fmt.Errorf("pqsl: %v", err)
	}

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

	currencyProcessors := map[string]map[app.Network]app.CurrencyProcessor{}

	dfcProcessor, err := processors.NewDfcProcessor(bscClient, common.HexToAddress(cfg.DFCBscContractAddress))
	if err != nil {
		return nil, err
	}

	currencyProcessors[app.DFC.Name] = map[app.Network]app.CurrencyProcessor{}
	currencyProcessors[app.DFC.Name][app.Networks.BSC] = dfcProcessor

	return app.Start(db, bscClient, polygonClient, currencyProcessors, cfg.BlockchainConfig,
		cfg.MailgunDomain, cfg.MailgunAPIKey)
}
