package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"time"

	"merryworld/metatradas/app/usdt"
	"merryworld/metatradas/app/util"
	"merryworld/metatradas/postgres/models"
	"merryworld/metatradas/web"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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

	dfcToken, err := usdt.NewUsdt(common.HexToAddress(os.Getenv("USDT_CONTRACT_ADDRESS")), m.client)
	if err != nil {
		log.Error("watchDeposit", err)
		return
	}

	var sink = make(chan *usdt.UsdtTransfer)

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
				log.Info("bsc watching...")
				defer sub.Unsubscribe()
				time.Sleep(5 * time.Minute)
			}()
		}
	}()

	for {
		tx := <-sink
		log.Info("processing deposit at " + tx.To.Hex())
		amount := tx.Value.Quo(tx.Value, big.NewInt(1e14)).Int64()
		// mi deposit is 20$
		if amount < 2*1e4 {
			log.Info("deposit amount too small")
			continue
		}

		ctx := context.Background()

		wallet, err := m.db.GetWellatByAddress(ctx, tx.To.String())
		if err == sql.ErrNoRows {
			log.Warn("strange, address not found", tx.To.Hex())
			continue
		}
		if err != nil {
			log.Critical("GetWalletByAddress", err)
			continue
		}

		_, err = m.moveBalanceToMaster(ctx, dfcToken, wallet)
		if err != nil {
			log.Error("moveBalanceToMaster", wallet.Address, err)
			continue
		}

		// $0.2 blockchain fee
		amount = amount - 2000

		if err := m.db.CreateDeposit(context.Background(), wallet.AccountID, tx.Raw.TxHash.Hex(), amount); err != nil {
			log.Critical("CreateDeposit", err)
			continue
		}
	}

}

func (m module) moveBalanceToMaster(ctx context.Context, token *usdt.Usdt, wallet *models.Wallet) (string, error) {

	bnbBal, err := m.checkBalance(ctx, wallet.Address)
	if err != nil {
		log.Errorf("moveBalanceToMaster->m.checkBalance %v", err)
		return "", errors.New("error in processing payment. Please try again later or contact the admin for help")
	}

	if bnbBal.Int64() < m.feeAmount().Int64() {
		if err := m.sendTokenTransferFee(ctx, wallet.Address); err != nil {
			log.Errorf("processDFCDeposit->m.sendTokenTransferFee %v", err)
			return "", err
		}
	}

	bal, err := token.BalanceOf(nil, common.HexToAddress(wallet.Address))
	if err != nil {
		log.Errorf("moveBalanceToMaster->BalanceOf %v", err)
		return "", err
	}

	return m.transferToken(ctx, wallet.PrivateKey, m.config.MasterAddress, bal)
}

func (m module) transferToken(ctx context.Context, privateKeyStr, to string, value *big.Int) (string, error) {
	if !util.IsValidAddress(to) {
		return "", errors.New("invalid address")
	}
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return "", err
	}

	dfcToken, err := usdt.NewUsdt(common.HexToAddress(m.config.USDTContractAddress), m.client)
	if err != nil {
		return "", err
	}

	toAddress := common.HexToAddress(to)

	chainID, err := m.client.ChainID(ctx)
	if err != nil {
		return "", fmt.Errorf("client.ChainID() %v", err)
	}
	opts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return "", fmt.Errorf("bind.NewKeyedTransactorWithChainID %v", err)
	}
	opts.GasLimit = 60000

	tx, err := dfcToken.Transfer(opts, toAddress, value)
	if err != nil {
		return "", fmt.Errorf("dfcToken.Transfer %v", err)
	}

	return tx.Hash().Hex(), nil
}
