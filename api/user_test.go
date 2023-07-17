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

func TestGetUserAPI(t *testing.T) {
	user := randomUser()

	testCases := []struct {
		name          string
		userId        int32
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userId: user.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.ID)).Times(1).Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, user)
			},
		},
		{
			name:   "NotFound",
			userId: user.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.ID)).Times(1).Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InternalError",
			userId: user.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.ID)).Times(1).Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "BadRequestInvalidID",
			userId: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(0)
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
			//создание заглушек
			//store.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.ID)).Times(1).Return(user, nil)
			tc.buildStubs(store)

			//start test server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/users/%d", tc.userId)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			//check response
			// require.Equal(t, http.StatusOK, recorder.Code)
			// requireBodyMatchAccount(t, recorder.Body, user)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestCreateUserAPI(t *testing.T) {
	user := randomUser()

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"login":    user.Login,
				"username": user.Username,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Login:    user.Login,
					Username: user.Username,
				}

				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Eq(arg)).Times(1).Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, user)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"login":    user.Login,
				"username": user.Username,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Login:    user.Login,
					Username: user.Username,
				}

				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "BadRequestInvalidLogin",
			body: gin.H{
				"login":    "",
				"username": "",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).Times(0)
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

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			//to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			fmt.Println("req", request)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}

func TestDeleteUserAPI(t *testing.T) {
	user := randomUser()

	testCases := []struct {
		name          string
		userID        int32
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
		//buildStubs1    func(store *mockdb.MockStore)
		//checkResponse1 func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userID: user.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteUser(gomock.Any(), gomock.Eq(user.ID)).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "InternalError",
			userID: user.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteUser(gomock.Any(), gomock.Eq(user.ID)).Times(1).Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "BadRequestInvalidID",
			userID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteUser(gomock.Any(), gomock.Any()).Times(0)
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

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/users/%d", tc.userID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}

// func TestUpdateUserAPI(t *testing.T) {
// 	user := randomUser()

// 	testCases := []struct {
// 		name          string
// 		userID        int32
// 		body          gin.H
// 		buildStubs    func(store *mockdb.MockStore)
// 		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name:   "OK",
// 			userID: user.ID,
// 			body: gin.H{
// 				"id":       user.ID,
// 				"username": user.Username,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				arg := db.UpdateUserParams{
// 					ID:       user.ID,
// 					Username: user.Username,
// 				}
// 				store.EXPECT().UpdateUser(gomock.Any(), gomock.Eq(arg)).Times(1).Return(user, nil)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 				//requireBodyMatchAccount(t, recorder.Body, user)
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			controller := gomock.NewController(t)
// 			defer controller.Finish()

// 			store := mockdb.NewMockStore(controller)
// 			tc.buildStubs(store)

// 			server := NewServer(store)
// 			recorder := httptest.NewRecorder()

// 			//to JSON
// 			data, err := json.Marshal(tc.body)
// 			require.NoError(t, err)

// 			//url := fmt.Sprintf("/users/%d", tc.userID)
// 			url := "/users"
// 			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
// 			require.NoError(t, err)

// 			server.router.ServeHTTP(recorder, request)

// 			tc.checkResponse(t, recorder)
// 		})
// 	}
// }

func TestUpdateUserAPI(t *testing.T) {
	user := randomUser()

	controller := gomock.NewController(t)
	defer controller.Finish()

	store := mockdb.NewMockStore(controller)

	arg := db.UpdateUserParams{
		ID:       int32(user.ID),
		Username: "bleat",
	}

	store.EXPECT().UpdateUser(gomock.Any(), gomock.Eq(arg)).Times(1).Return(user, nil)

	server := NewServer(store)
	recorder := httptest.NewRecorder()

	data, err := json.Marshal(arg)
	require.NoError(t, err)

	url := "/users"
	request, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(data))
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)

	fmt.Println("req", request)

	require.Equal(t, http.StatusOK, recorder.Code)
}

// func TestUpdateUser(t *testing.T) {
// 	user := randomUser()
// 	// user := db.User{
// 	// 	ID:       999,
// 	// 	Login:    "Yeban",
// 	// 	Username: "Sykaa",
// 	// 	Created:  time.Now().Local(),
// 	// }

// 	testCases := []struct {
// 		name          string
// 		body          gin.H
// 		buildStubs    func(store *mockdb.MockStore)
// 		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "OK",
// 			body: gin.H{
// 				"id":       user.ID,
// 				"username": user.Username,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				arg := db.UpdateUserParams{
// 					ID: int32(user.ID),
// 					//login:    user.Login,
// 					Username: user.Username,
// 				}
// 				store.EXPECT().UpdateUser(gomock.Any(), gomock.Eq(arg)).Times(1).Return(user, nil)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 				//requireBodyMatchAccount(t, recorder.Body, user)
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			store := mockdb.NewMockStore(ctrl)
// 			tc.buildStubs(store)

// 			server := NewServer(store)
// 			recorder := httptest.NewRecorder()

// 			//to JSON
// 			data, err := json.Marshal(tc.body)
// 			require.NoError(t, err)

// 			url := "/users/up"
// 			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
// 			require.NoError(t, err)

// 			server.router.ServeHTTP(recorder, request)

// 			tc.checkResponse(t, recorder)
// 		})
// 	}
// }

func randomUser() db.User {
	return db.User{
		ID:       int32(util.RandomInt(1, 1000)),
		Login:    util.RandomUser(),
		Username: util.RandomUser(),
		//Created:  time.Now().Local(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)
	require.Equal(t, user, gotUser)
}
