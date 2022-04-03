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

// ReferralPayout is an object representing the database table.
type ReferralPayout struct {
	ID             string `boil:"id" json:"id" toml:"id" yaml:"id"`
	AccountID      string `boil:"account_id" json:"account_id" toml:"account_id" yaml:"account_id"`
	FromAccountID  string `boil:"from_account_id" json:"from_account_id" toml:"from_account_id" yaml:"from_account_id"`
	SubscriptionID string `boil:"subscription_id" json:"subscription_id" toml:"subscription_id" yaml:"subscription_id"`
	Generation     int    `boil:"generation" json:"generation" toml:"generation" yaml:"generation"`
	Amount         int64  `boil:"amount" json:"amount" toml:"amount" yaml:"amount"`
	Date           int64  `boil:"date" json:"date" toml:"date" yaml:"date"`
	PaymentMethod  int    `boil:"payment_method" json:"payment_method" toml:"payment_method" yaml:"payment_method"`
	PaymentStatus  int    `boil:"payment_status" json:"payment_status" toml:"payment_status" yaml:"payment_status"`
	PaymentRef     string `boil:"payment_ref" json:"payment_ref" toml:"payment_ref" yaml:"payment_ref"`

	R *referralPayoutR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L referralPayoutL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ReferralPayoutColumns = struct {
	ID             string
	AccountID      string
	FromAccountID  string
	SubscriptionID string
	Generation     string
	Amount         string
	Date           string
	PaymentMethod  string
	PaymentStatus  string
	PaymentRef     string
}{
	ID:             "id",
	AccountID:      "account_id",
	FromAccountID:  "from_account_id",
	SubscriptionID: "subscription_id",
	Generation:     "generation",
	Amount:         "amount",
	Date:           "date",
	PaymentMethod:  "payment_method",
	PaymentStatus:  "payment_status",
	PaymentRef:     "payment_ref",
}

var ReferralPayoutTableColumns = struct {
	ID             string
	AccountID      string
	FromAccountID  string
	SubscriptionID string
	Generation     string
	Amount         string
	Date           string
	PaymentMethod  string
	PaymentStatus  string
	PaymentRef     string
}{
	ID:             "referral_payout.id",
	AccountID:      "referral_payout.account_id",
	FromAccountID:  "referral_payout.from_account_id",
	SubscriptionID: "referral_payout.subscription_id",
	Generation:     "referral_payout.generation",
	Amount:         "referral_payout.amount",
	Date:           "referral_payout.date",
	PaymentMethod:  "referral_payout.payment_method",
	PaymentStatus:  "referral_payout.payment_status",
	PaymentRef:     "referral_payout.payment_ref",
}

// Generated where

var ReferralPayoutWhere = struct {
	ID             whereHelperstring
	AccountID      whereHelperstring
	FromAccountID  whereHelperstring
	SubscriptionID whereHelperstring
	Generation     whereHelperint
	Amount         whereHelperint64
	Date           whereHelperint64
	PaymentMethod  whereHelperint
	PaymentStatus  whereHelperint
	PaymentRef     whereHelperstring
}{
	ID:             whereHelperstring{field: "\"referral_payout\".\"id\""},
	AccountID:      whereHelperstring{field: "\"referral_payout\".\"account_id\""},
	FromAccountID:  whereHelperstring{field: "\"referral_payout\".\"from_account_id\""},
	SubscriptionID: whereHelperstring{field: "\"referral_payout\".\"subscription_id\""},
	Generation:     whereHelperint{field: "\"referral_payout\".\"generation\""},
	Amount:         whereHelperint64{field: "\"referral_payout\".\"amount\""},
	Date:           whereHelperint64{field: "\"referral_payout\".\"date\""},
	PaymentMethod:  whereHelperint{field: "\"referral_payout\".\"payment_method\""},
	PaymentStatus:  whereHelperint{field: "\"referral_payout\".\"payment_status\""},
	PaymentRef:     whereHelperstring{field: "\"referral_payout\".\"payment_ref\""},
}

// ReferralPayoutRels is where relationship names are stored.
var ReferralPayoutRels = struct {
	Account     string
	FromAccount string
}{
	Account:     "Account",
	FromAccount: "FromAccount",
}

// referralPayoutR is where relationships are stored.
type referralPayoutR struct {
	Account     *Account `boil:"Account" json:"Account" toml:"Account" yaml:"Account"`
	FromAccount *Account `boil:"FromAccount" json:"FromAccount" toml:"FromAccount" yaml:"FromAccount"`
}

// NewStruct creates a new relationship struct
func (*referralPayoutR) NewStruct() *referralPayoutR {
	return &referralPayoutR{}
}

// referralPayoutL is where Load methods for each relationship are stored.
type referralPayoutL struct{}

var (
	referralPayoutAllColumns            = []string{"id", "account_id", "from_account_id", "subscription_id", "generation", "amount", "date", "payment_method", "payment_status", "payment_ref"}
	referralPayoutColumnsWithoutDefault = []string{"id", "account_id", "from_account_id", "subscription_id", "generation", "amount", "date", "payment_method", "payment_status", "payment_ref"}
	referralPayoutColumnsWithDefault    = []string{}
	referralPayoutPrimaryKeyColumns     = []string{"id"}
)

type (
	// ReferralPayoutSlice is an alias for a slice of pointers to ReferralPayout.
	// This should almost always be used instead of []ReferralPayout.
	ReferralPayoutSlice []*ReferralPayout

	referralPayoutQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	referralPayoutType                 = reflect.TypeOf(&ReferralPayout{})
	referralPayoutMapping              = queries.MakeStructMapping(referralPayoutType)
	referralPayoutPrimaryKeyMapping, _ = queries.BindMapping(referralPayoutType, referralPayoutMapping, referralPayoutPrimaryKeyColumns)
	referralPayoutInsertCacheMut       sync.RWMutex
	referralPayoutInsertCache          = make(map[string]insertCache)
	referralPayoutUpdateCacheMut       sync.RWMutex
	referralPayoutUpdateCache          = make(map[string]updateCache)
	referralPayoutUpsertCacheMut       sync.RWMutex
	referralPayoutUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single referralPayout record from the query.
func (q referralPayoutQuery) One(ctx context.Context, exec boil.ContextExecutor) (*ReferralPayout, error) {
	o := &ReferralPayout{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for referral_payout")
	}

	return o, nil
}

// All returns all ReferralPayout records from the query.
func (q referralPayoutQuery) All(ctx context.Context, exec boil.ContextExecutor) (ReferralPayoutSlice, error) {
	var o []*ReferralPayout

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to ReferralPayout slice")
	}

	return o, nil
}

// Count returns the count of all ReferralPayout records in the query.
func (q referralPayoutQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count referral_payout rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q referralPayoutQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if referral_payout exists")
	}

	return count > 0, nil
}

