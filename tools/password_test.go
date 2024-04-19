package tools

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetRandomHash(t *testing.T) {
	hash := GetRandomHash()
	require.NotEmpty(t, hash)
}

func TestGetPasswordHash(t *testing.T) {
	password := GetRandomString()
	hash, err := GetPasswordHash(password)

	require.NoError(t, err)
	require.NotEmpty(t, hash)
}

func TestCheckPassword(t *testing.T) {
	password := GetRandomString()
	hash, _ := GetPasswordHash(password)
	check := CheckPassword(password, hash)

	require.True(t, check)
	require.NotEmpty(t, hash)

	// Проверка что хеш не изменяется при том же пароле
	hashAgain, _ := GetPasswordHash(password)
	checkAgain := CheckPassword(password, hash)

	require.True(t, checkAgain)
	require.NotEmpty(t, hashAgain)

	otherPassword := GetRandomString()
	check = CheckPassword(otherPassword, hash)

	require.False(t, check)
	require.NotEmpty(t, otherPassword)

	hashOfHash, err := GetPasswordHash(hash)
	require.NoError(t, err)
	require.NotEmpty(t, hashOfHash)
	require.NotEqual(t, hash, hashOfHash)
}
