package app

import (
	"encoding/json"
	"merryworld/metatradas/web"
	"net/http"
)

type MakeTransferInput struct {
	ReceiverUsername string `json:"receiver_username"`
	Amount           int64  `json:"amount"`
}

type TransferViewModel struct {
	ID       string `json:"id"`
	Amount   int64  `json:"amount"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Date     int64  `json:"date"`
}

func (m module) makeTransfer(w http.ResponseWriter, r *http.Request) {
	var input MakeTransferInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Critical("makeTransfer", "json::Decode", err)
		web.SendErrorfJSON(w, "Error is decoding request. Please try again later")
		return
	}

	receiver, err := m.db.GetAccountByUsername(r.Context(), input.ReceiverUsername)
	if err != nil {
		web.SendErrorfJSON(w, "Invalid username")
		return
	}

	sender, err := m.currentAccount(r)
	if err != nil {
		log.Error("makeTransfer", err)
		web.SendErrorfJSON(w, "Something went wrong, please try again later")
		return
	}

	if receiver.ID == sender.ID {
		web.SendErrorfJSON(w, "The receiver must be an external account")
		return
	}

	if input.Amount > sender.Balance {
		web.SendErrorfJSON(w, "Insufficient balance")
		return
	}

	if err := m.db.Transfer(r.Context(), sender.ID, receiver.ID, input.Amount); err != nil {
		log.Error("makeTransfer", err)
		web.SendErrorfJSON(w, "Something went wrong, please try again later")
		return
	}

	web.SendJSON(w, true)
}

func (m module) transferHistory(w http.ResponseWriter, r *http.Request) {
	pagedReq := web.GetPanitionInfo(r)
	rec, total, err := m.db.TransferHistories(r.Context(), m.server.GetUserIDTokenCtx(r), pagedReq.Offset, pagedReq.Limit)
	if err != nil {
		log.Error("TransferHistories", err)
		web.SendErrorfJSON(w, "Something went wrong, please try again later")
		return
	}

	web.SendPagedJSON(w, rec, total)
}
