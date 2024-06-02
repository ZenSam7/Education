package token

import (
	"github.com/ZenSam7/Education/tools"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// Тесты точно такие же как для JWTMaker, но только тут Pasetomaker

func TestNewPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(tools.GetRandomString(minSecretKeySize)[:32]) // Строго 32!
	require.NoError(t, err)

	randomUsername := tools.GetRandomString()
	duration := time.Minute

	token, err := maker.CreateToken(randomUsername, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := NewPayload(randomUsername, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotEmpty(t, payload.ID)
	require.Equal(t, payload.Username, randomUsername)
	require.WithinDuration(t, payload.IssuedAt, time.Now(), time.Second)
	require.WithinDuration(t, payload.ExpiredAt, time.Now().Add(duration), time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(tools.GetRandomString(minSecretKeySize)[:32]) // Строго 32!
	require.NoError(t, err)

	token, err := maker.CreateToken(tools.GetRandomString(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrorExpiredToken.Error())
	require.Nil(t, payload)
}
