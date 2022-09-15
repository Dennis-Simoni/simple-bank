package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomEntry(t *testing.T, acc Account) Entry {

	arg := CreateEntryParams{
		AccountID: acc.ID,
		Amount:    acc.Balance,
	}

	e, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, e)
	require.Equal(t, acc.ID, e.AccountID)
	require.Equal(t, acc.Balance, e.Amount)
	require.NotEmpty(t, e.CreatedAt)
	require.NotEmpty(t, e.ID)
	return e
}

func TestCreateEntry(t *testing.T) {
	acc := createRandomAccount(t)
	createRandomEntry(t, acc)
}

func TestGetEntry(t *testing.T) {
	acc := createRandomAccount(t)
	e := createRandomEntry(t, acc)

	e2, err := testQueries.GetEntry(context.Background(), e.ID)
	require.NoError(t, err)
	require.Equal(t, e.ID, e2.ID)
	require.Equal(t, e.Amount, e2.Amount)
	require.Equal(t, e.AccountID, e2.AccountID)
	require.WithinDuration(t, e.CreatedAt, e2.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {

	acc := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomEntry(t, acc)
	}

	arg := ListEntriesParams{
		AccountID: acc.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, e := range entries {
		require.NotEmpty(t, e)
	}

}
