// Code generated by SQLBoiler 4.7.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Notification is an object representing the database table.
type Notification struct {
	ID         string `boil:"id" json:"id" toml:"id" yaml:"id"`
	AccountID  string `boil:"account_id" json:"account_id" toml:"account_id" yaml:"account_id"`
	Status     int    `boil:"status" json:"status" toml:"status" yaml:"status"`
	Title      string `boil:"title" json:"title" toml:"title" yaml:"title"`
	Content    string `boil:"content" json:"content" toml:"content" yaml:"content"`
	Date       int64  `boil:"date" json:"date" toml:"date" yaml:"date"`
	Type       int    `boil:"type" json:"type" toml:"type" yaml:"type"`
	ActionLink string `boil:"action_link" json:"action_link" toml:"action_link" yaml:"action_link"`
	ActionText string `boil:"action_text" json:"action_text" toml:"action_text" yaml:"action_text"`

	R *notificationR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L notificationL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var NotificationColumns = struct {
	ID         string
	AccountID  string
	Status     string
	Title      string
	Content    string
	Date       string
	Type       string
	ActionLink string
	ActionText string
}{
	ID:         "id",
	AccountID:  "account_id",
	Status:     "status",
	Title:      "title",
	Content:    "content",
	Date:       "date",
	Type:       "type",
	ActionLink: "action_link",
	ActionText: "action_text",
}

var NotificationTableColumns = struct {
	ID         string
	AccountID  string
	Status     string
	Title      string
	Content    string
	Date       string
	Type       string
	ActionLink string
	ActionText string
}{
	ID:         "notification.id",
	AccountID:  "notification.account_id",
	Status:     "notification.status",
	Title:      "notification.title",
	Content:    "notification.content",
	Date:       "notification.date",
	Type:       "notification.type",
	ActionLink: "notification.action_link",
	ActionText: "notification.action_text",
}

// Generated where

var NotificationWhere = struct {
	ID         whereHelperstring
	AccountID  whereHelperstring
	Status     whereHelperint
	Title      whereHelperstring
	Content    whereHelperstring
	Date       whereHelperint64
	Type       whereHelperint
	ActionLink whereHelperstring
	ActionText whereHelperstring
}{
	ID:         whereHelperstring{field: "\"notification\".\"id\""},
	AccountID:  whereHelperstring{field: "\"notification\".\"account_id\""},
	Status:     whereHelperint{field: "\"notification\".\"status\""},
	Title:      whereHelperstring{field: "\"notification\".\"title\""},
	Content:    whereHelperstring{field: "\"notification\".\"content\""},
	Date:       whereHelperint64{field: "\"notification\".\"date\""},
	Type:       whereHelperint{field: "\"notification\".\"type\""},
	ActionLink: whereHelperstring{field: "\"notification\".\"action_link\""},
	ActionText: whereHelperstring{field: "\"notification\".\"action_text\""},
}

// NotificationRels is where relationship names are stored.
var NotificationRels = struct {
	Account string
}{
	Account: "Account",
}

// notificationR is where relationships are stored.
type notificationR struct {
	Account *Account `boil:"Account" json:"Account" toml:"Account" yaml:"Account"`
}

// NewStruct creates a new relationship struct
func (*notificationR) NewStruct() *notificationR {
	return &notificationR{}
}

// notificationL is where Load methods for each relationship are stored.
type notificationL struct{}

var (
	notificationAllColumns            = []string{"id", "account_id", "status", "title", "content", "date", "type", "action_link", "action_text"}
	notificationColumnsWithoutDefault = []string{"id", "account_id", "status", "title", "content", "date"}
	notificationColumnsWithDefault    = []string{"type", "action_link", "action_text"}
	notificationPrimaryKeyColumns     = []string{"id"}
)

type (
	// NotificationSlice is an alias for a slice of pointers to Notification.
	// This should almost always be used instead of []Notification.
	NotificationSlice []*Notification

	notificationQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	notificationType                 = reflect.TypeOf(&Notification{})
	notificationMapping              = queries.MakeStructMapping(notificationType)
	notificationPrimaryKeyMapping, _ = queries.BindMapping(notificationType, notificationMapping, notificationPrimaryKeyColumns)
	notificationInsertCacheMut       sync.RWMutex
	notificationInsertCache          = make(map[string]insertCache)
	notificationUpdateCacheMut       sync.RWMutex
	notificationUpdateCache          = make(map[string]updateCache)
	notificationUpsertCacheMut       sync.RWMutex
	notificationUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single notification record from the query.
func (q notificationQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Notification, error) {
	o := &Notification{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for notification")
	}

	return o, nil
}

// All returns all Notification records from the query.
func (q notificationQuery) All(ctx context.Context, exec boil.ContextExecutor) (NotificationSlice, error) {
	var o []*Notification

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Notification slice")
	}

	return o, nil
}

// Count returns the count of all Notification records in the query.
func (q notificationQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count notification rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q notificationQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if notification exists")
	}

	return count > 0, nil
}

