package app

import (
	"context"
	"net/http"

	"deficonnect/defipayapi/web"
)

type Network string

type Currency struct {
	Name     string
	Symbol   string
	Networks []Network
}

type CurrencyProcessor interface {
	CheckBalance(ctx context.Context, walletAddress string, network Network) (int64, error);
	Transfer(ctx context.Context, fromAddressPK, toAddress string, network Network) error
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
