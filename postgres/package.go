package postgres

import (
	"context"
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
	packages, err := models.Packages().All(ctx, pg.Db)
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

func (pg PgDb) CreateSubscription(ctx context.Context, accountID, packageID string) error {
	pkg, err := pg.GetPackage(ctx, packageID)
	if err != nil {
		return err
	}

	acc, err := pg.GetAccount(ctx, accountID)
	if err != nil {
		return err
	}

	tx, err := pg.Db.Begin()
	if err != nil {
		return err
	}

	date := time.Now()
	note := "subscription to " + pkg.Name + " package"
	if err := pg.DebitAccountTx(ctx, tx, accountID, pkg.Price, date.Unix(), note); err != nil {
		tx.Rollback()
		return err
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
		return err
	}

	if acc.ReferralID.String != "" {
		refAmount := pkg.Price / 2
		if err := pg.CreditAccountTx(ctx, tx, acc.ReferralID.String, refAmount,
			date.Unix(), "referral earning from "+acc.Username); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
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
