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
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// PaymentLink is an object representing the database table.
type PaymentLink struct {
	Permalink     string      `boil:"permalink" json:"permalink" toml:"permalink" yaml:"permalink"`
	AccountID     null.String `boil:"account_id" json:"account_id,omitempty" toml:"account_id" yaml:"account_id,omitempty"`
	Email         string      `boil:"email" json:"email" toml:"email" yaml:"email"`
	Accountname   string      `boil:"accountname" json:"accountname" toml:"accountname" yaml:"accountname"`
	Accountnumber string      `boil:"accountnumber" json:"accountnumber" toml:"accountnumber" yaml:"accountnumber"`
	Bankname      string      `boil:"bankname" json:"bankname" toml:"bankname" yaml:"bankname"`
	Fixamount     int64       `boil:"fixamount" json:"fixamount" toml:"fixamount" yaml:"fixamount"`
	Title         string      `boil:"title" json:"title" toml:"title" yaml:"title"`
	Description   string      `boil:"description" json:"description" toml:"description" yaml:"description"`

	R *paymentLinkR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L paymentLinkL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var PaymentLinkColumns = struct {
	Permalink     string
	AccountID     string
	Email         string
	Accountname   string
	Accountnumber string
	Bankname      string
	Fixamount     string
	Title         string
	Description   string
}{
	Permalink:     "permalink",
	AccountID:     "account_id",
	Email:         "email",
	Accountname:   "accountname",
	Accountnumber: "accountnumber",
	Bankname:      "bankname",
	Fixamount:     "fixamount",
	Title:         "title",
	Description:   "description",
}

var PaymentLinkTableColumns = struct {
	Permalink     string
	AccountID     string
	Email         string
	Accountname   string
	Accountnumber string
	Bankname      string
	Fixamount     string
	Title         string
	Description   string
}{
	Permalink:     "payment_link.permalink",
	AccountID:     "payment_link.account_id",
	Email:         "payment_link.email",
	Accountname:   "payment_link.accountname",
	Accountnumber: "payment_link.accountnumber",
	Bankname:      "payment_link.bankname",
	Fixamount:     "payment_link.fixamount",
	Title:         "payment_link.title",
	Description:   "payment_link.description",
}

// Generated where

var PaymentLinkWhere = struct {
	Permalink     whereHelperstring
	AccountID     whereHelpernull_String
	Email         whereHelperstring
	Accountname   whereHelperstring
	Accountnumber whereHelperstring
	Bankname      whereHelperstring
	Fixamount     whereHelperint64
	Title         whereHelperstring
	Description   whereHelperstring
}{
	Permalink:     whereHelperstring{field: "\"payment_link\".\"permalink\""},
	AccountID:     whereHelpernull_String{field: "\"payment_link\".\"account_id\""},
	Email:         whereHelperstring{field: "\"payment_link\".\"email\""},
	Accountname:   whereHelperstring{field: "\"payment_link\".\"accountname\""},
	Accountnumber: whereHelperstring{field: "\"payment_link\".\"accountnumber\""},
	Bankname:      whereHelperstring{field: "\"payment_link\".\"bankname\""},
	Fixamount:     whereHelperint64{field: "\"payment_link\".\"fixamount\""},
	Title:         whereHelperstring{field: "\"payment_link\".\"title\""},
	Description:   whereHelperstring{field: "\"payment_link\".\"description\""},
}

// PaymentLinkRels is where relationship names are stored.
var PaymentLinkRels = struct {
	Account string
}{
	Account: "Account",
}

// paymentLinkR is where relationships are stored.
type paymentLinkR struct {
	Account *Account `boil:"Account" json:"Account" toml:"Account" yaml:"Account"`
}

// NewStruct creates a new relationship struct
func (*paymentLinkR) NewStruct() *paymentLinkR {
	return &paymentLinkR{}
}

// paymentLinkL is where Load methods for each relationship are stored.
type paymentLinkL struct{}

var (
	paymentLinkAllColumns            = []string{"permalink", "account_id", "email", "accountname", "accountnumber", "bankname", "fixamount", "title", "description"}
	paymentLinkColumnsWithoutDefault = []string{"permalink", "account_id", "email", "accountname", "accountnumber", "bankname", "fixamount", "title", "description"}
	paymentLinkColumnsWithDefault    = []string{}
	paymentLinkPrimaryKeyColumns     = []string{"permalink"}
)

type (
	// PaymentLinkSlice is an alias for a slice of pointers to PaymentLink.
	// This should almost always be used instead of []PaymentLink.
	PaymentLinkSlice []*PaymentLink

	paymentLinkQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	paymentLinkType                 = reflect.TypeOf(&PaymentLink{})
	paymentLinkMapping              = queries.MakeStructMapping(paymentLinkType)
	paymentLinkPrimaryKeyMapping, _ = queries.BindMapping(paymentLinkType, paymentLinkMapping, paymentLinkPrimaryKeyColumns)
	paymentLinkInsertCacheMut       sync.RWMutex
	paymentLinkInsertCache          = make(map[string]insertCache)
	paymentLinkUpdateCacheMut       sync.RWMutex
	paymentLinkUpdateCache          = make(map[string]updateCache)
	paymentLinkUpsertCacheMut       sync.RWMutex
	paymentLinkUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single paymentLink record from the query.
func (q paymentLinkQuery) One(ctx context.Context, exec boil.ContextExecutor) (*PaymentLink, error) {
	o := &PaymentLink{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for payment_link")
	}

	return o, nil
}

// All returns all PaymentLink records from the query.
func (q paymentLinkQuery) All(ctx context.Context, exec boil.ContextExecutor) (PaymentLinkSlice, error) {
	var o []*PaymentLink

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to PaymentLink slice")
	}

	return o, nil
}

// Count returns the count of all PaymentLink records in the query.
func (q paymentLinkQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count payment_link rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q paymentLinkQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if payment_link exists")
	}

	return count > 0, nil
}

