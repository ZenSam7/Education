package token

import (
	"github.com/ZenSam7/Education/tools"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// Тесты точно такие же как для JWTMaker, но только тут Pasetomaker

func TestNewPasetoMaker(t *testing.T) {
	maker := NewPasetoMaker(tools.GetRandomString(minSecretKeySize)[:minSecretKeySize])

	randomIDUser := tools.GetRandomInt()
	duration := time.Minute

	token, payload, err := maker.CreateToken(randomIDUser, tools.UsualRole, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = NewPayload(randomIDUser, tools.UsualRole, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotEmpty(t, payload.IDUser)
	require.Equal(t, payload.IDUser, randomIDUser)
	require.Equal(t, payload.Role, tools.UsualRole)
	require.WithinDuration(t, payload.IssuedAt, time.Now(), time.Second)
	require.WithinDuration(t, payload.ExpiredAt, time.Now().Add(duration), time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker := NewPasetoMaker(tools.GetRandomString(minSecretKeySize)[:32]) // Строго 32!

	token, payload, err := maker.CreateToken(tools.GetRandomInt(), tools.UsualRole, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrorExpiredToken.Error())
	require.Nil(t, payload)
}
