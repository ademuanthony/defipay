package app

import (
	"merryworld/metatradas/web"
	"net/http"
)

func (m module) GetDepositAddress(w http.ResponseWriter, r *http.Request) {
	wallet, err := m.db.GetDepositAddress(r.Context(), m.server.GetUserIDTokenCtx(r))
	if err != nil {
		log.Critical("GetDepositAddress", err)
		web.RenderErrorfJSON(w, "Something went wrong. Please try again later")
		return
	}

	web.RenderJSON(w, wallet.Address)
}
