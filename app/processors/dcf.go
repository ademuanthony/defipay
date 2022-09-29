package processors

import (
	"context"
	"crypto/ecdsa"
	"deficonnect/defipayapi/app/dfc"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type dfcProcessor struct {
	client   *ethclient.Client
	instance dfc.Dfc
}

func NewDfcProcessor(client *ethclient.Client, contractAddress common.Address) (*dfcProcessor, error) {
	instance, err := dfc.NewDfc(contractAddress, client)
	if err != nil {
		return nil, err
	}

	return &dfcProcessor{
		client:   client,
		instance: *instance,
	}, nil
}

func (p dfcProcessor) BalanceOf(ctx context.Context, walletAddress common.Address) (*big.Int, error) {
	return p.instance.BalanceOf(nil, walletAddress)
}

func (p dfcProcessor) DollarToToken(ctx context.Context, amount *big.Int) (*big.Int, error) {
	// todo use mainnet price
	return amount, nil
}

func (p dfcProcessor) Decimals(ctx context.Context) (uint8, error) {
	decimals, err := p.instance.Decimals(nil)
	if err != nil {
		return 0, err
	}
	return uint8(decimals.Int64()), nil
}

func (p dfcProcessor) Transfer(ctx context.Context, privateKey string, to common.Address,
	value *big.Int) (*types.Transaction, error) {
	opt, err := getAccountAuth(p.client, privateKey, 0)
	if err != nil {
		return nil, err
	}
	return p.instance.Transfer(opt, to, value)
}

func getAccountAuth(client *ethclient.Client, privateKeyString string, gasMultiplyer int64) (*bind.TransactOpts, error) {

	privateKey, err := crypto.HexToECDSA(privateKeyString)
	if err != nil {
		panic(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("invalid key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	//fetch the last use nonce of account
	nounce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		panic(err)
	}

	auth.Nonce = big.NewInt(int64(nounce))
	auth.Value = big.NewInt(0)                     // in wei 10:56
	auth.GasLimit = uint64(384696 * gasMultiplyer) // in units
	auth.GasPrice = gasPrice                       //big.NewInt(30001000047)

	return auth, nil
}
