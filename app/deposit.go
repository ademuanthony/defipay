package app

import (
	"merryworld/metatradas/web"
	"net/http"
)

func (m module) GetDepositAddress(w http.ResponseWriter, r *http.Request) {
	wallet, err := m.db.GetDepositAddress(r.Context(), m.server.GetUserIDTokenCtx(r))
	if err != nil {
		log.Critical("GetDepositAddress", err)
		web.SendErrorfJSON(w, "Something went wrong. Please try again later")
		return
	}

	web.SendJSON(w, wallet.Address)
}

func (m module) DepositHistories(w http.ResponseWriter, r *http.Request) {
	pageReq := web.GetPanitionInfo(r)
	deposits, totalCount, err := m.db.GetDeposits(r.Context(), m.server.GetUserIDTokenCtx(r), pageReq.Offset, pageReq.Limit)
	if err != nil {
		log.Error("DepositHistories", "GetDeposits", err)
		web.SendErrorfJSON(w, "Spmething went wrong. Please try again later")
		return
	}

	web.SendPagedJSON(w, deposits, totalCount)

}
