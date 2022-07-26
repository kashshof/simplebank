package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	acc1 := createTestAccount(t)
	acc2 := createTestAccount(t)
	fmt.Println(">> before:", acc1.Balance, acc2.Balance)

	//run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: acc1.ID,
				ToAccountId:   acc2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	//check results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, acc1.ID)
		require.Equal(t, transfer.ToAccountID, acc2.ID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//check from entry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.AccountID, acc1.ID)
		require.Equal(t, fromEntry.Amount, -amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		//check to entry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.AccountID, acc2.ID)
		require.Equal(t, toEntry.Amount, amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//check accounts
		fromAcc := result.FromAccount
		require.NotEmpty(t, fromAcc)
		require.Equal(t, fromAcc.ID, acc1.ID)

		toAcc := result.ToAccount
		require.NotEmpty(t, toAcc)
		require.Equal(t, toAcc.ID, acc2.ID)

		//check accounts' balances
		fmt.Println(">> tx:", fromAcc.Balance, toAcc.Balance)
		diff1 := acc1.Balance - fromAcc.Balance
		diff2 := toAcc.Balance - acc2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)
		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	//check final updated balance
	upAcc1, err := store.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	upAcc2, err := store.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", upAcc1.Balance, upAcc2.Balance)

	require.Equal(t, acc1.Balance-int64(n)*amount, upAcc1.Balance)
	require.Equal(t, acc2.Balance+int64(n)*amount, upAcc2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	acc1 := createTestAccount(t)
	acc2 := createTestAccount(t)
	fmt.Println(">> before:", acc1.Balance, acc2.Balance)

	//run n concurrent transfer transactions
	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccId := acc1.ID
		toAccId := acc2.ID
		if i%2 == 1 {
			fromAccId = acc2.ID
			toAccId = acc1.ID
		}
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: fromAccId,
				ToAccountId:   toAccId,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	//check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	//check final updated balance
	upAcc1, err := store.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	upAcc2, err := store.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", upAcc1.Balance, upAcc2.Balance)

	require.Equal(t, acc1.Balance, upAcc1.Balance)
	require.Equal(t, acc2.Balance, upAcc2.Balance)
}
