package postgres

import (
	"context"
	"merryworld/metatradas/postgres/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (pg PgDb) SetConfigValue(ctx context.Context, accountID, key, value string) error {
	exists, err := models.UserSettings(
		models.UserSettingWhere.AccountID.EQ(accountID),
		models.UserSettingWhere.ConfigKey.EQ(key),
	).Exists(ctx, pg.Db)
	if err != nil {
		return err
	}

	if exists {
		upCol := models.M{
			models.UserSettingColumns.ConfigValue: value,
		}
		_, err = models.UserSettings(
			models.UserSettingWhere.AccountID.EQ(accountID),
			models.UserSettingWhere.ConfigKey.EQ(key),
		).UpdateAll(ctx, pg.Db, upCol)
		return err
	}

	config := models.UserSetting{
		AccountID:   accountID,
		ConfigKey:   key,
		ConfigValue: value,
	}

	return config.Insert(ctx, pg.Db, boil.Infer())
}

func (pg PgDb) GetConfigValue(ctx context.Context, accountID, key string) (string, error) {
	config, err := models.UserSettings(
		qm.Select(models.UserSettingColumns.ConfigValue),
		models.UserSettingWhere.AccountID.EQ(accountID),
		models.UserSettingWhere.ConfigKey.EQ(key),
	).One(ctx, pg.Db)

	if err != nil {
		return "", err
	}

	return config.ConfigValue, nil
}

func (pg PgDb) GetConfigs(ctx context.Context, accountID string) (models.UserSettingSlice, error) {
	return models.UserSettings(
		models.UserSettingWhere.AccountID.EQ(accountID),
	).All(ctx, pg.Db)
}
