package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"merryworld/metatradas/app"
	"merryworld/metatradas/postgres/models"
	"time"

	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (pg PgDb) CreateAccount(ctx context.Context, input app.CreateAccountInput) error {
	account := models.Account{
		ID:          uuid.NewString(),
		ReferralID:  null.StringFrom(input.ReferralID),
		ReferralID2: null.StringFrom(input.ReferralID2),
		ReferralID3: null.StringFrom(input.ReferralID3),
		Username:    input.Username,
		Password:    input.Password,
		Email:       input.Email,
		PhoneNumber: input.PhoneNumber,
		FirstName:   input.Name,
		CreatedAt:   time.Now().Unix(),
	}

	tx, err := pg.Db.Begin()
	if err != nil {
		return err
	}

	err = account.Insert(ctx, tx, boil.Infer())
	if err != nil {
		tx.Rollback()
		return err
	}

	wallet := models.Wallet{
		ID:         uuid.NewString(),
		AccountID:  account.ID,
		Address:    input.DepositWalletAddress,
		PrivateKey: input.PrivateKey,
		CoinSymbol: "BEP20-USDT",
	}

	if err = wallet.Insert(ctx, tx, boil.Infer()); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pg PgDb) AddLogin(ctx context.Context, accountID, ip, platform string, date int64) error {
	info := models.LoginInfo{
		AccountID: accountID,
		Platform:  platform,
		IP:        ip,
		Date:      date,
	}

	return info.Insert(ctx, pg.Db, boil.Infer())
}

func (pg PgDb) LastLogin(ctx context.Context) (*models.LoginInfo, error) {
	maxDate := time.Now().Add(-1 * time.Minute).Unix()
	return models.LoginInfos(
		models.LoginInfoWhere.Date.LTE(maxDate),
		qm.OrderBy(models.LoginInfoColumns.Date+" desc"),
	).One(ctx, pg.Db)
}

func (pg PgDb) CreateDepositWallet(ctx context.Context, accountID, address, privateKey string) (*models.Wallet, error) {
	wallet := models.Wallet{
		ID:         uuid.NewString(),
		AccountID:  accountID,
		Address:    address,
		PrivateKey: privateKey,
		CoinSymbol: "BEP20-USDT",
	}

	err := wallet.Insert(ctx, pg.Db, boil.Infer())
	return &wallet, err
}

func (pg PgDb) GetAccount(ctx context.Context, id string) (*models.Account, error) {
	acc, err := models.FindAccount(ctx, pg.Db, id)
	if err != nil {
		return nil, err
	}
	bal, err := pg.AccountBalance(ctx, id)
	if err != nil {
		return nil, err
	}
	acc.Balance = bal

	return acc, nil

}

func (pg PgDb) GetPasswordResetCode(ctx context.Context, accountID string) (string, error) {
	// delete expired code
	minDate := time.Now().Add(-15 * time.Minute)
	if _, err := models.SecurityCodes(models.SecurityCodeWhere.Date.LTE(minDate.Unix())).DeleteAll(ctx, pg.Db); err != nil {
		return "", err
	}

	code, err := models.SecurityCodes(models.SecurityCodeWhere.Date.GT(minDate.Unix())).One(ctx, pg.Db)
	if err == sql.ErrNoRows {
		code = &models.SecurityCode{
			Code: randomCode(6),
			AccountID: accountID,
			Date: time.Now().Unix(),
		}
		if err = code.Insert(ctx, pg.Db, boil.Infer()); err != nil {
			return "", err
		}
	}

	if err != nil {
		return "", err
	}

	return code.Code, err
}

