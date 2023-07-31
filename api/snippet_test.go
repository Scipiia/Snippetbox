package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/scipiia/snippetbox/db/mock"
	db "github.com/scipiia/snippetbox/db/sqlc"
	"github.com/scipiia/snippetbox/util"
	"github.com/stretchr/testify/require"
)

func TestCreateSbippetAPI(t *testing.T) {
	snippet := randomSnippet()

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"account_id": snippet.AccountID,
				"title":      snippet.Title,
				"content":    snippet.Content,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateSnippetParams{
					AccountID: snippet.AccountID,
					Title:     snippet.Title,
					Content:   snippet.Content,
				}

				store.EXPECT().
					CreateSnippet(gomock.Any(), gomock.Eq(arg)).Times(1).Return(snippet, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchSnippet(t, recorder.Body, snippet)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"account_id": snippet.AccountID,
				"title":      snippet.Title,
				"content":    snippet.Content,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateSnippetParams{
					AccountID: snippet.AccountID,
					Title:     snippet.Title,
					Content:   snippet.Content,
				}

				store.EXPECT().
					CreateSnippet(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.Snippet{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "BadRequestInvalidAccountID",
			body: gin.H{
				"account_id": 0,
				"title":      snippet.Title,
				"content":    snippet.Content,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			store := mockdb.NewMockStore(controller)
			tc.buildStubs(store)

			server := NewTestServer(t, store)
			recorder := httptest.NewRecorder()

			//to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/accounts/snippet"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetSnippetAPI(t *testing.T) {
	snippet := randomSnippet()

	testCases := []struct {
		name          string
		snippetId     int32
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			snippetId: snippet.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetSnippet(gomock.Any(), gomock.Eq(snippet.ID)).Times(1).Return(snippet, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchSnippet(t, recorder.Body, snippet)
			},
		},
		{
			name:      "NotFound",
			snippetId: snippet.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetSnippet(gomock.Any(), gomock.Eq(snippet.ID)).Times(1).Return(db.Snippet{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			snippetId: snippet.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetSnippet(gomock.Any(), gomock.Eq(snippet.ID)).Times(1).Return(db.Snippet{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "BadRequestInvalidID",
			snippetId: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetSnippet(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			store := mockdb.NewMockStore(controller)
			tc.buildStubs(store)

			server := NewTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/snippet/%d", tc.snippetId)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}

func TestDeleteSnippetAPI(t *testing.T) {
	snippet := randomSnippet()

	testCases := []struct {
		name          string
		snippetID     int32
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			snippetID: snippet.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteSnippet(gomock.Any(), gomock.Eq(snippet.ID)).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			snippetID: snippet.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteSnippet(gomock.Any(), gomock.Eq(snippet.ID)).Times(1).Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "BadRequestInvalidID",
			snippetID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteSnippet(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			store := mockdb.NewMockStore(controller)
			tc.buildStubs(store)

			server := NewTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/snippet/%d", tc.snippetID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}

func TestListSnippets(t *testing.T) {
	user, _ := createRandomUser(t)
	account := randomAccount(user.Name)

	n := 5
	snippets := make([]db.Snippet, n)
	for i := 0; i < n; i++ {
		snippets[i] = randomSnippet()
	}

	type Query struct {
		accountID int
		pageID    int
		pageSize  int
	}

	testCases := []struct {
		name          string
		query         Query
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				accountID: int(account.ID),
				pageID:    1,
				pageSize:  n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListSnippetsParams{
					AccountID: int32(account.ID),
					Limit:     int32(n),
					Offset:    0,
				}
				store.EXPECT().
					ListSnippets(gomock.Any(), gomock.Eq(arg)).
					Times(1).Return(snippets, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchSnippets(t, recorder.Body, snippets)
			},
		},
		{
			name: "InternalError",
			query: Query{
				accountID: int(account.ID),
				pageID:    1,
				pageSize:  n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListSnippets(gomock.Any(), gomock.Any()).
					Times(1).Return([]db.Snippet{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			query: Query{
				accountID: int(account.ID),
				pageID:    -1,
				pageSize:  n,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListSnippets(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			query: Query{
				accountID: int(account.ID),
				pageID:    1,
				pageSize:  100000,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListSnippets(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			store := mockdb.NewMockStore(controller)
			tc.buildStubs(store)

			server := NewTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/accounts/snippet"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("account_id", fmt.Sprintf("%d", tc.query.accountID))
			q.Add("page_id", fmt.Sprintf("%d", tc.query.pageID))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}

func randomSnippet() db.Snippet {
	return db.Snippet{
		ID:        int32(util.RandomInt(1, 1000)),
		AccountID: int32(util.RandomInt(1, 1000)),
		Title:     util.RandomString(5),
		Content:   util.RandomString(10),
	}
}

func requireBodyMatchSnippet(t *testing.T, body *bytes.Buffer, snippet db.Snippet) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotSnippet db.Snippet
	err = json.Unmarshal(data, &gotSnippet)
	require.NoError(t, err)
	require.Equal(t, snippet, gotSnippet)
}

func requireBodyMatchSnippets(t *testing.T, body *bytes.Buffer, snippets []db.Snippet) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotSnippets []db.Snippet
	err = json.Unmarshal(data, &gotSnippets)
	require.NoError(t, err)
	require.Equal(t, snippets, gotSnippets)
}
