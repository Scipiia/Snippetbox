package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/scipiia/snippetbox/util"
	"github.com/stretchr/testify/require"
)

func createRandomSnippet(t *testing.T, user User) Snippet {
	arg := CreateSnippetParams{
		UserID:  user.ID,
		Title:   util.RandomTitle(),
		Content: util.RandomContent(),
	}

	snippet, err := testQueries.CreateSnippet(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, snippet)

	require.Equal(t, arg.UserID, snippet.UserID)
	require.Equal(t, arg.Title, snippet.Title)
	require.Equal(t, arg.Content, snippet.Content)

	require.NotZero(t, snippet.ID)
	require.NotZero(t, snippet.Created)

	return snippet
}

func TestCreateSnippet(t *testing.T) {
	user := createRandomUser(t)
	createRandomSnippet(t, user)
}

func TestGetSnippet(t *testing.T) {
	user := createRandomUser(t)
	snippet1 := createRandomSnippet(t, user)
	snippet2, err := testQueries.GetSnippet(context.Background(), snippet1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, snippet2)

	require.Equal(t, snippet1.ID, snippet2.ID)
	require.Equal(t, snippet1.UserID, snippet2.UserID)
	require.Equal(t, snippet1.Title, snippet2.Title)
	require.Equal(t, snippet1.Content, snippet2.Content)
	require.WithinDuration(t, snippet1.Created, snippet2.Created, time.Second)
}

func TestListSnippet(t *testing.T) {
	user := createRandomUser(t)
	for i := 0; i < 10; i++ {
		createRandomSnippet(t, user)
	}

	arg := ListSnippetsParams{
		UserID: user.ID,
		Limit:  5,
		Offset: 5,
	}

	snippets, err := testQueries.ListSnippets(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, snippets, 5)

	for _, snippet := range snippets {
		require.NotEmpty(t, snippet)
		require.Equal(t, arg.UserID, snippet.UserID)
	}
}

func TestDeleteSnippet(t *testing.T) {
	user := createRandomUser(t)
	snippet1 := createRandomSnippet(t, user)

	err := testQueries.DeleteSnippet(context.Background(), snippet1.ID)
	require.NoError(t, err)

	snippet2, err := testQueries.GetSnippet(context.Background(), snippet1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, snippet2)
}
