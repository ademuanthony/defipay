package processors

import (
	"context"
	"crypto/ecdsa"
	"deficonnect/defipayapi/app/dfc"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"

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

type pancakeSwapPriceOutput struct {
	UpdatedAt int64 `json:"updated_at"`
	Data      struct {
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	} `json:"data"`
}

func (p dfcProcessor) DollarToToken(ctx context.Context, amount *big.Int) (*big.Int, error) {
	dfcPrice, err := getTokenPrice(ctx, "0x97A143545c0F8200222C051aC0a2Fc93ACBE6bA2", 18)
	if err != nil {
		return nil, err
	}
	fmt.Println("amount", amount)
	fmt.Println("dfcPrice", dfcPrice)

	tokenAmount := amount.Div(amount.Mul(amount, big.NewInt(1e8)), dfcPrice)

	fmt.Println("tokenAmount", tokenAmount)
	// todo use mainnet price
	return tokenAmount, nil
}

func getTokenPrice(ctx context.Context, contractAddress string, dollarDecimals int) (*big.Int, error) {
	url := fmt.Sprintf("https://api.pancakeswap.info/api/v2/tokens/%s", contractAddress)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	var result pancakeSwapPriceOutput
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	priceParts := strings.Split(result.Data.Price, ".")
	if len(priceParts) == 1 {
		priceParts = append(priceParts, "0")
	}

	if len(priceParts[1]) < dollarDecimals {
		priceParts[1] = priceParts[1] + strings.Repeat("0", dollarDecimals-len(priceParts[1]))
	}

	priceParts[1] = priceParts[1][0:dollarDecimals]

	priceStr := strings.Join(priceParts, "")
	priceStr = strings.TrimLeft(priceStr, "0")

	price, valid := common.Big0.SetString(priceStr, 10)
	if !valid {
		return nil, errors.New("invalid price string")
	}
	return price, nil
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
