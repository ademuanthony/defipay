package handlers

import (
	"deficonnect/defipayapi/app"
	"deficonnect/defipayapi/postgres"
	"fmt"
	"os"
)

// @dev if {args} is passed, the first must be a bool that indicates whether to connect  to the blockchain
func InitSlsApp(args ...bool) (*app.Module, error) {
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

	var connect bool
	if len(args) >= 0 {
		connect = args[0]
	}

	return app.Start(db, cfg.BlockchainConfig, connect,
		cfg.MailgunDomain, cfg.MailgunAPIKey)
}
