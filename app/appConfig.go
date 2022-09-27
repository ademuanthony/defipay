package app

type AppConfig struct {
	// Blockchain
	BSCNode     string `long:"bsc-node" env:"BSC_NODE"`
	PolygonNode string `long:"polygon-node" env:"POLYGON_NODE"`

	MasterAddressKey string `long:"MASTER_ADDRESS_KEY" env:"MASTER_ADDRESS_KEY"`
	MasterAddress    string `long:"MASTER_ADDRESS" env:"MASTER_ADDRESS"`

	// PremiumWallet      string `long:"PREMIUM_ADDRESS" env:"PREMIUM_ADDRESS"`
	USDTBscContractAddress      string `env:"USDT_BSC"`
	USDTPolygonContractAddress  string `env:"USDT_POLYGON"`
	DFCBscContractAddress       string `env:"DFC_BSC"`
	CGoldPolygonContractAddress string `env:"CGOLD_POLYGON"`

	MastAccountID string `env:"MASTER_ACCOUNT_ID"`
}
