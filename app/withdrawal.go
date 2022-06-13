package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"merryworld/metatradas/app/util"
	"merryworld/metatradas/postgres/models"
	"merryworld/metatradas/web"
	"net/http"
	"net/url"
	"os"
	"time"
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

	if sender.Username == "mike91" {
		web.SendJSON(w, true)
		return
	}

	if input.Amount < 200000 {
		web.SendErrorfJSON(w, "Minimum withdrawal amount is $20")
		return
	}

	if input.Amount > sender.Balance {
		web.SendErrorfJSON(w, "Insufficient balance")
		return
	}

	if sender.WithdrawalAddresss == "" {
		web.SendErrorfJSON(w, "Please set your withdrawal wallet first")
		return
	}

	if err := m.db.Withdraw(r.Context(), sender.ID, input.Amount); err != nil {
		log.Error("Withdraw", err)
		web.SendErrorfJSON(w, "Something went wrong, please try again later")
		return
	}

	web.SendJSON(w, true)
}

func (m module) cancelWithdrawal(w http.ResponseWriter, r http.Request) {

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

func (m module) proccessPendingWithdrawal() {
	ctx := context.Background()
	withdrawals, err := m.db.GetWithdrawalsForProcessing(ctx)
	if err != nil {
		log.Error("proccessPendingWithdrawal->GetWithdrawalsForProcessing", err)
		return
	}
	for _, rec := range withdrawals {
		if err = m.proccessWithdrawal(ctx, rec); err != nil {
			log.Error("proccessPendingWithdrawal->proccessWithdrawal", err)
		}
	}
}

func (m module) proccessWithdrawal(ctx context.Context, withdarwal *models.Withdrawal) error {
	if withdarwal.ID == "22fb694e-dbbf-4282-94f0-d42585eac597" {
		return errors.New("mike91 tries to withdraw")
	}
	bnbAmount, err := m.convertClubDollarToBnb(ctx, withdarwal.Amount-5000) //blockchain fee
	if err != nil {
		return err
	}

	txHash, err := m.transfer(ctx, m.config.MasterAddressKey, withdarwal.Destination, bnbAmount)
	if err != nil {
		return fmt.Errorf("m.transfer %s %v", m.config.MasterAddress, err)
	}
	if err := m.db.SetWithdrawalTxHash(ctx, withdarwal.ID, txHash); err != nil {
		return fmt.Errorf("SetWithdrawalTxHash %v", err)
	}

	message := fmt.Sprintf("Your withdrawal request of %.2f has been processed successfully",
		float64(withdarwal.Amount)/float64(10000))
	title := "Withdrawal Proccessed"
	if err := m.db.CreateNotification(ctx, withdarwal.ID, title, message, "", "", NOTIFICATION_TYPE_TOPBAR); err != nil {
		return fmt.Errorf("CreateNotification %v", err)
	}

	return nil
}

func (m module) processReferralPayouts() {
	for {
		func() {
			defer time.Sleep(5 * time.Minute)
			ctx := context.Background()
			pendingPayouts, err := m.db.PendingReferralPayouts(ctx)
			if err != nil {
				log.Error("processReferralPayouts", "PendingReferralPayouts", err)
				return
			}
			for _, payout := range pendingPayouts {
				if err := m.proccessPayout(ctx, payout); err != nil {
					log.Error("processReferralPayouts", "proccessPayout", err)
				}
			}
		}()
	}
}

func (m module) proccessPayout(ctx context.Context, payout *models.ReferralPayout) error {
	account, err := m.db.GetAccount(ctx, payout.AccountID)
	if err != nil {
		return fmt.Errorf("GetAccount", err)
	}

	markCompletedAndNotify := func() error {
		payout.PaymentStatus = PAYMENTSTATUS_COMPLETED
		if err := m.db.UpdateReferralPayout(ctx, payout); err != nil {
			return err
		}

		notificationTitle := "Referral payment received"
		senderAccount, err := m.db.GetAccount(ctx, payout.FromAccountID)
		if err != nil {
			return fmt.Errorf("GetAccount %v", err)
		}
		destination := "wallet"
		if payout.PaymentMethod == PAYMENTMETHOD_C250D {
			destination = "Club250Cent backoffice"
		} else if account.WithdrawalAddresss == "" {
			destination = "available balance"
		}
		message := fmt.Sprintf("A referral commission of $%.4f was sent to your %s for %s",
			float64(payout.Amount)/float64(10000), destination, senderAccount.Username)

		return m.db.CreateNotification(ctx, payout.AccountID, notificationTitle, message, "", "", NOTIFICATION_TYPE_TOPBAR)
	}

	if payout.PaymentMethod == PAYMENTMETHOD_C250D {
		if err := m.transferC250Dollar(ctx, account.Username, payout.Amount); err != nil {
			// if strings.Contains(err.Error(), " transfer failed") {
			// 	payout.PaymentStatus = PAYMENTSTATUS_FAILED
			// 	if err := m.db.UpdateReferralPayout(ctx, payout); err != nil {
			// 		return err
			// 	}
			// }
			payout.PaymentStatus = PAYMENTSTATUS_FAILED
			if err := m.db.UpdateReferralPayout(ctx, payout); err != nil {
				return err
			}
			return fmt.Errorf("transferC250Dollar %v", err)
		}

		return markCompletedAndNotify()
	}

	if account.WithdrawalAddresss == "" || !util.IsValidAddress(account.WithdrawalAddresss) {
		if m.db.CreditAccount(ctx, payout.AccountID, payout.Amount, time.Now().Unix(), "referral payout from "+payout.FromAccountID); err != nil {
			return fmt.Errorf("CreditAccount %v", err)
		}
		return markCompletedAndNotify()
	}

	bnbAmount, err := m.convertClubDollarToBnb(ctx, payout.Amount-5000) //blockchain fee
	if err != nil {
		return fmt.Errorf("convertClubDollarToBnb %v", err)
	}

	txHash, err := m.transfer(ctx, m.config.MasterAddressKey, account.WithdrawalAddresss, bnbAmount)
	if err != nil {
		return fmt.Errorf("transfer %v", err)
	}
	payout.PaymentRef = txHash
	return markCompletedAndNotify()
}

type c250TransferInput struct {
	Username string `json:"username"`
	Amount   int64  `json:"amount"`
}

type c250TransferOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (m module) transferC250Dollar(ctx context.Context, username string, amount int64) error {
	input := c250TransferInput{Username: username, Amount: amount}
	var resp c250TransferOutput
	path := "/api/metatradas/transfer-payout?API_AUTH_KEY=" + url.QueryEscape(os.Getenv("ACCESS_SECRET"))
	if err := sendJsonRequest(ctx, &http.Client{}, http.MethodPost, path, input, &resp); err != nil {
		return err
	}
	if !resp.Success {
		msg := resp.Message
		if resp.Message == "" {
			msg = "transfer failed for " + username
		}
		return errors.New(msg)
	}
	return nil
}

// GetResponse attempts to collect json data from the given url string and decodes it into
// the destination
func sendJsonRequest(ctx context.Context, client *http.Client, method, url string, reqBody, destination interface{}) error {
	// if client has no timeout, set one
	if client.Timeout == time.Duration(0) {
		client.Timeout = 10 * time.Second
	}
	resp := new(http.Response)

	baseURL := "https://club250cent.com"
	url = baseURL + url

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("json.Marshal %v", err)
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req = req.WithContext(ctx)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	maxRetryAttempts := 1
	retryDelay := 0 * time.Second

	for i := 1; i <= maxRetryAttempts; i++ {
		res, err := client.Do(req)
		if err != nil {
			if res != nil {
				res.Body.Close()
			}
			if i == maxRetryAttempts {
				return err
			}
			time.Sleep(retryDelay)
			continue
		}
		resp = res
		break
	}

	err = json.NewDecoder(resp.Body).Decode(destination)
	if err != nil {
		return err
	}
	return nil
}
