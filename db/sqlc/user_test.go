package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"simplebank/util"
	"testing"
	"time"
)

func createRandomUser(t *testing.T) User {
	hp, err := util.HashedPassword(util.RandomString(6))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hp,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	u, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, u)

	require.Equal(t, arg.Username, u.Username)
	require.Equal(t, arg.HashedPassword, u.HashedPassword)
	require.Equal(t, arg.FullName, u.FullName)
	require.Equal(t, arg.Email, u.Email)

	require.True(t, u.PasswordChangedAt.IsZero())
	require.NotZero(t, u.CreatedAt)
	return u
}

func TestCreateAUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	want := createRandomUser(t)
	got, err := testQueries.GetUser(context.Background(), want.Username)
	require.NoError(t, err)
	require.NotEmpty(t, got)
	require.Equal(t, want.Username, got.Username)
	require.Equal(t, want.HashedPassword, got.HashedPassword)
	require.Equal(t, want.FullName, got.FullName)
	require.Equal(t, want.Email, got.Email)
	require.WithinDuration(t, want.PasswordChangedAt, got.PasswordChangedAt, time.Second)
	require.WithinDuration(t, want.CreatedAt, got.CreatedAt, time.Second)
}
