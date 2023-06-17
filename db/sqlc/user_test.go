package db

import (
	"context"
	"testing"
	"time"

	"github.com/binbomb/goapp/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := utils.HashPassword(utils.RandomString(6))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       utils.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       utils.RandomOwner(),
		Email:          utils.RandomEmail(),
	}
	newUser, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, newUser)
	require.Equal(t, arg.Username, newUser.Username)
	require.Equal(t, arg.HashedPassword, newUser.HashedPassword)
	require.Equal(t, arg.FullName, newUser.FullName)
	require.Equal(t, arg.Email, newUser.Email)
	require.True(t, newUser.PasswordChangedAt.IsZero())
	require.NotZero(t, newUser.CreatedAt)
	return newUser
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.NotZero(t, user2.CreatedAt)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)

}