// Account pointed to by the foreign key.
func (o *PaymentLink) Account(mods ...qm.QueryMod) accountQuery {
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
func (paymentLinkL) LoadAccount(ctx context.Context, e boil.ContextExecutor, singular bool, maybePaymentLink interface{}, mods queries.Applicator) error {
	var slice []*PaymentLink
	var object *PaymentLink

	if singular {
		object = maybePaymentLink.(*PaymentLink)
	} else {
		slice = *maybePaymentLink.(*[]*PaymentLink)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &paymentLinkR{}
		}
		if !queries.IsNil(object.AccountID) {
			args = append(args, object.AccountID)
		}

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &paymentLinkR{}
			}

			for _, a := range args {
				if queries.Equal(a, obj.AccountID) {
					continue Outer
				}
			}

			if !queries.IsNil(obj.AccountID) {
				args = append(args, obj.AccountID)
			}

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
		foreign.R.PaymentLinks = append(foreign.R.PaymentLinks, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if queries.Equal(local.AccountID, foreign.ID) {
				local.R.Account = foreign
				if foreign.R == nil {
					foreign.R = &accountR{}
				}
				foreign.R.PaymentLinks = append(foreign.R.PaymentLinks, local)
				break
			}
		}
	}

	return nil
}

// SetAccount of the paymentLink to the related item.
// Sets o.R.Account to related.
// Adds o to related.R.PaymentLinks.
func (o *PaymentLink) SetAccount(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Account) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"payment_link\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"account_id"}),
		strmangle.WhereClause("\"", "\"", 2, paymentLinkPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.Permalink}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	queries.Assign(&o.AccountID, related.ID)
	if o.R == nil {
		o.R = &paymentLinkR{
			Account: related,
		}
	} else {
		o.R.Account = related
	}

	if related.R == nil {
		related.R = &accountR{
			PaymentLinks: PaymentLinkSlice{o},
		}
	} else {
		related.R.PaymentLinks = append(related.R.PaymentLinks, o)
	}

	return nil
}

// RemoveAccount relationship.
// Sets o.R.Account to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *PaymentLink) RemoveAccount(ctx context.Context, exec boil.ContextExecutor, related *Account) error {
	var err error

	queries.SetScanner(&o.AccountID, nil)
	if _, err = o.Update(ctx, exec, boil.Whitelist("account_id")); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	if o.R != nil {
		o.R.Account = nil
	}
	if related == nil || related.R == nil {
		return nil
	}

	for i, ri := range related.R.PaymentLinks {
		if queries.Equal(o.AccountID, ri.AccountID) {
			continue
		}

		ln := len(related.R.PaymentLinks)
		if ln > 1 && i < ln-1 {
			related.R.PaymentLinks[i] = related.R.PaymentLinks[ln-1]
		}
		related.R.PaymentLinks = related.R.PaymentLinks[:ln-1]
		break
	}
	return nil
}

