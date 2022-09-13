package app

import (
	"math/big"
	"net/http"

	"deficonnect/defipayapi/web"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Network string

type Currency struct {
	Name     string
	Symbol   string
	Networks []Network
}

type CurrencyProcessor interface {
	BalanceOf(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error)
	Decimals(opts *bind.CallOpts) (uint8, error)
	Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error)
}

var (
	Networks = struct {
		BSC     Network
		Polygon Network
	}{
		BSC:     "Binance Smart Chain",
		Polygon: "Polygon",
	}

	USDT = Currency{
		Name:   "USDT",
		Symbol: "USDT",
		Networks: []Network{
			Networks.BSC, Networks.Polygon,
		},
	}

	DFC = Currency{
		Name:   "DefiConnect",
		Symbol: "DFC",
		Networks: []Network{
			Networks.BSC,
		},
	}

	CGold = Currency{
		Name:   "C250Gold",
		Symbol: "CGold",
		Networks: []Network{
			Networks.Polygon,
		},
	}
)

func (m module) supportedCurrencies(w http.ResponseWriter, r *http.Request) {
	currencies := []Currency{}
	for _, c := range []Currency{
		USDT, DFC, CGold,
	} {
		if m.currencyProcessors[c.Name] != nil {
			currencies = append(currencies, c)
		}
	}

	web.SendJSON(w, currencies)
}
