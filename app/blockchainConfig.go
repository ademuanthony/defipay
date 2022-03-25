package app

type BlockchainConfig struct {
	// Blockchain
	BSCNode            string `long:"mainnet" env:"MAINNET_NODE_ADDRESS"`
	MasterAddressKey   string `long:"MASTER_ADDRESS_KEY" env:"MASTER_ADDRESS_KEY"`
	MasterAddress      string `long:"MASTER_ADDRESS" env:"MASTER_ADDRESS"`
	// PremiumPrivateKey  string `long:"PREMIUM_PRIVATE_KEY" env:"PREMIUM_PRIVATE_KEY"`
	PrivateWallet      string `long:"RESERVE_WALLET" env:"RESERVE_WALLET"`
	// PremiumWallet      string `long:"PREMIUM_ADDRESS" env:"PREMIUM_ADDRESS"`
	USDTContractAddress string `long:"USDT_CONTRACT_ADDRESS" env:"USDT_CONTRACT_ADDRESS"`
}
