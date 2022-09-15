package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"simplebank/util"
	"testing"
	"time"
)

func createRandomTransfer(t *testing.T, acc1, acc2 Account) Transfer {
	args := CreateTransferParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        util.RandomBalance(),
	}

	tr, err := testQueries.CreateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, tr)

	require.Equal(t, tr.FromAccountID, acc1.ID)
	require.Equal(t, tr.ToAccountID, acc2.ID)
	require.Equal(t, tr.Amount, args.Amount)

	require.NotZero(t, tr.ID)
	require.NotZero(t, tr.CreatedAt)
	return tr
}

func TestCreateTransfer(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	createRandomTransfer(t, acc1, acc2)
}

func TestGetTransfer(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	tr := createRandomTransfer(t, acc1, acc2)

	tr2, err := testQueries.GetTransfer(context.Background(), tr.ID)
	require.NoError(t, err)
	require.NotEmpty(t, tr2)

	require.Equal(t, tr2.FromAccountID, tr.FromAccountID)
	require.Equal(t, tr2.ToAccountID, tr.ToAccountID)
	require.Equal(t, tr2.Amount, tr.Amount)

	require.WithinDuration(t, tr.CreatedAt, tr2.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {

	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomTransfer(t, acc1, acc2)
	}

	args := ListTransfersParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Limit:         5,
		Offset:        5,
	}

	transactions, err := testQueries.ListTransfers(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, transactions, 5)

	for _, transaction := range transactions {
		require.NotEmpty(t, transaction)
	}
}
