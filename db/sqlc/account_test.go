package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/scipiia/snippetbox/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	arg := CreateAccountParams{
		Login:    user.Name,
		Username: util.RandomUser(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Login, account.Login)
	require.Equal(t, arg.Username, account.Username)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.Created)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, acc2)

	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, acc1.Login, acc2.Login)
	require.Equal(t, acc1.Username, acc2.Username)

	require.WithinDuration(t, acc1.Created, acc2.Created, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	acc1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:       acc1.ID,
		Username: util.RandomUser(),
	}

	acc2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, acc2)

	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, acc1.Login, acc2.Login)
	require.Equal(t, arg.Username, acc2.Username)

	require.WithinDuration(t, acc1.Created, acc2.Created, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	acc1 := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.Error(t, err) //специально должны вернуться ошибка
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, acc2)
}
