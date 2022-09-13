package app

import (
	"net/http"
	
	"deficonnect/defipayapi/web"
)

type Network string

type Currency struct {
	Name     string
	Symbol   string
	Networks []Network
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
	web.SendJSON(w, []Currency{
		USDT, DFC, CGold,
	})
}
