package util

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestHashedPassword(t *testing.T) {
	psw := RandomString(6)

	hp, err := HashedPassword(psw)
	require.NoError(t, err)

	require.NotEmpty(t, hp)

	err = CheckPassword([]byte(hp), []byte(psw))
	require.NoError(t, err)

	wrongPsw := RandomString(6)
	err = CheckPassword([]byte(hp), []byte(wrongPsw))
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
