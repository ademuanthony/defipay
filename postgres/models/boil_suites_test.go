// Code generated by SQLBoiler 4.7.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import "testing"

// This test suite runs each operation test in parallel.
// Example, if your database has 3 tables, the suite will run:
// table1, table2 and table3 Delete in parallel
// table1, table2 and table3 Insert in parallel, and so forth.
// It does NOT run each operation group in parallel.
// Separating the tests thusly grants avoidance of Postgres deadlocks.
func TestParent(t *testing.T) {
	t.Run("Accounts", testAccounts)
	t.Run("AccountTransactions", testAccountTransactions)
	t.Run("DailyEarnings", testDailyEarnings)
	t.Run("Deposits", testDeposits)
	t.Run("Packages", testPackages)
	t.Run("Subscriptions", testSubscriptions)
	t.Run("Transfers", testTransfers)
	t.Run("Wallets", testWallets)
	t.Run("Withdrawals", testWithdrawals)
}

func TestDelete(t *testing.T) {
	t.Run("Accounts", testAccountsDelete)
	t.Run("AccountTransactions", testAccountTransactionsDelete)
	t.Run("DailyEarnings", testDailyEarningsDelete)
	t.Run("Deposits", testDepositsDelete)
	t.Run("Packages", testPackagesDelete)
	t.Run("Subscriptions", testSubscriptionsDelete)
	t.Run("Transfers", testTransfersDelete)
	t.Run("Wallets", testWalletsDelete)
	t.Run("Withdrawals", testWithdrawalsDelete)
}

func TestQueryDeleteAll(t *testing.T) {
	t.Run("Accounts", testAccountsQueryDeleteAll)
	t.Run("AccountTransactions", testAccountTransactionsQueryDeleteAll)
	t.Run("DailyEarnings", testDailyEarningsQueryDeleteAll)
	t.Run("Deposits", testDepositsQueryDeleteAll)
	t.Run("Packages", testPackagesQueryDeleteAll)
	t.Run("Subscriptions", testSubscriptionsQueryDeleteAll)
	t.Run("Transfers", testTransfersQueryDeleteAll)
	t.Run("Wallets", testWalletsQueryDeleteAll)
	t.Run("Withdrawals", testWithdrawalsQueryDeleteAll)
}

func TestSliceDeleteAll(t *testing.T) {
	t.Run("Accounts", testAccountsSliceDeleteAll)
	t.Run("AccountTransactions", testAccountTransactionsSliceDeleteAll)
	t.Run("DailyEarnings", testDailyEarningsSliceDeleteAll)
	t.Run("Deposits", testDepositsSliceDeleteAll)
	t.Run("Packages", testPackagesSliceDeleteAll)
	t.Run("Subscriptions", testSubscriptionsSliceDeleteAll)
	t.Run("Transfers", testTransfersSliceDeleteAll)
	t.Run("Wallets", testWalletsSliceDeleteAll)
	t.Run("Withdrawals", testWithdrawalsSliceDeleteAll)
}

func TestExists(t *testing.T) {
	t.Run("Accounts", testAccountsExists)
	t.Run("AccountTransactions", testAccountTransactionsExists)
	t.Run("DailyEarnings", testDailyEarningsExists)
	t.Run("Deposits", testDepositsExists)
	t.Run("Packages", testPackagesExists)
	t.Run("Subscriptions", testSubscriptionsExists)
	t.Run("Transfers", testTransfersExists)
	t.Run("Wallets", testWalletsExists)
	t.Run("Withdrawals", testWithdrawalsExists)
}

func TestFind(t *testing.T) {
	t.Run("Accounts", testAccountsFind)
	t.Run("AccountTransactions", testAccountTransactionsFind)
	t.Run("DailyEarnings", testDailyEarningsFind)
	t.Run("Deposits", testDepositsFind)
	t.Run("Packages", testPackagesFind)
	t.Run("Subscriptions", testSubscriptionsFind)
	t.Run("Transfers", testTransfersFind)
	t.Run("Wallets", testWalletsFind)
	t.Run("Withdrawals", testWithdrawalsFind)
}

func TestBind(t *testing.T) {
	t.Run("Accounts", testAccountsBind)
	t.Run("AccountTransactions", testAccountTransactionsBind)
	t.Run("DailyEarnings", testDailyEarningsBind)
	t.Run("Deposits", testDepositsBind)
	t.Run("Packages", testPackagesBind)
	t.Run("Subscriptions", testSubscriptionsBind)
	t.Run("Transfers", testTransfersBind)
	t.Run("Wallets", testWalletsBind)
	t.Run("Withdrawals", testWithdrawalsBind)
}

func TestOne(t *testing.T) {
	t.Run("Accounts", testAccountsOne)
	t.Run("AccountTransactions", testAccountTransactionsOne)
	t.Run("DailyEarnings", testDailyEarningsOne)
	t.Run("Deposits", testDepositsOne)
	t.Run("Packages", testPackagesOne)
	t.Run("Subscriptions", testSubscriptionsOne)
	t.Run("Transfers", testTransfersOne)
	t.Run("Wallets", testWalletsOne)
	t.Run("Withdrawals", testWithdrawalsOne)
}

