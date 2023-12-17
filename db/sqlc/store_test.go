package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	//testStore := NewStore(testDb)

	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)
	fmt.Println(">>before: ", acc1.Balance, acc2.Balance)
	// use concurrency to create transfer transaction
	n := 5
	amount := int64(10)
	errs := make(chan error)
	results := make(chan TransferTxResult)
	for i := 0; i < n; i++ {
		go func() {
			rs, err := testStore.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: acc1.ID,
				ToAccountID:   acc2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- rs
		}()
	}
	exists := make(map[int]bool)
	// check result get from channel
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		rs := <-results
		require.NotEmpty(t, rs)
		// check transfer
		transfer := rs.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, acc1.ID, transfer.FromAccountID)
		require.Equal(t, acc2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = testStore.GetTranfer(context.Background(), transfer.ID)
		require.NoError(t, err)
		// check entries
		fromEntry := rs.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, acc1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		_, err = testStore.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := rs.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, acc2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		_, err = testStore.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check account
		fromAccount := rs.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, acc1.ID)

		toAccount := rs.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, acc2.ID)

		fmt.Println(">>t	x: ", fromAccount.Balance, toAccount.Balance)
		diff1 := acc1.Balance - fromAccount.Balance
		diff2 := -acc2.Balance + toAccount.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, exists, k)
		exists[k] = true
	}

	updateAccount1, err := testStore.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)
	updateAccount2, err := testStore.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	fmt.Println(">>after: ", updateAccount1.Balance, updateAccount2.Balance)
	require.Equal(t, acc1.Balance-int64(n)*amount, updateAccount1.Balance)
	require.Equal(t, acc2.Balance+int64(n)*amount, updateAccount2.Balance)

}
func TestTransferTxDeadlock(t *testing.T) {
	//testStore := NewStore(testDb)

	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)
	fmt.Println(">>before: ", acc1.Balance, acc2.Balance)

	// use concurrency to create transfer transaction
	n := 10
	amount := int64(10)
	errs := make(chan error)
	for i := 0; i < n; i++ {

		// examples for deadlock happens

		fromAccountID := acc1.ID
		toAccountID := acc2.ID

		if i%2 == 1 {
			fromAccountID = acc2.ID
			toAccountID = acc1.ID
		}

		go func() {
			_, err := testStore.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}
	// check result get from channel
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updateAccount1, err := testStore.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)
	updateAccount2, err := testStore.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	fmt.Println(">>after: ", updateAccount1.Balance, updateAccount2.Balance)
	require.Equal(t, acc1.Balance, updateAccount1.Balance)
	require.Equal(t, acc2.Balance, updateAccount2.Balance)

}
