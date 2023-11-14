package db

import (
	"context"
	"database/sql"
	"simple_bank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccout(t *testing.T) {
	ac1 := createRandomAccount(t)
	ac2, err := testQueries.GetAccount(context.Background(), ac1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, ac2)

	require.Equal(t, ac1.ID, ac2.ID)
	require.Equal(t, ac1.Owner, ac2.Owner)
	require.Equal(t, ac1.Balance, ac2.Balance)
	require.Equal(t, ac1.Currency, ac2.Currency)
	require.WithinDuration(t, ac1.CreatedAt, ac2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	ac1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      ac1.ID,
		Balance: util.RandomMoney(),
	}
	ac2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, ac2)

	require.Equal(t, ac1.ID, ac2.ID)
	require.Equal(t, ac1.Owner, ac2.Owner)
	require.Equal(t, arg.Balance, ac2.Balance)
	require.Equal(t, ac1.Currency, ac2.Currency)
	require.WithinDuration(t, ac1.CreatedAt, ac2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccount(t *testing.T) {
	var ac Account
	for i := 0; i < 10; i++ {
		ac = createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Owner:  ac.Owner,
		Limit:  5,
		Offset: 0,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, ac.Owner, account.Owner)
	}
}
