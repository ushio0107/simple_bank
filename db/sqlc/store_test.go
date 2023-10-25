package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	ac1 := createRandomAccount(t)
	ac2 := createRandomAccount(t)

	n := 2
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			// Concurrency
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: ac1.ID,
				ToAccountID:   ac2.ID,
				Amount:        amount,
			})

			// Send the error and result to goroutine by channel.
			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)
	// Check results
	for i := 0; i < n; i++ {
		// Take out
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, ac1.ID, transfer.FromAccountID)
		require.Equal(t, ac2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)

		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// Check Entries.
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, ac1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, ac2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// Check accounts.
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, ac1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, ac2.ID, toAccount.ID)

		// org ac1 Balance - after transfer ac1 Balance = amount
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)
		diff1 := ac1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - ac2.Balance
		require.Equal(t, diff1, diff2)
		// The balance of ac1 before transfer should be larger than after transfer.
		require.True(t, diff1 > 0)

		// The balance of account 1 will be decreased by 1 times amount after the 1st transaction,
		// then 2 times amount after the 2nd transaction,
		// 3 times amount after the 3rd transaction,
		// and so on and so forth.
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedAc1, err := testQueries.GetAccount(context.Background(), ac1.ID)
	require.NoError(t, err)

	updatedAc2, err := testQueries.GetAccount(context.Background(), ac2.ID)
	require.NoError(t, err)

	// The account transfer amount(10) for 5 times, so the balance of ac1 should decrease 5*10.
	// FromAccount
	require.Equal(t, ac1.Balance-int64(n)*amount, updatedAc1.Balance)
	// ToAccount
	require.Equal(t, ac2.Balance+int64(n)*amount, updatedAc2.Balance)

}
