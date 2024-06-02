package token

import (
	"github.com/ZenSam7/Education/tools"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(tools.GetRandomString(minSecretKeySize))
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

func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(tools.GetRandomString(minSecretKeySize))
	require.NoError(t, err)

	token, err := maker.CreateToken(tools.GetRandomString(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrorExpiredToken.Error())
	require.Nil(t, payload)
}

func TestJWTTokenInvalidAlg(t *testing.T) {
	maker, err := NewJWTMaker(tools.GetRandomString(minSecretKeySize))
	require.NoError(t, err)

	payload, err := NewPayload(tools.GetRandomString(), time.Minute)
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