// Account pointed to by the foreign key.
func (o *Notification) Account(mods ...qm.QueryMod) accountQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.AccountID),
	}

	queryMods = append(queryMods, mods...)

	query := Accounts(queryMods...)
	queries.SetFrom(query.Query, "\"account\"")

	return query
}

// LoadAccount allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (notificationL) LoadAccount(ctx context.Context, e boil.ContextExecutor, singular bool, maybeNotification interface{}, mods queries.Applicator) error {
	var slice []*Notification
	var object *Notification

	if singular {
		object = maybeNotification.(*Notification)
	} else {
		slice = *maybeNotification.(*[]*Notification)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &notificationR{}
		}
		args = append(args, object.AccountID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &notificationR{}
			}

			for _, a := range args {
				if a == obj.AccountID {
					continue Outer
				}
			}

			args = append(args, obj.AccountID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`account`),
		qm.WhereIn(`account.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Account")
	}

	var resultSlice []*Account
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Account")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for account")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for account")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Account = foreign
		if foreign.R == nil {
			foreign.R = &accountR{}
		}
		foreign.R.Notifications = append(foreign.R.Notifications, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.AccountID == foreign.ID {
				local.R.Account = foreign
				if foreign.R == nil {
					foreign.R = &accountR{}
				}
				foreign.R.Notifications = append(foreign.R.Notifications, local)
				break
			}
		}
	}

	return nil
}

// SetAccount of the notification to the related item.
// Sets o.R.Account to related.
// Adds o to related.R.Notifications.
func (o *Notification) SetAccount(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Account) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"notification\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"account_id"}),
		strmangle.WhereClause("\"", "\"", 2, notificationPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.AccountID = related.ID
	if o.R == nil {
		o.R = &notificationR{
			Account: related,
		}
	} else {
		o.R.Account = related
	}

	if related.R == nil {
		related.R = &accountR{
			Notifications: NotificationSlice{o},
		}
	} else {
		related.R.Notifications = append(related.R.Notifications, o)
	}

	return nil
}

// Notifications retrieves all the records using an executor.
func Notifications(mods ...qm.QueryMod) notificationQuery {
	mods = append(mods, qm.From("\"notification\""))
	return notificationQuery{NewQuery(mods...)}
}

// FindNotification retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindNotification(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*Notification, error) {
	notificationObj := &Notification{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"notification\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, notificationObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from notification")
	}

	return notificationObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Notification) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no notification provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(notificationColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	notificationInsertCacheMut.RLock()
	cache, cached := notificationInsertCache[key]
	notificationInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			notificationAllColumns,
			notificationColumnsWithDefault,
			notificationColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(notificationType, notificationMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(notificationType, notificationMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"notification\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"notification\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into notification")
	}

	if !cached {
		notificationInsertCacheMut.Lock()
		notificationInsertCache[key] = cache
		notificationInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the Notification.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Notification) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	notificationUpdateCacheMut.RLock()
	cache, cached := notificationUpdateCache[key]
	notificationUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			notificationAllColumns,
			notificationPrimaryKeyColumns,
		)

		if len(wl) == 0 {
			return 0, errors.New("models: unable to update notification, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"notification\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, notificationPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(notificationType, notificationMapping, append(wl, notificationPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update notification row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for notification")
	}

	if !cached {
		notificationUpdateCacheMut.Lock()
		notificationUpdateCache[key] = cache
		notificationUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q notificationQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for notification")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for notification")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o NotificationSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), notificationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"notification\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, notificationPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in notification slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all notification")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Notification) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no notification provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(notificationColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	notificationUpsertCacheMut.RLock()
	cache, cached := notificationUpsertCache[key]
	notificationUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			notificationAllColumns,
			notificationColumnsWithDefault,
			notificationColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			notificationAllColumns,
			notificationPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert notification, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(notificationPrimaryKeyColumns))
			copy(conflict, notificationPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"notification\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(notificationType, notificationMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(notificationType, notificationMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert notification")
	}

	if !cached {
		notificationUpsertCacheMut.Lock()
		notificationUpsertCache[key] = cache
		notificationUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single Notification record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Notification) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Notification provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), notificationPrimaryKeyMapping)
	sql := "DELETE FROM \"notification\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from notification")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for notification")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q notificationQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no notificationQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from notification")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for notification")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o NotificationSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), notificationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"notification\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, notificationPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from notification slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for notification")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Notification) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindNotification(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *NotificationSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := NotificationSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), notificationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"notification\".* FROM \"notification\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, notificationPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in NotificationSlice")
	}

	*o = slice

	return nil
}

// NotificationExists checks if the Notification row exists.
func NotificationExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"notification\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if notification exists")
	}

	return exists, nil
}
