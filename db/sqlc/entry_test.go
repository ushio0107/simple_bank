package db

import (
	"context"
	"simple_bank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, ac Account) Entry {
	arg := CreateEntryParams{
		AccountID: ac.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.CreatedAt)
	require.NotZero(t, entry.ID)

	return entry
}

func TestCreateEntry(t *testing.T) {
	ac := createRandomAccount(t)
	createRandomEntry(t, ac)
}

func TestGetEntry(t *testing.T) {
	ac := createRandomAccount(t)
	e1 := createRandomEntry(t, ac)
	entry, err := testQueries.GetEntry(context.Background(), e1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, e1.ID, entry.ID)
	require.Equal(t, e1.AccountID, entry.AccountID)
	require.Equal(t, e1.Amount, entry.Amount)
	require.WithinDuration(t, e1.CreatedAt, entry.CreatedAt, time.Second)
}

func TestListEntry(t *testing.T) {
	ac := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t, ac)
	}
	arg := ListEntriesParams{
		AccountID: ac.ID,
		Limit:     5,
		Offset:    5,
	}
	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entries)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, arg.AccountID, entry.AccountID)
	}
}
