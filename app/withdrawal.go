package app

import (
	"encoding/json"
	"merryworld/metatradas/web"
	"net/http"
)

type MakeWithdrawalInput struct {
	Amount int64 `json:"amount"`
}

func (m module) makeWithdrawal(w http.ResponseWriter, r *http.Request) {
	var input MakeWithdrawalInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Critical("makeTransfer", "json::Decode", err)
		web.SendErrorfJSON(w, "Error is decoding request. Please try again later")
		return
	}

	sender, err := m.currentAccount(r)
	if err != nil {
		log.Error("makeTransfer", err)
		web.SendErrorfJSON(w, "Something went wrong, please try again later")
		return
	}

	if input.Amount > sender.Balance {
		web.SendErrorfJSON(w, "Insufficient balance")
		return
	}

	if err := m.db.Withdraw(r.Context(), sender.ID, input.Amount); err != nil {
		log.Error("Withdraw", err)
		web.SendErrorfJSON(w, "Something went wrong, please try again later")
		return
	}

	web.SendJSON(w, true)
}

func (m module) withdrawalHistory(w http.ResponseWriter, r *http.Request) {
	pagedReq := web.GetPanitionInfo(r)
	rec, total, err := m.db.Withdrawals(r.Context(), m.server.GetUserIDTokenCtx(r), pagedReq.Offset, pagedReq.Limit)
	if err != nil {
		log.Error("Withdrawals", err)
		web.SendErrorfJSON(w, "Something went wrong, please try again later")
		return
	}

	web.SendPagedJSON(w, rec, total)
}
