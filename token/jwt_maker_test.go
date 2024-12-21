package token

import (
	"github.com/ZenSam7/Education/tools"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewJWTMaker(t *testing.T) {
	maker := NewJWTMaker(tools.GetRandomString(minSecretKeySize)[:minSecretKeySize])

	randomIDUser := tools.GetRandomInt()
	duration := time.Minute

	token, payload, err := maker.CreateToken(randomIDUser, tools.AdminRole, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = NewPayload(randomIDUser, tools.AdminRole, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotEmpty(t, payload.IDUser)
	require.Equal(t, payload.IDUser, randomIDUser)
	require.Equal(t, payload.Role, tools.AdminRole)
	require.WithinDuration(t, payload.IssuedAt, time.Now(), time.Second)
	require.WithinDuration(t, payload.ExpiredAt, time.Now().Add(duration), time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	maker := NewJWTMaker(tools.GetRandomString(minSecretKeySize))

	token, payload, err := maker.CreateToken(tools.GetRandomInt(), tools.UsualRole, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrorExpiredToken.Error())
	require.Nil(t, payload)
}

func TestJWTTokenInvalidAlg(t *testing.T) {
	maker := NewJWTMaker(tools.GetRandomString(minSecretKeySize))

	payload, err := NewPayload(tools.GetRandomInt(), tools.UsualRole, time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	wrongTokenP := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	wrongToken, err := wrongTokenP.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)
	require.NotEmpty(t, wrongToken)

	wrongPayload, err := maker.VerifyToken(wrongToken)
	require.Error(t, err)
	require.EqualError(t, err, ErrorInvalidToken.Error())
	require.Nil(t, wrongPayload)
}
