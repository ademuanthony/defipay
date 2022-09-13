package app

import (
	"deficonnect/defipayapi/web"
	"net/http"
)

type Network string

type Currency struct {
	Name     string
	Symbol   string
	Networks map[Network]string
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
		Networks: map[Network]string{
			Networks.BSC:     "0x55d398326f99059fF775485246999027B3197955",
			Networks.Polygon: "0xc2132D05D31c914a87C6611C10748AEb04B58e8F",
		},
	}

	DFC = Currency{
		Name:   "DefiConnect",
		Symbol: "DFC",
		Networks: map[Network]string{
			Networks.BSC: "0x97a143545c0f8200222c051ac0a2fc93acbe6ba2",
		},
	}

	CGold = Currency{
		Name:   "C250Gold",
		Symbol: "CGold",
		Networks: map[Network]string{
			Networks.Polygon: "0xbC19aA6Ed11fD143D181999bFB2f972a62d91b57",
		},
	}
)

func (m module) supportedCurrencies(w http.ResponseWriter, r *http.Request) {
	web.SendJSON(w, []Currency{
		USDT, DFC, CGold,
	})
}
