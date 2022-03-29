package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"merryworld/metatradas/app"
	"merryworld/metatradas/postgres/models"
	"time"

	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (pg PgDb) CreatePackage(ctx context.Context, pkg models.Package) error {
	pkg.ID = uuid.NewString()
	return pkg.Insert(ctx, pg.Db, boil.Infer())
}

func (pg PgDb) PatchPackage(ctx context.Context, id string, input app.UpdatePackageInput) error {
	var upCol = models.M{}
	if input.Name != "" {
		upCol[models.PackageColumns.Name] = input.Name
	}
	if input.MinReturnPerMonth != 0 {
		upCol[models.PackageColumns.MinReturnPerMonth] = input.MinReturnPerMonth
	}
	if input.MaxReturnPerMonth != 0 {
		upCol[models.PackageColumns.MaxReturnPerMonth] = input.MaxReturnPerMonth
	}
	if input.Accuracy != 0 {
		upCol[models.PackageColumns.Accuracy] = input.Accuracy
	}
	if input.TradesPerDay != 0 {
		upCol[models.PackageColumns.TradesPerDay] = input.TradesPerDay
	}

	_, err := models.Packages(models.PackageWhere.ID.EQ(id)).UpdateAll(ctx, pg.Db, upCol)
	return err
}

func (pg PgDb) GetPackages(ctx context.Context) ([]*models.Package, error) {
	packages, err := models.Packages(
		qm.OrderBy(models.PackageColumns.Price+" desc"),
	).All(ctx, pg.Db)
	if err != nil {
		return nil, err
	}

	return packages, nil
}

func (pg PgDb) GetPackage(ctx context.Context, id string) (*models.Package, error) {
	return models.FindPackage(ctx, pg.Db, id)
}

func (pg PgDb) GetPackageByName(ctx context.Context, name string) (*models.Package, error) {
	return models.Packages(models.PackageWhere.Name.EQ(name)).One(ctx, pg.Db)
}

func (pg PgDb) CreateSubscription(ctx context.Context, accountID, packageID string, c250 bool) error {
	pkg, err := pg.GetPackage(ctx, packageID)
	if err != nil {
		return fmt.Errorf("GetPackage %v - %s", err, packageID)
	}

	acc, err := pg.GetAccount(ctx, accountID)
	if err != nil {
		return fmt.Errorf("GetAccount %v", err)
	}

	tx, err := pg.Db.Begin()
	if err != nil {
		return fmt.Errorf("Begin %v", err)
	}

	date := time.Now()

	if !c250 {
		note := "subscription to " + pkg.Name + " package"
		if err := pg.DebitAccountTx(ctx, tx, accountID, pkg.Price, date.Unix(), note); err != nil {
			tx.Rollback()
			return fmt.Errorf("DebitAccountTx %v", err)
		}
	}

	sub := models.Subscription{
		ID:        uuid.NewString(),
		PackageID: packageID,
		AccountID: accountID,
		StartDate: date.Unix(),
		EndDate:   date.Add(365 * 24 * time.Hour).Unix(),
	}

	if err := sub.Insert(ctx, tx, boil.Infer()); err != nil {
		tx.Rollback()
		return fmt.Errorf("Insert %v", err)
	}

	if acc.ReferralID.String != "" {
		// TODO: indicate c250 ref earnings
		if err := pg.payReferrer(ctx, tx, sub.ID, acc.ID, date.Unix(), acc.ReferralID.String, pkg.Price, 1, c250); err != nil {
			tx.Rollback()
			return fmt.Errorf("payReferrer %v", err)
		}
	}

	return tx.Commit()
}

func (pg PgDb) payReferrer(ctx context.Context, tx *sql.Tx, subscriptionID, payerID string, 
	date int64, refId string, subAmount int64, level int, c250 bool) error {
	// first level is 15%
	if level > 3 {
		return nil
	}
	var percentage int64

	switch level {
	case 1:
		percentage = 15
	case 2:
		percentage = 10
	case 3:
		percentage = 5
	}

	amount := subAmount * percentage / 100

	method := app.PAYMENTMETHOD_BNB
	if c250 {
		method = app.PAYMENTMETHOD_C250D
	}

	refPayout := models.ReferralPayout {
		ID: uuid.NewString(),
		AccountID: refId,
		Generation: level,
		Date: date,
		Amount: amount,
		SubscriptionID: subscriptionID,
		FromAccountID: payerID,
		PaymentStatus: app.PAYMENTSTATUS_PENDING,
		PaymentRef: "",
		PaymentMethod: method,
	}

	if err := refPayout.Insert(ctx, tx, boil.Infer()); err != nil {
		return err
	}
	// if err := pg.CreditAccountTx(ctx, tx, refId, amount,
	// 	date, "referral earning from "+payerUsername); err != nil {
	// 	return err
	// }
	acc, err := models.FindAccount(ctx, tx, refId)
	if err == sql.ErrNoRows || acc.ReferralID.String == "" {
		return nil
	}
	if err != nil {
		return err
	}
	return pg.payReferrer(ctx, tx, subscriptionID, payerID, date, acc.ReferralID.String, subAmount, level+1, c250)
}

func (pg PgDb) PendingReferralPayouts(ctx context.Context) (models.ReferralPayoutSlice, error) {
	return models.ReferralPayouts(
		models.ReferralPayoutWhere.PaymentStatus.EQ(app.PAYMENTSTATUS_PENDING),
	).All(ctx, pg.Db)
}

func (pg PgDb) UpdateReferralPayout(ctx context.Context, payout *models.ReferralPayout) error {
	_, err := payout.Update(ctx, pg.Db, boil.Infer())
	return err
}

func (pg PgDb) ActiveSubscription(ctx context.Context, accountID string) (*models.Subscription, error) {
	date := time.Now().Unix()
	return models.Subscriptions(
		models.SubscriptionWhere.StartDate.LTE(date),
		models.SubscriptionWhere.EndDate.GTE(date),
		models.SubscriptionWhere.AccountID.EQ(accountID),
		qm.Load(models.SubscriptionRels.Package),
	).One(ctx, pg.Db)
}
