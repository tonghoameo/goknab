package token

import (
	"fmt"
	"testing"
	"time"

	"github.com/binbomb/goapp/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(utils.RandomString(32))
	require.NoError(t, err)

	username := utils.RandomOwner()
	//issueAt := time.Now()
	duration := time.Minute
	//expiredAt := issueAt.Add(duration)

	token, payload, err := maker.CreateToken(username, -duration)
	fmt.Println("token is created ", token)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)
	require.EqualError(t, err, ErrExpiredToken.Error())
	fmt.Println("payload is created ", payload)
}
func TestExpiredPasetoToken(t *testing.T) {

	maker, err := NewPasetoMaker(utils.RandomString(32))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(utils.RandomOwner(), -time.Minute)
	fmt.Println("token is created ", token)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)
	require.EqualError(t, err, ErrExpiredToken.Error())
}
func TestInvalidPasetoToken(t *testing.T) {

	maker, err := NewPasetoMaker(utils.RandomString(32))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(utils.RandomOwner(), time.Minute)
	fmt.Println("token is created ", token)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)
	require.EqualError(t, err, ErrInvalidToken.Error())
}
