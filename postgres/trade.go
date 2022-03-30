package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"merryworld/metatradas/app"
	"merryworld/metatradas/postgres/models"

	"github.com/google/uuid"
	"github.com/jinzhu/now"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

const (
	ADAY = 86400
)

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

func (pg PgDb) ActiveTrades(ctx context.Context, accountID string) ([]app.Trade, error) {
	trades, err := models.Trades(
		models.TradeWhere.AccountID.EQ(accountID),
		models.TradeWhere.Date.EQ(now.BeginningOfDay().Unix()),
		models.TradeWhere.StartDate.LTE(time.Now().Unix()),
		qm.OrderBy(models.TradeColumns.StartDate+" desc"),
	).All(ctx, pg.Db)
	if err != nil {
		return nil, err
	}

	randomNumber := func(min, max int64) int64 {
		rand.Seed(time.Now().UnixNano())
		return int64(rand.Int63n(max-min+1) + (min))
	}

	var tradeView []app.Trade

	currentTime := time.Now().Unix()
	for _, t := range trades {
		trade := app.Trade{
			ID: t.ID,
			AccountID: t.AccountID,
			TradeNo: t.TradeNo,
			Date: t.Date,
			StartDate: t.StartDate,
			EndDate: t.EndDate,
			Amount: t.Amount,
			Profit: t.Profit,
		}

		if t.EndDate <= time.Now().Unix() {
			tradeView = append(tradeView, trade)
			continue
		}

		trade.EndDate = 0

		if currentTime-t.LastViewTime > ADAY/(24*12) { //every 5 mins
			if (currentTime - t.StartDate) < ADAY/24 { // just started
				trade.Profit = randomNumber((5*t.Profit)/100, (15*t.Profit)/100)
			} else if t.EndDate-currentTime <= ADAY/24 { // almost ended
				trade.Profit = randomNumber((95*t.Profit)/100, t.Profit)
			} else {
				trade.Profit = randomNumber((25*t.Profit)/100, (350*t.Profit)/100)
			}
			col := models.M{
				models.TradeColumns.LastViewProfit: trade.Profit,
				models.TradeColumns.LastViewTime:   currentTime,
			}
			if _, err := models.Trades(models.TradeWhere.ID.EQ(t.ID)).UpdateAll(ctx, pg.Db, col); err != nil {
				log.Error("ActiveTrades", "UpdateAll", err)
			}
		} else {
			trade.Profit = t.LastViewProfit
		}

		tradeView = append(tradeView, trade)

	}

	return tradeView, nil
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

func (pg PgDb) BuildTradingSchedule(ctx context.Context) error {
	// no trading on sundays
	if now.BeginningOfDay().Weekday() == time.Sunday {
		return nil
	}

	date := now.BeginningOfDay().Unix()
	minStartDate := now.BeginningOfDay().Add(time.Minute * 5).Unix()

	count, err := models.TradeSchedules(models.TradeScheduleWhere.Date.EQ(date)).Count(ctx, pg.Db)
	if err != nil {
		return err
	}
	// we build schedule once in a day
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

		for tradeNo := 1; tradeNo <= p.TradesPerDay; tradeNo += 1 {
			seed := (20 / p.TradesPerDay) * tradeNo
			maxStartDate := now.BeginningOfDay().Add(time.Hour * time.Duration(seed)).Unix()
			statement := `
				insert into trade_schedule(account_id, trade_no, total_trades, date, target_profit_percentage, start_date)
					select
						DISTINCT account.id as account_id,
						%d,
						%d,
						%d,
						COALESCE((floor(random()*(%d-%d+1))+%d), 0) as target_profit_percentage,
						COALESCE((floor(random()*(%d-%d+1))+%d), 0) as start_date
					from account 
						inner join subscription on account.id = subscription.account_id
					where 
						account.principal > 0 and
						subscription.start_date <= %d and subscription.end_date >= %d and subscription.package_id = '%s'
			`
			pDivisor := 30 * p.TradesPerDay
			if _, err := models.TradeSchedules(
				qm.SQL(fmt.Sprintf(statement,
					tradeNo,
					p.TradesPerDay,
					date,
					p.MaxReturnPerMonth*1000/pDivisor, p.MinReturnPerMonth*1000/pDivisor, p.MinReturnPerMonth*1000/pDivisor,
					maxStartDate,
					minStartDate,
					minStartDate,
					date,
					date,
					p.ID),
				),
			).ExecContext(ctx, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}

func (pg PgDb) PopulateTrades(ctx context.Context) error {
	// no earnings on sundays
	if now.BeginningOfDay().Weekday() == time.Sunday {
		return nil
	}

	date := now.BeginningOfDay().Unix()
	maxEndDate := date + ADAY - 60

	count, err := models.Trades(models.TradeWhere.Date.EQ(date)).Count(ctx, pg.Db)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	statement := `
			insert into trade (account_id, trade_no, date, start_date, end_date, amount, profit)
				select 
					account_id,
					trade_no,
					date,
					start_date,
					COALESCE((floor(random()*(%d-(start_date + 3600)+1))+(start_date + 3600)), 0) as end_date, 
					COALESCE((floor(random()*(
						(account.principal/total_trades) - 
						(account.principal/(total_trades*2)) +1) + (account.principal/(total_trades*2))
						)), 0) as amount,
					(account.principal * target_profit_percentage)/100000 as profit
				from trade_schedule 
				inner join account on account.id = trade_schedule.account_id
			 where 
			 date = %d
		`
	if _, err := models.DailyEarnings(
		qm.SQL(fmt.Sprintf(statement,
			maxEndDate,
			date),
		),
	).ExecContext(ctx, pg.Db); err != nil {
		return err
	}

	return nil
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

	tx, err := pg.Db.Begin()
	if err != nil {
		return err
	}

	statement := `
			insert into daily_earning (account_id, date, percentage, principal)
				select 
				distinct account_id, 
					date, 
					((sum(profit)*100000)/account.principal) as percentage,
					account.principal
				from trade
				inner join account on account.id = trade.account_id
				where date = %d group by account_id, date, account.principal
		`
	if _, err := models.DailyEarnings(
		qm.SQL(fmt.Sprintf(statement, date)),
	).ExecContext(ctx, tx); err != nil {
		tx.Rollback()
		return err
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