// Account pointed to by the foreign key.
func (o *ReferralPayout) Account(mods ...qm.QueryMod) accountQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.AccountID),
	}

	queryMods = append(queryMods, mods...)

	query := Accounts(queryMods...)
	queries.SetFrom(query.Query, "\"account\"")

	return query
}

// FromAccount pointed to by the foreign key.
func (o *ReferralPayout) FromAccount(mods ...qm.QueryMod) accountQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.FromAccountID),
	}

	queryMods = append(queryMods, mods...)

	query := Accounts(queryMods...)
	queries.SetFrom(query.Query, "\"account\"")

	return query
}

// LoadAccount allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (referralPayoutL) LoadAccount(ctx context.Context, e boil.ContextExecutor, singular bool, maybeReferralPayout interface{}, mods queries.Applicator) error {
	var slice []*ReferralPayout
	var object *ReferralPayout

	if singular {
		object = maybeReferralPayout.(*ReferralPayout)
	} else {
		slice = *maybeReferralPayout.(*[]*ReferralPayout)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &referralPayoutR{}
		}
		args = append(args, object.AccountID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &referralPayoutR{}
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
		foreign.R.ReferralPayouts = append(foreign.R.ReferralPayouts, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.AccountID == foreign.ID {
				local.R.Account = foreign
				if foreign.R == nil {
					foreign.R = &accountR{}
				}
				foreign.R.ReferralPayouts = append(foreign.R.ReferralPayouts, local)
				break
			}
		}
	}

	return nil
}

// LoadFromAccount allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (referralPayoutL) LoadFromAccount(ctx context.Context, e boil.ContextExecutor, singular bool, maybeReferralPayout interface{}, mods queries.Applicator) error {
	var slice []*ReferralPayout
	var object *ReferralPayout

	if singular {
		object = maybeReferralPayout.(*ReferralPayout)
	} else {
		slice = *maybeReferralPayout.(*[]*ReferralPayout)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &referralPayoutR{}
		}
		args = append(args, object.FromAccountID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &referralPayoutR{}
			}

			for _, a := range args {
				if a == obj.FromAccountID {
					continue Outer
				}
			}

			args = append(args, obj.FromAccountID)

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
		object.R.FromAccount = foreign
		if foreign.R == nil {
			foreign.R = &accountR{}
		}
		foreign.R.FromAccountReferralPayouts = append(foreign.R.FromAccountReferralPayouts, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.FromAccountID == foreign.ID {
				local.R.FromAccount = foreign
				if foreign.R == nil {
					foreign.R = &accountR{}
				}
				foreign.R.FromAccountReferralPayouts = append(foreign.R.FromAccountReferralPayouts, local)
				break
			}
		}
	}

	return nil
}