func (pg PgDb) ValidatePasswordResetCode(ctx context.Context, accountID, code string) (bool, error) {
	minDate := time.Now().Add(-15 * time.Minute)
	if _, err := models.SecurityCodes(models.SecurityCodeWhere.Date.LTE(minDate.Unix())).DeleteAll(ctx, pg.Db); err != nil {
		return false, err
	}

	lastCode, err := models.SecurityCodes(
		models.SecurityCodeWhere.AccountID.EQ(accountID),
		models.SecurityCodeWhere.Date.GT(minDate.Unix()),
		qm.OrderBy(models.SecurityCodeColumns.Date + " desc"),
	).One(ctx, pg.Db)
	if err != nil {
		return false, err
	}

	return code == lastCode.Code, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("0123456789")

func randomCode(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (pg PgDb) GetAllAccountsCount(ctx context.Context) (int64, error) {
	return models.Accounts().Count(ctx, pg.Db)
}

func (pg PgDb) GetAccounts(ctx context.Context, skip, limit int) ([]*models.Account, error) {
	return models.Accounts(
		qm.Offset(skip),
		qm.Limit(limit),
	).All(ctx, pg.Db)
}

func (pg PgDb) GetAccountIDs(ctx context.Context) ([]string, error) {
	accounts, err := models.Accounts(
		qm.Select(models.AccountColumns.ID),
	).All(ctx, pg.Db)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, acc := range accounts {
		ids = append(ids, acc.ID)
	}

	return ids, nil
}

func (pg PgDb) GetAccountByUsername(ctx context.Context, username string) (*models.Account, error) {
	return models.Accounts(
		models.AccountWhere.Username.EQ(username),
	).One(ctx, pg.Db)
}

func (pg PgDb) UpdateAccountDetail(ctx context.Context, accountID string, input app.UpdateDetailInput) error {
	var upCol = models.M{}
	if input.FirstName != "" {
		upCol[models.AccountColumns.FirstName] = input.FirstName
	}
	if input.PhoneNumber != "" {
		upCol[models.AccountColumns.PhoneNumber] = input.PhoneNumber
	}
	if input.LastName != "" {
		upCol[models.AccountColumns.LastName] = input.LastName
	}
	if input.WithdrawalAddress != "" {
		upCol[models.AccountColumns.WithdrawalAddresss] = input.WithdrawalAddress
	}

	_, err := models.Accounts(models.AccountWhere.ID.EQ(accountID)).UpdateAll(ctx, pg.Db, upCol)
	return err
}

func (pg PgDb) ChangePassword(ctx context.Context, accountID, password string) error {
	colUp := models.M{
		models.AccountColumns.Password: password,
	}
	_, err := models.Accounts(models.AccountWhere.ID.EQ(accountID)).UpdateAll(ctx, pg.Db, colUp)
	return err
}

func (pg PgDb) GetRefferalCount(ctx context.Context, accountID string) (int64, error) {
	return models.Accounts(
		models.AccountWhere.ReferralID.EQ(null.StringFrom(accountID)),
	).Count(ctx, pg.Db)
}

func (pg PgDb) GetTeamInformation(ctx context.Context, accountID string) (*app.TeamInfo, error) {
	g1, err := models.Accounts(
		models.AccountWhere.ReferralID.EQ(null.StringFrom(accountID)),
	).Count(ctx, pg.Db)
	if err != nil {
		return nil, err
	}

	g2, err := models.Accounts(
		models.AccountWhere.ReferralID2.EQ(null.StringFrom(accountID)),
	).Count(ctx, pg.Db)
	if err != nil {
		return nil, err
	}

	g3, err := models.Accounts(
		models.AccountWhere.ReferralID3.EQ(null.StringFrom(accountID)),
	).Count(ctx, pg.Db)
	if err != nil {
		return nil, err
	}

	statement := "select coalesce(sum(principal), 0) as principal from account where referral_id = $1"
	acc, err := models.Accounts(qm.SQL(statement, accountID)).One(ctx, pg.Db)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	var p1 int64
	if acc != nil {
		p1 = acc.Principal
	}

	statement = "select coalesce(sum(principal), 0) as principal from account where referral_id_2 = $1"
	acc, err = models.Accounts(qm.SQL(statement, accountID)).One(ctx, pg.Db)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	var p2 int64
	if acc != nil {
		p2 = acc.Principal
	}

	statement = "select coalesce(sum(principal), 0) as principal from account where referral_id_3 = $1"
	acc, err = models.Accounts(qm.SQL(statement, accountID)).One(ctx, pg.Db)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	var p3 int64
	if acc != nil {
		p3 = acc.Principal
	}

	return &app.TeamInfo{
		FirstGeneration:   g1,
		SecoundGeneration: g2,
		ThirdGeneration:   g3,
		Pool1:             p1,
		Pool2:             p2,
		Pool3:             p3,
	}, nil
}

func (pg PgDb) GetDepositAddress(ctx context.Context, accountID string) (*models.Wallet, error) {
	return models.Wallets(models.WalletWhere.AccountID.EQ(accountID)).One(ctx, pg.Db)
}

func (pg PgDb) GetDeposits(ctx context.Context, accountID string, offset, limit int) ([]*models.Deposit, int64, error) {
	deposits, err := models.Deposits(
		models.DepositWhere.AccountID.EQ(accountID),
		qm.OrderBy(models.DepositColumns.Date+" desc"),
		qm.Limit(limit), qm.Offset(offset),
	).All(ctx, pg.Db)

	if err != nil {
		return nil, 0, err
	}

	totalCount, err := models.Deposits(models.DepositWhere.AccountID.EQ(accountID)).Count(ctx, pg.Db)

	if err != nil {
		return nil, 0, err
	}

	return deposits, totalCount, nil
}

func (pg PgDb) CreditAccount(ctx context.Context, accountID string, amount, date int64, ref string) error {
	tx, err := pg.Db.Begin()
	if err != nil {
		return err
	}
	if err := pg.CreditAccountTx(ctx, tx, accountID, amount, date, ref); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pg PgDb) CreditAccountTx(ctx context.Context, tx *sql.Tx, accountID string, amount, date int64, ref string) error {
	transaction := models.AccountTransaction{
		AccountID:   accountID,
		Amount:      amount,
		TXType:      app.TxTypeCredit,
		Date:        date,
		Description: ref,
	}

	if err := transaction.Insert(ctx, tx, boil.Infer()); err != nil {
		return err
	}

	return nil
}

func (pg PgDb) DebitAccountTx(ctx context.Context, tx *sql.Tx, accountID string, amount, date int64, ref string) error {
	acc, err := models.Accounts(qm.Select(models.AccountColumns.Balance), models.AccountWhere.ID.EQ(accountID)).One(ctx, tx)
	if err != nil {
		return err
	}
	if acc.Balance < amount {
		return errors.New("insufficient balance")
	}
	transaction := models.AccountTransaction{
		AccountID:   accountID,
		Amount:      amount,
		TXType:      app.TxTypeDebit,
		Date:        date,
		Description: ref,
	}

	if err := transaction.Insert(ctx, tx, boil.Infer()); err != nil {
		return err
	}
	
	return nil
}

func (pg PgDb) AccountBalance(ctx context.Context, accountId string) (int64, error) {
	var statement = `SELECT 
	SUM(amount) AS balance FROM (
		SELECT
			CASE WHEN tx.tx_type = 'credit' THEN tx.amount ELSE -1 * tx.amount END AS amount 
		FROM account_transaction tx
		WHERE tx.account_id = $1
	) res`

	var result null.Int64
	err := pg.Db.QueryRow(statement, accountId).Scan(&result)
	if err != nil && err.Error() == sql.ErrNoRows.Error() {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("pg.Db.QueryRow %v", err)
	}
	return result.Int64, err
}

func (pg PgDb) MyDownlines(ctx context.Context, accountID string, generation int64, offset, limit int) ([]app.DownlineInfo, int64, error) {
	query := []qm.QueryMod{}
	switch generation {
	case 1:
		query = append(query, models.AccountWhere.ReferralID.EQ(null.StringFrom(accountID)))
	case 2:
		models.AccountWhere.ReferralID2.EQ(null.StringFrom(accountID))
	case 3:
		models.AccountWhere.ReferralID3.EQ(null.StringFrom(accountID))
	}

	totalCount, err := models.Accounts(query...).Count(ctx, pg.Db)
	if err != nil {
		return nil, 0, err
	}

	query = append(query, qm.Load(models.AccountRels.Subscriptions),
		qm.OrderBy(models.AccountColumns.CreatedAt+" desc"),
		qm.Offset(offset),
		qm.Limit(limit))

	accounts, err := models.Accounts(
		query...,
	).All(ctx, pg.Db)

	if err != nil {
		return nil, 0, err
	}

	var downlines []app.DownlineInfo
	for _, acc := range accounts {
		downline := app.DownlineInfo{
			ID:        acc.ID,
			Username:  acc.Username,
			FirstName: acc.FirstName,
			LastName:  acc.LastName,
			Date:      acc.CreatedAt,
		}
		currcentData := time.Now().Unix()
		for _, s := range acc.R.Subscriptions {
			if s.StartDate <= currcentData && s.EndDate >= currcentData {
				pkg, err := models.FindPackage(ctx, pg.Db, s.PackageID)
				if err != nil {
					return nil, 0, err
				}
				downline.PackageName = pkg.Name
			}
		}
		downlines = append(downlines, downline)
	}

	return downlines, totalCount, nil
}
