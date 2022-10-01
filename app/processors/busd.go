package processors

import (
	"context"
	"deficonnect/defipayapi/app/usdt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type busdProcessor struct {
	client   *ethclient.Client
	instance usdt.Usdt
}

func NewBusdProcessor(client *ethclient.Client, contractAddress common.Address) (*busdProcessor, error) {
	instance, err := usdt.NewUsdt(contractAddress, client)
	if err != nil {
		return nil, err
	}

	return &busdProcessor{
		client:   client,
		instance: *instance,
	}, nil
}

func (p busdProcessor) BalanceOf(ctx context.Context, walletAddress common.Address) (*big.Int, error) {
	return p.instance.BalanceOf(nil, walletAddress)
}

func (p busdProcessor) DollarToToken(ctx context.Context, amount *big.Int) (*big.Int, error) {
	dfcPrice, err := getTokenPrice(ctx, "0xe9e7CEA3DedcA5984780Bafc599bD69ADd087D56", 18)
	if err != nil {
		return nil, err
	}

	tokenAmount := amount.Div(amount.Mul(amount, big.NewInt(1e18)), dfcPrice)

	return tokenAmount, nil
}

func (p busdProcessor) Decimals(ctx context.Context) (uint8, error) {
	decimals, err := p.instance.Decimals(nil)
	if err != nil {
		return 0, err
	}
	return decimals, nil
}

func (p busdProcessor) Transfer(ctx context.Context, privateKey string, to common.Address,
	value *big.Int) (*types.Transaction, error) {
	opt, err := getAccountAuth(p.client, privateKey, 0)
	if err != nil {
		return nil, err
	}
	return p.instance.Transfer(opt, to, value)
}