// SetAccount of the referralPayout to the related item.
// Sets o.R.Account to related.
// Adds o to related.R.ReferralPayouts.
func (o *ReferralPayout) SetAccount(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Account) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"referral_payout\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"account_id"}),
		strmangle.WhereClause("\"", "\"", 2, referralPayoutPrimaryKeyColumns),
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
		o.R = &referralPayoutR{
			Account: related,
		}
	} else {
		o.R.Account = related
	}

	if related.R == nil {
		related.R = &accountR{
			ReferralPayouts: ReferralPayoutSlice{o},
		}
	} else {
		related.R.ReferralPayouts = append(related.R.ReferralPayouts, o)
	}

	return nil
}

// SetFromAccount of the referralPayout to the related item.
// Sets o.R.FromAccount to related.
// Adds o to related.R.FromAccountReferralPayouts.
func (o *ReferralPayout) SetFromAccount(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Account) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"referral_payout\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"from_account_id"}),
		strmangle.WhereClause("\"", "\"", 2, referralPayoutPrimaryKeyColumns),
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

	o.FromAccountID = related.ID
	if o.R == nil {
		o.R = &referralPayoutR{
			FromAccount: related,
		}
	} else {
		o.R.FromAccount = related
	}

	if related.R == nil {
		related.R = &accountR{
			FromAccountReferralPayouts: ReferralPayoutSlice{o},
		}
	} else {
		related.R.FromAccountReferralPayouts = append(related.R.FromAccountReferralPayouts, o)
	}

	return nil
}

// ReferralPayouts retrieves all the records using an executor.
func ReferralPayouts(mods ...qm.QueryMod) referralPayoutQuery {
	mods = append(mods, qm.From("\"referral_payout\""))
	return referralPayoutQuery{NewQuery(mods...)}
}

// FindReferralPayout retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindReferralPayout(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*ReferralPayout, error) {
	referralPayoutObj := &ReferralPayout{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"referral_payout\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, referralPayoutObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from referral_payout")
	}

	return referralPayoutObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *ReferralPayout) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no referral_payout provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(referralPayoutColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	referralPayoutInsertCacheMut.RLock()
	cache, cached := referralPayoutInsertCache[key]
	referralPayoutInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			referralPayoutAllColumns,
			referralPayoutColumnsWithDefault,
			referralPayoutColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(referralPayoutType, referralPayoutMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(referralPayoutType, referralPayoutMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"referral_payout\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"referral_payout\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into referral_payout")
	}

	if !cached {
		referralPayoutInsertCacheMut.Lock()
		referralPayoutInsertCache[key] = cache
		referralPayoutInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the ReferralPayout.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *ReferralPayout) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	referralPayoutUpdateCacheMut.RLock()
	cache, cached := referralPayoutUpdateCache[key]
	referralPayoutUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			referralPayoutAllColumns,
			referralPayoutPrimaryKeyColumns,
		)

		if len(wl) == 0 {
			return 0, errors.New("models: unable to update referral_payout, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"referral_payout\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, referralPayoutPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(referralPayoutType, referralPayoutMapping, append(wl, referralPayoutPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update referral_payout row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for referral_payout")
	}

	if !cached {
		referralPayoutUpdateCacheMut.Lock()
		referralPayoutUpdateCache[key] = cache
		referralPayoutUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q referralPayoutQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for referral_payout")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for referral_payout")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ReferralPayoutSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), referralPayoutPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"referral_payout\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, referralPayoutPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in referralPayout slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all referralPayout")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *ReferralPayout) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no referral_payout provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(referralPayoutColumnsWithDefault, o)

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

	referralPayoutUpsertCacheMut.RLock()
	cache, cached := referralPayoutUpsertCache[key]
	referralPayoutUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			referralPayoutAllColumns,
			referralPayoutColumnsWithDefault,
			referralPayoutColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			referralPayoutAllColumns,
			referralPayoutPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert referral_payout, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(referralPayoutPrimaryKeyColumns))
			copy(conflict, referralPayoutPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"referral_payout\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(referralPayoutType, referralPayoutMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(referralPayoutType, referralPayoutMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert referral_payout")
	}

	if !cached {
		referralPayoutUpsertCacheMut.Lock()
		referralPayoutUpsertCache[key] = cache
		referralPayoutUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single ReferralPayout record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *ReferralPayout) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no ReferralPayout provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), referralPayoutPrimaryKeyMapping)
	sql := "DELETE FROM \"referral_payout\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from referral_payout")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for referral_payout")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q referralPayoutQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no referralPayoutQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from referral_payout")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for referral_payout")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ReferralPayoutSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), referralPayoutPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"referral_payout\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, referralPayoutPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from referralPayout slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for referral_payout")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *ReferralPayout) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindReferralPayout(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ReferralPayoutSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ReferralPayoutSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), referralPayoutPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"referral_payout\".* FROM \"referral_payout\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, referralPayoutPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in ReferralPayoutSlice")
	}

	*o = slice

	return nil
}

// ReferralPayoutExists checks if the ReferralPayout row exists.
func ReferralPayoutExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"referral_payout\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if referral_payout exists")
	}

	return exists, nil
}
