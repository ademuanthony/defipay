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
		ID:                 uuid.NewString(),
		ReferralID:         null.StringFrom(input.ReferralID),
		ReferralID2:        null.StringFrom(input.ReferralID2),
		ReferralID3:        null.StringFrom(input.ReferralID3),
		Username:           input.Username,
		Password:           input.Password,
		Email:              input.Email,
		PhoneNumber:        input.PhoneNumber,
		WithdrawalAddresss: input.WalletAddress,
		FirstName:          input.Name,
		CreatedAt:          time.Now().Unix(),
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

func (pg PgDb) CreditAccountTx(ctx context.Context, tx *sql.Tx, accountID string, amount, date int64, ref string) error {
	// transaction := models.AccountTransaction{
	// 	AccountID:   accountID,
	// 	Amount:      amount,
	// 	TXType:      app.TxTypeCredit,
	// 	Date:        date,
	// 	Description: ref,
	// }

	// if err := transaction.Insert(ctx, tx, boil.Infer()); err != nil {
	// 	return err
	// }

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
		models.InvestmentWhere.Status.EQ(0),
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

func (pg PgDb) Investment(ctx context.Context, id string) (*models.Investment, error) {
	return models.FindInvestment(ctx, pg.Db, id)
}

func (pg PgDb) ReleaseInvestment(ctx context.Context, id string) error {
	investment, err := pg.Investment(ctx, id)
	if err != nil {
		return err
	}

	tx, err := pg.Db.Begin()
	if err != nil {
		return err
	}

	no, err := investment.Delete(ctx, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if no < 1 {
		tx.Rollback()
		return errors.New("no investment was released")
	}

	statement := `update account set balance = balance + $1, principal = principal - $1, matured_principal = matured_principal - $1 where id = $2`
	if _, err = models.Accounts(qm.SQL(statement, investment.Amount, investment.AccountID)).ExecContext(ctx, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
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
	// no earnings on sundays
	if now.BeginningOfDay().Weekday() == time.Sunday {
		return nil
	}

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
					COALESCE((floor(random()*(%d-%d+1))+%d), 0) as percentage, 
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
	COALESCE(sum((daily_earning.principal * daily_earning.percentage)/100000), 0) as principal 
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
		 COALESCE(sum((daily_earning.principal * daily_earning.percentage)/100000), 0) as total from 
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