func TestAll(t *testing.T) {
	t.Run("Accounts", testAccountsAll)
	t.Run("AccountTransactions", testAccountTransactionsAll)
	t.Run("DailyEarnings", testDailyEarningsAll)
	t.Run("Deposits", testDepositsAll)
	t.Run("Packages", testPackagesAll)
	t.Run("Subscriptions", testSubscriptionsAll)
	t.Run("Transfers", testTransfersAll)
	t.Run("Wallets", testWalletsAll)
	t.Run("Withdrawals", testWithdrawalsAll)
}

func TestCount(t *testing.T) {
	t.Run("Accounts", testAccountsCount)
	t.Run("AccountTransactions", testAccountTransactionsCount)
	t.Run("DailyEarnings", testDailyEarningsCount)
	t.Run("Deposits", testDepositsCount)
	t.Run("Packages", testPackagesCount)
	t.Run("Subscriptions", testSubscriptionsCount)
	t.Run("Transfers", testTransfersCount)
	t.Run("Wallets", testWalletsCount)
	t.Run("Withdrawals", testWithdrawalsCount)
}

func TestInsert(t *testing.T) {
	t.Run("Accounts", testAccountsInsert)
	t.Run("Accounts", testAccountsInsertWhitelist)
	t.Run("AccountTransactions", testAccountTransactionsInsert)
	t.Run("AccountTransactions", testAccountTransactionsInsertWhitelist)
	t.Run("DailyEarnings", testDailyEarningsInsert)
	t.Run("DailyEarnings", testDailyEarningsInsertWhitelist)
	t.Run("Deposits", testDepositsInsert)
	t.Run("Deposits", testDepositsInsertWhitelist)
	t.Run("Packages", testPackagesInsert)
	t.Run("Packages", testPackagesInsertWhitelist)
	t.Run("Subscriptions", testSubscriptionsInsert)
	t.Run("Subscriptions", testSubscriptionsInsertWhitelist)
	t.Run("Transfers", testTransfersInsert)
	t.Run("Transfers", testTransfersInsertWhitelist)
	t.Run("Wallets", testWalletsInsert)
	t.Run("Wallets", testWalletsInsertWhitelist)
	t.Run("Withdrawals", testWithdrawalsInsert)
	t.Run("Withdrawals", testWithdrawalsInsertWhitelist)
}

// TestToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestToOne(t *testing.T) {
	t.Run("AccountTransactionToAccountUsingAccount", testAccountTransactionToOneAccountUsingAccount)
	t.Run("DepositToAccountUsingAccount", testDepositToOneAccountUsingAccount)
	t.Run("SubscriptionToAccountUsingAccount", testSubscriptionToOneAccountUsingAccount)
	t.Run("SubscriptionToPackageUsingPackage", testSubscriptionToOnePackageUsingPackage)
	t.Run("TransferToAccountUsingReceiver", testTransferToOneAccountUsingReceiver)
	t.Run("TransferToAccountUsingSender", testTransferToOneAccountUsingSender)
	t.Run("WalletToAccountUsingAccount", testWalletToOneAccountUsingAccount)
	t.Run("WithdrawalToAccountUsingAccount", testWithdrawalToOneAccountUsingAccount)
}

// TestOneToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOne(t *testing.T) {}

// TestToMany tests cannot be run in parallel
// or deadlocks can occur.
func TestToMany(t *testing.T) {
	t.Run("AccountToAccountTransactions", testAccountToManyAccountTransactions)
	t.Run("AccountToDeposits", testAccountToManyDeposits)
	t.Run("AccountToSubscriptions", testAccountToManySubscriptions)
	t.Run("AccountToReceiverTransfers", testAccountToManyReceiverTransfers)
	t.Run("AccountToSenderTransfers", testAccountToManySenderTransfers)
	t.Run("AccountToWallets", testAccountToManyWallets)
	t.Run("AccountToWithdrawals", testAccountToManyWithdrawals)
	t.Run("PackageToSubscriptions", testPackageToManySubscriptions)
}

// TestToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneSet(t *testing.T) {
	t.Run("AccountTransactionToAccountUsingAccountTransactions", testAccountTransactionToOneSetOpAccountUsingAccount)
	t.Run("DepositToAccountUsingDeposits", testDepositToOneSetOpAccountUsingAccount)
	t.Run("SubscriptionToAccountUsingSubscriptions", testSubscriptionToOneSetOpAccountUsingAccount)
	t.Run("SubscriptionToPackageUsingSubscriptions", testSubscriptionToOneSetOpPackageUsingPackage)
	t.Run("TransferToAccountUsingReceiverTransfers", testTransferToOneSetOpAccountUsingReceiver)
	t.Run("TransferToAccountUsingSenderTransfers", testTransferToOneSetOpAccountUsingSender)
	t.Run("WalletToAccountUsingWallets", testWalletToOneSetOpAccountUsingAccount)
	t.Run("WithdrawalToAccountUsingWithdrawals", testWithdrawalToOneSetOpAccountUsingAccount)
}

// TestToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneRemove(t *testing.T) {
	t.Run("TransferToAccountUsingReceiverTransfers", testTransferToOneRemoveOpAccountUsingReceiver)
	t.Run("TransferToAccountUsingSenderTransfers", testTransferToOneRemoveOpAccountUsingSender)
}

// TestOneToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneSet(t *testing.T) {}

// TestOneToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneRemove(t *testing.T) {}

// TestToManyAdd tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyAdd(t *testing.T) {
	t.Run("AccountToAccountTransactions", testAccountToManyAddOpAccountTransactions)
	t.Run("AccountToDeposits", testAccountToManyAddOpDeposits)
	t.Run("AccountToSubscriptions", testAccountToManyAddOpSubscriptions)
	t.Run("AccountToReceiverTransfers", testAccountToManyAddOpReceiverTransfers)
	t.Run("AccountToSenderTransfers", testAccountToManyAddOpSenderTransfers)
	t.Run("AccountToWallets", testAccountToManyAddOpWallets)
	t.Run("AccountToWithdrawals", testAccountToManyAddOpWithdrawals)
	t.Run("PackageToSubscriptions", testPackageToManyAddOpSubscriptions)
}

// TestToManySet tests cannot be run in parallel
// or deadlocks can occur.
func TestToManySet(t *testing.T) {
	t.Run("AccountToReceiverTransfers", testAccountToManySetOpReceiverTransfers)
	t.Run("AccountToSenderTransfers", testAccountToManySetOpSenderTransfers)
}

// TestToManyRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyRemove(t *testing.T) {
	t.Run("AccountToReceiverTransfers", testAccountToManyRemoveOpReceiverTransfers)
	t.Run("AccountToSenderTransfers", testAccountToManyRemoveOpSenderTransfers)
}

func TestReload(t *testing.T) {
	t.Run("Accounts", testAccountsReload)
	t.Run("AccountTransactions", testAccountTransactionsReload)
	t.Run("DailyEarnings", testDailyEarningsReload)
	t.Run("Deposits", testDepositsReload)
	t.Run("Packages", testPackagesReload)
	t.Run("Subscriptions", testSubscriptionsReload)
	t.Run("Transfers", testTransfersReload)
	t.Run("Wallets", testWalletsReload)
	t.Run("Withdrawals", testWithdrawalsReload)
}

func TestReloadAll(t *testing.T) {
	t.Run("Accounts", testAccountsReloadAll)
	t.Run("AccountTransactions", testAccountTransactionsReloadAll)
	t.Run("DailyEarnings", testDailyEarningsReloadAll)
	t.Run("Deposits", testDepositsReloadAll)
	t.Run("Packages", testPackagesReloadAll)
	t.Run("Subscriptions", testSubscriptionsReloadAll)
	t.Run("Transfers", testTransfersReloadAll)
	t.Run("Wallets", testWalletsReloadAll)
	t.Run("Withdrawals", testWithdrawalsReloadAll)
}

func TestSelect(t *testing.T) {
	t.Run("Accounts", testAccountsSelect)
	t.Run("AccountTransactions", testAccountTransactionsSelect)
	t.Run("DailyEarnings", testDailyEarningsSelect)
	t.Run("Deposits", testDepositsSelect)
	t.Run("Packages", testPackagesSelect)
	t.Run("Subscriptions", testSubscriptionsSelect)
	t.Run("Transfers", testTransfersSelect)
	t.Run("Wallets", testWalletsSelect)
	t.Run("Withdrawals", testWithdrawalsSelect)
}

func TestUpdate(t *testing.T) {
	t.Run("Accounts", testAccountsUpdate)
	t.Run("AccountTransactions", testAccountTransactionsUpdate)
	t.Run("DailyEarnings", testDailyEarningsUpdate)
	t.Run("Deposits", testDepositsUpdate)
	t.Run("Packages", testPackagesUpdate)
	t.Run("Subscriptions", testSubscriptionsUpdate)
	t.Run("Transfers", testTransfersUpdate)
	t.Run("Wallets", testWalletsUpdate)
	t.Run("Withdrawals", testWithdrawalsUpdate)
}

func TestSliceUpdateAll(t *testing.T) {
	t.Run("Accounts", testAccountsSliceUpdateAll)
	t.Run("AccountTransactions", testAccountTransactionsSliceUpdateAll)
	t.Run("DailyEarnings", testDailyEarningsSliceUpdateAll)
	t.Run("Deposits", testDepositsSliceUpdateAll)
	t.Run("Packages", testPackagesSliceUpdateAll)
	t.Run("Subscriptions", testSubscriptionsSliceUpdateAll)
	t.Run("Transfers", testTransfersSliceUpdateAll)
	t.Run("Wallets", testWalletsSliceUpdateAll)
	t.Run("Withdrawals", testWithdrawalsSliceUpdateAll)
}