// PaymentLinks retrieves all the records using an executor.
func PaymentLinks(mods ...qm.QueryMod) paymentLinkQuery {
	mods = append(mods, qm.From("\"payment_link\""))
	return paymentLinkQuery{NewQuery(mods...)}
}

// FindPaymentLink retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindPaymentLink(ctx context.Context, exec boil.ContextExecutor, permalink string, selectCols ...string) (*PaymentLink, error) {
	paymentLinkObj := &PaymentLink{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"payment_link\" where \"permalink\"=$1", sel,
	)

	q := queries.Raw(query, permalink)

	err := q.Bind(ctx, exec, paymentLinkObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from payment_link")
	}

	return paymentLinkObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *PaymentLink) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no payment_link provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(paymentLinkColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	paymentLinkInsertCacheMut.RLock()
	cache, cached := paymentLinkInsertCache[key]
	paymentLinkInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			paymentLinkAllColumns,
			paymentLinkColumnsWithDefault,
			paymentLinkColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(paymentLinkType, paymentLinkMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(paymentLinkType, paymentLinkMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"payment_link\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"payment_link\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into payment_link")
	}

	if !cached {
		paymentLinkInsertCacheMut.Lock()
		paymentLinkInsertCache[key] = cache
		paymentLinkInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the PaymentLink.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *PaymentLink) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	paymentLinkUpdateCacheMut.RLock()
	cache, cached := paymentLinkUpdateCache[key]
	paymentLinkUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			paymentLinkAllColumns,
			paymentLinkPrimaryKeyColumns,
		)

		if len(wl) == 0 {
			return 0, errors.New("models: unable to update payment_link, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"payment_link\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, paymentLinkPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(paymentLinkType, paymentLinkMapping, append(wl, paymentLinkPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update payment_link row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for payment_link")
	}

	if !cached {
		paymentLinkUpdateCacheMut.Lock()
		paymentLinkUpdateCache[key] = cache
		paymentLinkUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q paymentLinkQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for payment_link")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for payment_link")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o PaymentLinkSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), paymentLinkPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"payment_link\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, paymentLinkPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in paymentLink slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all paymentLink")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *PaymentLink) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no payment_link provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(paymentLinkColumnsWithDefault, o)

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

	paymentLinkUpsertCacheMut.RLock()
	cache, cached := paymentLinkUpsertCache[key]
	paymentLinkUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			paymentLinkAllColumns,
			paymentLinkColumnsWithDefault,
			paymentLinkColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			paymentLinkAllColumns,
			paymentLinkPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert payment_link, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(paymentLinkPrimaryKeyColumns))
			copy(conflict, paymentLinkPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"payment_link\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(paymentLinkType, paymentLinkMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(paymentLinkType, paymentLinkMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert payment_link")
	}

	if !cached {
		paymentLinkUpsertCacheMut.Lock()
		paymentLinkUpsertCache[key] = cache
		paymentLinkUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single PaymentLink record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *PaymentLink) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no PaymentLink provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), paymentLinkPrimaryKeyMapping)
	sql := "DELETE FROM \"payment_link\" WHERE \"permalink\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from payment_link")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for payment_link")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q paymentLinkQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no paymentLinkQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from payment_link")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for payment_link")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o PaymentLinkSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), paymentLinkPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"payment_link\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, paymentLinkPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from paymentLink slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for payment_link")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *PaymentLink) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindPaymentLink(ctx, exec, o.Permalink)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *PaymentLinkSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := PaymentLinkSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), paymentLinkPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"payment_link\".* FROM \"payment_link\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, paymentLinkPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in PaymentLinkSlice")
	}

	*o = slice

	return nil
}

// PaymentLinkExists checks if the PaymentLink row exists.
func PaymentLinkExists(ctx context.Context, exec boil.ContextExecutor, permalink string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"payment_link\" where \"permalink\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, permalink)
	}
	row := exec.QueryRowContext(ctx, sql, permalink)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if payment_link exists")
	}

	return exists, nil
}
