package app

import (
	"context"
	"database/sql"
	"math/big"
	"merryworld/metatradas/app/dfc"
	"merryworld/metatradas/web"
	"net/http"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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

func (m module) watchDeposit() {
	dfcToken, err := dfc.NewDfc(common.HexToAddress(os.Getenv("USDT_CONTRACT")), m.client)
	if err != nil {
		log.Error("watchDeposit", err)
		return
	}

	var sink = make(chan *dfc.DfcTransfer)

	go func() {
		for {
			func() {
				addresses, err := m.db.GetWalletByAddresses(context.Background())
				if err != nil {
					log.Error("GetWalletByAddresses", err)
				}

				var toAddresses []common.Address
				for _, add := range addresses {
					toAddresses = append(toAddresses, common.HexToAddress(add))
				}

				sub, err := dfcToken.WatchTransfer(&bind.WatchOpts{}, sink, nil, toAddresses)
				if err != nil {
					log.Error("watchTranfer", err)
					return
				}
				defer sub.Unsubscribe()
				time.Sleep(5*time.Minute)
			}()
		}
	}()

	for {
		tx := <-sink
		log.Info("processing deposit at" + tx.To.Hex())
		// mi deposit is 20$
		if tx.Value.Div(tx.Value, big.NewInt(1e18)).Int64() < 20 {
			continue
		}

		wallet, err := m.db.GetWellatByAddress(context.Background(), tx.To.String())
		if err == sql.ErrNoRows {
			continue
		}
		if err != nil {
			log.Critical("GetWalletByAddress", err)
			continue
		}
		// TODO: move the fund to the main wallet

		//get wallet address
		// process deposit
		divisor := big.NewInt(1e14)
		amountBig := tx.Value.Div(tx.Value, divisor)

		if err := m.db.CreateDeposit(context.Background(), wallet.AccountID, tx.Raw.BlockHash.Hex(), amountBig.Int64()); err != nil {
			log.Critical("CreateDeposit", err)
			continue
		}
	}

}
