package postgres

import (
	"context"
	"database/sql"
	"merryworld/metatradas/app"
	"merryworld/metatradas/postgres/models"
	"time"

	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (pg PgDb) CreateNotification(ctx context.Context, accountID, title, message string) error {
	tx, err := pg.Db.Begin()
	if err != nil {
		return err
	}

	if err := pg.createNotificationTx(ctx, tx, accountID, title, message); err != nil {
		return tx.Rollback()
	}

	return tx.Commit()
}

func (pg PgDb) createNotificationTx(ctx context.Context, tx *sql.Tx, accountID, title, message string) error {
	notification := models.Notification{
		ID:        uuid.NewString(),
		AccountID: accountID,
		Date:      time.Now().Unix(),
		Status:    app.NOTIFICATION_STATUS_NEW,
		Title:     title,
		Content:   message,
	}

	return notification.Insert(ctx, tx, boil.Infer())
}

func (pg PgDb) UnReadNotificationCount(ctx context.Context, accountID string) (int64, error) {
	return models.Notifications(
		models.NotificationWhere.AccountID.EQ(accountID),
		models.NotificationWhere.Status.EQ(app.NOTIFICATION_STATUS_NEW),
	).Count(ctx, pg.Db)
}

func (pg PgDb) GetNotifications(ctx context.Context, accountID string, offset, limit int) (models.NotificationSlice, int64, error) {
	query := []qm.QueryMod{
		models.NotificationWhere.AccountID.EQ(accountID),
	}
	totalCount, err := models.Notifications(query...).Count(ctx, pg.Db)
	if err != nil {
		return nil, 0, err
	}

	query = append(query, qm.Offset(offset), qm.Limit(limit), 
		qm.OrderBy(models.NotificationColumns.Date+" desc"),
		qm.Select(models.NotificationColumns.ID, models.NotificationColumns.Title, models.NotificationColumns.Date),
	)

	notifications, err := models.Notifications(query...).All(ctx, pg.Db)

	return notifications, totalCount, err
}

func (pg PgDb) GetNewNotifications(ctx context.Context, accountID string, offset, limit int) (models.NotificationSlice, int64, error) {
	query := []qm.QueryMod{
		models.NotificationWhere.AccountID.EQ(accountID),
		models.NotificationWhere.Status.EQ(app.NOTIFICATION_STATUS_NEW),
	}
	totalCount, err := models.Notifications(query...).Count(ctx, pg.Db)
	if err != nil {
		return nil, 0, err
	}

	query = append(query, qm.Offset(offset), qm.Limit(limit), 
		qm.OrderBy(models.NotificationColumns.Date+" desc"),
		qm.Select(models.NotificationColumns.ID, models.NotificationColumns.Title, models.NotificationColumns.Date),
	)

	notifications, err := models.Notifications(query...).All(ctx, pg.Db)

	return notifications, totalCount, err
}

func (pg PgDb) GetNotification(ctx context.Context, id string) (*models.Notification, error) {
	notification, err := models.FindNotification(ctx, pg.Db, id)
	if err != nil {
		return nil, err
	}
	if notification.Status != app.NOTIFICATION_STATUS_READ {
		col := models.M{
			models.NotificationColumns.Status: app.NOTIFICATION_STATUS_READ,
		}
		if _, err = models.Notifications(models.NotificationWhere.ID.EQ(id)).UpdateAll(ctx, pg.Db, col); err != nil {
			return nil, err
		}
	}
	
	return notification, nil
}
