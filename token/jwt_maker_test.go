package token

import (
	"fmt"
	"testing"
	"time"

	"github.com/binbomb/goapp/simplebank/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)

	username := utils.RandomOwner()
	//issueAt := time.Now()
	duration := time.Minute
	//expiredAt := issueAt.Add(duration)

	token, err := maker.CreateToken(username, -duration)
	fmt.Println("token is created ", token)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)
	require.EqualError(t, err, ErrExpiredToken.Error())
	fmt.Println("payload is created ", payload)
	//
	//require.NotZero(t, payload.RegisteredClaims.ID)
	//require.Equal(t, username, payload.Username)

	// fmt.Println("issueAt is created ", payload.RegisteredClaims.IssuedAt.Time)
	// require.WithinDuration(t, issueAt, payload.RegisteredClaims.IssuedAt.Time, time.Minute)
	// require.WithinDuration(t, expiredAt, payload.RegisteredClaims.ExpiresAt.Time, time.Second)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	payload, err := NewPayload(utils.RandomOwner(), time.Hour)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)
	payload, err = maker.VerifyToken(token)
	fmt.Println("error print test valid ", err)
	require.Error(t, err)

	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)

}
