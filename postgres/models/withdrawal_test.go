// Code generated by SQLBoiler 4.7.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/volatiletech/randomize"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/strmangle"
)

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testWithdrawals(t *testing.T) {
	t.Parallel()

	query := Withdrawals()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testWithdrawalsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Withdrawal{}
	if err = randomize.Struct(seed, o, withdrawalDBTypes, true, withdrawalColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Withdrawal struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := o.Delete(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Withdrawals().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testWithdrawalsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Withdrawal{}
	if err = randomize.Struct(seed, o, withdrawalDBTypes, true, withdrawalColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Withdrawal struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := Withdrawals().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Withdrawals().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testWithdrawalsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Withdrawal{}
	if err = randomize.Struct(seed, o, withdrawalDBTypes, true, withdrawalColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Withdrawal struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := WithdrawalSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Withdrawals().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testWithdrawalsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Withdrawal{}
	if err = randomize.Struct(seed, o, withdrawalDBTypes, true, withdrawalColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Withdrawal struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := WithdrawalExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if Withdrawal exists: %s", err)
	}
	if !e {
		t.Errorf("Expected WithdrawalExists to return true, but got false.")
	}
}

func testWithdrawalsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Withdrawal{}
	if err = randomize.Struct(seed, o, withdrawalDBTypes, true, withdrawalColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Withdrawal struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	withdrawalFound, err := FindWithdrawal(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if withdrawalFound == nil {
		t.Error("want a record, got nil")
	}
}

func testWithdrawalsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Withdrawal{}
	if err = randomize.Struct(seed, o, withdrawalDBTypes, true, withdrawalColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Withdrawal struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = Withdrawals().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testWithdrawalsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Withdrawal{}
	if err = randomize.Struct(seed, o, withdrawalDBTypes, true, withdrawalColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Withdrawal struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := Withdrawals().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testWithdrawalsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	withdrawalOne := &Withdrawal{}
	withdrawalTwo := &Withdrawal{}
	if err = randomize.Struct(seed, withdrawalOne, withdrawalDBTypes, false, withdrawalColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Withdrawal struct: %s", err)
	}
	if err = randomize.Struct(seed, withdrawalTwo, withdrawalDBTypes, false, withdrawalColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Withdrawal struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = withdrawalOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = withdrawalTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Withdrawals().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testWithdrawalsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	withdrawalOne := &Withdrawal{}
	withdrawalTwo := &Withdrawal{}
	if err = randomize.Struct(seed, withdrawalOne, withdrawalDBTypes, false, withdrawalColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Withdrawal struct: %s", err)
	}
	if err = randomize.Struct(seed, withdrawalTwo, withdrawalDBTypes, false, withdrawalColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Withdrawal struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = withdrawalOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = withdrawalTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Withdrawals().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testWithdrawalsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Withdrawal{}
	if err = randomize.Struct(seed, o, withdrawalDBTypes, true, withdrawalColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Withdrawal struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Withdrawals().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testWithdrawalsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Withdrawal{}
	if err = randomize.Struct(seed, o, withdrawalDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Withdrawal struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(withdrawalColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := Withdrawals().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testWithdrawalToOneAccountUsingAccount(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local Withdrawal
	var foreign Account

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, withdrawalDBTypes, false, withdrawalColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Withdrawal struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, accountDBTypes, false, accountColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Account struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.AccountID = foreign.ID
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.Account().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := WithdrawalSlice{&local}
	if err = local.L.LoadAccount(ctx, tx, false, (*[]*Withdrawal)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Account == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Account = nil
	if err = local.L.LoadAccount(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Account == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testWithdrawalToOneSetOpAccountUsingAccount(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Withdrawal
	var b, c Account

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, withdrawalDBTypes, false, strmangle.SetComplement(withdrawalPrimaryKeyColumns, withdrawalColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, accountDBTypes, false, strmangle.SetComplement(accountPrimaryKeyColumns, accountColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, accountDBTypes, false, strmangle.SetComplement(accountPrimaryKeyColumns, accountColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Account{&b, &c} {
		err = a.SetAccount(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Account != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.Withdrawals[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.AccountID != x.ID {
			t.Error("foreign key was wrong value", a.AccountID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.AccountID))
		reflect.Indirect(reflect.ValueOf(&a.AccountID)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.AccountID != x.ID {
			t.Error("foreign key was wrong value", a.AccountID, x.ID)
		}
	}
}

func testWithdrawalsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Withdrawal{}
	if err = randomize.Struct(seed, o, withdrawalDBTypes, true, withdrawalColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Withdrawal struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = o.Reload(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testWithdrawalsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Withdrawal{}
	if err = randomize.Struct(seed, o, withdrawalDBTypes, true, withdrawalColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Withdrawal struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := WithdrawalSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testWithdrawalsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Withdrawal{}
	if err = randomize.Struct(seed, o, withdrawalDBTypes, true, withdrawalColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Withdrawal struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Withdrawals().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	withdrawalDBTypes = map[string]string{`ID`: `character varying`, `AccountID`: `character varying`, `Amount`: `bigint`, `Date`: `bigint`, `Destination`: `character varying`, `Ref`: `character varying`, `Status`: `character varying`}
	_                 = bytes.MinRead
)

func testWithdrawalsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(withdrawalPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(withdrawalAllColumns) == len(withdrawalPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Withdrawal{}
	if err = randomize.Struct(seed, o, withdrawalDBTypes, true, withdrawalColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Withdrawal struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Withdrawals().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, withdrawalDBTypes, true, withdrawalPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Withdrawal struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testWithdrawalsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(withdrawalAllColumns) == len(withdrawalPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Withdrawal{}
	if err = randomize.Struct(seed, o, withdrawalDBTypes, true, withdrawalColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Withdrawal struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Withdrawals().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, withdrawalDBTypes, true, withdrawalPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Withdrawal struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(withdrawalAllColumns, withdrawalPrimaryKeyColumns) {
		fields = withdrawalAllColumns
	} else {
		fields = strmangle.SetComplement(
			withdrawalAllColumns,
			withdrawalPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	typ := reflect.TypeOf(o).Elem()
	n := typ.NumField()

	updateMap := M{}
	for _, col := range fields {
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.Tag.Get("boil") == col {
				updateMap[col] = value.Field(i).Interface()
			}
		}
	}

	slice := WithdrawalSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testWithdrawalsUpsert(t *testing.T) {
	t.Parallel()

	if len(withdrawalAllColumns) == len(withdrawalPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := Withdrawal{}
	if err = randomize.Struct(seed, &o, withdrawalDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Withdrawal struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Withdrawal: %s", err)
	}

	count, err := Withdrawals().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, withdrawalDBTypes, false, withdrawalPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Withdrawal struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Withdrawal: %s", err)
	}

	count, err = Withdrawals().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
