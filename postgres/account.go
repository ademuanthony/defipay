package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"merryworld/metatradas/app"
	"merryworld/metatradas/postgres/models"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/now"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (pg PgDb) CreateAccount(ctx context.Context, input app.CreateAccountInput) error {
	account := models.Account{
		ID:          uuid.NewString(),
		ReferralID:  null.StringFrom(input.ReferralID),
		Username:    input.Username,
		Password:    input.Password,
		Email:       input.Email,
		PhoneNumber: input.PhoneNumber,
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
		AccountID:  account.ID,
		Address:    input.WalletAddress,
		PrivateKey: input.PrivateKey,
		CoinSymbol: "BEP20-USDT",
	}

	if err = wallet.Insert(ctx, tx, boil.Infer()); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pg PgDb) GetAccount(ctx context.Context, id string) (*models.Account, error) {
	return models.FindAccount(ctx, pg.Db, id)
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

func (pg PgDb) GetDepositAddress(ctx context.Context, accountID string) (*models.Wallet, error) {
	return models.Wallets(models.WalletWhere.AccountID.EQ(accountID)).One(ctx, pg.Db)
}

func (pg PgDb) GetDeposits(ctx context.Context, accountID string, offset, limit int) ([]*models.Deposit, int64, error) {
	deposits, err := models.Deposits(
		models.DepositWhere.AccountID.EQ(accountID),
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

	statement := `update account set balance = balance + $1 where id = $2`
	_, err := models.Accounts(qm.SQL(statement, amount, accountID)).ExecContext(ctx, pg.Db)

	return err
}

func (pg PgDb) DebitAccountTx(ctx context.Context, tx *sql.Tx, accountID string, amount, date int64, ref string) error {
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

	statement := `update account set balance = balance - $1 where id = $2`
	_, err := models.Accounts(qm.SQL(statement, amount, accountID)).ExecContext(ctx, pg.Db)

	return err
}

func (pg PgDb) Invest(ctx context.Context, accountID string, amount int64) error {
	tx, err := pg.Db.Begin()
	if err != nil {
		return err
	}

	acc, err := models.FindAccount(ctx, tx, accountID)
	if err != nil {
		return err
	}
	if acc.Balance < amount {
		tx.Rollback()
		return errors.New("insufficient fund")
	}

	investment := models.Investment{
		ID:             uuid.NewString(),
		AccountID:      accountID,
		Amount:         amount,
		Date:           time.Now().Unix(),
		ActivationDate: now.BeginningOfDay().UTC().Add(24 * time.Hour).Unix(),
	}

	if err := investment.Insert(ctx, tx, boil.Infer()); err != nil {
		tx.Rollback()
		return err
	}

	statement := `update account set balance = balance - $1, principal = principal + $1 where id = $2`
	if _, err = models.Accounts(qm.SQL(statement, amount, accountID)).ExecContext(ctx, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()

}

func (pg PgDb) Investments(ctx context.Context, accountId string, offset, limit int) ([]*models.Investment, int64, error) {
	rec, err := models.Investments(
		models.InvestmentWhere.AccountID.EQ(accountId),
		qm.Offset(offset),
		qm.Limit(limit),
		qm.OrderBy(models.InvestmentColumns.Date+" desc"),
	).All(ctx, pg.Db)

	if err != nil {
		return nil, 0, err
	}

	count, err := models.Investments(models.InvestmentWhere.AccountID.EQ(accountId)).Count(ctx, pg.Db)
	if err != nil {
		return nil, 0, err
	}

	return rec, count, nil
}

func (pg PgDb) DailyEarnings(ctx context.Context, accountId string, offset, limit int) ([]*models.DailyEarning, int64, error) {
	today := now.BeginningOfDay().Unix()

	rec, err := models.DailyEarnings(
		models.DailyEarningWhere.AccountID.EQ(accountId),
		models.DailyEarningWhere.Date.LT(today),
		qm.Offset(offset),
		qm.Limit(limit),
		qm.OrderBy(models.DailyEarningColumns.Date+" desc"),
	).All(ctx, pg.Db)

	if err != nil {
		return nil, 0, err
	}

	count, err := models.DailyEarnings(
		models.DailyEarningWhere.AccountID.EQ(accountId),
		models.DailyEarningWhere.Date.LT(today),
	).Count(ctx, pg.Db)
	if err != nil {
		return nil, 0, err
	}

	return rec, count, nil
}

func (pg PgDb) PopulateEarnings(ctx context.Context) error {

	date := now.BeginningOfDay().Unix()
	count, err := models.DailyEarnings(models.DailyEarningWhere.Date.EQ(date)).Count(ctx, pg.Db)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	packages, err := pg.GetPackages(ctx)
	if err != nil {
		return err
	}
	tx, err := pg.Db.Begin()
	if err != nil {
		return err
	}

	for _, p := range packages {
		count, _ := models.Subscriptions(
			models.SubscriptionWhere.PackageID.EQ(p.ID),
		).Count(ctx, tx)
		if count == 0 {
			continue
		}

		statement := `
			insert into daily_earning (account_id, date, percentage, principal)
				select 
					DISTINCT account.id as account_id,
					%d, 
					(floor(random()*(%d-%d+1))+%d) as percentage, 
					account.principal
				from account 
				inner join subscription on account.id = subscription.account_id
			 where 
			 	account.principal > 0 and
			 	subscription.start_date <= %d and subscription.end_date >= %d and subscription.package_id = '%s'
		`
		if _, err := models.DailyEarnings(
			qm.SQL(fmt.Sprintf(statement, date,
				p.MaxReturnPerMonth*1000/30, p.MinReturnPerMonth*1000/30, p.MinReturnPerMonth*1000/30, date, date, p.ID),
			),
		).ExecContext(ctx, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (pg PgDb) ProcessWeeklyPayout(ctx context.Context) error {
	date := now.BeginningOfDay()
	if date.Weekday() != time.Sunday {
		return nil
	}

	lastPayout, err := models.WeeklyPayouts(
		qm.OrderBy(models.WeeklyPayoutColumns.Date+" desc"),
	).One(ctx, pg.Db)
	if err == sql.ErrNoRows {
		lastPayout = &models.WeeklyPayout{}
	}
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	lastPayDate := time.Unix(lastPayout.Date, 0)
	if time.Since(lastPayDate).Hours() < 24*7 {
		return nil
	}

	tx, err := pg.Db.Begin()
	if err != nil {
		return err
	}

	statement := `select
	 sum((daily_earning.principal * daily_earning.percentage)/100000) as principal 
	 from daily_earning where date >= $1 and date < $2`
	totalDaily, err := models.DailyEarnings(qm.SQL(statement, lastPayout.Date, today)).One(ctx, tx)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return err
	}

	weeklyPay := models.WeeklyPayout{
		ID:     uuid.NewString(),
		Date:   today,
		Amount: totalDaily.Principal,
	}

	if err := weeklyPay.Insert(ctx, tx, boil.Infer()); err != nil {
		tx.Rollback()
		return err
	}

	statement = `
	update account set balance  = balance + floor(sub.total) FROM (
		select distinct
		 daily_earning.account_id,
		 sum((daily_earning.principal * daily_earning.percentage)/100000) as total from 
		 daily_earning
		where daily_earning.date >= $1 and daily_earning.date < $2
		group by daily_earning.account_id
		) sub where id = sub.account_id`

	if _, err := models.DailyEarnings(qm.SQL(statement, lastPayout.Date, today)).ExecContext(ctx, pg.Db); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
