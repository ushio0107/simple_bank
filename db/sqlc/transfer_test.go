package db

import (
	"context"
	"simple_bank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, ac1 Account, ac2 Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: ac1.ID,
		ToAccountID:   ac2.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	ac1 := createRandomAccount(t)
	ac2 := createRandomAccount(t)

	createRandomTransfer(t, ac1, ac2)
}

func TestGetTransfer(t *testing.T) {
	ac1 := createRandomAccount(t)
	ac2 := createRandomAccount(t)

	transfer1 := createRandomTransfer(t, ac1, ac2)
	transfer, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, transfer1.ID, transfer.ID)
	require.Equal(t, transfer1.FromAccountID, transfer.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer.Amount)
	require.WithinDuration(t, transfer.CreatedAt, transfer.CreatedAt, time.Second)
}

func TestListTransfer(t *testing.T) {
	ac1 := createRandomAccount(t)
	ac2 := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomTransfer(t, ac1, ac2)
	}

	arg := ListTransfersParams{
		FromAccountID: ac1.ID,
		ToAccountID:   ac2.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
		require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	}

}
